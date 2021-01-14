/*
Copyright 2021 Intuit Inc.

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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// RollingUpgradeSpec defines the desired state of RollingUpgrade
type RollingUpgradeSpec struct {
	PostDrainDelaySeconds int                 `json:"postDrainDelaySeconds,omitempty"`
	NodeIntervalSeconds   int                 `json:"nodeIntervalSeconds,omitempty"`
	AsgName               string              `json:"asgName,omitempty"`
	PreDrain              PreDrainSpec        `json:"preDrain,omitempty"`
	PostDrain             PostDrainSpec       `json:"postDrain,omitempty"`
	PostTerminate         PostTerminateSpec   `json:"postTerminate,omitempty"`
	Strategy              UpdateStrategy      `json:"strategy,omitempty"`
	IgnoreDrainFailures   bool                `json:"ignoreDrainFailures,omitempty"`
	ForceRefresh          bool                `json:"forceRefresh,omitempty"`
	ReadinessGates        []NodeReadinessGate `json:"readinessGates,omitempty"`
}

// RollingUpgradeStatus defines the observed state of RollingUpgrade
type RollingUpgradeStatus struct {
	CurrentStatus           string                    `json:"currentStatus,omitempty"`
	StartTime               string                    `json:"startTime,omitempty"`
	EndTime                 string                    `json:"endTime,omitempty"`
	TotalProcessingTime     string                    `json:"totalProcessingTime,omitempty"`
	NodesProcessed          int                       `json:"nodesProcessed,omitempty"`
	TotalNodes              int                       `json:"totalNodes,omitempty"`
	Conditions              []RollingUpgradeCondition `json:"conditions,omitempty"`
	LastNodeTerminationTime metav1.Time               `json:"lastTerminationTime,omitempty"`
	LastNodeDrainTime       metav1.Time               `json:"lastDrainTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=rollingupgrades,scope=Namespaced,shortName=ru
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.currentStatus",description="current status of the rollingupgarde"
// +kubebuilder:printcolumn:name="TotalNodes",type="string",JSONPath=".status.totalNodes",description="total nodes involved in the rollingupgarde"
// +kubebuilder:printcolumn:name="NodesProcessed",type="string",JSONPath=".status.nodesProcessed",description="current number of nodes processed in the rollingupgarde"

// RollingUpgrade is the Schema for the rollingupgrades API
type RollingUpgrade struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RollingUpgradeSpec   `json:"spec,omitempty"`
	Status RollingUpgradeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RollingUpgradeList contains a list of RollingUpgrade
type RollingUpgradeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RollingUpgrade `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RollingUpgrade{}, &RollingUpgradeList{})
}

// PreDrainSpec contains the fields for actions taken before draining the node.
type PreDrainSpec struct {
	Script string `json:"script,omitempty"`
}

// PostDrainSpec contains the fields for actions taken after draining the node.
type PostDrainSpec struct {
	Script         string `json:"script,omitempty"`
	WaitSeconds    int64  `json:"waitSeconds,omitempty"`
	PostWaitScript string `json:"postWaitScript,omitempty"`
}

// PostTerminateSpec contains the fields for actions taken after terminating the node.
type PostTerminateSpec struct {
	Script string `json:"script,omitempty"`
}

type NodeReadinessGate struct {
	MatchLabels map[string]string `json:"matchLabels,omitempty" protobuf:"bytes,1,rep,name=matchLabels"`
}

const (
	// Status
	StatusInit     = "init"
	StatusRunning  = "running"
	StatusComplete = "completed"
	StatusError    = "error"

	// Conditions
	UpgradeComplete UpgradeConditionType = "Complete"
)

var (
	FiniteStates = []string{StatusComplete, StatusError}
)

// RollingUpgradeCondition describes the state of the RollingUpgrade
type RollingUpgradeCondition struct {
	Type   UpgradeConditionType   `json:"type,omitempty"`
	Status corev1.ConditionStatus `json:"status,omitempty"`
}

type UpdateStrategyType string
type UpdateStrategyMode string
type UpgradeConditionType string

const (
	RandomUpdateStrategy          UpdateStrategyType = "randomUpdate"
	UniformAcrossAzUpdateStrategy UpdateStrategyType = "uniformAcrossAzUpdate"

	UpdateStrategyModeLazy  UpdateStrategyMode = "lazy"
	UpdateStrategyModeEager UpdateStrategyMode = "eager"
)

// UpdateStrategy holds the information needed to perform update based on different update strategies
type UpdateStrategy struct {
	Type           UpdateStrategyType `json:"type,omitempty"`
	Mode           UpdateStrategyMode `json:"mode,omitempty"`
	MaxUnavailable intstr.IntOrString `json:"maxUnavailable,omitempty"`
	DrainTimeout   int                `json:"drainTimeout"`
}

func (c UpdateStrategyMode) String() string {
	return string(c)
}

// NamespacedName returns namespaced name of the object.
func (r *RollingUpgrade) NamespacedName() string {
	return fmt.Sprintf("%s/%s", r.Namespace, r.Name)
}

func (r *RollingUpgrade) ScalingGroupName() string {
	return r.Spec.AsgName
}

func (r *RollingUpgrade) CurrentStatus() string {
	return r.Status.CurrentStatus
}

func (r *RollingUpgrade) UpdateStrategyType() UpdateStrategyType {
	return r.Spec.Strategy.Type
}

func (r *RollingUpgrade) MaxUnavailable() intstr.IntOrString {
	return r.Spec.Strategy.MaxUnavailable
}

func (r *RollingUpgrade) LastNodeTerminationTime() metav1.Time {
	return r.Status.LastNodeTerminationTime
}

func (r *RollingUpgrade) LastNodeDrainTime() metav1.Time {
	return r.Status.LastNodeDrainTime
}

func (r *RollingUpgrade) NodeIntervalSeconds() int {
	return r.Spec.NodeIntervalSeconds
}

func (r *RollingUpgrade) PostDrainDelaySeconds() int {
	return r.Spec.PostDrainDelaySeconds
}

func (r *RollingUpgrade) SetCurrentStatus(status string) {
	r.Status.CurrentStatus = status
}

func (r *RollingUpgrade) SetStartTime(t string) {
	r.Status.StartTime = t
}

func (r *RollingUpgrade) StartTime() string {
	return r.Status.StartTime
}

func (r *RollingUpgrade) SetEndTime(t string) {
	r.Status.EndTime = t
}

func (r *RollingUpgrade) EndTime() string {
	return r.Status.EndTime
}

func (r *RollingUpgrade) SetTotalNodes(n int) {
	r.Status.TotalNodes = n
}

func (r *RollingUpgrade) SetNodesProcessed(n int) {
	r.Status.NodesProcessed = n
}

func (r *RollingUpgrade) IsForceRefresh() bool {
	return r.Spec.ForceRefresh
}

// Migrate r.setDefaultsForRollingUpdateStrategy & r.validateRollingUpgradeObj into v1alpha1 RollingUpgrade.Validate()
func (r *RollingUpgrade) Validate() (bool, error) {
	return true, nil
}
