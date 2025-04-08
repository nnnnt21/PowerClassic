package session

import (
	"PowerClassic/event"
	"PowerClassic/events"
	"PowerClassic/network"
	"PowerClassic/packets"
	"fmt"
	"github.com/anthdm/hollywood/actor"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"sync"
	"time"
)

// Promoter promotes a session once handshake is complete
type Promoter interface {
	PromoteSession(session *Session, evt *events.PlayerIdentificationEvent) error
}

// TODO: sloppy probably not needing
type ActorEngineProvider interface {
	GetEngine() *actor.Engine
}

type Session struct {
	tcpConn         net.Conn
	udpConn         *net.UDPConn
	udpAddr         *net.UDPAddr
	incomingChan    chan []byte
	outgoingChan    chan network.SerializablePacket
	handler         PacketHandler
	sessionPromoter Promoter
	engineProvider  ActorEngineProvider
	done            chan bool
	eventBus        *event.EventBus

	data *SessionData
}

type SessionData struct {
	Username        string
	ProtocolVersion byte
	VerificationKey string
}

func NewSession(tcpConn net.Conn, sessionPromoter Promoter, engineProvider ActorEngineProvider) *Session {
	s := &Session{
		tcpConn:         tcpConn,
		incomingChan:    make(chan []byte, 64),
		outgoingChan:    make(chan network.SerializablePacket, 64),
		sessionPromoter: sessionPromoter,
		engineProvider:  engineProvider,
		done:            make(chan bool),
		data:            &SessionData{},
		eventBus:        event.NewEventBus(),
	}
	s.SetHandler(NewAuthHandler(s))
	return s
}

func (s *Session) GetEventBus() *event.EventBus {
	return s.eventBus
}

func (s *Session) setUDPConnection(udpAddr *net.UDPAddr, udpConn *net.UDPConn) {
	s.udpAddr = udpAddr
	s.udpConn = udpConn
}

func (s *Session) SendPacketNow(sp network.SerializablePacket) error {
	buf := network.NewBuffer(2048)
	err := network.SerializePacket(sp, buf)
	if err != nil {
		return err
	}
	_, err = s.tcpConn.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("TCP outbound network failed to send: %v", err)
	}
	return nil
}

func (s *Session) Data() *SessionData {
	return s.data
}

func (s *Session) AddOutgoingPacket(sp network.SerializablePacket) {
	s.outgoingChan <- sp
}

func (s *Session) PromoteSession(evt *events.PlayerIdentificationEvent) error {
	return s.sessionPromoter.PromoteSession(s, evt)
}

func (s *Session) SetHandler(handler PacketHandler) {
	s.handler = handler
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return network.NewBuffer(2048)
	},
}

var incomingDataPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 2048)
	},
}

func (s *Session) Start() {
	incomingBuffer := make([]byte, 2097151)
	defer func() {
		err := s.tcpConn.Close()
		if err != nil {
			panic(fmt.Errorf("error closing TCP connection: %v", err))
		}
		s.done <- true
		close(s.done)
		close(s.outgoingChan)
		close(s.incomingChan)
	}()

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.AddOutgoingPacket(&packets.PingPacket{})
			case <-s.done:
				ticker.Stop()
				return
			}
		}
	}()

	go s.processIncomingPackets()
	go s.processOutgoingPackets()

	for {
		incomingLen, err := s.tcpConn.Read(incomingBuffer)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Error().Msgf("error reading incoming packet: %v", err)
		}
		packetData := incomingDataPool.Get().([]byte)
		if cap(packetData) < incomingLen {
			packetData = make([]byte, incomingLen)
		}
		packetData = packetData[:incomingLen]
		copy(packetData, incomingBuffer[:incomingLen])
		s.incomingChan <- packetData
	}
}

func (s *Session) processOutgoingPackets() {
	for {
		select {
		case sp, ok := <-s.outgoingChan:
			if !ok {
				return
			}
			if err := s.SendPacketNow(sp); err != nil {
				log.Err(err).Msg("failed to send outgoing packet")
			}
		case <-s.done:
			return
		}
	}
}

func (s *Session) processIncomingPackets() {
	for raw := range s.incomingChan {
		buf := bufferPool.Get().(*network.Buffer)
		buf.ResetWithData(raw)

		pid, err := network.PeekPacketID(buf)

		if err != nil {
			log.Err(err).Msg("failed to deserialize base packet")
			bufferPool.Put(buf)
			incomingDataPool.Put(raw)
			continue
		}

		switch pid {
		case packets.PLAYER_IDENTIFICATION:
			var pidPacket packets.PlayerIdentificationPacket
			if err := pidPacket.Deserialize(buf); err != nil {
				panic(fmt.Errorf("failed to deserialize PLAYER_IDENTIFICATION packet: %v", err))
			}
			s.handler.HandlePlayerIdentificationPacket(&pidPacket)
		case packets.PLAYER_MOVEMENT:
			var tpPkt packets.PlayerTeleportPacket
			if err := tpPkt.Deserialize(buf); err != nil {
				panic(fmt.Errorf("failed to deserialize PLAYER_MOVEMENT packet: %v", err))
			}
			s.handler.HandleMovement(s.engineProvider.GetEngine(), &tpPkt)
		default:
			log.Info().Msgf("unknown packet type: %v", pid)
		}

		bufferPool.Put(buf)
		incomingDataPool.Put(raw)
	}
}
