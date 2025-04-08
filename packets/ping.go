package packets

import "PowerClassic/network"

type PingPacket struct {
}

func (p *PingPacket) PacketID() byte {
	return PING
}

func (p *PingPacket) SerializePayload(b *network.Buffer) error {
	return nil
}

func (p *PingPacket) Deserialize(b *network.Buffer) error {
	return nil
}
