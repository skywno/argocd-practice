# Argo CD Overlay Prerequisite

Before applying this overlay, create the `argocd-sops-gpg` secret first.

This overlay expects the secret to already exist in the Argo CD namespace.

Create it from your private key file (`ksops/private.asc`) with:

```bash
kubectl -n argocd create secret generic argocd-sops-gpg \
  --from-file=private.asc=ksops/private.asc
```

How this is used by the overlay:
- `argocd-gpg-key-patch.yaml` mounts secret `argocd-sops-gpg` at `/secret`.
- An init container runs `gpg --import /secret/private.asc`.
- The imported keyring is mounted into repo-server and used by KSOPS/SOPS during manifest decryption.
