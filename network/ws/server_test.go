package ws_test

import (
	"github.com/rs/zerolog/log"
	"github.com/tieing/lemon/network/inet"
	"github.com/tieing/lemon/network/ws"
	"testing"
)

func TestServer(t *testing.T) {
	server := ws.NewServer("1")
	server.OnStart(func() {
		log.Info().Msg("server is started")
	})
	server.OnConnect(func(conn inet.Conn) {
		log.Info().Msgf("connection is opened, connection id: %d", conn.ID())
	})
	server.OnDisconnect(func(conn inet.Conn) {
		log.Info().Msgf("connection is closed, connection id: %d", conn.ID())
	})
	server.OnReceive(func(conn inet.Conn, msg []byte) {
		log.Info().Msgf("receive msg from client, connection id: %d, msg: %s", conn.ID(), string(msg))

		if err := conn.Push([]byte("I'm fine~~")); err != nil {
			log.Error().Msgf("push message failed: %v", err)
		}
	})

	if err := server.Start(); err != nil {
		log.Fatal().Msgf("start server failed: %v", err)
	}
}
