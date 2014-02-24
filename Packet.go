package amf

type Packet struct {
	version uint16
	headers []*Header
	messages []*Message
}

func NewPacket(numHeaders, numMessages int) *Packet {
	p := Packet{}
	p.headers = make([]*Header, numHeaders)
	p.messages = make([]*Message, numMessages)
	
	return &p
}
