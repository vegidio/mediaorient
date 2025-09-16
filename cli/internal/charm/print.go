package charm

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/samber/lo"
	"github.com/vegidio/mediaorient"
)

func PrintDownloadModel() {
	fmt.Printf("\n🌟 %s\n", yellow.Render("Downloading the neural network model..."))
}

func PrintCalculateFiles(amount int) {
	fmt.Printf("⏳ Calculating orientation in %s files\n", green.Render(strconv.Itoa(amount)))
}

func PrintCalculateDirectory(dir string) {
	fmt.Printf("⏳ Calculating orientation in the directory %s\n", green.Render(dir))
}

func PrintError(message string, a ...interface{}) {
	format := fmt.Sprintf(message, a...)
	fmt.Printf("🧨 %s\n", red.Render(format))
}

func PrintReport(media []mediaorient.Media, includeZero bool) {
	groups := lo.GroupBy(media, func(m mediaorient.Media) int {
		return m.Rotation
	})

	if !includeZero {
		delete(groups, 0)
	}

	angles := lo.Keys(groups)
	sort.Ints(angles)

	for _, k := range angles {
		fmt.Printf("\nMedia rotated %s clockwise:\n", getRotationColor(k))

		for _, m := range groups[k] {
			fmt.Printf("  -> %s\n", bold.Render(m.Name))
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
