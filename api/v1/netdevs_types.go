/*
Copyright 2022.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NetDevsSpec defines the desired state of NetDevs
type NetDevsSpec struct {
	Username      string `json:"username"`      // Username for the network device.
	Password      string `json:"password"`      // Password for the device.
	Host          string `json:"host"`          // Address for the device.
	Port          int    `json:"port"`          // Port the switch is running on
	Transport     string `json:"transport"`     // http, https , socket , http_local for goeapi
	NetworkOs     string `json:"networkos"`     // eos
	RunningConfig string `json:"runningconfig"` // Running-config

}

// NetDevsStatus defines the observed state of NetDevs
type NetDevsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// NetDevs is the Schema for the netdevs API
type NetDevs struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetDevsSpec   `json:"spec,omitempty"`
	Status NetDevsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NetDevsList contains a list of NetDevs
type NetDevsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetDevs `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NetDevs{}, &NetDevsList{})
}
