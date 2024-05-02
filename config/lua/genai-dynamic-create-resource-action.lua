data:
 resource.customizations.actions.argoproj.io_Rollout: |
  discovery.lua: |
    actions = {}
    actions["genai-analysis"] = {}
    return actions
  definitions:
  - name: genai-analysis
    action.lua: |
      local genaianalysis = {}
      genaianalysis.apiVersion = "argoproj.extensions.io/v1alpha1"
      genaianalysis.kind = "ArgoOperationRun"
      genaianalysis.metadata = {}
      local os = require("os")
      genaianalysis.metadata.name = "genai-analysis-run"
      genaianalysis.metadata.namespace = obj.metadata.namespace
      genaianalysis.metadata.labels = {}
      genaianalysis.metadata.labels["app"] = obj.metadata.labels["app"]
      genaianalysis.metadata.labels["rollout-type"] = "Background"
      genaianalysis.metadata.labels["rollouts-pod-template-hash"] = obj.status.currentPodHash
      genaianalysis.metadata.annotations = {}
      genaianalysis.metadata.annotations["rollout.argoproj.io/revision"] = obj.metadata.annotations["rollout.argoproj.io/revision"]
      local ownerRef = {}
      ownerRef.apiVersion = obj.apiVersion
      ownerRef.kind = obj.kind
      ownerRef.name = obj.metadata.name
      ownerRef.uid = obj.metadata.uid
      ownerRef.blockOwnerDeletion = true
      ownerRef.controller = true
      genaianalysis.metadata.ownerReferences = {}
      genaianalysis.metadata.ownerReferences[1] = ownerRef
      impactedResource = {}
      impactedResource.operation = "create"
      impactedResource.resource = genaianalysis
      local result = {}
      result[1] = impactedResource
      return result
