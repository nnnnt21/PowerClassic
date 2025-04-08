package packets

import "PowerClassic/network"

type LevelDataChunkPacket struct {
	ChunkLength     int16
	ChunkData       []byte
	PercentComplete byte
}

func (p *LevelDataChunkPacket) PacketID() byte {
	return LEVEL_DATA_CHUNK
}

func (p *LevelDataChunkPacket) SerializePayload(b *network.Buffer) error {
	b.WriteShort(p.ChunkLength)
	err := b.WriteByteArray(p.ChunkData)
	if err != nil {
		return err
	}
	err = b.WriteByte(p.PercentComplete)
	if err != nil {
		return err
	}
	return nil
}

func (p *LevelDataChunkPacket) Deserialize(b *network.Buffer) error {
	cl, err := b.ReadShort()
	if err != nil {
		return err
	}
	cd, err := b.ReadByteArray()
	if err != nil {
		return err
	}
	pc, err := b.ReadByte()
	if err != nil {
		return err
	}

	p.ChunkLength = cl
	p.ChunkData = cd
	p.PercentComplete = pc

	return nil
}
