package main

import (
	"PowerClassic/events"
	"PowerClassic/messages"
	"PowerClassic/player"
	"PowerClassic/session"
	"PowerClassic/world"
	"fmt"
	"github.com/anthdm/hollywood/actor"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"time"
)

type Server struct {
	pid *actor.PID

	Config *Config

	actorEngine *actor.Engine

	world *world.World

	tcpServer *TCPServer
}

func (s *Server) Receive(ctx *actor.Context) {
	switch ctx.Message().(type) {
	case actor.Started:
		s.world = world.NewWorld(ctx.Engine(), 1024, 64, 1024)

		worldPid := ctx.SpawnChild(func() actor.Receiver { return s.world }, "world")
		s.world.SetPID(worldPid)

		log.Debug().Msgf("creating flat world")
		resp := ctx.Request(worldPid, &world.WorldRunnable{func(ctx *actor.Context, w *world.World) {
			err := w.CreateFlatWorld()
			ctx.Respond(&messages.Response{E: err})
		}}, time.Second*60)

		res, err := resp.Result()
		log.Debug().Msgf("flat world created")

		worldResp, ok := res.(*messages.Response)

		if err != nil || !ok || worldResp.E != nil {
			panic(fmt.Errorf("failed to create flat world %v %v", err, worldResp.E))
		}

		tcp_port, err := strconv.Atoi(os.Getenv("TCP_PORT"))
		if err != nil {
			log.Err(err).Msg("Error reading TCP port")
		}

		s.tcpServer = NewTCPServer(s, tcp_port)

		go s.tcpServer.Start()
	}
}

func NewServer(config *Config) *Server {
	eng, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}
	server := &Server{
		Config:      config,
		actorEngine: eng,
	}
	server.pid = eng.Spawn(func() actor.Receiver {
		return server
	}, "server")

	return server
}

func (s *Server) GetEngine() *actor.Engine {
	return s.actorEngine
}

func (s *Server) PromoteSession(session *session.Session, evt *events.PlayerIdentificationEvent) error {

	wp := evt.SpawnWorld()

	w, ok := wp.(*world.World)
	if !ok {
		return fmt.Errorf("expected world")
	}

	p := player.NewPlayer(session, w.GetNextEntityId(), *evt.SpawnX(), *evt.SpawnY(), *evt.SpawnZ(), s.actorEngine)

	s.actorEngine.Send(w.GetPID(), &world.WorldRunnable{func(ctx *actor.Context, w *world.World) {
		w.AddEntity(ctx, p, evt)
	}})

	return nil
}

var _ actor.Receiver = (*Server)(nil)
