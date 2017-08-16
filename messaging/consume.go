package messaging

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	libs "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/messaging/kafka"

	"encoding/json"
	"reflect"

	"github.com/k8guard/k8guard-action/actions"

	"github.com/k8guard/k8guard-action/db"

	"github.com/k8guard/k8guardlibs/violations"
)

func ConsumeMessages() {

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, syscall.SIGTERM)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	topic := libs.Cfg.KafkaActionTopic

	master, err := kafka.NewConsumer(kafka.ACTION_CLIENTID, libs.Cfg)
	if err != nil {
		panic(err)

	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	partitions, _ := master.Partitions(topic)

	messages := make(chan *sarama.ConsumerMessage)
	for _, partition := range partitions {

		libs.Log.Info("Creating Consumer ", topic, " on partition ", partition)
		consumer, err := master.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}

		go func(consumer sarama.PartitionConsumer) {
			for {
				select {
				case err := <-consumer.Errors():
					fmt.Println(err)
				case msg := <-consumer.Messages():
					messages <- msg
				case <-signals:
					libs.Log.Debug("Interrupt is detected")
				}
			}
		}(consumer)
	}

	libs.Log.Info("Waiting for messages")
	for {
		message := <-messages
		parseViolationMessage(message)
	}
}

func parseViolationMessage(msg *sarama.ConsumerMessage) {
	libs.Log.Info("Taking Action ...")

	messageData := map[string]interface{}{}
	err := json.Unmarshal(msg.Value, &messageData)
	if err != nil {
		libs.Log.Fatal(err)
	}

	dataBytes, _ := json.Marshal(messageData["data"])

	entityViolations := []violations.Violation{}
	var actionableEntity actions.ActionableEntity

	switch messageData["kind"] {
	case string(kafka.POD_MESSAGE):
		libs.Log.Debug("Parsing Pod Message")

		pod := actions.ActionPod{}
		json.Unmarshal(dataBytes, &pod)
		actionableEntity = pod

		entityViolations = append(entityViolations, pod.Violations...)

		break
	case string(kafka.DEPLOYMENT_MESSAGE):
		libs.Log.Debug("Parsing Deployment Message")

		deployment := actions.ActionDeployment{}
		json.Unmarshal(dataBytes, &deployment)
		actionableEntity = deployment

		entityViolations = append(entityViolations, deployment.Violations...)

		break
	case string(kafka.DAEMONSET_MESSAGE):
		libs.Log.Debug("Parsing DaemonSet Message")

		daemonSet := actions.ActionDaemonSet{}
		json.Unmarshal(dataBytes, &daemonSet)
		actionableEntity = daemonSet

		entityViolations = append(entityViolations, daemonSet.Violations...)

		break
	case string(kafka.INGRESS_MESSAGE):
		libs.Log.Debug("Parsing Ingress Message")

		ingress := actions.ActionIngress{}
		json.Unmarshal(dataBytes, &ingress)
		actionableEntity = ingress

		entityViolations = append(entityViolations, ingress.Violations...)

		break
	case string(kafka.JOB_MESSAGE):
		libs.Log.Debug("Parsing Job Message")

		job := actions.ActionJob{}
		json.Unmarshal(dataBytes, &job)
		actionableEntity = job

		entityViolations = append(entityViolations, job.Violations...)

		break
	case string(kafka.CRONJOB_MESSAGE):
		libs.Log.Debug("Parsing CronJob Message")

		cronjob := actions.ActionCronJob{}
		json.Unmarshal(dataBytes, &cronjob)
		actionableEntity = cronjob

		entityViolations = append(entityViolations, cronjob.Violations...)
		break
	default:
		libs.Log.Error("Unknown Message Kind: ", messageData["kind"])
		return
	}

	for _, violation := range entityViolations {

		vEntity, err := actions.ConvertActionableEntityToViolatableEntity(actionableEntity)
		if err != nil {
			libs.Log.Fatal(err)
		}

		// Insert violation into log
		db.InsertVLOGRow(vEntity, violation, reflect.TypeOf(actionableEntity).Name())
		action := createAction(violation)
		vActionRow := db.SelectVActionRow(vEntity, violation, reflect.TypeOf(actionableEntity).Name())
		doneActions := actions.DoAction(action, actionableEntity, vEntity, vActionRow.Actions, libs.Cfg.ActionDryRun)

		if len(doneActions) == 0 {
			// If we did no actions don't insert anything
			return
		}

		for actionName, t := range doneActions {

			// Insert action into log
			db.InsertActionLogRow(vEntity.Namespace, reflect.TypeOf(actionableEntity).Name(), vEntity.Name, string(violation.Type), violation.Source, actionName)

			if _, ok := vActionRow.Actions[actionName]; ok {
				vActionRow.Actions[actionName] = append(vActionRow.Actions[actionName], t...)
			} else {
				vActionRow.Actions[actionName] = t
			}
		}

		// Insert violation state
		db.InsertVactionRow(vEntity.Namespace, reflect.TypeOf(actionableEntity).Name(), vEntity.Name, string(violation.Type), violation.Source, vActionRow.Actions)

	}

}

func createAction(violation violations.Violation) actions.Action {
	var action actions.Action

	switch vType := violation.Type; vType {
	case violations.SINGLE_REPLICA_TYPE:
		action = actions.SingleReplicaAction{Violation: violation}
		break
	case violations.IMAGE_SIZE_TYPE:
		action = actions.ImageSizeAction{Violation: violation}
		break
	case violations.IMAGE_REPO_TYPE:
		action = actions.ImageRepoAction{Violation: violation}
		break
	case violations.INGRESS_HOST_INVALID_TYPE:
		action = actions.IngressAction{Violation: violation}
		break
	case violations.CAPABILITIES_TYPE:
		action = actions.CapabilitiesAction{Violation: violation}
		break
	case violations.PRIVILEGED_TYPE:
		action = actions.PrivilegedAction{Violation: violation}
		break
	case violations.HOST_VOLUMES_TYPE:
		action = actions.HostVolumesAction{Violation: violation}
		break
	case violations.REQUIRED_POD_ANNOTATIONS_TYPE:
		action = actions.RequiredPodAnnotationAction{Violation: violation}
		break
	case violations.REQUIRED_DAEMONSETS_TYPE:
		action = actions.RequiredDaemonSetAction{Violation: violation}
		break
	default:
		libs.Log.Fatal("Unknown Violation Type ", vType)
	}

	return action
}
