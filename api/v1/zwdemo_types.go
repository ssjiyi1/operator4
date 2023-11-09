/*
Copyright 2023.

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

// ZwdemoSpec defines the desired state of Zwdemo
type ZwdemoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Zwdemo. Edit zwdemo_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// ZwdemoStatus defines the observed state of Zwdemo
type ZwdemoStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Zwdemo is the Schema for the zwdemoes API
type Zwdemo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZwdemoSpec   `json:"spec,omitempty"`
	Status ZwdemoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ZwdemoList contains a list of Zwdemo
type ZwdemoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Zwdemo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Zwdemo{}, &ZwdemoList{})
}
