package gapi

import (
	"testing"
	"time"

	db "github.com/morka17/shiny_bank/v1/src/db/sqlc"
	"github.com/morka17/shiny_bank/v1/src/utils"
	"github.com/morka17/shiny_bank/v1/src/worker"
	"github.com/stretchr/testify/assert"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	assert.NoError(t, err)

	return server
}
