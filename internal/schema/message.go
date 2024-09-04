package schema

type EventMessage interface {
	JSON() ([]byte, error)
}
