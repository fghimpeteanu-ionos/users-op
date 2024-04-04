Commands for a k8s operator:
```bash
# create a project
kubebuilder init --domain filip.org --repo github.com/fghimpeteanu-ionos/user-op

# create an API
kubebuilder create api --group user --version v1 --kind User


```