# user-op
This is a simple example of a Kubernetes operator written in Go. The sole purpose of this operator is
to add a user in a Database table called `users`

## Description
The milestones of this project are:
- [x] Create a Custom Resource Definition (CRD) for the user in order to create CRs (e.g. User CR)
- [ ] Create a database table called `users` in a PostGreSQL database. The database should be installed via Helm.
- [ ] Add logic in the controller to add a user in the `users` table when a User CR is created
- [ ] Deploy the controller to a Kubernetes cluster and integrate it with the local database.
- [ ] Add a second controller that will take care of creating a PostgresSQL Database

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -k config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/user-op:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/user-op:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the DB helm chart
```sh
make install-db
````

2. Connect to the DB
```sh
kubectl port-forward --namespace default svc/my-db-postgresql 25432:5432
# authenticate with postgres user
PGPASSWORD=pg-pretest psql --host 127.0.0.1 -U postgres -d user-db -p 25432
```

3. Install the CRDs into the cluster:
```sh
make install
```

4. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):
```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

5. Create a User CR:
```sh
kubectl apply -f config/samples/user.yaml 
```

6. List resourceds
```sh
watch "kubectl get all,configmap,secrets,pv,pvc"
```

7. Cleanup
```sh
# uninstall PostgreSQL
make uninstall-db

# Remove the PVC and PV created for the PostgreSQL
```

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

