package packets

import "PowerClassic/network"

type LevelFinalizePacket struct {
	X int16
	Y int16
	Z int16
}

func (p *LevelFinalizePacket) PacketID() byte {
	return LEVEL_FINALIZE
}

func (p *LevelFinalizePacket) SerializePayload(b *network.Buffer) error {
	b.WriteShort(p.X)
	b.WriteShort(p.Y)
	b.WriteShort(p.Z)
	return nil
}

func (p *LevelFinalizePacket) Deserialize(b *network.Buffer) error {
	x, err := b.ReadShort()
	if err != nil {
		return err
	}
	y, err := b.ReadShort()
	if err != nil {
		return err
	}
	z, err := b.ReadShort()
	if err != nil {
		return err
	}

	p.X = x
	p.Y = y
	p.Z = z

	return nil
}
