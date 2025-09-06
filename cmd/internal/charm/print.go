package charm

import (
	"fmt"
	"strconv"
)

func PrintError(message string, a ...interface{}) {
	format := fmt.Sprintf(message, a...)
	fmt.Printf("\n🧨 %s\n", red.Render(format))
}

func PrintCalculateFiles(amount int) {
	fmt.Printf("\n⏳ Determining the orientation in %s files\n", green.Render(strconv.Itoa(amount)))
}

func PrintCalculateDirectory(dir string) {
	fmt.Printf("\n⏳ Determining the orientation in the directory %s\n", green.Render(dir))
}
