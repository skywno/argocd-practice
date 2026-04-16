# Notes

- Gatekeeper templates must be applied before constraints.
- This repository enforces ordering with sync waves:
  - `opa/policies/templates/kustomization.yaml`: `argocd.argoproj.io/sync-wave: "1"`
  - `opa/policies/constraints/kustomization.yaml`: `argocd.argoproj.io/sync-wave: "5"`
- Constraints also set `SkipDryRunOnMissingResource=true` to avoid dry-run failures while template CRDs are becoming available.

- During teardown, deleting `root-app.yaml` can get stuck if `idp-control-plane` (`AppProject`) is deleted first.
- Current best practice in this repo:
  - manage `AppProject` separately (`argocd/bootstrap/root-project.yaml`)
  - delete app-of-apps first, then delete the project
- If deletion is already stuck, re-apply `root-project.yaml`, then delete `idp-root` with cascade again.
