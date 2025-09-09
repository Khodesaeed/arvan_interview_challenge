# GitOps with ArgoCD

This directory contains the GitOps configuration for deploying applications using the **ArgoCD App of Apps pattern**.  
The root application manages multiple child applications defined under [`argo-apps/`](./argo-apps).

## Structure

``` text
gitops/
├── README.md # This file
└── argo-apps/ # Application manifests managed by ArgoCD
    ├── ingress-app.yaml
    ├── geo-ip-app.yaml
    ├── kube-prom-app.yaml
    └── cnpg-app.yaml
```


Each file in `argo-apps/` represents an **ArgoCD Application** that manages a specific component of the platform.

---

## Applications

### 1. Ingress Controller (`ingress-app.yaml`)

- **Purpose**: Deploys and manages the NGINX Ingress Controller.
- **Role in the cluster**:  
  Provides external access to services running inside Kubernetes. All HTTP/HTTPS traffic flows through this component.
- **Notes**:
  - Typically deployed in `ingress-nginx` namespace.
  - Required before exposing other applications.

---

### 2. Geo-IP Service (`geo-ip-app.yaml`)

- **Purpose**: Deploys the sample **Geo-IP application**.
- **Role in the cluster**:  
  Demonstrates a simple microservice that likely resolves or processes geographic IP information.
- **Notes**:
  - Serves as the custom application for this challenge.
  - Will be accessible through the Ingress controller once deployed.

---

### 3. Monitoring Stack (`kube-prom-app.yaml`)

- **Purpose**: Installs **Kube-Prometheus stack** for observability.
- **Role in the cluster**:  
  Provides monitoring, alerting, and dashboards for cluster and applications.
- **Notes**:
  - Includes Prometheus, Alertmanager, and Grafana.
  - Can be extended with custom dashboards and alerts.

---

### 4. PostgreSQL Database (`cnpg-app.yaml`)

- **Purpose**: Provisions a **CloudNativePG (CNPG)** PostgreSQL cluster.
- **Role in the cluster**:  
  Provides persistent storage and database services for applications.
- **Notes**:
  - Manages HA Postgres with backups and failover.
  - Used by the Geo-IP service (or other apps requiring a database).

---

## Deployment Flow

1. Apply the **root ArgoCD Application** (App of Apps) pointing to this repository.  
2. ArgoCD will recursively deploy each child application:
   - First the **Ingress Controller**  
   - Then supporting services like **CNPG** and **Kube-Prometheus**  
   - Finally the **Geo-IP application**

---

## References

- [ArgoCD App of Apps Pattern](https://argo-cd.readthedocs.io/en/stable/operator-manual/cluster-bootstrapping/)
- [Kube-Prometheus Stack](https://github.com/prometheus-operator/kube-prometheus)
- [CloudNativePG](https://cloudnative-pg.io/)
- [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
