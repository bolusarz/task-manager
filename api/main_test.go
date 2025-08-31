package api

import (
	"os"
	"testing"

	db "github.com/bolusarz/task-manager/db/sqlc"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	return NewServer(store)
}

func TestMain(m *testing.M) {
	validate.RegisterValidation("strong", IsPasswordStrong)

	exitCode := m.Run()

	os.Exit(exitCode)
}
