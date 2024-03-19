// Copyright 2024 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespacedCloudProfile represents certain properties about a provider environment.
type NamespacedCloudProfile struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Spec defines the provider environment properties.
	Spec NamespacedCloudProfileSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// Most recently observed status of the NamespacedCloudProfile.
	Status NamespacedCloudProfileStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespacedCloudProfileList is a collection of CloudProfiles.
type NamespacedCloudProfileList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list object metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is the list of CloudProfiles.
	Items []NamespacedCloudProfile `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// NamespacedCloudProfileSpec is the specification of a NamespacedCloudProfile.
// It must contain exactly one of its defined keys.
type NamespacedCloudProfileSpec struct {
	// CABundle is a certificate bundle which will be installed onto every host machine of shoot cluster targeting this profile.
	// +optional
	CABundle *string `json:"caBundle,omitempty" protobuf:"bytes,1,opt,name=caBundle"`
	// Kubernetes contains constraints regarding allowed values of the 'kubernetes' block in the Shoot specification.
	// +optional
	Kubernetes *KubernetesSettings `json:"kubernetes" protobuf:"bytes,2,opt,name=kubernetes"`
	// MachineImages contains constraints regarding allowed values for machine images in the Shoot specification.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +optional
	MachineImages []MachineImage `json:"machineImages" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,3,opt,name=machineImages"`
	// MachineTypes contains constraints regarding allowed values for machine types in the 'workers' block in the Shoot specification.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +optional
	MachineTypes []MachineType `json:"machineTypes" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,4,opt,name=machineTypes"`
	// Regions contains constraints regarding allowed values for regions and zones.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +optional
	Regions []Region `json:"regions" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,6,opt,name=regions"`
	// SeedSelector contains an optional list of labels on `Seed` resources that marks those seeds whose shoots may use this provider profile.
	// An empty list means that all seeds of the same provider type are supported.
	// This is useful for environments that are of the same type (like openstack) but may have different "instances"/landscapes.
	// Optionally a list of possible providers can be added to enable cross-provider scheduling. By default, the provider
	// type of the seed must match the shoot's provider.
	// +optional
	SeedSelector *SeedSelector `json:"seedSelector,omitempty" protobuf:"bytes,7,opt,name=seedSelector"`
	// VolumeTypes contains constraints regarding allowed values for volume types in the 'workers' block in the Shoot specification.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +optional
	VolumeTypes []VolumeType `json:"volumeTypes,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,9,opt,name=volumeTypes"`
	// A pointer to the NamespacedCloudProfile's parent CloudProfile.
	Parent CloudProfileReference `json:"parent" protobuf:"bytes,10,req,name=parent"`
}

// NamespacedCloudProfileStatus holds the most recently observed status of the project.
type NamespacedCloudProfileStatus struct {
	// CloudProfile is the most recently generated CloudProfile of the NamespacedCloudProfile
	CloudProfileSpec CloudProfileSpec `json:"cloudProfileSpec,omitempty" protobuf:"bytes,1,opt,name=cloudProfileSpec"`
}

// CloudProfileReference holds the information about the parent of the NamespacedCloudProfile.
type CloudProfileReference struct {
	Kind string `json:"kind" protobuf:"bytes,1,req,name=kind"`
	Name string `json:"name" protobuf:"bytes,2,req,name=name"`
}