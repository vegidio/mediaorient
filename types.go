package mediaorient

import (
	"fmt"
	"image"
)

// Media represents a media object.
type Media struct {
	// Name of the media.
	Name string `json:"name"`
	// Type of the media (e.g., image, video).
	Type string `json:"type"`
	// Rotation represents the rotation of the media.
	Rotation int `json:"rotation"`
	// Frames contain the image data of the media.
	Frames []image.Image `json:"-"`
	// Width represents the width of the media in pixels.
	Width int `json:"width"`
	// Height represents the height of the media in pixels.
	Height int `json:"height"`
}

func (m *Media) String() string {
	return fmt.Sprintf(`{Name: %s, Type: %s, Rotation: %d, Width: %d, Height: %d}`,
		m.Name, m.Type, m.Rotation, m.Width, m.Height)
}
