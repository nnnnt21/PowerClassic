package packets

import "PowerClassic/network"

type PlayerIdentificationPacket struct {
	ProtocolVersion byte
	Username        string
	VerificationKey string
	Unused          byte
}

func (p *PlayerIdentificationPacket) PacketID() byte {
	return PLAYER_IDENTIFICATION
}

func (p *PlayerIdentificationPacket) SerializePayload(b *network.Buffer) error {
	err := b.WriteByte(p.ProtocolVersion)
	if err != nil {
		return err
	}
	err = b.WriteString(p.Username)
	if err != nil {
		return err
	}
	err = b.WriteString(p.VerificationKey)
	if err != nil {
		return err
	}
	err = b.WriteByte(p.Unused)
	if err != nil {
		return err
	}
	return nil
}

func (p *PlayerIdentificationPacket) Deserialize(b *network.Buffer) error {
	pv, err := b.ReadByte()
	if err != nil {
		return err
	}
	un, err := b.ReadString()
	if err != nil {
		return err
	}
	vk, err := b.ReadString()
	if err != nil {
		return err
	}
	us, err := b.ReadByte()
	if err != nil {
		return err
	}

	p.ProtocolVersion = pv
	p.Username = un
	p.VerificationKey = vk
	p.Unused = us

	return nil
}
