package message

type Object interface {
	Get(key string) (obj any, found bool)
}

type Command interface {
	Message
	GetCommand() string
	GetTransactionId() float64
	GetObject() Object
	GetParameters() []any
}
