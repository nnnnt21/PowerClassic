package packets

import "PowerClassic/network"

type DisconnectPacket struct {
	Reason string
}

func (p *DisconnectPacket) PacketID() byte {
	return DISCONNECT
}

func (p *DisconnectPacket) SerializePayload(b *network.Buffer) error {
	err := b.WriteString(p.Reason)
	if err != nil {
		return err
	}
	return nil
}

func (p *DisconnectPacket) Deserialize(b *network.Buffer) error {
	reason, err := b.ReadString()
	if err != nil {
		return err
	}

	p.Reason = reason

	return nil
}
