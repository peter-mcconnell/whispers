package whispers

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

func byteArrayToString(b []byte) string {
	n := -1
	for i, v := range b {
		if v == 0 {
			n = i
			break
		}
	}
	if n == -1 {
		n = len(b)
	}
	return string(b[:n])
}

func parseEventData(data []byte) *eventT {
	var event eventT

	if len(data) >= int(unsafe.Sizeof(eventT{})) {
		if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &event); err == nil {
			return &event
		}
	}

	return nil
}
