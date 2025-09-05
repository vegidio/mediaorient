package mediaorient

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	downloader "github.com/vegidio/ffmpeg-downloader"
)

// Holds the file path to the FFmpeg binary. Defaults to the system-installed path if not explicitly set.
var ffmpegPath = getFFmpegPath("mediaorient")

func getFFmpegPath(configName string) string {
	installed := downloader.IsSystemInstalled()
	if installed {
		return ""
	}

	path, installed := downloader.IsStaticallyInstalled(configName)
	if installed {
		return path
	}

	path, err := downloader.Download(configName)
	if err != nil {
		return ""
	}

	return path
}

func listFiles(directory string, mediaTypes []string, recursive bool) ([]string, error) {
	files := make([]string, 0)

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip on error
		}

		// If this is a directory below the root, and we're not in recursive mode, skip it
		if d.IsDir() && !recursive && path != directory {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if slices.Contains(mediaTypes, ext) {
				files = append(files, path)
			}
		}

		return nil
	})

	if err != nil {
		return files, err
	}

	return files, nil
}

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
