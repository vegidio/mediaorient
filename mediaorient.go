package mediaorient

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/disintegration/imaging"
	"github.com/samber/lo"
	ort "github.com/yalue/onnxruntime_go"
)

var session *ort.DynamicAdvancedSession

func CalculateImageOrientation(name string, images []image.Image) (*Media, error) {
	var err error
	size := 384
	rotations := make([]int, 0)

	for _, img := range images {
		input := preprocess(img, size)

		inputTensor, tErr := ort.NewTensor[float32](ort.NewShape(1, 3, int64(size), int64(size)), input)
		if tErr != nil {
			return nil, tErr
		}

		outputTensor, tErr := ort.NewEmptyTensor[float32](ort.NewShape(1, 4))
		if tErr != nil {
			inputTensor.Destroy()
			return nil, tErr
		}

		if err = session.Run([]ort.Value{inputTensor}, []ort.Value{outputTensor}); err != nil {
			inputTensor.Destroy()
			outputTensor.Destroy()
			return nil, err
		}

		logits := outputTensor.GetData()
		probs := softmax(logits)
		bestIdx, best := 0, probs[0]
		for i := 1; i < len(probs); i++ {
			if probs[i] > best {
				best = probs[i]
				bestIdx = i
			}
		}

		rotation := []int{0, 90, 180, 270}[bestIdx]
		rotations = append(rotations, rotation)

		inputTensor.Destroy()
		outputTensor.Destroy()
	}

	media := createMedia(name, images, rotations)
	return media, nil
}

func CalculateFileOrientation(filePath string) (*Media, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file '%s': %w", filePath, err)
	}

	defer file.Close()

	ext := strings.ToLower(filepath.Ext(file.Name()))
	images := make([]image.Image, 0)

	if slices.Contains(validImageTypes, ext) {
		img, _, imgErr := image.Decode(file)
		if imgErr != nil {
			return nil, imgErr
		}

		images = append(images, img)

	} else if slices.Contains(validVideoTypes, ext) {
		videos, vidErr := extractFrames(file.Name())
		if vidErr != nil {
			return nil, vidErr
		}

		images = append(images, videos...)
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("no valid images found in file '%s'", filePath)
	}

	return CalculateImageOrientation(filePath, images)
}

// region - Private functions

func softmax(x []float32) []float32 {
	maxValue := x[0]
	for _, v := range x {
		if v > maxValue {
			maxValue = v
		}
	}
	sum := float64(0)
	out := make([]float32, len(x))
	for i, v := range x {
		e := math.Exp(float64(v - maxValue))
		out[i] = float32(e)
		sum += e
	}
	for i := range out {
		out[i] /= float32(sum)
	}
	return out
}

func preprocess(img image.Image, size int) []float32 {
	// Center-crop to square and resize (closest to common ImageNet eval)
	square := imaging.Fill(img, size, size, imaging.Center, imaging.Lanczos)

	b := square.Bounds()
	w, h := b.Dx(), b.Dy()
	mean := [3]float32{0.485, 0.456, 0.406} // ImageNet mean
	std := [3]float32{0.229, 0.224, 0.225}  // ImageNet std

	data := make([]float32, 3*w*h) // NCHW
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, bl, _ := square.At(x, y).RGBA()
			rf := float32(uint8(r>>8)) / 255.0
			gf := float32(uint8(g>>8)) / 255.0
			bf := float32(uint8(bl>>8)) / 255.0
			i := y*w + x
			data[0*w*h+i] = (rf - mean[0]) / std[0]
			data[1*w*h+i] = (gf - mean[1]) / std[1]
			data[2*w*h+i] = (bf - mean[2]) / std[2]
		}
	}

	return data
}

func createMedia(name string, images []image.Image, rotations []int) *Media {
	mediaType := "image"
	if len(images) > 1 {
		mediaType = "video"
	}

	rotation := lo.MaxBy(lo.Entries(lo.CountValues(rotations)), func(a, b lo.Entry[int, int]) bool {
		return a.Value > b.Value
	})

	return &Media{
		Name:     name,
		Type:     mediaType,
		Rotation: rotation.Key,
		Frames:   images,
		Width:    images[0].Bounds().Dx(),
		Height:   images[0].Bounds().Dy(),
	}
}

// endregion
