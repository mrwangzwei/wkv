package queue

type Queue interface {
	Pull() (string, error)
	Push(cont string) error
	Len() int
}
