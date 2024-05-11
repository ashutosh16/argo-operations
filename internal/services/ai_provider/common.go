package ai_provider

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Failure struct {
	Context string `json:"context"`
}

type Failures struct {
	Failures []Failure `json:"failures"`
}

type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Status            ApplicationStatus `json:"status,omitempty"`
}

// ApplicationStatus contains status information for the application
type ApplicationStatus struct {
	// Resources is a list of Kubernetes resources managed by this application
	Resources []ResourceStatus `json:"resources,omitempty"`
	// Sync contains information about the application's current sync status
	Sync SyncStatus `json:"sync,omitempty"`
	// Health contains information about the application's current health status
	Health     HealthStatus           `json:"health,omitempty"`
	Conditions []ApplicationCondition `json:"conditions,omitempty"`
	// ReconciledAt indicates when the application state was reconciled using the latest git version
	ReconciledAt *metav1.Time `json:"reconciledAt,omitempty"`
}

type ResourceStatus struct {
	Group           string         `json:"group,omitempty" `
	Version         string         `json:"version,omitempty" `
	Kind            string         `json:"kind,omitempty" `
	Namespace       string         `json:"namespace,omitempty" `
	Name            string         `json:"name,omitempty" `
	Status          SyncStatusCode `json:"status,omitempty" `
	Health          *HealthStatus  `json:"health,omitempty" `
	Hook            bool           `json:"hook,omitempty" `
	RequiresPruning bool           `json:"requiresPruning,omitempty" `
	SyncWave        int64          `json:"syncWave,omitempty" `
}

// SyncStatus contains information about the currently observed live and desired states of an application
type SyncStatus struct {
	Status SyncStatusCode `json:"status" `
}

// SyncStatusCode is a type which represents possible comparison results
type SyncStatusCode string

// Possible comparison results
const (
	// SyncStatusCodeUnknown indicates that the status of a sync could not be reliably determined
	SyncStatusCodeUnknown SyncStatusCode = "Unknown"
	// SyncStatusCodeOutOfSync indicates that desired and live states match
	SyncStatusCodeSynced SyncStatusCode = "Synced"
	// SyncStatusCodeOutOfSync indicates that there is a drift between desired and live states
	SyncStatusCodeOutOfSync SyncStatusCode = "OutOfSync"
)

// HealthStatus contains information about the currently observed health state of an application or resource
type HealthStatus struct {
	// Status holds the status code of the application or resource
	Status HealthStatusCode `json:"status,omitempty" `
	// Message is a human-readable informational message describing the health status
	Message string `json:"message,omitempty"`
}
type ApplicationConditionType = string

// ApplicationCondition contains details about an application condition, which is usually an error or warning
type ApplicationCondition struct {
	// Type is an application condition type
	Type ApplicationConditionType `json:"type" `
	// Message contains human-readable message indicating details about condition
	Message string `json:"message" `
	// LastTransitionTime is the time the condition was last observed
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`
}

const (
	// ApplicationConditionDeletionError indicates that controller failed to delete application
	ApplicationConditionDeletionError = "DeletionError"
	// ApplicationConditionInvalidSpecError indicates that application source is invalid
	ApplicationConditionInvalidSpecError = "InvalidSpecError"
	// ApplicationConditionComparisonError indicates controller failed to compare application state
	ApplicationConditionComparisonError = "ComparisonError"
	// ApplicationConditionSyncError indicates controller failed to automatically sync the application
	ApplicationConditionSyncError = "SyncError"
	// ApplicationConditionUnknownError indicates an unknown controller error
	ApplicationConditionUnknownError = "UnknownError"
	// ApplicationConditionSharedResourceWarning indicates that controller detected resources which belongs to more than one application
	ApplicationConditionSharedResourceWarning = "SharedResourceWarning"
	// ApplicationConditionRepeatedResourceWarning indicates that application source has resource with same Group, Kind, Name, Namespace multiple times
	ApplicationConditionRepeatedResourceWarning = "RepeatedResourceWarning"
	// ApplicationConditionExcludedResourceWarning indicates that application has resource which is configured to be excluded
	ApplicationConditionExcludedResourceWarning = "ExcludedResourceWarning"
	// ApplicationConditionOrphanedResourceWarning indicates that application has orphaned resources
	ApplicationConditionOrphanedResourceWarning = "OrphanedResourceWarning"
)

type HealthStatusCode string

const (
	// Indicates that health assessment failed and actual health status is unknown
	HealthStatusUnknown HealthStatusCode = "Unknown"
	// Progressing health status means that resource is not healthy but still have a chance to reach healthy state
	HealthStatusProgressing HealthStatusCode = "Progressing"
	// Resource is 100% healthy
	HealthStatusHealthy HealthStatusCode = "Healthy"
	// Assigned to resources that are suspended or paused. The typical example is a
	// [suspended](https://kubernetes.io/docs/tasks/job/automated-tasks-with-cron-jobs/#suspend) CronJob.
	HealthStatusSuspended HealthStatusCode = "Suspended"
	// Degrade status is used if resource status indicates failure or resource could not reach healthy state
	// within some timeout.
	HealthStatusDegraded HealthStatusCode = "Degraded"
	// Indicates that resource is missing in the cluster.
	HealthStatusMissing HealthStatusCode = "Missing"
)
