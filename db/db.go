package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type DBConn struct {
	isDB        bool
	conn        *sql.DB
	mq          *amqp.Channel
	mqQueueName string
	dataTable   string
	metaTable   string
	StartBlock  uint64
	Network     string
	ChainID     uint
	store       [][]byte
}

var zeroFloat = big.NewFloat(0.0)

func SetupConnection() (DBConn, error) {
	dbType := viper.GetString("general.persistence")

	switch dbType {
	case "postgres":
		db, err := setupPostgres()
		return DBConn{isDB: true,
			conn:        db,
			mq:          nil,
			mqQueueName: "",
			dataTable:   viper.GetString("postgres.datatable"),
			metaTable:   viper.GetString("postgres.metatable"),
			StartBlock:  viper.GetUint64("general.startBlock"),
			Network:     viper.GetString("general.network"),
			ChainID:     viper.GetUint("general.chainid"),
		}, err
	case "mq":
		mq, err := setupRabbitMQ()
		return DBConn{isDB: false,
			conn:        nil,
			mq:          mq,
			mqQueueName: viper.GetString("mq.queue"),
			dataTable:   "",
			metaTable:   "",
			StartBlock:  viper.GetUint64("general.startBlock"),
			Network:     viper.GetString("general.network"),
			ChainID:     viper.GetUint("general.chainid"),
		}, err
	default:
		break
	}

	return DBConn{}, errors.New("unsupported db: " + dbType)
}

func setupPostgres() (*sql.DB, error) {
	var (
		host   = viper.GetString("postgres.host")
		port   = viper.GetUint64("postgres.port")
		user   = viper.GetString("postgres.user")
		pass   = viper.GetString("postgres.pass")
		dbname = viper.GetString("postgres.dbname")
		ssl    = viper.GetString("postgres.sslmode")
	)

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, ssl)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return &sql.DB{}, err
	}

	err = db.Ping()
	if err != nil {
		return &sql.DB{}, err
	}

	log.Info("connected to the postgres database")
	return db, nil
}

func setupRabbitMQ() (*amqp.Channel, error) {
	var (
		host = viper.GetString("mq.host")
		port = viper.GetUint64("mq.port")
		user = viper.GetString("mq.user")
		pass = viper.GetString("mq.pass")
	)
	mqConnStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", user, pass, host, port)

	connectRabbitMQ, err := amqp.Dial(mqConnStr)
	if err != nil {
		return &amqp.Channel{}, err
	}

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		return &amqp.Channel{}, err
	}

	_, err = channelRabbitMQ.QueueDeclare(
		viper.GetString("mq.queue"), // queue name
		true,                        // durable
		false,                       // auto delete
		false,                       // exclusive
		false,                       // no wait
		nil,                         // arguments
	)
	if err != nil {
		return &amqp.Channel{}, err
	}
	return channelRabbitMQ, nil
}

func (d *DBConn) GetMostRecentPostedBlockHeight() uint64 {
	if !d.isDB {
		log.Warn("Resume feature unavailable in non postgres database. Assuming new DB")
		log.Warn("Transaction support unavailable for the given persistence backend")
		return d.StartBlock
	}

	query := fmt.Sprintf("SELECT height FROM %s WHERE nwtype='%s' AND network=%d ORDER BY height DESC LIMIT 1",
		d.metaTable, d.Network, d.ChainID)

	rows, err := d.conn.Query(query)
	util.ENOK(err)
	defer rows.Close()

	mostRecent := d.StartBlock - 1
	foundRow := false
	for rows.Next() {
		err = rows.Scan(&mostRecent)
		util.ENOK(err)
		foundRow = true
	}

	if !foundRow {
		log.Warn("no recent blocks found in db. assuming new db")
	}
	return mostRecent
}

func (d *DBConn) BeginTx() (context.Context, *sql.Tx) {
	if !d.isDB {
		return nil, nil
	}
	ctx := context.Background()
	tx, err := d.conn.BeginTx(ctx, nil)
	util.ENOK(err)
	return ctx, tx
}

