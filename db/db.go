package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DBConn struct {
	conn       *sql.DB
	dataTable  string
	metaTable  string
	StartBlock uint64
	Network    string
	ChainID    uint
}

func SetupConnection() (DBConn, error) {
	dbType := viper.GetString("db.type")

	switch dbType {
	case "postgres":
		db, err := setupPostgres()
		return DBConn{conn: db,
			dataTable:  viper.GetString("db.datatable"),
			metaTable:  viper.GetString("db.metatable"),
			StartBlock: viper.GetUint64("general.start_block"),
			Network:    viper.GetString("general.network"),
			ChainID:    viper.GetUint("general.chainid"),
		}, err
	default:
		break
	}

	return DBConn{}, errors.New("unsupported db: " + dbType)
}

func setupPostgres() (*sql.DB, error) {
	var (
		host   = viper.GetString("db.host")
		port   = viper.GetInt64("db.port")
		user   = viper.GetString("db.user")
		pass   = viper.GetString("db.pass")
		dbname = viper.GetString("db.dbname")
		ssl    = viper.GetString("db.sslmode")
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

func (d *DBConn) GetMostRecentPostedBlockHeight() uint64 {
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
	ctx := context.Background()
	tx, err := d.conn.BeginTx(ctx, nil)
	util.ENOK(err)
	return ctx, tx
}

func (d *DBConn) AddToTx(dbCtx *context.Context, dbTx *sql.Tx, items []interface{}, bm itypes.BlockSynopsis, blockHeight uint64) {
	currentTime := time.Now().Unix()
	for _, item := range items {
		query := ""
		switch item.(type) {
		case itypes.Mint:
			query = d.getQueryStringMint(item.(itypes.Mint), currentTime)
		case itypes.Burn:
			query = d.getQueryStringBurn(item.(itypes.Burn), currentTime)
		case *itypes.Swap:
			query = d.getQueryStringSwap(item.(itypes.Swap), currentTime)
		}
		_, err := dbTx.ExecContext(*dbCtx, query)
		util.ENOK(err)
	}

	// Add block synopsis
	query := d.getQueryStringBlockSynopsis(blockHeight, currentTime, bm)
	_, err := dbTx.ExecContext(*dbCtx, query)
	util.ENOK(err)
}

func (d *DBConn) getQueryStringMint(item itypes.Mint, currentTime int64) string {
	const insquery string = "INSERT INTO %s "
	const fields string = "(nwtype, network, 	time, 			  inserted_at, 		token0, token1, pair, amount0, amount1, amountusd, reserves0, reserves1, reservesusd, type, sender, transaction, slippage, height) "
	const valuesfmt string = "VALUES ('%s', %d, TO_TIMESTAMP(%d), TO_TIMESTAMP(%d), '%s',   '%s',   '%s', %f,      %f,      %f,        %f,        %f,        %f,          '%s', '%s',   '%s',        %f,       %d    );"
	return fmt.Sprintf(insquery+fields+valuesfmt, d.dataTable, // table to insert to
		d.Network,                      // nwtype
		d.ChainID,                      // network
		item.Time,                      // time
		currentTime,                    // inserted_at
		item.Token0.String()[2:],       // token0 (removed 0x prefix)
		item.Token1.String()[2:],       // token1 (removed 0x prefix)
		item.PairContract.String()[2:], // pair
		item.Amount0,                   // amount0
		item.Amount1,                   // amount1
		0.0,                            // amountusd, FIXME
		item.Reserve0,                  // reserves0
		item.Reserve1,                  // reserves1
		0.0,                            // reservesusd, FIXME
		"mint",                         // type
		"",                             // sender FIXME
		item.Transaction.String()[2:],  // transaction (removed 0x prefix)
		0.0,                            // slippage
		item.Height,                    // height
	)
}

func (d *DBConn) getQueryStringBurn(item itypes.Burn, currentTime int64) string {
	const insquery string = "INSERT INTO %s "
	const fields string = "(nwtype, network, 	time, 			  inserted_at, 		token0, token1, pair, amount0, amount1, amountusd, reserves0, reserves1, reservesusd, type, sender, transaction, slippage, height) "
	const valuesfmt string = "VALUES ('%s', %d, TO_TIMESTAMP(%d), TO_TIMESTAMP(%d), '%s',   '%s',   '%s', %f,      %f,      %f,        %f,        %f,        %f,          '%s', '%s',   '%s',        %f,       %d    );"
	return fmt.Sprintf(insquery+fields+valuesfmt, d.dataTable, // table to insert to
		d.Network,                      // nwtype
		d.ChainID,                      // network
		item.Time,                      // time
		currentTime,                    // inserted_at
		item.Token0.String()[2:],       // token0 (removed 0x prefix)
		item.Token1.String()[2:],       // token1 (removed 0x prefix)
		item.PairContract.String()[2:], // pair
		item.Amount0,                   // amount0
		item.Amount1,                   // amount1
		0.0,                            // amountusd, FIXME
		item.Reserve0,                  // reserves0
		item.Reserve1,                  // reserves1
		0.0,                            // reservesusd, FIXME
		"burn",                         // type
		"",                             // sender FIXME
		item.Transaction.String()[2:],  // transaction (removed 0x prefix)
		0.0,                            // slippage
		item.Height,                    // height
	)
}

func (d *DBConn) getQueryStringSwap(item itypes.Swap, currentTime int64) string {
	const insquery string = "INSERT INTO %s "
	const fields string = "(nwtype, network, 	time, 			  inserted_at, 		token0, token1, pair, amount0, amount1, amountusd, reserves0, reserves1, reservesusd, type, sender, recipient, transaction, slippage, height) "
	const valuesfmt string = "VALUES ('%s', %d, TO_TIMESTAMP(%d), TO_TIMESTAMP(%d), '%s',   '%s',   '%s', %f,      %f,      %f,        %f,        %f,        %f,          '%s', '%s',   '%s',      '%s',        %f,       %d    );"
	return fmt.Sprintf(insquery+fields+valuesfmt, d.dataTable, // table to insert to
		d.Network,                      // nwtype
		d.ChainID,                      // network
		item.Time,                      // time
		currentTime,                    // inserted_at
		item.Token0.String()[2:],       // token0 (removed 0x prefix)
		item.Token1.String()[2:],       // token1 (removed 0x prefix)
		item.PairContract.String()[2:], // pair
		item.Amount0,                   // amount0
		item.Amount1,                   // amount1
		0.0,                            // amountusd, FIXME
		item.Reserve0,                  // reserves0
		item.Reserve1,                  // reserves1
		0.0,                            // reservesusd, FIXME
		"swap",                         // type
		"",                             // sender FIXME
		"",                             // recipient FIXME
		item.Transaction.String()[2:],  // transaction (removed 0x prefix)
		0.0,                            // slippage
		item.Height,                    // height
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
