i learned that OPA templates need to be deployed before the OPA constraints, in order to do that, `argocd.argoproj.io/sync-wave: "1"` annotation has been used. 


when deleting `root-app.yaml`, deletion of the entire application gets stuck due to error message that unable to delete application resources:error getting app project: "idp-ontrol-plane": appproject.argoproj.io "idp-control-plane" not found.

To fix this issue, Best practices is that AppProject needs to be separately manged from app-of-apps. (bottsrap it once, don't prune it with root app). If they are kept together, then give AppProject a lower sync wave so it's created first and pruned last, reducing this race/order issue. 
[../.images/error.png]
