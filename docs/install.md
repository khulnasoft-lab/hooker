To run Hooker you will first need to configure the [Hooker Configuration File](/hooker/config), which contains all the message routing logic.
After the configuration file is ready, you can run the official Hooker container image: **khulnasoft/hooker:latest**, or compile it from source.

There are different options to mount your customize configuration file to Hooker - if running as a Docker container, then you simply mount the configuration files as a volume mount. If running as a Kubernetes deployment, you will need to mount it as a ConfigMap. See the below usage examples for how to run Hooker on different scenarios.

After Hooker will run, it will expose two endpoints, HTTP and HTTPS. You can send your JSON messages to these endpoints, where they will be delivered to their target system based on the defined rules.

### Docker
To run Hooker as a Docker container, you mount the cfg.yaml to '/config/cfg.yaml' path in the Hooker container.


```bash
docker run -d --name=hooker -v /<path to configuration file>/cfg.yaml:/config/cfg.yaml \
    -e HOOKER_CFG=/config/cfg.yaml -e HOOKER_HTTP=0.0.0.0:8084 -e HOOKER_HTTPS=0.0.0.0:8444 \
    -p 8084:8084 -p 8444:8444 khulnasoft/hooker:latest
```

### Kubernetes
When running Hooker on Kubernetes, the configuration file is passed as a ConfigMap that is mounted to the Hooker pod.


#### Cloud Providers

``` bash
kubectl create -f https://raw.githubusercontent.com/khulnasoft-lab/hooker/main/deploy/kubernetes/hooker.yaml
```

#### Using HostPath

``` bash
kubectl create -f https://raw.githubusercontent.com/khulnasoft-lab/hooker/main/deploy/kubernetes/hostPath/hooker-pv.yaml
```

!!! Note "Persistent Volumes Explained"
    - `hooker-db`: persistent storage directory `/server/database`
    - `hooker-config`: mount the cfg.yaml to a writable directory `/config/cfg.yaml`
    - `hooker-rego-templates`: mount custom rego templates
    - `hooker-rego-filters`: mount custom rego filters
To edit the default Hooker-UI user

```
kubectl -n hooker set env deployment/my-hookerui -e HOOKER_ADMIN_USER=testabc -e HOOKER_ADMIN_PASSWORD=password
```

The Hooker endpoints
```
http://hooker-svc.default.svc.cluster.local:8082
```
```
https://hooker-svc.default.svc.cluster.local:8445
```

The Hooker-UI endpoint
````
http://hooker-ui-svc.default.svc.cluster.local:8000
````

#### Controller/Runner
To use Controller/Runner functionality within Kubernetes, you can follow a reference manifest implementation:
- [Controller](https://github.com/khulnasoft-lab/hooker/blob/main/deploy/kubernetes/hooker-controller.yaml)
- [Runner](https://github.com/khulnasoft-lab/hooker/blob/main/deploy/kubernetes/hooker-runner.yaml)

### Helm
When running Hooker on Kubernetes, the configuration file is passed as a ConfigMap that is mounted to the Hooker pod.

This chart bootstraps a Hooker deployment on a [Kubernetes](https://kubernetes.io/) cluster using the [Helm package manager](https://helm.sh/).

#### Prerequisites
- Kubernetes 1.17+
- Helm 3+

#### Test the Chart Repository

```bash
cd deploy/helm
helm install my-hooker -n hooker --dry-run --set-file applicationConfigPath="../../cfg.yaml" ./hooker
```

#### Installing the Chart from the Source Code

```bash
cd deploy/helm
helm install app --create-namespace -n hooker ./hooker
```

#### Installing from the the Khulnasoft Chart Repository

Let's add the Helm chart and deploy Hooker executing:


```bash
helm repo add khulnasoft-lab https://khulnasoft-lab.github.io/helm-charts/
helm repo update
helm search repo hooker
helm install app --create-namespace -n hooker khulnasoft-lab/hooker
```

Check that all the pods are in Running state:

`kubectl get pods -n hooker`

We check the logs:

```
kubectl logs deployment/my-hookerui -n hooker | head
```

```
kubectl logs statefulsets/my-hooker -n hooker | head
```

#### Delete Chart

```bash
helm -n hooker delete my-hooker
```

#### From Source
Clone and build the project:
```bash
git clone git@github.com:khulnasoft-lab/hooker.git
make build
```
After that, modify the cfg.yaml file and set the 'HOOKER_CFG' environment variable to point to it.
```bash
export HOOKER_CFG=<path to cfg.yaml>
./bin/hooker
```
