package actions

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	libs "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/violations"
)

type Action interface {
	DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string
}

type SingleReplicaAction struct {
	violations.Violation
}

type CapabilitiesAction struct {
	violations.Violation
}

type PrivilegedAction struct {
	violations.Violation
}

type HostVolumesAction struct {
	violations.Violation
}

type ImageSizeAction struct {
	violations.Violation
}

type ImageRepoAction struct {
	violations.Violation
}

type IngressAction struct {
	violations.Violation
}

type RequiredPodAction struct {
	violations.Violation
}

type RequiredPodAnnotationAction struct {
	violations.Violation
}

type RequiredPodLabelAction struct {
	violations.Violation
}

type RequiredNamespaceAction struct {
	violations.Violation
}

type RequiredNamespaceAnnotationAction struct {
	violations.Violation
}

type RequiredNamespaceLabelAction struct {
	violations.Violation
}

type RequiredDeploymentAction struct {
	violations.Violation
}

type RequiredDeploymentAnnotationAction struct {
	violations.Violation
}

type RequiredDeploymentLabelAction struct {
	violations.Violation
}

type RequiredDaemonSetAction struct {
	violations.Violation
}

type RequiredDaemonSetAnnotationAction struct {
	violations.Violation
}

type RequiredDaemonSetLabelAction struct {
	violations.Violation
}

type RequiredResourceQuotaAction struct {
	violations.Violation
}

type NoOwnerAction struct {
	violations.Violation
}

// action for containers with extra capablities.
func (a CapabilitiesAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Extra Capabilities", a.Violation.Source, a.Type)
}

// Action for privileged mode containers
func (a PrivilegedAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Privileged Mode", a.Violation.Source, a.Type)
}

// Action for any pod with a hostVolume
func (a HostVolumesAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Host Volumes Mounted", a.Violation.Source, a.Type)
}

// action for pods with single replica , currently action is supressed.
func (a SingleReplicaAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processSupressedAction(entity, vEntity, lastActions, "Single Replica", a.Source, a.Type)
}

// action for a container with a big image size
func (a ImageSizeAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processSupressedAction(entity, vEntity, lastActions, "Invalid Image Size", a.Source, a.Type)
}

// action for invalid repo for an image
func (a ImageRepoAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Invalid Image Repo", a.Violation.Source, a.Type)
}

// action for ingress, a special kind that we don't warn.
func (a IngressAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	// While in safe mode last warning = false
	actMessage := createActionMessage(vEntity.Namespace, reflect.TypeOf(entity).Name(), vEntity.Name, "Invalid Ingress", a.Violation.Source, len(lastActions["notify"]), libs.Cfg.ActionSafeMode == false)
	NotifyOfViolation(actMessage)
	if libs.Cfg.ActionSafeMode == false {
		entity.DoAction()
	} else {
		libs.Log.Debug("Skipping action for ", vEntity.Name, " ", a.Type, " due to safe mode.")
		return []string{"notify"}
	}

	return []string{"notify", "entity_action"}
}

// action for missing mandatory namespace
func (a RequiredNamespaceAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Missing required namespace", a.Violation.Source, a.Type)
}

// action for missing namespace annotation
func (a RequiredNamespaceAnnotationAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Missing namespace annotation", a.Violation.Source, a.Type)
}

// action for missing namespace label
func (a RequiredNamespaceLabelAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Missing namespace label", a.Violation.Source, a.Type)
}

// action for missing mandatory deployment
func (a RequiredDeploymentAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Missing required deployment", a.Violation.Source, a.Type)
}

// action for missing namespace annotation
func (a RequiredDeploymentAnnotationAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Missing deployment annotation", a.Violation.Source, a.Type)
}

// action for missing namespace label
func (a RequiredDeploymentLabelAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Missing deployment label", a.Violation.Source, a.Type)
}

// action for missing mandatory pod
func (a RequiredPodAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Missing required pod", a.Violation.Source, a.Type)
}

// action for missing pod annotation
func (a RequiredPodAnnotationAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Missing pod annotation", a.Violation.Source, a.Type)
}

// action for missing pod label
func (a RequiredPodLabelAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Missing pod label", a.Violation.Source, a.Type)
}

// action for missing mandatory daemonset
func (a RequiredDaemonSetAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Missing required daemonset", a.Violation.Source, a.Type)
}

// action for missing daemonset annotation
func (a RequiredDaemonSetAnnotationAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Missing daemonset annotation", a.Violation.Source, a.Type)
}

