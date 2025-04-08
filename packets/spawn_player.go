package packets

import "PowerClassic/network"

type SpawnPlayerPacket struct {
	PlayerId   uint8
	PlayerName string
	X          float32
	Y          float32
	Z          float32
	Yaw        byte
	Pitch      byte
}

func (p *SpawnPlayerPacket) PacketID() byte {
	return SPAWN_PLAYER
}

func (p *SpawnPlayerPacket) SerializePayload(b *network.Buffer) error {
	err := b.WriteByte(p.PlayerId)
	if err != nil {
		return err
	}
	err = b.WriteString(p.PlayerName)
	if err != nil {
		return err
	}
	err = b.WriteFShort(p.X)
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
	err = b.WriteByte(p.Yaw)
	if err != nil {
		return err
	}
	err = b.WriteByte(p.Pitch)
	if err != nil {
		return err
	}
	return nil
}

func (p *SpawnPlayerPacket) Deserialize(b *network.Buffer) error {
	pid, err := b.ReadByte()
	if err != nil {
		return err
	}
	pn, err := b.ReadString()
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
	p.PlayerName = pn
	p.X = x
	p.Y = y
	p.Z = z
	p.Yaw = yaw
	p.Pitch = pitch

	return nil
}
