//go:build !oracle

package store

import "errors"

// OracleStore is only compiled with -tags oracle, because godror needs cgo
// and Oracle Instant Client. The default build is pure Go and uses SQLite.
type OracleStore struct{ Store }

func NewOracleStore(string) (*OracleStore, error) {
	return nil, errors.New("built without Oracle support; rebuild with -tags oracle")
}
