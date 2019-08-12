package amf

/**
 * Decoder for AMF Protocol
 */
import (
	"bufio"
	"encoding/binary"
	"io"
	"github.com/marcuswu/amf/amf0"
)

type Decoder struct {
	r       io.Reader
}

// should use io.LimitedReader
func NewDecoder(r io.Reader) *Decoder {
	if _, ok := r.(*bufio.Reader); ok {
		return &Decoder{r: r}
	}
	return &Decoder{r: bufio.NewReader(r)}
}

func (dec *Decoder) Decode() (p *Packet, err error) {
	p, err = dec.decodePacket()
	if err != nil {
		return nil, err
	}
	
	return
}

func (dec *Decoder) decodePacket() (p *Packet, err error) {
	p = &Packet{}
	u16 := make([]byte, 2)

	_, err = dec.r.Read(u16)
	if err != nil {
		return nil, err
	}
	p.Version = binary.BigEndian.Uint16(u16)

	_, err = dec.r.Read(u16)
	if err != nil {
		return nil, err
	}
	headerCount := binary.BigEndian.Uint16(u16)

	p.Headers = make([]*Header, headerCount)
	for i := 0; i < len(p.Headers); i++ {
		p.Headers[i], err = dec.decodeHeader()
		if err != nil {
			return nil, err
		}
	}

	_, err = dec.r.Read(u16)
	if err != nil {
		return nil, err
	}
	messageCount := binary.BigEndian.Uint16(u16)

	p.Messages = make([]*Message, messageCount)
	for i := 0; i < len(p.Messages); i++ {
		p.Messages[i], err = dec.decodeMessage()
		if err != nil {
			return nil, err
		}
	}
	
	return
}

func (dec *Decoder) decodeHeader() (h *Header, err error) {
	h = &Header{}
	u8 := make([]byte, 1)
	u16 := make([]byte, 2)
	u32 := make([]byte, 4)

	_, err = dec.r.Read(u16)
	if err != nil {
		return nil, err
	}
	headerNameLen := binary.BigEndian.Uint16(u16)

	headerNameBytes := make([]byte, headerNameLen)
	_, err = dec.r.Read(headerNameBytes)
	if err != nil {
		return nil, err
	}
	h.name = string(headerNameBytes)

	_, err = dec.r.Read(u8)
	if err != nil {
		return nil, err
	}
	h.mustUnderstand = u8[0] != 0

	_, err = dec.r.Read(u32)
	if err != nil {
		return nil, err
	}

	var amf0Decoder *amf0.Decoder = amf0.NewDecoder(dec.r)
	h.data, err = amf0Decoder.Decode()
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (dec *Decoder) decodeMessage() (m *Message, err error) {
	m = &Message{}
	u16 := make([]byte, 2)
	u32 := make([]byte, 4)

	_, err = dec.r.Read(u16)
	if err != nil {
		return nil, err
	}
	targetUriLen := binary.BigEndian.Uint16(u16)

	targetUriBytes := make([]byte, targetUriLen)
	_, err = dec.r.Read(targetUriBytes)
	if err != nil {
		return nil, err
	}
	m.TargetURI = string(targetUriBytes)

	_, err = dec.r.Read(u16)
	if err != nil {
		return nil, err
	}
	responseUriLen := binary.BigEndian.Uint16(u16)

	responseUriBytes := make([]byte, responseUriLen)
	_, err = dec.r.Read(responseUriBytes)
	if err != nil {
		return nil, err
	}
	m.ResponseURI = string(responseUriBytes)

	_, err = dec.r.Read(u32)
	if err != nil {
		return nil, err
	}

	var amf0Decoder *amf0.Decoder = amf0.NewDecoder(dec.r)
	m.Data, err = amf0Decoder.Decode()
	if err != nil {
		return nil, err
	}

	return m, nil
}