func (d *DBConn) CommitTx(dbTx *sql.Tx) error {
	if d.isDB {
		return dbTx.Commit()
	}
	// In case of mq
	for _, item := range d.store {
		err := d.mq.Publish(
			"",            // exchange
			d.mqQueueName, // queue name
			false,         // mandatory
			false,         // immediate
			amqp.Publishing{
				ContentType:     "application/json",
				ContentEncoding: "application/json",
				Timestamp:       time.Now(),
				Body:            item,
			}, // message to publish
		)
		if err != nil {
			return err
		}
	}
	// reset the store
	d.store = [][]byte{}
	return nil
}

func (d *DBConn) AddToTx(dbCtx *context.Context, dbTx *sql.Tx, items []interface{}, bm itypes.BlockSynopsis, blockHeight uint64) {
	currentTime := time.Now().Unix()
	for _, item := range items {
		query := ""
		switch it := item.(type) {
		case itypes.Mint:
			if d.isDB {
				query = d.getQueryStringMint(it, currentTime)
			} else {
				mqMessage, err := json.Marshal(it)
				util.ENOK(err)
				d.store = append(d.store, mqMessage)
			}
		case itypes.Burn:
			if d.isDB {
				query = d.getQueryStringBurn(it, currentTime)
			} else {
				mqMessage, err := json.Marshal(it)
				util.ENOK(err)
				d.store = append(d.store, mqMessage)
			}
		case itypes.Swap:
			if d.isDB {
				query = d.getQueryStringSwap(it, currentTime)
			} else {
				mqMessage, err := json.Marshal(it)
				util.ENOK(err)
				d.store = append(d.store, mqMessage)
			}
		}
		if d.isDB {
			_, err := dbTx.ExecContext(*dbCtx, query)
			util.ENOKF(err, query)
		}
	}

	// Add block synopsis
	if d.isDB {
		query := d.getQueryStringBlockSynopsis(blockHeight, currentTime, bm)
		_, err := dbTx.ExecContext(*dbCtx, query)
		util.ENOK(err)
	} else {
		mqMessage, err := json.Marshal(bm)
		util.ENOK(err)
		d.store = append(d.store, mqMessage)
	}
}

func (d *DBConn) getQueryStringMint(item itypes.Mint, currentTime int64) string {
	const insquery string = "INSERT INTO %s "
	const fields string = "(nwtype, network, 	time, 			  inserted_at, 		token0, token1, pair, amount0, amount1, amountusd, reserves0, reserves1, reservesusd, type, sender, transaction, slippage, height) "
	const valuesfmt string = "VALUES ('%s', %d, TO_TIMESTAMP(%d), TO_TIMESTAMP(%d), '%s',   '%s',   '%s', %f,      %f,      %f,        %f,        %f,        %f,          '%s', '%s',   '%s',        %f,       %d    );"
	return fmt.Sprintf(insquery+fields+valuesfmt, d.dataTable, // table to insert to
		d.Network,   // nwtype
		d.ChainID,   // network
		item.Time,   // time
		currentTime, // inserted_at
		strings.ToLower(item.Token0.String()[2:]),       // token0 (removed 0x prefix)
		strings.ToLower(item.Token1.String()[2:]),       // token1 (removed 0x prefix)
		strings.ToLower(item.PairContract.String()[2:]), // pair
		item.Amount0,                           // amount0
		zeroFloat,                              // amount1
		0.0,                                    // amountusd, FIXME
		item.Reserve0,                          // reserves0
		item.Reserve1,                          // reserves1
		0.0,                                    // reservesusd, FIXME
		"mint",                                 // type
		strings.ToLower(item.Sender.Hex()[2:]), // sender FIXME (removed 0x prefix)
		strings.ToLower(item.Transaction.String()[2:]), // transaction (removed 0x prefix)
		0.0,         // slippage
		item.Height, // height
	)
}

