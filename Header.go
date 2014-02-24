package amf

type Header struct {
	name string
	mustUnderstand bool
	data interface{}
}
