package stmts

/*
   This package we need to store CQL statements
*/

const (
	/*
	   ========================================
	   Init statements
	   ========================================
	*/

	CREATE_KEYSPACE = "CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }"

	// Violation Log Table
	CREATE_VIOLATION_LOG_TABLE = `
		CREATE TABLE IF NOT EXISTS %s.vlog_namespace_type (
			namespace varchar,
			cluster varchar,
			type varchar,
			source varchar,
			vType varchar,
			vSource varchar,
			created_at timestamp,
			PRIMARY KEY((namespace,cluster,type,source),created_at))
			WITH CLUSTERING ORDER BY (created_at DESC)
	`

	CREATE_ACTION_LOG_NAMESPACE_TYPE_TABLE = `
		CREATE TABLE IF NOT EXISTS %s.alog_namespace_type (
			namespace varchar,
			cluster varchar,
			type varchar,
			source varchar,
			vType varchar,
			vSource varchar,
			action varchar,
			created_at timestamp,
			PRIMARY KEY((namespace,type),created_at))
			WITH CLUSTERING ORDER BY (created_at DESC)
	`

	CREATE_ACTION_LOG_TYPE_TABLE = `
		CREATE TABLE IF NOT EXISTS %s.alog_type (
			namespace varchar,
			cluster varchar,
			type varchar,
			source varchar,
			vType varchar,
			vSource varchar,
			action varchar,
			created_at timestamp,
			PRIMARY KEY((type),created_at))
			WITH CLUSTERING ORDER BY (created_at DESC)
	`
	CREATE_ACTION_LOG_VTYPE_TABLE = `
		CREATE TABLE IF NOT EXISTS %s.alog_vType (
			namespace varchar,
			cluster varchar,
			type varchar,
			source varchar,
			vType varchar,
			vSource varchar,
			action varchar,
			created_at timestamp,
			PRIMARY KEY((vType),created_at))
			WITH CLUSTERING ORDER BY (created_at DESC)
	`

	CREATE_ACTION_LOG_ACTION_TABLE = `
		CREATE TABLE IF NOT EXISTS %s.alog_action (
			namespace varchar,
			cluster varchar,
			type varchar,
			source varchar,
			vType varchar,
			vSource varchar,
			action varchar,
			created_at timestamp,
			PRIMARY KEY((action),created_at))
			WITH CLUSTERING ORDER BY (created_at DESC)
	`

	// Tracks the status of a violation
	CREATE_VACTION_TABLE = `
		CREATE TABLE IF NOT EXISTS %s.vaction (
			namespace varchar,
			cluster varchar,
			type varchar,
			source varchar,
			vType varchar,
			vSource varchar,
			actions frozen<map<varchar, list<timestamp>>>,
			created_at timestamp,
			expire_at timestamp,
			PRIMARY KEY((namespace,cluster,type,source,vtype,vsource),created_at))
			WITH CLUSTERING ORDER BY (created_at desc)
	`

	INSERT_TO_VLOG = `INSERT INTO %s.vlog_namespace_type (namespace, cluster, type, source, vType, vSource, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`

	INSERT_TO_ALOG_NAMESPACE_TYPE = `INSERT INTO %s.alog_namespace_type (namespace, cluster, type, source, vType, vSource, action, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	INSERT_TO_ALOG_TYPE           = `INSERT INTO %s.alog_type (namespace, cluster, type, source, vType, vSource, action, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	INSERT_TO_ALOG_VTYPE          = `INSERT INTO %s.alog_vType (namespace, cluster, type, source, vType, vSource, action, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	INSERT_TO_ALOG_ACTION         = `INSERT INTO %s.alog_action (namespace, cluster, type, source, vType, vSource, action, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	INSERT_TO_VACTION = `INSERT INTO %s.vaction (namespace, cluster, type, source, vType, vSource, actions, created_at ,expire_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	SELECT_ENTITY_FROM_VACTION = `SELECT namespace, type, source, vType, vSource, actions, created_at, expire_at FROM %s.vaction WHERE namespace = ? AND cluster = ? AND type = ? AND source = ? AND vType = ? AND vSource = ? LIMIT 1`
)