func (d *DBConn) getQueryStringBurn(item itypes.Burn, currentTime int64) string {
	const insquery string = "INSERT INTO %s "
	const fields string = "(nwtype, network, 	time, 			  inserted_at, 		token0, token1, pair, amount0, amount1, amountusd, reserves0, reserves1, reservesusd, type, sender, recipient, transaction, slippage, height) "
	const valuesfmt string = "VALUES ('%s', %d, TO_TIMESTAMP(%d), TO_TIMESTAMP(%d), '%s',   '%s',   '%s', %f,      %f,      %f,        %f,        %f,        %f,          '%s', '%s',   '%s',  '%s',      %f,       %d    );"
	return fmt.Sprintf(insquery+fields+valuesfmt, d.dataTable, // table to insert to
		d.Network,   // nwtype
		d.ChainID,   // network
		item.Time,   // time
		currentTime, // inserted_at
		strings.ToLower(item.Token0.String()[2:]),       // token0 (removed 0x prefix)
		strings.ToLower(item.Token1.String()[2:]),       // token1 (removed 0x prefix)
		strings.ToLower(item.PairContract.String()[2:]), // pair
		item.Amount0,                             // amount0
		item.Amount1,                             // amount1
		0.0,                                      // amountusd, FIXME
		item.Reserve0,                            // reserves0
		item.Reserve1,                            // reserves1
		0.0,                                      // reservesusd, FIXME
		"burn",                                   // type
		strings.ToLower(item.Sender.Hex()[2:]),   // sender FIXME (removed 0x prefix)
		strings.ToLower(item.Receiver.Hex()[2:]), // recipient (removed 0x prefix)
		strings.ToLower(item.Transaction.String()[2:]), // transaction (removed 0x prefix)
		0.0,         // slippage
		item.Height, // height
	)
}

func (d *DBConn) getQueryStringSwap(item itypes.Swap, currentTime int64) string {
	const insquery string = "INSERT INTO %s "
	const fields string = "(nwtype, network, 	time, 			  inserted_at, 		token0, token1, pair, amount0, amount1, amountusd, reserves0, reserves1, reservesusd, type, sender, recipient, transaction, slippage, height) "
	const valuesfmt string = "VALUES ('%s', %d, TO_TIMESTAMP(%d), TO_TIMESTAMP(%d), '%s',   '%s',   '%s', %f,      %f,      %f,        %f,        %f,        %f,          '%s', '%s',   '%s',      '%s',        %f,       %d    );"
	return fmt.Sprintf(insquery+fields+valuesfmt, d.dataTable, // table to insert to
		d.Network,   // nwtype
		d.ChainID,   // network
		item.Time,   // time
		currentTime, // inserted_at
		strings.ToLower(item.Token0.String()[2:]),       // token0 (removed 0x prefix)
		strings.ToLower(item.Token1.String()[2:]),       // token1 (removed 0x prefix)
		strings.ToLower(item.PairContract.String()[2:]), // pair
		item.Amount0,                             // amount0
		item.Amount1,                             // amount1
		0.0,                                      // amountusd, FIXME
		item.Reserve0,                            // reserves0
		item.Reserve1,                            // reserves1
		0.0,                                      // reservesusd, FIXME
		"swap",                                   // type
		strings.ToLower(item.Sender.Hex()[2:]),   // sender (removed 0x prefix)
		strings.ToLower(item.Receiver.Hex()[2:]), // recipient (removed 0x prefix)
		strings.ToLower(item.Transaction.String()[2:]), // transaction (removed 0x prefix)
		0.0,         // slippage
		item.Height, // height
	)
}

func (d *DBConn) getQueryStringBlockSynopsis(blockHeight uint64, currentTime int64, bm itypes.BlockSynopsis) string {
	if bm.TotalLogs != bm.MintLogs+bm.BurnLogs+bm.SwapLogs {
		log.Fatal("arithmetic error for block synopsis: ", bm)
	}
	const insquery string = "INSERT INTO %s "
	const fields string = "(nwtype, network, height, inserted_at, mint_logs, burn_logs, swap_logs, total_logs) "
	const valuesfmt string = "VALUES ('%s', %d, %d,  TO_TIMESTAMP(%d), %d,   %d,        %d,        %d);"
	return fmt.Sprintf(insquery+fields+valuesfmt, d.metaTable, // table to insert to
		d.Network,    // nwtype
		d.ChainID,    // network
		blockHeight,  // height
		currentTime,  // inserted_at
		bm.MintLogs,  // mint_logs
		bm.BurnLogs,  // burn_logs
		bm.SwapLogs,  // swap_logs
		bm.TotalLogs, // total_logs
	)
}
