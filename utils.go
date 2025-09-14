package mediaorient

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	downloader "github.com/vegidio/ffmpeg-downloader"
)

// Holds the file path to the FFmpeg binary. Defaults to the system-installed path if not explicitly set.
var ffmpegPath = downloader.GetFFmpegPath("mediaorient")

func loadFrames(directory string) ([]image.Image, error) {
	images := make([]image.Image, 0)

	files, err := os.ReadDir(directory)
	if err != nil {
		return images, err
	}

	for _, file := range files {
		if !file.IsDir() {
			fullPath := filepath.Join(directory, file.Name())

			f, fErr := os.Open(fullPath)
			if fErr != nil {
				return images, fErr
			}

			img, _, imgErr := image.Decode(f)
			if imgErr != nil {
				return images, imgErr
			}

			f.Close()
			images = append(images, img)
		}
	}

	return images, nil
}

func extractFrames(filePath string) ([]image.Image, error) {
	images := make([]image.Image, 0)

	tempDir, err := os.MkdirTemp("", "mediaorient-*")
	if err != nil {
		return images, fmt.Errorf("error creating temp directory: %w", err)
	}

	defer os.RemoveAll(tempDir)

	// Export 1 frame per second
	path := filepath.Join(tempDir, "frame_%04d.jpg")
	command := ffmpeg.Input(filePath).
		Filter("fps", ffmpeg.Args{"1"}).
		Output(path).
		Silent(true)

	if ffmpegPath == "" {
		_ = command.Run()
	} else {
		_ = command.SetFfmpegPath(ffmpegPath).Run()
	}

	images, _ = loadFrames(tempDir)
	if len(images) > 0 {
		return images, nil
	}

	// Failed to export multiple frames, so let's try to export a single frame
	path = filepath.Join(tempDir, "frame.jpg")
	command = ffmpeg.Input(filePath).
		Output(path, ffmpeg.KwArgs{"vframes": 1}).
		Silent(true)

	if ffmpegPath == "" {
		err = command.Run()
	} else {
		err = command.SetFfmpegPath(ffmpegPath).Run()
	}

	if err != nil {
		return images, fmt.Errorf("error exporting video frames from '%s': %w", filePath, err)
	}

	images, err = loadFrames(tempDir)
	if err != nil {
		return images, fmt.Errorf("error loading videos frames from '%s': %w", filePath, err)
	}

	return images, nil
}
