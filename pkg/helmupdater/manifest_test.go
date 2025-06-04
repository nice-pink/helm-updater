package helmupdater

import (
	"strings"
	"testing"

	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/helm-updater/pkg/models"
	"github.com/nice-pink/helm-updater/testdata"
)

const (
	TEST_VERSION_CURRENT string = "8.0.1"
	TEST_VERSION_NEW     string = "100.0.0"
	BASE_FOLDER_TEST     string = "../../testdata"
)

var (
	TERRAFORM_APP = models.App{
		AutoUpdate:             false,
		System:                 models.Terraform,
		Repo:                   "https://charts.dexidp.io",
		Name:                   "nginx-ingress-controller",
		Version:                "*",
		ContainerImage:         "ignored",
		ContainerVersionPrefix: "",
		RepoUsername:           "",
		RepoPassword:           "",
		Path:                   "main.tf",
	}

	KUSTOMIZE_APP = models.App{
		AutoUpdate:             false,
		System:                 models.Kustomize,
		Repo:                   "https://charts.dexidp.io",
		Name:                   "nginx-ingress-controller",
		Version:                "*",
		ContainerImage:         "ignored",
		ContainerVersionPrefix: "",
		RepoUsername:           "",
		RepoPassword:           "",
		Path:                   "kustomization.yaml",
	}

	DEPLOYMENT_APP = models.App{
		AutoUpdate:             false,
		System:                 models.Deployment,
		Repo:                   "https://charts.dexidp.io",
		Name:                   "nginx-ingress-controller",
		Version:                "*",
		ContainerImage:         "quay.io/oauth2-proxy/oauth2-proxy",
		ContainerVersionPrefix: "v",
		RepoUsername:           "",
		RepoPassword:           "",
		Path:                   "deployment.yaml",
	}

	ARGO_APP = models.App{
		AutoUpdate:             false,
		System:                 models.ArgoCd,
		Repo:                   "https://charts.dexidp.io",
		Name:                   "nginx-ingress-controller",
		Version:                "*",
		ContainerImage:         "quay.io/oauth2-proxy/oauth2-proxy",
		ContainerVersionPrefix: "",
		RepoUsername:           "",
		RepoPassword:           "",
		Path:                   "argo.yaml",
	}
)

// update

func TestUpdateTerraform(t *testing.T) {
	data, replaced, err := Update(TERRAFORM_APP, TEST_VERSION_NEW, "", testdata.TERRAFORM_APP)
	if err != nil {
		t.Fatal(err, "update error")
	}

	if !replaced {
		t.Fatal("did not update content")
	} //else {
	// 	log.Info(string(data))
	// 	t.Fatal("no")
	// }

	sData := string(data)
	if !strings.Contains(sData, TEST_VERSION_NEW) {
		log.Error("does not contain new version")
	}
}

func TestUpdateKustomize(t *testing.T) {
	data, replaced, err := Update(KUSTOMIZE_APP, TEST_VERSION_NEW, "", testdata.KUSTOMIZE_APP)
	if err != nil {
		t.Fatal(err, "update error")
	}

	if !replaced {
		t.Fatal("did not update content")
	} //else {
	// 	log.Info(string(data))
	// 	t.Fatal("no")
	// }

	sData := string(data)
	if !strings.Contains(sData, TEST_VERSION_NEW) {
		log.Error("does not contain new version")
	}
}

func TestUpdateDeployment(t *testing.T) {
	data, replaced, err := Update(DEPLOYMENT_APP, TEST_VERSION_NEW, "", testdata.DEPLOYMENT_APP)
	if err != nil {
		t.Fatal(err, "update error")
	}

	if !replaced {
		t.Fatal("did not update content")
	} //else {
	// 	log.Info(string(data))
	// 	t.Fatal("no")
	// }

	sData := string(data)
	if !strings.Contains(sData, TEST_VERSION_NEW) {
		log.Error("does not contain new version")
	}
}

func TestUpdateArgo(t *testing.T) {
	data, replaced, err := Update(ARGO_APP, TEST_VERSION_NEW, "", testdata.ARGO_APP)
	if err != nil {
		t.Fatal(err, "update error")
	}

	if !replaced {
		t.Fatal("did not update content")
	} //else {
	// 	log.Info(string(data))
	// 	t.Fatal("no")
	// }

	sData := string(data)
	if !strings.Contains(sData, TEST_VERSION_NEW) {
		log.Error("does not contain new version")
	}
}

// get version

func TestGetCurrentVersionTerraform(t *testing.T) {
	v := getCurrentVersion(TERRAFORM_APP, testdata.TERRAFORM_APP)
	if v != TEST_VERSION_CURRENT {
		t.Error("not requested version", v, TEST_VERSION_CURRENT)
	}
}

func TestGetCurrentVersionDeployment(t *testing.T) {
	v := getCurrentVersion(DEPLOYMENT_APP, testdata.DEPLOYMENT_APP)
	if v != TEST_VERSION_CURRENT {
		t.Error("not requested version", v, TEST_VERSION_CURRENT)
	}
}

func TestGetCurrentVersionKustomize(t *testing.T) {
	v := getCurrentVersion(KUSTOMIZE_APP, testdata.KUSTOMIZE_APP)
	if v != TEST_VERSION_CURRENT {
		t.Error("not requested version", v, TEST_VERSION_CURRENT)
	}
}

func TestGetCurrentVersionArgo(t *testing.T) {
	v := getCurrentVersion(ARGO_APP, testdata.ARGO_APP)
	if v != TEST_VERSION_CURRENT {
		t.Error("not requested version", v, TEST_VERSION_CURRENT)
	}
}

// update version

func TestUpdateVersionSegment(t *testing.T) {
	manifest := testdata.TERRAFORM_APP

	// same
	app := TERRAFORM_APP
	v, updated := updateVersionSegment(app, TEST_VERSION_CURRENT, manifest)
	if updated {
		t.Error("the versions are already the same")
	}
	if v != TEST_VERSION_CURRENT {
		t.Error("the versions should be equal")
	}

	// all
	v, updated = updateVersionSegment(app, TEST_VERSION_NEW, manifest)
	if !updated {
		t.Error("update all versions")
	}
	if v != TEST_VERSION_NEW {
		t.Error("the version should be the new one")
	}

	// fix major
	newVersion := "8.1.3"
	app.Version = "8"
	v, updated = updateVersionSegment(app, newVersion, manifest)
	if !updated {
		t.Error("update all versions")
	}
	if v != newVersion {
		t.Error("wrong version. app version:", app.Version, ", new version:", newVersion)
	}
	// star
	app.Version = "8.*"
	v, updated = updateVersionSegment(app, newVersion, manifest)
	if !updated {
		t.Error("update all versions")
	}
	if v != newVersion {
		t.Error("wrong version. app version:", app.Version, ", new version:", newVersion)
	}

	// fix minor
	app.Version = "8.0"
	v, updated = updateVersionSegment(app, newVersion, manifest)
	if !updated {
		t.Error("update all versions")
	}
	if v != "8.0.3" {
		t.Error("wrong version. app version:", app.Version, ", new version:", newVersion)
	}
	// star
	app.Version = "8.0.*"
	v, updated = updateVersionSegment(app, newVersion, manifest)
	if !updated {
		t.Error("update all versions")
	}
	if v != "8.0.3" {
		t.Error("wrong version. app version:", app.Version, ", new version:", newVersion)
	}
}
