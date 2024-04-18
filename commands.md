Commands for a k8s operator:
```bash
# create a project
kubebuilder init --domain filip.org --repo github.com/fghimpeteanu-ionos/user-op

# create an API
kubebuilder create api --group user --version v1 --kind User

# how to test the Postgresql Helm chart
helm template -f helm/postgresql-values.yaml test bitnami/postgresql
```