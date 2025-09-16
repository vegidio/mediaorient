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

// Initialize sets up the media orientation detection system by ensuring all required dependencies are available.
//
// This function performs a three-step initialization process:
//  1. Installs the ONNX runtime library if not already present in the user's config directory
//  2. Downloads the latest image orientation model if needed or outdated
//  3. Initializes the ONNX runtime session with the downloaded model
//
// The name parameter specifies the application name used to create a dedicated config directory under the user's
// standard configuration path (e.g., ~/.config/name on Linux). It's important that you reuse the same name on later
// calls to Initialize() to ensure that the same config directory is used.
//
// The onDownload callback function is called when the model needs to be downloaded, allowing the caller to notify users
// about the download process (e.g., show a progress indicator).
//
// Returns an error if any step fails, including:
//   - Unable to create config directories or files
//   - Network errors during model download
//   - ONNX runtime initialization failures
//
// # Example:
//
//	err := mediaorient.Initialize("myapp", func() {
//	    fmt.Println("Downloading model, please wait...")
//	})
//	if err != nil {
//	    log.Fatal("Failed to initialize:", err)
//	}
//	defer mediaorient.Destroy() // Clean up resources
func Initialize(name string, onDownload func()) error {
	// Install the ONNX runtime if it's not already installed
	if yes := shouldInstallRuntime(name); yes {
		if err := installRuntime(name); err != nil {
			return err
		}
	}

	// Download the model if it's not already present
	if url, yes := shouldDownloadModel(name); yes {
		// Notify the user that the model is being downloaded
		onDownload()

		if err := downloadModel(url, name); err != nil {
			return err
		}
	}

	// Initialize the ONNX runtime
	if err := startRuntime(name); err != nil {
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

func shouldInstallRuntime(name string) bool {
	configDir, err := fs.MkUserConfigDir(name)
	if err != nil {
		log.Fatalf("error getting user config directory: %v\n", err)
	}

	runtimePath := filepath.Join(configDir, libOnnxName)
	_, err = os.Stat(runtimePath)
	return os.IsNotExist(err)
}

func installRuntime(name string) error {
	file, err := fs.MkUserConfigFile(name, libOnnxName)
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

func shouldDownloadModel(name string) (string, bool) {
	configDir, err := fs.MkUserConfigDir(name)
	if err != nil {
		log.Fatalf("error getting user config directory: %v\n", err)
	}

	url, remoteHash, err := getLatestModel(name)
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

func downloadModel(url string, name string) error {
	file, err := fs.MkUserConfigFile(name, modelName)
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

func startRuntime(name string) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	runtimePath := filepath.Join(configDir, name, libOnnxName)
	modelPath := filepath.Join(configDir, name, modelName)

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

func getLatestModel(name string) (string, string, error) {
	cachePath, err := fs.MkUserConfigDir(name, "cache")
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
