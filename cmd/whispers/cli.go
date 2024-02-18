//go:build amd64

package main

import (
	"context"
	"log"

	"github.com/peter-mcconnell/whispers/pkg/whispers"
)

func main() {
	cfg := cfgFromFlags()
	ctx := context.Background()

	if err := whispers.EnvSetup(ctx, cfg); err != nil {
		log.Fatal(err)
	}

	if err := whispers.Listen(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
