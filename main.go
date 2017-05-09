package main

import (
	"k8guard-action/db"
	"k8guard-action/messaging"

	libs "github.com/k8guard/k8guardlibs"
)

func main() {
	libs.Log.Info("Hello From k8guard-action")

	err := db.Connect(libs.Cfg.CassandraHosts)
	if err != nil {
		panic(err.Error())
	}
	messaging.ConsumeMessages()

}
