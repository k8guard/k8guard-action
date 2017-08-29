package messaging

import (
	libs "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/messaging"
	"github.com/k8guard/k8guardlibs/messaging/types"

	"encoding/json"
	"reflect"

	"github.com/k8guard/k8guard-action/actions"

	"github.com/k8guard/k8guard-action/db"

	"github.com/k8guard/k8guardlibs/violations"
)

func ConsumeMessages() {

	c, err := messaging.CreateMessageConsumer(
		types.MessageBrokerType(libs.Cfg.MessageBroker), types.ACTION_CLIENTID, libs.Cfg)
	if err != nil {
		panic(err)
	}

	defer func() {
		c.Close()
	}()

	messages := make(chan []byte)
	c.ConsumeMessages(messages)

	libs.Log.Info("Waiting for messages")
	for {
		message := <-messages
		parseViolationMessage(message)
	}
}

func parseViolationMessage(msg []byte) {
	libs.Log.Info("Processing violation message ...")

	messageData := map[string]interface{}{}
	err := json.Unmarshal(msg, &messageData)
	if err != nil {
		libs.Log.Fatal(err)
	}

	dataBytes, _ := json.Marshal(messageData["data"])

	entityViolations := []violations.Violation{}
	var actionableEntity actions.ActionableEntity

	switch messageData["kind"] {
	case string(types.POD_MESSAGE):
		libs.Log.Debug("Parsing Pod Message")

		pod := actions.ActionPod{}
		json.Unmarshal(dataBytes, &pod)
		actionableEntity = pod

		entityViolations = append(entityViolations, pod.Violations...)

		break
	case string(types.NAMESPACE_MESSAGE):
		libs.Log.Debug("Parsing Namespace Message")

		namespace := actions.ActionNamespace{}
		json.Unmarshal(dataBytes, &namespace)
		actionableEntity = namespace

		entityViolations = append(entityViolations, namespace.Violations...)

		break
	case string(types.DEPLOYMENT_MESSAGE):
		libs.Log.Debug("Parsing Deployment Message")

		deployment := actions.ActionDeployment{}
		json.Unmarshal(dataBytes, &deployment)
		actionableEntity = deployment

		entityViolations = append(entityViolations, deployment.Violations...)

		break
	case string(types.DAEMONSET_MESSAGE):
		libs.Log.Debug("Parsing DaemonSet Message")

		daemonSet := actions.ActionDaemonSet{}
		json.Unmarshal(dataBytes, &daemonSet)
		actionableEntity = daemonSet

		entityViolations = append(entityViolations, daemonSet.Violations...)

		break
	case string(types.INGRESS_MESSAGE):
		libs.Log.Debug("Parsing Ingress Message")

		ingress := actions.ActionIngress{}
		json.Unmarshal(dataBytes, &ingress)
		actionableEntity = ingress

		entityViolations = append(entityViolations, ingress.Violations...)

		break
	case string(types.JOB_MESSAGE):
		libs.Log.Debug("Parsing Job Message")

		job := actions.ActionJob{}
		json.Unmarshal(dataBytes, &job)
		actionableEntity = job

		entityViolations = append(entityViolations, job.Violations...)

		break
	case string(types.CRONJOB_MESSAGE):
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
	switch vType := violation.Type; vType {
	case violations.SINGLE_REPLICA_TYPE:
		return actions.SingleReplicaAction{Violation: violation}
	case violations.IMAGE_SIZE_TYPE:
		return actions.ImageSizeAction{Violation: violation}
	case violations.IMAGE_REPO_TYPE:
		return actions.ImageRepoAction{Violation: violation}
	case violations.INGRESS_HOST_INVALID_TYPE:
		return actions.IngressAction{Violation: violation}
	case violations.CAPABILITIES_TYPE:
		return actions.CapabilitiesAction{Violation: violation}
	case violations.PRIVILEGED_TYPE:
		return actions.PrivilegedAction{Violation: violation}
	case violations.HOST_VOLUMES_TYPE:
		return actions.HostVolumesAction{Violation: violation}
	case violations.REQUIRED_NAMESPACES_TYPE:
		return actions.RequiredNamespaceAction{Violation: violation}
	case violations.REQUIRED_NAMESPACE_ANNOTATIONS_TYPE:
		return actions.RequiredNamespaceAnnotationAction{Violation: violation}
	case violations.REQUIRED_NAMESPACE_LABELS_TYPE:
		return actions.RequiredNamespaceLabelAction{Violation: violation}
	case violations.REQUIRED_DEPLOYMENTS_TYPE:
		return actions.RequiredDeploymentAction{Violation: violation}
	case violations.REQUIRED_DEPLOYMENT_ANNOTATIONS_TYPE:
		return actions.RequiredDeploymentAnnotationAction{Violation: violation}
	case violations.REQUIRED_DEPLOYMENT_LABELS_TYPE:
		return actions.RequiredDeploymentLabelAction{Violation: violation}
	case violations.REQUIRED_PODS_TYPE:
		return actions.RequiredPodAction{Violation: violation}
	case violations.REQUIRED_POD_ANNOTATIONS_TYPE:
		return actions.RequiredPodAnnotationAction{Violation: violation}
	case violations.REQUIRED_POD_LABELS_TYPE:
		return actions.RequiredPodLabelAction{Violation: violation}
	case violations.REQUIRED_DAEMONSETS_TYPE:
		return actions.RequiredDaemonSetAction{Violation: violation}
	case violations.REQUIRED_DAEMONSET_ANNOTATIONS_TYPE:
		return actions.RequiredDaemonSetAnnotationAction{Violation: violation}
	case violations.REQUIRED_DAEMONSET_LABELS_TYPE:
		return actions.RequiredDaemonSetLabelAction{Violation: violation}
	case violations.REQUIRED_RESOURCEQUOTA_TYPE:
		return actions.RequiredResourceQuotaAction{Violation: violation}
	case violations.NO_OWNER_ANNOTATION_TYPE:
		return actions.NoOwnerAction{Violation: violation}
	default:
		libs.Log.Fatal("Unknown Violation Type ", vType)
	}

	return nil
}
