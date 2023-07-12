package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/morka17/shiny_bank/v1/src/db/sqlc"
	"github.com/morka17/shiny_bank/v1/src/token"
	"github.com/morka17/shiny_bank/v1/src/utils"
	"github.com/morka17/shiny_bank/v1/src/worker"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
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


func newContextWithBearer(t *testing.T, tokenMaker token.Maker, username string, duration time.Duration) context.Context {

	accessToken, _, err := tokenMaker.CreateToken(username, duration)
	assert.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorization: []string{
			bearerToken,
		},
	}
	return metadata.NewIncomingContext(context.Background(), md)


}