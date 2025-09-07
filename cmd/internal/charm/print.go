package charm

import (
	"fmt"
	"strconv"

	"github.com/vegidio/mediaorient"
)

func TextFilesMessage(amount int) string {
	return fmt.Sprintf("⏳ Analysing %s files to calculate orientation...", green.Render(strconv.Itoa(amount)))
}

func TextDirectoryMessage(directory string) string {
	return fmt.Sprintf("⏳ Analysing directory %s to calculate orientation...", green.Render(directory))
}

func PrintError(message string, a ...interface{}) {
	format := fmt.Sprintf(message, a...)
	fmt.Printf("\n🧨 %s\n", red.Render(format))
}

func PrintReport(media []mediaorient.Media) {
	for _, m := range media {
		if m.Rotation > 0 {
			if m.Type == "image" {
				fmt.Printf("📸 The %s %s is rotated clockwise by %s\n",
					orange.Render(m.Type), bold.Render(m.Name), getRotationColor(m.Rotation))
			} else {
				fmt.Printf("🎬 The %s %s is rotated clockwise by %s\n",
					magenta.Render(m.Type), bold.Render(m.Name), getRotationColor(m.Rotation))
			}
		}
	}
}

// region - Private functions

func getRotationColor(rotation int) string {
	switch rotation {
	case 90:
		return yellow.Render("90º")
	case 180:
		return blue.Render("180º")
	case 270:
		return red.Render("270º")
	default:
		return bold.Render("0º")
	}
}

// endregion
