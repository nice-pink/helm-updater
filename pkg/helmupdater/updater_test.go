package helmupdater

import (
	"os"
	"testing"

	"github.com/nice-pink/helm-updater/pkg/models"
	"github.com/nice-pink/helm-updater/testdata"
)

const (
	TEST_FILE_BASE_PATH string = "../../testdata"
)

func TestUpdateVersionKustomize(t *testing.T) {
	app := KUSTOMIZE_APP

	for _, path := range app.Paths {
		// write test file
		outputPath := TEST_FILE_BASE_PATH + "/" + path
		err := os.WriteFile(outputPath, []byte(testdata.KUSTOMIZE_APP), 0666)
		if err != nil {
			t.Error("cannot write test file", outputPath)
		}

		validateUpdate(app, path, outputPath, t)
	}
}

func TestUpdateVersionDeployment(t *testing.T) {
	app := K8S_APP

	for _, path := range app.Paths {
		// write test file
		outputPath := TEST_FILE_BASE_PATH + "/" + path
		err := os.WriteFile(outputPath, []byte(testdata.K8S_APP), 0666)
		if err != nil {
			t.Error("cannot write test file", outputPath)
		}

		validateUpdate(app, path, outputPath, t)
	}

}

func validateUpdate(app models.App, appPath, outputPath string, t *testing.T) {
	// new test version
	version := "100.0.0"

	// don't update manifest
	replaced, newAvailable, err := updateVersionInPath(app, appPath, version, TEST_FILE_BASE_PATH)
	if err != nil {
		t.Error("update version error", err)
	}
	if replaced {
		t.Error("version was replaced but AutoUpdate=false")
	}
	if !newAvailable {
		t.Error("new version is available but was not identified")
	}

	// write manifest (autoupdate)
	app.AutoUpdate = true
	replaced, newAvailable, err = updateVersionInPath(app, appPath, version, TEST_FILE_BASE_PATH)
	if err != nil {
		t.Error("update version error", err)
	}
	if !replaced {
		t.Error("version was not replaced but AutoUpdate=true")
	}
	if !newAvailable {
		t.Error("new version is available but was not identified")
	}

	// check file
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Error("cannot read manifest file", outputPath)
	}

	sData := string(data)
	v := getCurrentVersion(app, sData)
	if v != version {
		t.Error("versions in manifest don't match")
	}

	err = os.Remove(outputPath)
	if err != nil {
		t.Error("cannot delete manifest file", outputPath)
	}
}
