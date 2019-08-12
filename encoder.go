package amf

/**
 * Encoder for AMF Protocol
 */
import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"github.com/marcuswu/amf/amf3"
)

type Encoder struct {
	w       *bufio.Writer
}

// should use io.LimitedReader
func NewEncoder(w io.Writer) *Encoder {
	if buf, ok := w.(*bufio.Writer); ok {
		return &Encoder{w: buf}
	}
	return &Encoder{w: bufio.NewWriter(w)}
}

func (enc *Encoder) Encode(p *Packet) (err error) {
	err = enc.encodePacket(p)
	if err != nil {
		return nil
	}
	enc.w.Flush()

	return
}

func (enc *Encoder) encodePacket(p *Packet) (err error) {
	u16 := make([]byte, 2)

	binary.BigEndian.PutUint16(u16, p.Version)
	_, err = enc.w.Write(u16)
	if err != nil {
		return nil
	}

	binary.BigEndian.PutUint16(u16, uint16(len(p.Headers)))
	_, err = enc.w.Write(u16)
	if err != nil {
		return err
	}
	for i := range p.Headers {
		err = enc.encodeHeader(p.Headers[i])
		if err != nil {
			return err
		}
	}

	binary.BigEndian.PutUint16(u16, uint16(len(p.Messages)))
	_, err = enc.w.Write(u16)
	if err != nil {
		return err
	}
	for i := range p.Messages {
		err = enc.encodeMessage(p.Messages[i])
		if err != nil {
			return err
		}
	}

	return
}

func (enc *Encoder) encodeHeader(h *Header) (err error) {
	h = &Header{}
	u8 := make([]byte, 1)
	u16 := make([]byte, 2)
	u32 := make([]byte, 4)

	binary.BigEndian.PutUint16(u16, uint16(len(h.name)))
	_, err = enc.w.Write(u16)
	if err != nil {
		return nil
	}

	_, err = enc.w.Write([]byte(h.name))
	if err != nil {
		return nil
	}

	if h.mustUnderstand {
		u8[0] = 1
	} else {
		u8[0] = 0
	}
	_, err = enc.w.Write(u8)
	if err != nil {
		return nil
	}

	var messageBuffer *bytes.Buffer = &bytes.Buffer{}
	var amf3Encoder *amf3.Encoder = amf3.NewEncoder(messageBuffer)
	amf3Encoder.Encode(h.data)
	var messageLen uint32 = uint32(messageBuffer.Len()) + 1

	binary.BigEndian.PutUint32(u32, messageLen)
	_, err = enc.w.Write(u32)
	if err != nil {
		return nil
	}

	//Switch to AMF3
	_, err = enc.w.Write([]byte{ 0x11 })
	if err != nil {
		return nil
	}

	_, err = enc.w.Write(messageBuffer.Bytes())
	if err != nil {
		return err
	}

	return
}

func (enc *Encoder) encodeMessage(m *Message) (err error) {
	u16 := make([]byte, 2)
	u32 := make([]byte, 4)

	binary.BigEndian.PutUint16(u16, uint16(len(m.TargetURI)))
	_, err = enc.w.Write(u16)
	if err != nil {
		return err
	}

	_, err = enc.w.Write([]byte(m.TargetURI))
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint16(u16, uint16(len(m.ResponseURI)))
	_, err = enc.w.Write(u16)
	if err != nil {
		return err
	}

	_, err = enc.w.Write([]byte(m.ResponseURI))
	if err != nil {
		return err
	}

	var messageBuffer *bytes.Buffer = &bytes.Buffer{}
	var amf3Encoder *amf3.Encoder = amf3.NewEncoder(messageBuffer)
	amf3Encoder.Encode(m.Data)
//	messageLen := uint32(messageBuffer.Len())

	binary.BigEndian.PutUint32(u32, 1)
	_, err = enc.w.Write(u32)
	if err != nil {
		return nil
	}

	//Switch to AMF3
	_, err = enc.w.Write([]byte{ 0x11 })
	if err != nil {
		return nil
	}

	_, err = enc.w.Write(messageBuffer.Bytes())
	if err != nil {
		return err
	}

	return
}
