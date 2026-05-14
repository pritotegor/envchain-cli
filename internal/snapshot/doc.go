// Package snapshot provides point-in-time capture and restore for
// environment chains managed by envchain-cli.
//
// A Snapshot records all key-value pairs belonging to a chain at the
// moment Capture is called, along with the chain name and a UTC
// timestamp. The snapshot can later be passed to Restore to rewrite
// the chain's contents, optionally requiring an explicit overwrite
// flag when the destination chain already exists.
//
// Typical usage:
//
//	sn := snapshot.NewSnapshotter(store)
//	snap, err := sn.Capture("production")
//	// ... persist snap as JSON, share, or diff ...
//	err = sn.Restore(snap, false)
package snapshot
