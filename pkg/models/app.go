package models

type SystemType string

const (
	Kustomize SystemType = "kustomize" // update kustomize manifest
	ArgoCd    SystemType = "argocd"    // update argocd application manifest
	K8s       SystemType = "k8s"       // update container image in deployment, statefulsets, daemonsets
	Terraform SystemType = "terraform" // terraform module
)

type App struct {
	AutoUpdate             bool
	Name                   string
	Repo                   string
	Private                bool
	Version                string
	Paths                  []string
	System                 SystemType
	ContainerImage         string
	ContainerVersionPrefix string
	RepoUsername           string
	RepoPassword           string
}
