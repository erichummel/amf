package amf

type Packet struct {
	Version  uint16
	Headers  []*Header
	Messages []*Message
}

func NewPacket(numHeaders, numMessages int) *Packet {
	p := Packet{}
	p.Headers = make([]*Header, numHeaders)
	p.Messages = make([]*Message, numMessages)
	
	return &p
}
