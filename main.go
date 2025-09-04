package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"math"
	"os"

	"github.com/disintegration/imaging"
	ort "github.com/yalue/onnxruntime_go"
)

func main() {
	size := 384
	model := "model/efficient_net_v2.onnx"

	ort.SetSharedLibraryPath("libs/darwin_arm64/libonnxruntime.dylib")

	if err := ort.InitializeEnvironment(); err != nil {
		log.Fatalf("Initialize ORT failed: %v", err)
	}
	defer ort.DestroyEnvironment()

	f, err := os.Open("assets/image_0.jpg")
	if err != nil {
		log.Fatalf("open image: %v", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("decode image: %v", err)
	}

	input := preprocess(img, size)
	inputTensor, err := ort.NewTensor[float32](ort.NewShape(1, 3, int64(size), int64(size)), input)
	if err != nil {
		log.Fatalf("make input tensor: %v", err)
	}
	defer inputTensor.Destroy()

	outputTensor, err := ort.NewEmptyTensor[float32](ort.NewShape(1, 4))
	if err != nil {
		log.Fatalf("make output tensor: %v", err)
	}
	defer outputTensor.Destroy()

	sess, err := ort.NewAdvancedSession(model,
		[]string{"input"}, []string{"output"},
		[]ort.Value{inputTensor}, []ort.Value{outputTensor}, nil)

	if err != nil {
		log.Fatalf("create session: %v", err)
	}
	defer sess.Destroy()

	if err := sess.Run(); err != nil {
		log.Fatalf("inference failed: %v", err)
	}

	logits := outputTensor.GetData() // length 4
	probs := softmax(logits)
	bestIdx, best := 0, probs[0]
	for i := 1; i < len(probs); i++ {
		if probs[i] > best {
			best = probs[i]
			bestIdx = i
		}
	}

	rotation := []int{0, 90, 180, 270}[bestIdx]
	fmt.Printf("%d\n", rotation)
}

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
	// Center-crop to square and resize (closest to common ImageNet eval).
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
