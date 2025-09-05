package main

import (
	"log"

	"github.com/vegidio/mediaorient"
)

func main() {
	if err := mediaorient.Initialize(); err != nil {
		log.Fatal("Failed to initialize media orientation detection:", err)
	}
	defer mediaorient.Destroy()

	media, err := mediaorient.CalculateFileOrientation("../assets/image_270.jpg")
	if err != nil {
		log.Println("Error:", err)
	}

	log.Println(media)
}
