/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ArgoSupportPhase string

// Possible ArgoSupportPhase values
const (
	ArgoSupportPhaseRunning   ArgoSupportPhase = "Running"
	ArgoSupportPhaseCompleted ArgoSupportPhase = "Completed"
	ArgoSupportPhaseFailed    ArgoSupportPhase = "Failed"
)

// ArgoAISupportSpec defines the desired state of ArgoAISupport
type ArgoAISupportSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +kubebuilder:validation:Required
	Workflows []Workflow `json:"workflows,omitempty"`
}

// ArgoAISupportStatus defines the observed state of ArgoAISupport
type ArgoAISupportStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Results            []Result     `json:"results,omitempty"`
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`
	// The generation observed by the  controller from metadata.generation
	// +kubebuilder:validation:Optional
	ObservedGeneration int64            `json:"observedGeneration,omitempty"`
	Phase              ArgoSupportPhase `json:"phase,omitempty"`
}

type Feedback struct {
	DownVote    bool   `json:"downVote,omitempty"`
	FeedbackMsg string `json:"feedbackMsg,omitempty"`
	UpVote      bool   `json:"upVote,omitempty"`
}

type Help struct {
	Links        []string `json:"links,omitempty"`
	SlackChannel string   `json:"slackChannel,omitempty"`
}

type Summary struct {
	MainSummary string `json:"mainSummary,omitempty"`
}

type Result struct {
	Feedback   Feedback         `json:"feedback,omitempty"`
	FinishedAt *metav1.Time     `json:"finishedAt,omitempty"`
	Help       Help             `json:"help,omitempty"`
	Name       string           `json:"name,omitempty"`
	StartedAt  *metav1.Time     `json:"startedAt,omitempty"`
	Summary    Summary          `json:"summary,omitempty"`
	Message    string           `json:"message,omitempty"`
	Phase      ArgoSupportPhase `json:"phase,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// ArgoAISupport is the Schema for the argoaisupports API
type ArgoAISupport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArgoAISupportSpec   `json:"spec,omitempty"`
	Status ArgoAISupportStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ArgoAISupportList contains a list of ArgoAISupport
type ArgoAISupportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ArgoAISupport `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ArgoAISupport{}, &ArgoAISupportList{})
}
