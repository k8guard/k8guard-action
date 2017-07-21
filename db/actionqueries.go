package db

import (
	"fmt"
	"time"

	"github.com/k8guard/k8guard-action/db/stmts"

	"github.com/gocql/gocql"

	libs "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/violations"
)

func InsertVLOGRow(vEntity libs.ViolatableEntity, violation violations.Violation, entityType string) {
	err := Sess.Query(fmt.Sprintf(stmts.INSERT_TO_VLOG, libs.Cfg.CassandraKeyspace), vEntity.Namespace, libs.Cfg.ClusterName, entityType, vEntity.Name, string(violation.Type), violation.Source, time.Now()).Exec()
	if err != nil {
		panic(err)
	}
}

func InsertActionLogRow(namespace string, entityType string, entitySource string, violationType string, violationSource string, action string) {
	b := Sess.NewBatch(gocql.LoggedBatch)

	now := time.Now()

	b.Query(fmt.Sprintf(stmts.INSERT_TO_ALOG_NAMESPACE_TYPE, libs.Cfg.CassandraKeyspace), namespace, libs.Cfg.ClusterName, entityType, entitySource, violationType, violationSource, action, now)
	b.Query(fmt.Sprintf(stmts.INSERT_TO_ALOG_TYPE, libs.Cfg.CassandraKeyspace), namespace, libs.Cfg.ClusterName, entityType, entitySource, violationType, violationSource, action, now)
	b.Query(fmt.Sprintf(stmts.INSERT_TO_ALOG_VTYPE, libs.Cfg.CassandraKeyspace), namespace, libs.Cfg.ClusterName, entityType, entitySource, violationType, violationSource, action, now)
	b.Query(fmt.Sprintf(stmts.INSERT_TO_ALOG_ACTION, libs.Cfg.CassandraKeyspace), namespace, libs.Cfg.ClusterName, entityType, entitySource, violationType, violationSource, action, now)

	err := Sess.ExecuteBatch(b)
	if err != nil {
		panic(err)
	}
}

func SelectVActionRow(vEntity libs.ViolatableEntity, violation violations.Violation, entityType string) VActionRow {
	vActionQuery := VActionRow{
		Namespace: vEntity.Namespace,
		Type:      entityType,
		Source:    vEntity.Name,
		VType:     string(violation.Type),
		VSource:   violation.Source,
	}

	iter := Sess.Query(fmt.Sprintf(stmts.SELECT_ENTITY_FROM_VACTION, libs.Cfg.CassandraKeyspace),
		vActionQuery.Namespace, libs.Cfg.ClusterName, vActionQuery.Type, vActionQuery.Source, vActionQuery.VType, vActionQuery.VSource).Iter()

	vActionRow := VActionRow{Actions: map[string][]time.Time{}}

	for iter.Scan(&vActionRow.Namespace, &vActionRow.Type, &vActionRow.Source, &vActionRow.VType,
		&vActionRow.VSource, &vActionRow.Actions, &vActionRow.CreatedAt, &vActionRow.ExpiresAt) {
		break
	}

	if err := iter.Close(); err != nil {
		panic(err)
	}

	if vActionRow.ExpiresAt.IsZero() == false {
		//If it is expired create a new row with new data
		if vActionRow.ExpiresAt.Before(time.Now()) {
			vActionRow = VActionRow{Actions: map[string][]time.Time{}}
		}
	}

	return vActionRow
}

func InsertVactionRow(namespace string, entityType string, entitySource string, violationType string, violationSource string, actions map[string][]time.Time) {
	err := Sess.Query(fmt.Sprintf(stmts.INSERT_TO_VACTION, libs.Cfg.CassandraKeyspace), namespace, libs.Cfg.ClusterName, entityType, entitySource, violationType, violationSource, actions, time.Now(), time.Now().Add(libs.Cfg.DurationViolationExpires)).Exec()
	if err != nil {
		panic(err)
	}
}
