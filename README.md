# K0rdent Observability and FinOps
This repo contains 4 charts to deploy an observability stack using [k0rdent](https://github.com/K0rdent/kcm) and get [OpenTelemetry](https://opentelemetry.io/) data into storage clusters aggregated into single grafana interface.

![alt text](docs/otel.png)

## Mothership chart
* central grafana interface
* promxy to forward calls to multiple downstream regional metrics servers
* local victoriametrics storage for alerting record rules
* k0rdent helmchart definitions and service templates to deploy storage and collectors charts into managedclusters

### Demo deployment
In `demo/demo-mothership-values.yaml` set your target ingress names that you are going to use for your storage clusters, but they can always be changed after the fact

By default the secrets defined in the `values.yaml` are created automatically and propagated to managed clusters using Sveltos cluster profile.

You can retrieve grafana password and username using the following command

```bash
kubectl get secret grafana-admin-credentials -o jsonpath="{.data.GF_SECURITY_ADMIN_USER}" -n kof | base64 -d; echo

kubectl get secret grafana-admin-credentials -o jsonpath="{.data.GF_SECURITY_ADMIN_PASSWORD}" -n kof | base64 -d; echo
```

```bash
helm repo add kof https://k0rdent.github.io/kof/
helm repo update
helm upgrade -i kof-mothership kof/kof-mothership -n kof -f demo/demo-mothership-values.yaml
```

## Storage chart

Deploys metrics and logs [VictoriaMetrics](https://victoriametrics.com/) storages.

* Grafana - storage-cluster scoped Grafana instance, deployed and configured with grafana-operator
* vmcluster - metrics storage, ingestion, querying
* vmlogs - logs storage
* vmauth - auth frontend for metrics and logs ingestion and query services

#### Cluster requirements
- cert-manager
- ingress-nginx

To deploy storage `clusterdeployment` configure desired ingress names for vmauth and regional Grafana in it's values for the `kof-storage` template.
`demo/cluster/aws-storage.yaml` contains example definitions

```bash
kubectl apply -f demo/cluster/aws-storage.yaml
# you can check helm chart deployment status using ClusterSummary object:
kubectl get clustersummaries.config.projectsveltos.io -n kcm-system
```
Once the storage clusterdeploymet is ready - retrieve its kubeconfig and get loadbalancer IP/DNS name for your ingress-nginx service.

```bash
kubectl get secret -n kcm-system aws-storage-kubeconfig -o jsonpath={.data.value} | base64 -d  > /tmp/kcm-aws-storage-kubeconfig.yaml
export KUBECONFIG=/tmp/kcm-aws-storage-kubeconfig.yaml
kubectl get svc -n ingress-nginx ingress-nginx-controller
```

Create secrets for grafana and vmauth according to the names provided in helm values.

With your preffered DNS hosting, set your ingress domains to resolve to that IP/DNS name, that's how the traffic will flow to/from regional cluster. 
To simplify this process it is posssible to enable [external-dns](https://kubernetes-sigs.github.io/external-dns/) helm chart deployment in values.

Once your domain is resolvable your Grafana and vmauth should be accessible.

## Operators chart
* opentelemetry-operator - [OpenTelemetry Operator](https://opentelemetry.io/docs/kubernetes/operator/)
* prometheus-operator-crds - [Prometheus Operator](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus-operator-crds)

This chart pre-installs all required CRDs to create Opentelemetry Collectors for metrics and logs

## Collectors chart
* opentelemetry-collectors - [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) configured to monitor logs and metrics and send them to a storage cluster

To deploy operators and collectors to a `clusterdeployment` configure ingress names for storage vmauth in its values for the `kof-collectors` template.

```
kubectl apply -f demo/cluster/aws-managed.yaml
# you can check helm chart deployment status using ClusterSummary object:
kubectl get clustersummaries.config.projectsveltos.io -n kcm-system
```

Once your managed clusters are up, create secrets for storage cluster authentication, it should start pushing metrics and logs to your storage one, through ingress domain you've configured.
Check your storage cluster's Grafana for results first, then you should be able to see the same cluster in Grafana on the "mothership".

### Scaling up
* Deploy more managed clusters in a single region and point them to the existing storage victoria stack.
* Repeat the previous two steps for each desired region
* Update mothership chart configuration with every deployed regional stack to aggregate the data
