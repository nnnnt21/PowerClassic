package packets

import "PowerClassic/network"

type LevelInitializePacket struct {
}

func (p *LevelInitializePacket) PacketID() byte {
	return LEVEL_INITIALIZE
}

func (p *LevelInitializePacket) SerializePayload(b *network.Buffer) error {
	return nil
}

func (p *LevelInitializePacket) Deserialize(b *network.Buffer) error {
	return nil
}