// action for missing daemonset label
func (a RequiredDaemonSetLabelAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Missing daemonset label", a.Violation.Source, a.Type)
}

// action for missing mandatory resourcequota
func (a RequiredResourceQuotaAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "Missing required resourcequota", a.Violation.Source, a.Type)
}

// action for missing owner
func (a NoOwnerAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	return processAction(entity, vEntity, lastActions, "No owner", a.Violation.Source, a.Type)
}

func ConvertActionableEntityToViolatableEntity(entity ActionableEntity) (libs.ViolatableEntity, error) {

	var vEntity libs.ViolatableEntity

	switch t := entity.(type) {
	case ActionDeployment:
		vEntity = t.ViolatableEntity
		break
	case ActionNamespace:
		vEntity = t.ViolatableEntity
		break
	case ActionDaemonSet:
		vEntity = t.ViolatableEntity
		break
	case ActionPod:
		vEntity = t.ViolatableEntity
		break
	case ActionIngress:
		vEntity = t.ViolatableEntity
		break
	case ActionJob:
		vEntity = t.ViolatableEntity
		break
	case ActionCronJob:
		vEntity = t.ViolatableEntity
		break
	default:
		return libs.ViolatableEntity{}, errors.New(fmt.Sprintf("Unknown Actionable Entity Type %s", t))
	}

	return vEntity, nil
}

func processAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time, violationMessage string, violationSource string, violationType violations.ViolationType) []string {
	lastTimeWarned, doIt := getLastTimeWarnedAndifToDoAction(lastActions)
	if doIt {
		entity.DoAction()
		return []string{"entity_action"}
	}

	if canSkipNotification(lastTimeWarned) {
		libs.Log.Debug("Skipping notification for ", vEntity.Name, " ", violationType, " it was notified less than ", libs.Cfg.DurationBetweenNotifyingAgain, " ago.")
		return []string{}
	}

	aMessage := createActionMessage(vEntity.Namespace, reflect.TypeOf(entity).Name(), vEntity.Name, violationMessage, violationSource, len(lastActions["notify"]), isLastWarning(lastActions))
	NotifyOfViolation(aMessage)
	return []string{"notify"}

}

func processSupressedAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time, violationMessage string, violationSource string, violationType violations.ViolationType) []string {
	lastTimeWarned, _ := getLastTimeWarnedAndifToDoAction(lastActions)

	if canSkipNotification(lastTimeWarned) {
		libs.Log.Debug("Skipping notification for ", vEntity.Name, " ", violationType, " it was notified less than ", libs.Cfg.DurationBetweenNotifyingAgain, " ago.")
		return []string{}
	}

	aMessage := createActionMessage(vEntity.Namespace, reflect.TypeOf(entity).Name(), vEntity.Name, violationMessage, violationSource, len(lastActions["notify"]), isLastWarning(lastActions))
	NotifyOfViolation(aMessage)
	return []string{"notify"}

}

func createActionMessage(namespace string, entityType string, sourceName string, violationType string, violationSource string, warningCount int, lastWarning bool) actionMessage {

	aMessage := actionMessage{
		Namespace:       namespace,
		Cluster:         libs.Cfg.ClusterName,
		EntityType:      strings.Replace(entityType, "Action", "", 1) + " Name",
		EntitySource:    sourceName,
		ViolationType:   violationType,
		ViolationSource: violationSource,
		WarningCount:    warningCount + 1,
		// There will be no last warning in safe mode.
		LastWarning: lastWarning,
	}

	return aMessage

}

func canSkipNotification(lastTimeWarned time.Time) bool {
	return lastTimeWarned.IsZero() == false && time.Now().Sub(lastTimeWarned) < libs.Cfg.DurationBetweenNotifyingAgain

}

func getLastTimeWarnedAndifToDoAction(lastActions map[string][]time.Time) (time.Time, bool) {
	var lastTimeWarned time.Time
	var doAction = false

	if t, ok := lastActions["notify"]; ok {
		lastTimeWarned = t[len(t)-1]
		if libs.Cfg.ActionSafeMode == false {
			if len(t) >= libs.Cfg.WarningCountBeforeAction {
				doAction = true
			}
		}
	}
	return lastTimeWarned, doAction
}

func isLastWarning(lastActions map[string][]time.Time) bool {
	return libs.Cfg.ActionSafeMode == false && len(lastActions["notify"]) >= libs.Cfg.WarningCountBeforeAction-1
}
