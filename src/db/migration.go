package db

import "github.com/shaileshhb/equisplit/src/server"

// MigrateTables will migrate all the tables.
func MigrateTables(ser *server.Server) {
	ser.MigrateTables([]server.ModuleConfig{})
}
