package db

import (
	"database/sql"
	"errors"
	"fmt"

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
	query := fmt.Sprintf("SELECT height FROM %s WHERE nwtype='%s' AND network=%x ORDER BY height DESC LIMIT 1",
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
