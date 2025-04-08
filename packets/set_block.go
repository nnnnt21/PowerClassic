package packets

import "PowerClassic/network"

type SetBlockPacket struct {
	X         int16
	Y         int16
	Z         int16
	BlockType byte
}

func (p *SetBlockPacket) PacketID() byte {
	return SET_BLOCK
}

func (p *SetBlockPacket) SerializePayload(b *network.Buffer) error {
	b.WriteShort(p.X)
	b.WriteShort(p.Y)
	b.WriteShort(p.Z)
	if err := b.WriteByte(p.BlockType); err != nil {
		return err
	}
	return nil
}

func (p *SetBlockPacket) Deserialize(b *network.Buffer) error {
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
	bt, err := b.ReadByte()
	if err != nil {
		return err
	}
	p.X = x
	p.Y = y
	p.Z = z
	p.BlockType = bt
	return nil
}
