package internal

type SubAction uint8

const (
	SAInvalid     SubAction = 0
	SASubscribe   SubAction = 1
	SAUnsubscribe SubAction = 2
)

func (ac SubAction) string() string {
	switch ac {
	case SASubscribe:
		return "subscribe"
	case SAUnsubscribe:
		return "unsubscribe"
	default:
		panic("should never happen")
	}
}
