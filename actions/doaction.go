package actions

import (
	"reflect"
	"time"

	libs "github.com/k8guard/k8guardlibs"
)

//  actionable is interface, violatable is struct
func DoAction(action Action, entity ActionableEntity, violatableEntity libs.ViolatableEntity, lastActions map[string][]time.Time, dryRun bool) map[string][]time.Time {
	doneActions := map[string][]time.Time{}

	if dryRun {
		libs.Log.Info("Running dry run for action ", reflect.TypeOf(action).Name())
		return doneActions
	}

	for _, doneAction := range action.DoAction(entity, violatableEntity, lastActions) {
		if _, ok := doneActions[doneAction]; ok {
			doneActions[doneAction] = append(doneActions[doneAction], time.Now())
		} else {
			doneActions[doneAction] = []time.Time{time.Now()}
		}
	}

	return doneActions
}
