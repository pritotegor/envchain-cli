// Package rollback provides the ability to restore a chain to a previously
// captured snapshot.
//
// A Rollbacker wraps a chain.Store and a snapshot.Snapshotter. Given a chain
// name and a snapshot identifier it replaces the current chain contents with
// the values recorded at snapshot time.
//
// Example:
//
//	rb := rollback.NewRollbacker(store)
//	if err := rb.Rollback("myapp", "before-deploy", false); err != nil {
//		log.Fatal(err)
//	}
package rollback
