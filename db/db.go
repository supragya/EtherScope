package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/Blockpour/Blockpour-Geth-Indexer/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type DBConn struct {
	isDB        bool
	conn        *sql.DB
	mq          *amqp.Channel
	doResume    bool
	resumeURL   string
	mqQueueName string
	dataTable   string
	metaTable   string
	StartBlock  uint64
	Network     string
	ChainID     uint
	store       [][]byte
}

type VersionWrapper struct {
	Version uint8
	Message any
}

func VersionWrapped(message any) VersionWrapper {
	return VersionWrapper{
		Version: version.PersistenceVersion,
		Message: message,
	}
}

func SetupConnection() (DBConn, error) {
	dbType := viper.GetString("general.persistence")

	switch dbType {
	case "mq":
		mq, err := setupRabbitMQ()
		return DBConn{isDB: false,
			conn:        nil,
			mq:          mq,
			doResume:    !viper.GetBool("mq.skipResume"),
			resumeURL:   viper.GetString("mq.resumeURL"),
			mqQueueName: viper.GetString("mq.queue"),
			dataTable:   "",
			metaTable:   "",
			StartBlock:  viper.GetUint64("general.startBlock"),
			Network:     viper.GetString("general.network"),
			ChainID:     viper.GetUint("general.chainID"),
		}, err
	default:
		break
	}

	return DBConn{}, errors.New("unsupported db: " + dbType)
}

func setupRabbitMQ() (*amqp.Channel, error) {
	var (
		host = viper.GetString("mq.host")
		port = viper.GetUint64("mq.port")
		user = viper.GetString("mq.user")
		pass = viper.GetString("mq.pass")
	)
	connPrefix := "amqp"
	if viper.GetBool("mq.secureConnection") {
		connPrefix = "amqps"
	}
	mqConnStr := fmt.Sprintf("%s://%s:%s@%s:%d/", connPrefix, user, pass, host, port)

	connectRabbitMQ, err := amqp.Dial(mqConnStr)
	if err != nil {
		return &amqp.Channel{}, err
	}

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		return &amqp.Channel{}, err
	}

	_, err = channelRabbitMQ.QueueDeclare(
		viper.GetString("mq.queue"),         // queue name
		viper.GetBool("mq.queueIsDurable"),  // durable
		viper.GetBool("mq.queueAutoDelete"), // auto delete
		viper.GetBool("mq.queueExclusive"),  // exclusive
		viper.GetBool("mq.queueNoWait"),     // no wait
		nil,                                 // arguments
	)
	if err != nil {
		return &amqp.Channel{}, err
	}
	return channelRabbitMQ, nil
}

func (d *DBConn) GetMostRecentPostedBlockHeight() uint64 {

	type ResumeAt struct {
		Data struct {
			Height int `json:"block_height"`
		} `json:"data"`
	}

	if !d.isDB {
		log.Warn("transaction support unavailable for the given persistence backend")
		if !d.doResume {
			log.Warn("resume feature skipped in non postgres database. assuming new DB")
			return d.StartBlock
		} else {
			resp, err := http.Get(d.resumeURL)
			util.ENOK(err)
			body, err := ioutil.ReadAll(resp.Body)
			var responseObject ResumeAt
			json.Unmarshal(body, &responseObject)
			ResumeBlock := int64(responseObject.Data.Height)
			util.ENOK(err)
			var respBody struct{ Resume uint64 }
			util.ENOK(json.Unmarshal(body, &respBody))
			fmt.Println(respBody)
			return uint64(ResumeBlock)
		}
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
	for _, item := range items {
		switch it := item.(type) {
		case itypes.Transfer:
			mqMessage, err := json.Marshal(VersionWrapped(it))
			util.ENOK(err)
			d.store = append(d.store, mqMessage)
		case itypes.Mint:
			mqMessage, err := json.Marshal(VersionWrapped(it))
			util.ENOK(err)
			d.store = append(d.store, mqMessage)
		case itypes.Burn:
			mqMessage, err := json.Marshal(VersionWrapped(it))
			util.ENOK(err)
			d.store = append(d.store, mqMessage)
		case itypes.Swap:
			mqMessage, err := json.Marshal(VersionWrapped(it))
			util.ENOK(err)
			d.store = append(d.store, mqMessage)
		}
	}

	// Add block synopsis
	mqMessage, err := json.Marshal(bm)
	util.ENOK(err)
	d.store = append(d.store, mqMessage)
}
