// Package audit provides an append-only structured audit log for envchain-cli.
//
// Every mutation command (add, delete, rename, copy, merge) can record an Entry
// to a newline-delimited JSON file so operators can review what changed, when,
// and to which chain.
//
// Usage:
//
//	logger, err := audit.NewLogger("/var/log/envchain/audit.log")
//	if err != nil { ... }
//
//	// record an action
//	_ = logger.Record("add", "myapp", "DB_URL", "")
//
//	// read back all entries
//	entries, err := logger.ReadAll()
package audit
