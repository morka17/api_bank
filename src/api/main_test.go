package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/morka17/shiny_bank/v1/src/db/sqlc"
	"github.com/morka17/shiny_bank/v1/src/utils"
	"github.com/stretchr/testify/assert"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenSymmetricKey: utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	assert.NoError(t, err)

	return server 
}


func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
