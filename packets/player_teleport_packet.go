package packets

import "PowerClassic/network"

type PlayerTeleportPacket struct {
	PlayerId uint8
	X        float32
	Y        float32
	Z        float32
	Yaw      byte
	Pitch    byte
}

func (p *PlayerTeleportPacket) PacketID() byte {
	return PLAYER_TELEPORT
}

func (p *PlayerTeleportPacket) SerializePayload(b *network.Buffer) error {
	if err := b.WriteByte(p.PlayerId); err != nil {
		return err
	}
	err := b.WriteFShort(p.X)
	if err != nil {
		return err
	}
	err = b.WriteFShort(p.Y)
	if err != nil {
		return err
	}
	err = b.WriteFShort(p.Z)
	if err != nil {
		return err
	}
	if err := b.WriteByte(p.Yaw); err != nil {
		return err
	}
	if err := b.WriteByte(p.Pitch); err != nil {
		return err
	}
	return nil
}

func (p *PlayerTeleportPacket) Deserialize(b *network.Buffer) error {
	pid, err := b.ReadByte()
	if err != nil {
		return err
	}
	x, err := b.ReadFShort()
	if err != nil {
		return err
	}
	y, err := b.ReadFShort()
	if err != nil {
		return err
	}
	z, err := b.ReadFShort()
	if err != nil {
		return err
	}
	yaw, err := b.ReadByte()
	if err != nil {
		return err
	}
	pitch, err := b.ReadByte()
	if err != nil {
		return err
	}

	p.PlayerId = pid
	p.X = x
	p.Y = y
	p.Z = z
	p.Yaw = yaw
	p.Pitch = pitch
	return nil
}
