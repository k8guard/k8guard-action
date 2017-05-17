package db

import (
	"github.com/k8guard/k8guard-action/db/stmts"

	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gocql/gocql"
	libs "github.com/k8guard/k8guardlibs"
)

// Here we'll store connection
var Sess *gocql.Session

// This is wrapper for gocql
func Connect(hosts []string) error {

	log.Info("Connecting to db")

	// Creating cluster for cassandra
	cluster := gocql.NewCluster(hosts...)
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = time.Second * 5

	// Initializing session
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	// Storring cassandra session
	Sess = session

	// Initializing database
	err = initDB()
	if err != nil {
		return err
	}

	return nil
}

// This we need to initialize database scheme
func initDB() error {

	log.Info("Initing DB")

	// Creating KEYSPACE (if not exists)
	err := Sess.Query(fmt.Sprintf(stmts.CREATE_KEYSPACE, libs.Cfg.CassandraKeyspace)).Exec()
	if err != nil {
		return err
	}

	err = Sess.Query(fmt.Sprintf(stmts.CREATE_VACTION_TABLE, libs.Cfg.CassandraKeyspace)).Exec()
	if err != nil {
		return err
	}

	err = Sess.Query(fmt.Sprintf(stmts.CREATE_VIOLATION_LOG_TABLE, libs.Cfg.CassandraKeyspace)).Exec()
	if err != nil {
		return err
	}

	err = Sess.Query(fmt.Sprintf(stmts.CREATE_ACTION_LOG_NAMESPACE_TYPE_TABLE, libs.Cfg.CassandraKeyspace)).Exec()
	if err != nil {
		return err
	}
	err = Sess.Query(fmt.Sprintf(stmts.CREATE_ACTION_LOG_TYPE_TABLE, libs.Cfg.CassandraKeyspace)).Exec()
	if err != nil {
		return err
	}
	err = Sess.Query(fmt.Sprintf(stmts.CREATE_ACTION_LOG_VTYPE_TABLE, libs.Cfg.CassandraKeyspace)).Exec()
	if err != nil {
		return err
	}
	err = Sess.Query(fmt.Sprintf(stmts.CREATE_ACTION_LOG_ACTION_TABLE, libs.Cfg.CassandraKeyspace)).Exec()
	if err != nil {
		return err
	}

	return nil
}
