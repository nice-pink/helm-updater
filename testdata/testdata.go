package testdata

const (
	DEPLOYMENT_APP string = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-ingress-controller
spec:
  template:
    spec:
      # Add OAuth2 proxy sidecar container configuration
      containers:
      - name: nginx-ingress-controller
        image: quay.io/oauth2-proxy/oauth2-proxy:v8.0.1
`

	KUSTOMIZE_APP string = `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

helmCharts:
- name: nginx-ingress-controller
  repo: https://charts.bitnami.com/bitnami
  version: 8.0.1
  releaseName: nginx-ingress-controller
  valuesFile: values.yaml
`

	TERRAFORM_APP string = `resource "helm_release" "nginx" {
  name       = "my-nginx"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "nginx-ingress-controller"
  version    = "8.0.1"

  set {
    name  = "service.type"
    value = "ClusterIP"
  }
}
`

	ARGO_APP string = `apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: nginx-ingress-controller
  namespace: argocd
spec:
  project: ops
  source:
    chart: nginx-ingress-controller
    repoURL: https://charts.bitnami.com/bitnami
    targetRevision: 8.0.1
    helm:
      releaseName: nginx-ingress-controller
  destination:
    server: "https://kubernetes.default.svc"
    namespace: ops
`
)
