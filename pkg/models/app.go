package models

type SystemType string

const (
	Kustomize  SystemType = "kustomize"  // update kustomize manifest
	ArgoCd     SystemType = "argocd"     // update argocd application manifest
	Deployment SystemType = "deployment" // update container image in deployment
	Terraform  SystemType = "terraform"  // terraform module
	Helm       SystemType = "helm"       // direct helm release - Still not implemented!!!
)

type App struct {
	AutoUpdate             bool
	Name                   string
	Repo                   string
	Private                bool
	Version                string
	Path                   string
	System                 SystemType
	ContainerImage         string
	ContainerVersionPrefix string
	RepoUsername           string
	RepoPassword           string
}
