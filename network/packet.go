package network

type SerializablePacket interface {
	PacketID() byte
	SerializePayload(b *Buffer) error
	Deserialize(b *Buffer) error
}

func SerializePacket(sp SerializablePacket, b *Buffer) error {
	err := b.WriteByte(sp.PacketID())
	if err != nil {
		return err
	}
	err = sp.SerializePayload(b)
	if err != nil {
		return err
	}
	return nil
}

func PeekPacketID(b *Buffer) (byte, error) {
	pid, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	return pid, nil
}
