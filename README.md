# k8s-simple-app

Simple Go application that shows how to use tools: [ytt](https://get-ytt.io), [kbld](https://get-kbld.io), [kapp](https://get-kapp.io) and [kwt](https://github.com/k14s/kapp) from k14s org.

## Install k14s Tools

Head over to [k14s.io](https://k14s.io/) for installation instructions.

## Deploy

### Step 1

Introduces [kapp](https://get-kapp.io) for deploying k8s resources.

```bash
kapp deploy -a simple-app -R -f config-step-1-minimal/
kapp inspect -a simple-app --tree
kapp logs -f -a simple-app
```

Modify `metadata.name` for Deployment resource in `config-step-1-minimal/config.yml`, and run:

```bash
kapp deploy -a simple-app -R -f config-step-1-minimal/ --diff-changes
```

### Step 2

Introduces [ytt](https://get-ytt.io) templating for more flexible configuration.

```bash
ytt template -R -f config-step-2-template/ | kapp deploy -a simple-app -f- --diff-changes -y
```

### Step 3

Introduces [kbld](https://get-kbld.io) functionality for building images from source code. This step requires Minikube. If Minikube is not available, skip to the next step.

```bash
eval $(minikube docker-env)
ytt template -R -f config-step-3-build-local/ | kbld -f- | kapp deploy -a simple-app -f- --diff-changes -y
```

### Step 4

Introduces [kbld](https://get-kbld.io) functionality to push to remote registries. This step can works against Minikube or remote cluster.

```bash
docker login -u dkalinin -p ...
ytt template -R -f config-step-4-build-and-push/ -v push_images=true -v push_images_repo=docker.io/dkalinin/k8s-simple-app | kbld -f- | kapp deploy -a simple-app -f- --diff-changes -y
```

## Directory Layout

- [`app.go`](app.go): simple Go HTTP server
- [`Dockerfile`](Dockerfile): Dockerfile to build Go app
- `config-step-1-minimal/`
  - [`config.yml`](config-step-1-minimal/config.yml): basic k8s Service and Deployment configuration for the app
- `config-step-2-template/`
  - [`config.yml`](config-step-2-template/config.yml): slightly modified configuration to use `ytt` features, such as data module and functions
  - [`values.yml`](config-step-2-template/values.yml): defines two values used in `config.yml`
- `config-step-3-build-local/`
  - [`config.yml`](config-step-3-build-local/build.yml): tells `kbld` about how to build container image from source (app.go + Dockerfile)
  - [`config.yml`](config-step-3-build-local/config.yml): ~same as prev step~
  - [`values.yml`](config-step-3-build-local/values.yml): ~same as prev step~
- `config-step-4-build-and-push/`
  - [`config.yml`](config-step-4-build-and-push/build.yml): tells `kbld` about how to build container image and push to remote registry
  - [`config.yml`](config-step-4-build-and-push/config.yml): ~same as prev step~
  - [`values.yml`](config-step-4-build-and-push/values.yml): defines shared configuration, including configuration for pushing container images
