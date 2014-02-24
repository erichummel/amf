package amf

type Message struct {
	targetUri string
	responseUri string
	data interface{}
}
