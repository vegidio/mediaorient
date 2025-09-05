package mediaorient

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"

	ort "github.com/yalue/onnxruntime_go"
)

//go:embed model/efficient_net_v2.onnx
var modelBinary []byte
var modelName = "efficient_net_v2.onnx"

func init() {
	// Check if OnnxRuntime and the model are already saved in the user's config directory
	onnxPath, modelPath, exists := hasBinaries("mediaorient")

	if !exists {
		if err := saveBinaries(onnxPath, modelPath); err != nil {
			log.Fatalf("error initializing the app: %v\n", err)
		}
	}

	if err := startRuntime(onnxPath, modelPath); err != nil {
		log.Fatalf("error initializing the app: %v\n", err)
	}
}

func hasBinaries(configName string) (string, string, bool) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("error getting user config directory: %v\n", err)
	}

	fullConfigDir := filepath.Join(configDir, configName)
	onnxPath := filepath.Join(fullConfigDir, libOnnxName)
	modelPath := filepath.Join(fullConfigDir, modelName)

	if _, fErr := os.Stat(onnxPath); fErr != nil {
		return onnxPath, modelPath, false
	}
	if _, fErr := os.Stat(modelPath); fErr != nil {
		return onnxPath, modelPath, false
	}

	return onnxPath, modelPath, true
}

func saveBinaries(onnxPath, modelPath string) error {
	directory := filepath.Dir(onnxPath)
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}

	// Copy the OnnxRuntime library
	f1, err := os.Create(onnxPath)
	if err != nil {
		return err
	}
	defer f1.Close()

	if fErr := os.WriteFile(onnxPath, libOnnxBinary, 0755); fErr != nil {
		return fErr
	}

	// Copy the orientation model
	f2, err := os.Create(modelPath)
	if err != nil {
		return err
	}
	defer f2.Close()

	if fErr := os.WriteFile(modelPath, modelBinary, 0755); fErr != nil {
		return fErr
	}

	return nil
}

func startRuntime(onnxPath, modelPath string) error {
	var err error
	ort.SetSharedLibraryPath(onnxPath)

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
