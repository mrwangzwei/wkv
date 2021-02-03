package queue

type Queue interface {
	Pull() ([]byte, error)
	Push(cont []byte) error
	Len() int
}
