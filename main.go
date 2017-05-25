package main

import (
	"github.com/k8guard/k8guard-action/messaging"

	libs "github.com/k8guard/k8guardlibs"

	"github.com/revel/modules/db/app"
)

func main() {
	libs.Log.Info("Hello From k8guard-action")
	db.Init()
	messaging.ConsumeMessages()

}
