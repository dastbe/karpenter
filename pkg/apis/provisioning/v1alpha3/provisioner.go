/*
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

package v1alpha3

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ProvisionerSpec is the top level provisioner specification. Provisioners
// launch nodes in response to pods where status.conditions[type=unschedulable,
// status=true]. Node configuration is driven by through a combination of
// provisioner specification (defaults) and pod scheduling constraints
// (overrides). A single provisioner is capable of managing highly diverse
// capacity within a single cluster and in most cases, only one should be
// necessary. For advanced use cases like workload separation and sharding, it's
// possible to define multiple provisioners. These provisioners may have
// different defaults and can be specifically targeted by pods using
// pod.spec.nodeSelector["karpenter.sh/provisioner-name"]=$PROVISIONER_NAME.
type ProvisionerSpec struct {
	// Cluster that launched nodes connect to.
	Cluster Cluster `json:"cluster"`
	// Constraints are applied to all nodes launched by this provisioner.
	// +optional
	Constraints `json:",inline"`
	// TTLSecondsAfterEmpty is the number of seconds the controller will wait
	// before attempting to terminate a node, measured from when the node is
	// detected to be empty. A Node is considered to be empty when it does not
	// have pods scheduled to it, excluding daemonsets.
	//
	// Termination due to underutilization is disabled if this field is not set.
	// +optional
	TTLSecondsAfterEmpty *int64 `json:"ttlSecondsAfterEmpty,omitempty"`
	// TTLSecondsUntilExpired is the number of seconds the controller will wait
	// before terminating a node, measured from when the node is created. This
	// is useful to implement features like eventually consistent node upgrade,
	// memory leak protection, and disruption testing.
	//
	// Termination due to expiration is disabled if this field is not set.
	// +optional
	TTLSecondsUntilExpired *int64 `json:"ttlSecondsUntilExpired,omitempty"`
}

// Cluster configures the cluster that the provisioner operates against. If
// not specified, it will default to using the controller's kube-config.
type Cluster struct {
	// Endpoint is required for nodes to connect to the API Server.
	// +required
	Endpoint string `json:"endpoint"`
	// CABundle used by nodes to verify API Server certificates. If omitted (nil),
	// it will be dynamically loaded at runtime from the in-cluster configuration
	// file /var/run/secrets/kubernetes.io/serviceaccount/ca.crt.
	// An empty value ("") can be used to signal that no CABundle should be used.
	// +optional
	CABundle *string `json:"caBundle,omitempty"`
	// Name may be required to detect implementing cloud provider resources.
	// +optional
	Name *string `json:"name,omitempty"`
}

// Constraints are applied to all nodes created by the provisioner. They can be
// overriden by NodeSelectors at the pod level.
type Constraints struct {
	// Taints will be applied to every node launched by the Provisioner. If
	// specified, the provisioner will not provision nodes for pods that do not
	// have matching tolerations.
	// +optional
	Taints []v1.Taint `json:"taints,omitempty"`
	// Labels will be applied to every node launched by the Provisioner unless
	// overriden by pod node selectors. Well known labels control provisioning
	// behavior. Additional labels may be supported by your cloudprovider.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Zones constrains where nodes will be launched by the Provisioner. If
	// unspecified, defaults to all zones in the region. Cannot be specified if
	// label "topology.kubernetes.io/zone" is specified.
	// +optional
	Zones []string `json:"zones,omitempty"`
	// InstanceTypes constrains which instances types will be used for nodes
	// launched by the Provisioner. If unspecified, it will support all types.
	// Cannot be specified if label "node.kubernetes.io/instance-type" is specified.
	// +optional
	InstanceTypes []string `json:"instanceTypes,omitempty"`
	// Architecture constrains the underlying node architecture
	// +optional
	Architecture *string `json:"architecture,omitempty"`
	// OperatingSystem constrains the underlying node operating system
	// +optional
	OperatingSystem *string `json:"operatingSystem,omitempty"`
}

var (
	ArchitectureAmd64 = "amd64"
	ArchitectureArm64 = "arm64"
)

var (
	OperatingSystemLinux = "linux"
)

var (
	ProvisionerNameLabelKey         = SchemeGroupVersion.Group + "/provisioner-name"
	NotReadyTaintKey                = SchemeGroupVersion.Group + "/not-ready"
	DoNotEvictPodAnnotationKey      = SchemeGroupVersion.Group + "/do-not-evict"
	EmptinessTimestampAnnotationKey = SchemeGroupVersion.Group + "/emptiness-timestamp"
	TerminationFinalizer            = SchemeGroupVersion.Group + "/termination"
	DefaultProvisioner              = types.NamespacedName{Name: "default"}
)

// RestrictedLabels prevent usage of specific labels.
var RestrictedLabels = []string{
	// Use strongly typed fields instead
	v1.LabelArchStable,
	v1.LabelOSStable,
	v1.LabelTopologyZone,
	v1.LabelInstanceTypeStable,
	// Used internally by provisioning logic
	ProvisionerNameLabelKey,
	EmptinessTimestampAnnotationKey,
	v1.LabelHostname,
}

// Provisioner is the Schema for the Provisioners API
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=provisioners,scope=Cluster
// +kubebuilder:subresource:status
type Provisioner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProvisionerSpec   `json:"spec,omitempty"`
	Status ProvisionerStatus `json:"status,omitempty"`
}

// ProvisionerList contains a list of Provisioner
// +kubebuilder:object:root=true
type ProvisionerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Provisioner `json:"items"`
}
