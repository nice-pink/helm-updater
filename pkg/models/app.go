package models

type SystemType string

const (
	Kustomize SystemType = "kustomize" // update kustomize manifest
	ArgoCd    SystemType = "argocd"    // update argocd application manifest
	Helm      SystemType = "helm"      // direct helm release
)

type App struct {
	Name         string
	Repo         string
	Private      bool
	Version      string
	Path         string
	System       SystemType
	RepoUsername string
	RepoPassword string
}

type Config struct {
	Apps       []App
	BaseFolder string
}
