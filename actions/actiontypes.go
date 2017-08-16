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

type RequiredPodAnnotationAction struct {
	violations.Violation
}

// action for containers with extra capablities.
func (a CapabilitiesAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	lastTimeWarned, doIt := getLastTimeWarnedAndifToDoAction(lastActions)
	if doIt {
		entity.DoAction()
		return []string{"entity_action"}
	}

	if canSkipNotification(lastTimeWarned) {
		libs.Log.Debug("Skipping notification for ", vEntity.Name, " ", a.Type, " it was notified less than ", libs.Cfg.DurationBetweenNotifyingAgain, " ago.")
		return []string{}
	}

	aMessage := createActionMessage(vEntity.Namespace, reflect.TypeOf(entity).Name(), vEntity.Name, "Extra Capabilities", a.Violation.Source, len(lastActions["notify"]), isLastWarning(lastActions))
	NotifyOfViolation(aMessage)
	return []string{"notify"}

}

// Action for privileged mode containers
func (a PrivilegedAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	lastTimeWarned, doIt := getLastTimeWarnedAndifToDoAction(lastActions)
	if doIt {
		entity.DoAction()
		return []string{"entity_action"}
	}
	if canSkipNotification(lastTimeWarned) {
		libs.Log.Info("Skipping notification for ", vEntity.Name, " ", a.Type, " it was notified less than ", libs.Cfg.DurationBetweenNotifyingAgain, " ago.")
		return []string{}
	}
	actMessage := createActionMessage(vEntity.Namespace, reflect.TypeOf(entity).Name(), vEntity.Name, "Privileged Mode", a.Violation.Source, len(lastActions["notify"]), isLastWarning(lastActions))
	NotifyOfViolation(actMessage)

	return []string{"notify"}

}

// Action for any pod with a hostVolume
func (a HostVolumesAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	lastTimeWarned, doIt := getLastTimeWarnedAndifToDoAction(lastActions)
	if doIt {
		entity.DoAction()
		return []string{"entity_action"}
	}
	if canSkipNotification(lastTimeWarned) {
		libs.Log.Debug("Skipping notification for ", vEntity.Name, " ", a.Type, " it was notified less than ", libs.Cfg.DurationBetweenNotifyingAgain, " ago.")
		return []string{}
	}

	actMessage := createActionMessage(vEntity.Namespace, reflect.TypeOf(entity).Name(), vEntity.Name, "Host Volumes Mounted", a.Violation.Source, len(lastActions["notify"]), isLastWarning(lastActions))
	NotifyOfViolation(actMessage)

	return []string{"notify"}

}

// action for pods with single replica , currently action is supressed.
func (a SingleReplicaAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {

	lastTimeWarned, _ := getLastTimeWarnedAndifToDoAction(lastActions)

	if canSkipNotification(lastTimeWarned) {
		libs.Log.Debug("Skipping notification for ", vEntity.Name, " ", a.Type, " it was notified less than ", libs.Cfg.DurationBetweenNotifyingAgain, " ago.")
		return []string{}
	}

	aMessage := createActionMessage(vEntity.Namespace, reflect.TypeOf(entity).Name(), vEntity.Name, "Single Replica", a.Source, len(lastActions), false)

	NotifyOfViolation(aMessage)

	return []string{"notify"}
}

// action for a container with a big image size
func (a ImageSizeAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {
	lastTimeWarned, _ := getLastTimeWarnedAndifToDoAction(lastActions)
	if canSkipNotification(lastTimeWarned) {
		libs.Log.Debug("Skipping notification for ", vEntity.Name, " ", a.Type, " it was notified less than ", libs.Cfg.DurationBetweenNotifyingAgain, " ago.")
		return []string{}
	}
	actMessage := createActionMessage(vEntity.Namespace, reflect.TypeOf(entity).Name(), vEntity.Name, "Invalid Image Size", a.Source, len(lastActions), false)
	NotifyOfViolation(actMessage)

	return []string{"notify"}
}

// action for invalid repo for an image
func (a ImageRepoAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {

	lastTimeWarned, doIt := getLastTimeWarnedAndifToDoAction(lastActions)
	if doIt {
		entity.DoAction()
		return []string{"entity_action"}
	}
	if canSkipNotification(lastTimeWarned) {
		libs.Log.Debug("Skipping notification for ", vEntity.Name, " ", a.Type, " it was notified less than ", libs.Cfg.DurationBetweenNotifyingAgain, " ago.")
		return []string{}
	}

	actMessage := createActionMessage(vEntity.Namespace, reflect.TypeOf(entity).Name(), vEntity.Name, "Invalid Image Repo", a.Violation.Source, len(lastActions["notify"]), isLastWarning(lastActions))
	NotifyOfViolation(actMessage)
	return []string{"notify"}

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

// action for missing pod annotation
func (a RequiredPodAnnotationAction) DoAction(entity ActionableEntity, vEntity libs.ViolatableEntity, lastActions map[string][]time.Time) []string {

	lastTimeWarned, doIt := getLastTimeWarnedAndifToDoAction(lastActions)
	if doIt {
		entity.DoAction()
		return []string{"entity_action"}
	}
	if canSkipNotification(lastTimeWarned) {
		libs.Log.Debug("Skipping notification for ", vEntity.Name, " ", a.Type, " it was notified less than ", libs.Cfg.DurationBetweenNotifyingAgain, " ago.")
		return []string{}
	}

	actMessage := createActionMessage(vEntity.Namespace, reflect.TypeOf(entity).Name(), vEntity.Name, "Missing pod annotation", a.Violation.Source, len(lastActions["notify"]), isLastWarning(lastActions))
	NotifyOfViolation(actMessage)
	return []string{"notify"}

}

func ConvertActionableEntityToViolatableEntity(entity ActionableEntity) (libs.ViolatableEntity, error) {

	var vEntity libs.ViolatableEntity

	switch t := entity.(type) {
	case ActionDeployment:
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
