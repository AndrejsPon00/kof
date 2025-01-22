# Development

## Prerequisites

* Make sure that you have a [K0rdent](https://github.com/K0rdent/kcm/blob/main/docs/dev.md) installed.
It's your "mothership" cluster.

* DNS to test kof with managed clusters installations

Install cli tools

```bash
make cli-install
```

## Local deployment (without K0rdent)

Install into local clusters these helm charts using Makefile

```bash
make dev-storage-deploy
make dev-operators-deploy
make dev-collectors-deploy
```

When everything up and running you can connect to grafana using port-forwarding

```bash
kubectl --namespace kof port-forward svc/grafana-vm-service 3000:3000
```

## Managed clusters deployment with K0rdent in AWS

Install helm charts into a local registry

```bash
make helm-push
```

Define your DNS zone (automatically managed by external-dns)

```bash
KOF_DNS="dev.example.net"
```

Install "mothership" helm chart into your "mothership" cluster


```bash
make dev-ms-deploy-aws
```

Create "storage" managed cluster using KCM

```bash
make dev-storage-deploy-cloud  # to deploy using AWS standalone template
make dev-storage-deploy-cloud CLOUD_CLUSTER_TEMPLATE=aws-eks # to deploy using EKS cluster template
```

Create "managed" managed cluster using KCM

```bash
make dev-managed-deploy-cloud  # to deploy using AWS standalone template
make dev-managed-deploy-cloud CLOUD_CLUSTER_TEMPLATE=aws-eks # to deploy using EKS cluster template
```

When everything up and running you can connect to grafana using port-forwarding from your "mothership" cluster

```bash
kubectl --namespace kof port-forward svc/grafana-vm-service 3000:3000
```
