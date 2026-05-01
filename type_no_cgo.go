//go:build !cgo

package rcl

type typePackage string

type MessageTypePackage string
type ServiceTypePackage string

func (r TypeLoader) LoadMessage(name string) (MessageTypePackage, error) {
	return MessageTypePackage(name), nil
}

func (r TypeLoader) LoadService(name string) (ServiceTypePackage, error) {
	return ServiceTypePackage(name), nil
}

type MessageTypeSupport string

func (p MessageTypePackage) Get(name string) (MessageTypeSupport, error) {
	return MessageTypeSupport(string(p) + "/" + name), nil
}

type ServiceTypeSupport string

func (p ServiceTypePackage) Get(name string) (ServiceTypeSupport, error) {
	return ServiceTypeSupport(string(p) + "/" + name), nil
}

func (s ServiceTypeSupport) Request() MessageTypeSupport {
	return MessageTypeSupport(s + "+Request")
}

func (s ServiceTypeSupport) Response() MessageTypeSupport {
	return MessageTypeSupport(s + "+Response")
}

type Message struct {
	name string
}

func (m Message) Close() error { return nil }

func (s MessageTypeSupport) Create() (Message, error) {
	return Message{name: string(s)}, nil
}

func (s MessageTypeSupport) Destroy(_ interface{}) error {
	return nil
}

func Serialize(_ *SerializedMessage, _ Message) error   { return nil }
func Deserialize(_ Message, _ *SerializedMessage) error { return nil }
