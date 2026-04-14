# Backend-First Internal Developer Portal (MicroK8s)

This sample implements a backend-first Internal Developer Portal control plane with:

- Argo CD as GitOps reconciler
- Crossplane as platform API engine
- OPA Gatekeeper as policy enforcement
- Git as the single source of truth

No frontend portal (such as Backstage) is required for this sample.

## Repository layout

- `argocd/bootstrap`: Argo CD installation and root app bootstrap
- `argocd/apps`: app-of-apps children with sync-wave ordering
- `crossplane/config`: providers, `ProviderConfig`, XRD, and Composition
- `opa/policies`: Gatekeeper constraint templates and constraints
- `platform`: sample tenant, claim, and policy test workloads

## Prerequisites

- MicroK8s installed
- `kubectl` and `argocd` CLIs installed
- Access to this Git repository from the cluster

Enable required MicroK8s add-ons:

```bash
microk8s status --wait-ready
microk8s enable dns storage ingress metrics-server
microk8s config > ~/.kube/config
```

## Bootstrap sequence

1. Update repository URL placeholders in:
   - `argocd/bootstrap/root-app.yaml`
   - `argocd/apps/project.yaml`
   - `argocd/apps/crossplane-config.yaml`
   - `argocd/apps/opa-policies.yaml`
   - `argocd/apps/platform.yaml`
2. Install Argo CD:

```bash
kubectl apply -k crossplane-argocd-opa/argocd/bootstrap
kubectl -n argocd rollout status deploy/argocd-server --timeout=5m
```

3. Bootstrap app-of-apps:

```bash
kubectl apply -f crossplane-argocd-opa/argocd/bootstrap/root-app.yaml
```

## Validation checks

Argo CD sync and health:

```bash
argocd app get idp-root
argocd app get crossplane-core
argocd app get crossplane-config
argocd app get opa-gatekeeper-core
argocd app get opa-gatekeeper-policies
argocd app get idp-platform-sample
```

Crossplane API readiness:

```bash
kubectl get xrd xservices.platform.example.org
kubectl get composition xservice.microk8s.platform.example.org
kubectl get providers.pkg.crossplane.io
kubectl get services.platform.example.org -A
kubectl get xservices.platform.example.org
```

Gatekeeper and policy checks:

```bash
kubectl get constrainttemplates
kubectl get k8srequiredlabels,k8sdisallowlatest,k8srequiredresources,k8sdisallowprivileged
kubectl apply -f crossplane-argocd-opa/platform/policy-test-valid-deployment.yaml
kubectl apply -f crossplane-argocd-opa/platform/policy-test-invalid-latest.yaml
```

Expected behavior:

- `policy-test-valid-deployment.yaml` is admitted.
- `policy-test-invalid-latest.yaml` is denied by `K8sDisallowLatest`.

Disaster recovery reconciliation check:

```bash
kubectl -n tenant-a delete deploy compliant-sample
argocd app sync idp-platform-sample
kubectl -n tenant-a get deploy compliant-sample
```

## Production defaults included

- Argo CD auto-sync (`prune` and `selfHeal`)
- Argo CD retry backoff
- Application finalizers for resource cleanup
- Sync-wave ordering for CRDs/controllers before dependent resources
- Namespace and workload policy guardrails via OPA Gatekeeper

## Placeholder secrets and cloud provider migration

Current sample uses local in-cluster Crossplane provider auth (`InjectedIdentity`) and does not require cloud credentials.

To move toward cloud production:

1. Add cloud provider package and `ProviderConfig` in `crossplane/config/providers`.
2. Replace local Composition resources with managed cloud resources.
3. Introduce secret management (SOPS + KSOPS or External Secrets) and remove placeholders.
4. Add environment overlays for staging/production and stricter AppProject source/destination controls.
