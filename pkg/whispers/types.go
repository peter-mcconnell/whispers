package whispers

type eventT struct {
	Pid      int32
	Comm     [16]byte
	Username [80]byte
	Password [80]byte
}
