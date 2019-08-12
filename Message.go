package amf

type Message struct {
	TargetURI   string
	ResponseURI string
	Data        interface{}
}
