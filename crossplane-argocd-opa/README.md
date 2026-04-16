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
- `argocd/bootstrap/root-project.yaml`: dedicated `AppProject` (`idp-control-plane`) managed separately from app-of-apps
- `crossplane/config`: providers, `ProviderConfig`, XRD, and Composition
- `opa/core`: reserved for Gatekeeper core manifests (current deployment uses Helm via `argocd/apps/opa-core.yaml`)
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
   - `argocd/apps/opa-core.yaml`
   - `argocd/apps/crossplane-config.yaml`
   - `argocd/apps/opa-policies.yaml`
   - `argocd/apps/platform.yaml`
2. Install Argo CD and bootstrap root-app:

```bash
kubectl create ns argocd
kubectl apply -n argocd --server-side --force-conflicts -k argocd/bootstrap

kubectl apply -f argocd/bootstrap/root-project.yaml
kubectl apply -f argocd/bootstrap/root-app.yaml
```

## Safe teardown sequence

Delete workloads first, then remove the `AppProject`:

```bash
argocd app delete idp-root --cascade --yes
kubectl wait application/idp-root -n argocd --for=delete --timeout=300s

kubectl delete -f argocd/bootstrap/root-project.yaml
```

Why this ordering matters:

- Child apps (`crossplane-core`, `opa-gatekeeper-core`, etc.) use `spec.project: idp-control-plane`.
- If the `AppProject` is deleted first, Argo CD cannot finish child app finalization and deletion may get stuck.

Recovery when deletion is already stuck:

```bash
kubectl apply -f argocd/bootstrap/root-project.yaml
argocd app delete idp-root --cascade --yes
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
kubectl get xrd myapps.example.crossplane.io
kubectl get composition app-yaml
kubectl get providers.pkg.crossplane.io
kubectl get myapps.example.crossplane.io -A
```

Gatekeeper and policy checks:

```bash
kubectl get constrainttemplates
kubectl get k8srequiredlabels,k8sdisallowlatest,k8srequiredresources,k8sdisallowprivileged
kubectl apply -f crossplane-argocd-opa/platform/policy-test-valid-deployment.yaml
kubectl apply -f crossplane-argocd-opa/platform/policy-test-invalid-latest.yaml
```

If custom constraint kinds are missing, wait for templates first:

```bash
kubectl wait --for=condition=Created constrainttemplates.templates.gatekeeper.sh/k8srequiredlabels --timeout=120s
kubectl wait --for=condition=Created constrainttemplates.templates.gatekeeper.sh/k8sdisallowlatest --timeout=120s
kubectl wait --for=condition=Created constrainttemplates.templates.gatekeeper.sh/k8srequiredresources --timeout=120s
kubectl wait --for=condition=Created constrainttemplates.templates.gatekeeper.sh/k8sdisallowprivileged --timeout=120s
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
