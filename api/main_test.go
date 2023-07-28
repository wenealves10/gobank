package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/wenealves10/gobank/db/sqlc"
	"github.com/wenealves10/gobank/utils"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenPassetoKey:     utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
