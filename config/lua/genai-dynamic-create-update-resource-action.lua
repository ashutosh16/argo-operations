data:
  resource.customizations.actions.argoproj.io_Rollout: |
    discovery.lua: |
    actions = {}
    actions["create-genai"] = {}
    return actions
    definitions:
    - name: create-genai
    action.lua: |
    local genaianalysis = {}
    genaianalysis.apiVersion = "argoproj.extensions.io/v1alpha1"
    genaianalysis.kind = "ArgoSupport"
    genaianalysis.metadata = {}
    local os = require("os")
    genaianalysis.metadata.name = "gen-ai"
    genaianalysis.metadata.namespace = obj.metadata.namespace
    genaianalysis.metadata.labels = {}
    genaianalysis.metadata.labels["app"] = obj.metadata.labels["app"]
    genaianalysis.metadata.labels["app.kubernetes.io/instance"] = obj.metadata.labels["app.kubernetes.io/instance"]
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
    impactedResource.operation = "update"
    impactedResource.resource = genaianalysis
    local result = {}
    result[1] = impactedResource
    return result
    discovery.lua: |
    actions = {}
    actions["update-genai"] = {}
    return actions
    definitions:
    - name: update-genai
    action.lua: |
    local genaianalysis = {}
    genaianalysis.apiVersion = "argoproj.extensions.io/v1alpha1"
    genaianalysis.kind = "ArgoSupport"
    genaianalysis.metadata = {}
    local os = require("os")
    genaianalysis.metadata.name = "gen-ai"
    genaianalysis.metadata.namespace = obj.metadata.namespace
    genaianalysis.metadata.labels = {}
    genaianalysis.metadata.labels["app"] = obj.metadata.labels["app"]
    genaianalysis.metadata.labels["app.kubernetes.io/instance"] = obj.metadata.labels["app.kubernetes.io/instance"]
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
    impactedResource.operation = "update"
    impactedResource.resource = genaianalysis
    local result = {}
    result[1] = impactedResource
    return result
