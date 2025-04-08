package packets

import "PowerClassic/network"

type ServerIdentificationPacket struct {
	ProtocolVersion byte
	ServerName      string
	MOTD            string
	UserType        byte
}

func (p *ServerIdentificationPacket) PacketID() byte {
	return SERVER_IDENTIFICATION
}

func (p *ServerIdentificationPacket) SerializePayload(b *network.Buffer) error {
	err := b.WriteByte(p.ProtocolVersion)
	if err != nil {
		return err
	}
	err = b.WriteString(p.ServerName)
	if err != nil {
		return err
	}
	err = b.WriteString(p.MOTD)
	if err != nil {
		return err
	}
	err = b.WriteByte(p.UserType)
	if err != nil {
		return err
	}
	return nil
}

func (p *ServerIdentificationPacket) Deserialize(b *network.Buffer) error {
	pv, err := b.ReadByte()
	if err != nil {
		return err
	}
	sn, err := b.ReadString()
	if err != nil {
		return err
	}
	motd, err := b.ReadString()
	if err != nil {
		return err
	}
	ut, err := b.ReadByte()
	if err != nil {
		return err
	}

	p.ProtocolVersion = pv
	p.ServerName = sn
	p.MOTD = motd
	p.UserType = ut

	return nil
}
