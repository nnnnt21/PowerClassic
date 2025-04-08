package packets

import "PowerClassic/network"

type PositionOrientationUpdatePacket struct {
	PlayerId uint8
	ChangeX  int8
	ChangeY  int8
	ChangeZ  int8
	Yaw      byte
	Pitch    byte
}

func (p *PositionOrientationUpdatePacket) PacketID() byte {
	return POSITION_UPDATE
}

func (p *PositionOrientationUpdatePacket) SerializePayload(b *network.Buffer) error {
	if err := b.WriteByte(p.PlayerId); err != nil {
		return err
	}
	if err := b.WriteSByte(p.ChangeX); err != nil {
		return err
	}
	if err := b.WriteSByte(p.ChangeY); err != nil {
		return err
	}
	if err := b.WriteSByte(p.ChangeZ); err != nil {
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

func (p *PositionOrientationUpdatePacket) Deserialize(b *network.Buffer) error {
	playerId, err := b.ReadByte()
	if err != nil {
		return err
	}
	changeX, err := b.ReadSByte()
	if err != nil {
		return err
	}
	changeY, err := b.ReadSByte()
	if err != nil {
		return err
	}
	changeZ, err := b.ReadSByte()
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

	p.PlayerId = playerId
	p.ChangeX = changeX
	p.ChangeY = changeY
	p.ChangeZ = changeZ
	p.Yaw = yaw
	p.Pitch = pitch
	return nil
}
