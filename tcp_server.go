package main

import (
	"PowerClassic/event"
	"PowerClassic/events"
	"PowerClassic/session"
	"github.com/rs/zerolog/log"
	"net"
	"strconv"
)

type TCPServer struct {
	tcpCon net.Listener
	port   int
	server *Server
}

func NewTCPServer(server *Server, port int) *TCPServer {
	listener, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	return &TCPServer{tcpCon: listener, port: port, server: server}
}

func (s *TCPServer) GetPort() int {
	return s.port
}

func (s *TCPServer) Start() {
	log.Debug().Msgf("Starting TCP Listener on port: %d", s.port)

	for {
		c, err := s.tcpCon.Accept()
		if err != nil {
			log.Err(err).Msg("Error getting tcp connection")
			return
		}
		sess := session.NewSession(c, s.server, s.server)

		sess.GetEventBus().Register("PlayerIdentification", func(evt event.Event) {
			pid, ok := evt.(*events.PlayerIdentificationEvent)
			if ok {
				var spawnX = float32(10)
				var spawnY = float32(10)
				var spawnZ = float32(10)
				pid.SetSpawnX(&spawnX)
				pid.SetSpawnY(&spawnY)
				pid.SetSpawnZ(&spawnZ)
				pid.SetSpawnWorld(s.server.world)
			}
		})

		go sess.Start()
	}
}
