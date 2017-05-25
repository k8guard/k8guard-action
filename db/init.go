package db

import (
	"github.com/k8guard/k8guard-action/db/stmts"

	"fmt"

	"github.com/gocql/gocql"
	libs "github.com/k8guard/k8guardlibs"
	libsdb "github.com/k8guard/k8guardlibs/db"
)

// Here we'll store connection
var Sess *gocql.Session
var err error

func InitDB() {
	Sess = libsdb.Connect(libs.Cfg.CassandraHosts)

}

func createDB() error {

	libs.Log.Info("Initializing the DB schema ")

	// Creating KEYSPACE (if not exists)
	// some cassandras config may not allow Create keyspace
	if libs.Cfg.CassandraCreateKeyspace {

		err = Sess.Query(fmt.Sprintf(stmts.CREATE_KEYSPACE, libs.Cfg.CassandraKeyspace)).Exec()
		if err != nil {
			return err
		}
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
