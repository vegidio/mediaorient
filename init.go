package mediaorient

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v74/github"
	"github.com/samber/lo"
	"github.com/vegidio/go-sak/crypto"
	"github.com/vegidio/go-sak/fs"
	gitsak "github.com/vegidio/go-sak/github"
	"github.com/vegidio/go-sak/memo"
	ort "github.com/yalue/onnxruntime_go"
)

var modelName = "image_orientation.onnx"

// Initialize sets up the media orientation detection system by configuring the ONNX runtime and ensuring all required
// binaries are available.
//
// # Returns:
//   - error: nil on success, or an error describing what went wrong during initialization.
//
// # Example:
//
//	if err := mediaorient.Initialize(); err != nil {
//	    log.Fatal("Failed to initialize media orientation detection:", err)
//	}
func Initialize() error {
	// Install the ONNX runtime if it's not already installed
	if yes := shouldInstallRuntime(); yes {
		if err := installRuntime(); err != nil {
			return err
		}
	}

	// Download the model if it's not already present
	if url, yes := shouldDownloadModel(); yes {
		if err := downloadModel(url); err != nil {
			return err
		}
	}

	// Initialize the ONNX runtime
	if err := startRuntime(); err != nil {
		return err
	}

	return nil
}

// Destroy cleans up resources used by the media orientation detection system.
//
// This function performs cleanup operations to free memory and resources allocated during initialization. It should be
// called when the application is shutting down or when orientation detection functionality is no longer needed.
//
// It's recommended to use this function with defer for proper cleanup:
//
// # Example:
//
//	if err := mediaorient.Initialize(); err != nil {
//	    log.Fatal("Initialization failed:", err)
//	}
//	defer mediaorient.Destroy() // Ensure cleanup on exit
func Destroy() {
	session.Destroy()
	ort.DestroyEnvironment()
}

// region - Private functions

func shouldInstallRuntime() bool {
	configDir, err := fs.MkUserConfigDir("mediaorient")
	if err != nil {
		log.Fatalf("error getting user config directory: %v\n", err)
	}

	runtimePath := filepath.Join(configDir, libOnnxName)
	_, err = os.Stat(runtimePath)
	return os.IsNotExist(err)
}

func installRuntime() error {
	file, err := fs.MkUserConfigFile("mediaorient", libOnnxName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(libOnnxBinary)
	if err != nil {
		return err
	}

	return nil
}

func shouldDownloadModel() (string, bool) {
	configDir, err := fs.MkUserConfigDir("mediaorient")
	if err != nil {
		log.Fatalf("error getting user config directory: %v\n", err)
	}

	url, remoteHash, err := getLatestModel()
	if err != nil {
		log.Fatalf("error downloading the latest model: %v\n", err)
	}

	modelPath := filepath.Join(configDir, modelName)
	if _, fErr := os.Stat(modelPath); os.IsNotExist(fErr) {
		// The model is not present, so we must download it
		return url, true
	}

	localHash, err := crypto.Sha256File(modelPath)
	if err != nil {
		log.Fatalf("error getting the model signature: %v\n", err)
	}

	return url, localHash != remoteHash
}

func downloadModel(url string) error {
	file, err := fs.MkUserConfigFile("mediaorient", modelName)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func startRuntime() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	runtimePath := filepath.Join(configDir, "mediaorient", libOnnxName)
	modelPath := filepath.Join(configDir, "mediaorient", modelName)

	ort.SetSharedLibraryPath(runtimePath)
	if err = ort.InitializeEnvironment(); err != nil {
		return err
	}

	session, err = ort.NewDynamicAdvancedSession(modelPath, []string{"input"}, []string{"output"}, nil)
	if err != nil {
		ort.DestroyEnvironment()
		return err
	}

	return nil
}

func getLatestModel() (string, string, error) {
	cachePath, err := fs.MkUserConfigDir("mediaorient", "cache")
	if err != nil {
		return "", "", err
	}

	opts := memo.CacheOpts{MaxEntries: 100, MaxCapacity: 1024 * 1024}
	m, err := memo.NewDiskOnly(cachePath, opts)
	if err != nil {
		return "", "", err
	}
	defer m.Close()

	ctx := context.Background()
	key := memo.KeyFrom("getLatestModel")
	ttl := time.Hour * 24 * 7 // Cache for 1 week

	release, err := memo.Do(m, ctx, key, ttl, func(ctx context.Context) (*github.RepositoryRelease, error) {
		r, gErr := gitsak.GetLatestRelease("vegidio", "mediaorient")
		if gErr != nil {
			return nil, gErr
		}

		return r, nil
	})

	if err != nil {
		return "", "", err
	}

	asset, found := lo.Find(release.Assets, func(item *github.ReleaseAsset) bool {
		return item.GetName() == modelName
	})

	if !found {
		return "", "", fmt.Errorf("model not found in latest release")
	}

	url := fmt.Sprintf("https://github.com/vegidio/mediaorient/releases/download/%s/image_orientation.onnx",
		release.GetName())
	hash := strings.TrimPrefix(asset.GetDigest(), "sha256:")

	return url, hash, nil
}

// endregion
