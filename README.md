![logo](logos/CarvelLogo.png)

# k8s-simple-app-example

Example repo shows how to use tools from carvel dev: [ytt](https://carvel.dev/ytt), [kbld](https://carvel.dev/kbld), [kapp](https://carvel.dev/kapp) and [kwt](https://github.com/vmware-tanzu/carvel-kwt) to work with a simple Go app on Kubernetes.

Associated blog post: [Deploying Kubernetes Applications with ytt, kbld, and kapp](https://carvel.dev/blog/deploying-apps-with-ytt-kbld-kapp/).

## Install Carvel Tools

Head over to [carvel.dev](https://carvel.dev/) for installation instructions.

## Deploying Application

Each top level step has an associated `config-step-*` directory. Refer to [Directory Layout](#directory-layout) for details about files.

### Step 1: Deploying application

Introduces [kapp](https://carvel.dev/kapp) for deploying k8s resources.

```bash
kapp deploy -a simple-app -f config-step-1-minimal/
kapp inspect -a simple-app --tree
kapp logs -f -a simple-app
```

### Step 1a: Viewing application

Once deployed successfully, you can access frontend service at `127.0.0.1:8080` in your browser via `kubectl port-forward` command:

```bash
kubectl port-forward svc/simple-app 8080:80
```

You will have to restart port forward command after making any changes as pods are recreated. Alternatively consider using [kwt](https://github.com/vmware-tanzu/carvel-kwt) which exposes cluser IP subnets and cluster DNS to your machine and does not require any restarts:

```bash
sudo -E kwt net start
```

and open [`http://simple-app.default.svc.cluster.local/`](http://simple-app.default.svc.cluster.local/).

### Step 1b: Modifying application configuration

Modify `HELLO_MSG` environment value from `stranger` to something else in `config-step-1-minimal/config.yml`, and run:

```bash
kapp deploy -a simple-app -f config-step-1-minimal/ --diff-changes
```

In following steps we'll use `-c` shorthand for `--diff-changes`.

### Step 2: Configuration templating

Introduces [ytt](https://carvel.dev/ytt) templating for more flexible configuration.

```bash
kapp deploy -a simple-app -c -f <(ytt -f config-step-2-template/)
```

ytt provides a way to configure data values from command line as well:

```bash
kapp deploy -a simple-app -c -f <(ytt -f config-step-2-template/ -v hello_msg=another-stranger)
```

New message should be returned from the app in the browser.

### Step 2a: Configuration patching

Introduces [ytt overlays](https://carvel.dev/ytt/docs/latest/lang-ref-ytt-overlay/) to patch configuration without modifying original `config.yml`.

```bash
kapp deploy -a simple-app -c -f <(ytt -f config-step-2-template/ -f config-step-2a-overlays/custom-scale.yml)
```

### Step 2b: Customizing configuration data values per environment

Requires ytt v0.13.0+.

Introduces [use of multiple data values](https://carvel.dev/ytt/docs/latest/ytt-data-values/) to show layering of configuration for different environment without modifying default `values.yml`.

```bash
kapp deploy -a simple-app -c -f <(ytt -f config-step-2-template/ -f config-step-2b-multiple-data-values/)
```

### Step 3: Building container images locally

Introduces [kbld](https://carvel.dev/kbld) functionality for building images from source code. This step requires Minikube. If Minikube is not available, skip to the next step.

```bash
eval $(minikube docker-env)
kapp deploy -a simple-app -c -f <(ytt -f config-step-3-build-local/ | kbld -f-)
```

Note that rerunning above command again should be a noop, given that nothing has changed.

### Step 3a: Modifying application source code

Uncomment `fmt.Fprintf(w, "<p>local change</p>")` line in `app.go`, and re-run above command:

```bash
kapp deploy -a simple-app -c -f <(ytt -f config-step-3-build-local/ | kbld -f-)
```

Observe that new container was built, and deployed. This change should be returned from the app in the browser.

### Step 4: Building and pushing container images to registry

Introduces [kbld](https://carvel.dev/kbld) functionality to push to remote registries. This step can work with Minikube or any remote cluster.

```bash
docker login -u dkalinin -p ...
kapp deploy -a simple-app -c -f <(ytt -f config-step-4-build-and-push/ -v push_images_repo=gcr.io/projectX/k8s-simple-app | kbld -f-)
```

### Step 5: Clean up cluster resources

```bash
kapp delete -a simple-app
```

There is currently no functionality in kbld to remove pushed images from registry.

## Directory Layout

- [`app.go`](app.go): simple Go HTTP server
- [`Dockerfile`](Dockerfile): Dockerfile to build Go app
- `config-step-1-minimal/`
  - [`config.yml`](config-step-1-minimal/config.yml): basic k8s Service and Deployment configuration for the app
- `config-step-2-template/`
  - [`config.yml`](config-step-2-template/config.yml): slightly modified configuration to use `ytt` features, such as data module and functions
  - [`values.yml`](config-step-2-template/values.yml): defines extracted data values used in `config.yml`
- `config-step-2a-overlays/`
  - [`custom-scale.yml`](config-step-2a-overlays/custom-scale.yml): ytt overlay to set number of deployment replicas to 3
- `config-step-3-build-local/`
  - [`build.yml`](config-step-3-build-local/build.yml): tells `kbld` about how to build container image from source (app.go + Dockerfile)
  - [`config.yml`](config-step-3-build-local/config.yml): _same as prev step_
  - [`values.yml`](config-step-3-build-local/values.yml): _same as prev step_
- `config-step-4-build-and-push/`
  - [`build.yml`](config-step-4-build-and-push/build.yml): _same as prev step_
  - [`push.yml`](config-step-4-build-and-push/push.yml): tells `kbld` about how to push container image to remote registry
  - [`config.yml`](config-step-4-build-and-push/config.yml): _same as prev step_
  - [`values.yml`](config-step-4-build-and-push/values.yml): defines shared configuration, including configuration for pushing container images

### Join the Community and Make Carvel Better
Carvel is better because of our contributors and maintainers. It is because of you that we can bring great software to the community.
Please join us during our online community meetings. Details can be found on our [Carvel website](https://carvel.dev/community/).

You can chat with us on Kubernetes Slack in the #carvel channel and follow us on Twitter at @carvel_dev.

Check out which organizations are using and contributing to Carvel: [Adopter's list](https://github.com/vmware-tanzu/carvel/blob/master/ADOPTERS.md)
