package main

import (
	"cli/internal/charm"
	"context"
	"fmt"
	"os"

	_ "github.com/vegidio/avif-go"
	_ "github.com/vegidio/heif-go"

	"github.com/vegidio/mediaorient"
)

func main() {
	if err := mediaorient.Initialize(); err != nil {
		charm.PrintError(fmt.Sprintf("Failed to initialize media orientation detection: %v\n", err))
		return
	}
	defer mediaorient.Destroy()

	// Add support for AVIF and HEIC images
	mediaorient.AddImageType(".avif", ".heic")

	fmt.Print("\n")
	cmd := buildCliCommands()

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		charm.PrintError(err.Error())
	}
}
