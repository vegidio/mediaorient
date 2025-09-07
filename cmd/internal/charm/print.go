package charm

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/samber/lo"
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
	groups := lo.GroupBy(media, func(m mediaorient.Media) int {
		return m.Rotation
	})

	delete(groups, 0)
	angles := lo.Keys(groups)
	sort.Ints(angles)

	for _, k := range angles {
		fmt.Printf("\nMedia rotated %s clockwise:\n", getRotationColor(k))

		for _, m := range groups[k] {
			fmt.Printf("  -> %s is rotated %s\n", bold.Render(m.Name), getRotationColor(m.Rotation))
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
