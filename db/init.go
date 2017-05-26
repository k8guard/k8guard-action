package db

import (
	"k8guard-action/db"
	"reflect"

	"github.com/k8guard/k8guard-action/db/stmts"

	"fmt"
	"time"

	"github.com/gocql/gocql"
	libs "github.com/k8guard/k8guardlibs"
)

// Here we'll store connection
var Sess *gocql.Session
var err error

// This is wrapper for gocql
func Connect(hosts []string) error {

	libs.Log.Info("Connecting to db")

	// Creating cluster for cassandra
	cluster := gocql.NewCluster(hosts...)
	cluster.Consistency = gocql.LocalQuorum
	cluster.Timeout = time.Second * 15
	// Auth if username is set
	if libs.Cfg.CassandraUsername != "" {
		libs.Log.Debug("Connecting with username ", libs.Cfg.CassandraUsername)
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: libs.Cfg.CassandraUsername,
			Password: libs.Cfg.CassandraPassword,
		}
	}
	if libs.Cfg.CassandraCaPath != "" {
		libs.Log.Debug("Using Ca")
		cluster.SslOpts = &gocql.SslOptions{
			CaPath:                 libs.Cfg.CassandraCaPath,
			EnableHostVerification: libs.Cfg.CassandraSslHostValidation,
		}
	}

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

	libs.Log.Info("Initing DB")
	// creating keyspace if configured true
	if libs.Cfg.CassandraCreateKeyspace {
		// Creating KEYSPACE (if not exists)
		libs.Log.Info("Creating the keyspace ", libs.Cfg.CassandraKeyspace)
		err = Sess.Query(fmt.Sprintf(stmts.CREATE_KEYSPACE, libs.Cfg.CassandraKeyspace)).Exec()
		if err != nil {
			return err
		}
	} else {
		libs.Log.Info("Skipping creating the keyspace")
	}

	// creating tables if configured true
	if libs.Cfg.CassandraCreateTables {
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
	} else {
		libs.Log.Info("Skipping creating tables")
	}
	vActionRow := db.SelectVActionRow(vEntity, violation, reflect.TypeOf(actionableEntity).Name())

	return nil
}
