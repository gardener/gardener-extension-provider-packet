// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudProfile represents certain properties about a provider environment.
type CloudProfile struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Spec defines the provider environment properties.
	// +optional
	Spec CloudProfileSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudProfileList is a collection of CloudProfiles.
type CloudProfileList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list object metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is the list of CloudProfiles.
	Items []CloudProfile `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// CloudProfileSpec is the specification of a CloudProfile.
// It must contain exactly one of its defined keys.
type CloudProfileSpec struct {
	// CABundle is a certificate bundle which will be installed onto every host machine of shoot cluster targeting this profile.
	// +optional
	CABundle *string `json:"caBundle,omitempty" protobuf:"bytes,1,opt,name=caBundle"`
	// Kubernetes contains constraints regarding allowed values of the 'kubernetes' block in the Shoot specification.
	Kubernetes KubernetesSettings `json:"kubernetes" protobuf:"bytes,2,opt,name=kubernetes"`
	// MachineImages contains constraints regarding allowed values for machine images in the Shoot specification.
	// +patchMergeKey=name
	// +patchStrategy=merge
	MachineImages []MachineImage `json:"machineImages" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,3,rep,name=machineImages"`
	// MachineTypes contains constraints regarding allowed values for machine types in the 'workers' block in the Shoot specification.
	// +patchMergeKey=name
	// +patchStrategy=merge
	MachineTypes []MachineType `json:"machineTypes" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,4,rep,name=machineTypes"`
	// ProviderConfig contains provider-specific configuration for the profile.
	// +optional
	ProviderConfig *ProviderConfig `json:"providerConfig,omitempty" protobuf:"bytes,5,opt,name=providerConfig"`
	// Regions contains constraints regarding allowed values for regions and zones.
	// +patchMergeKey=name
	// +patchStrategy=merge
	Regions []Region `json:"regions" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,6,rep,name=regions"`
	// SeedSelector contains an optional list of labels on `Seed` resources that marks those seeds whose shoots may use this provider profile.
	// An empty list means that all seeds of the same provider type are supported.
	// This is useful for environments that are of the same type (like openstack) but may have different "instances"/landscapes.
	// +optional
	SeedSelector *metav1.LabelSelector `json:"seedSelector,omitempty" protobuf:"bytes,7,opt,name=seedSelector"`
	// Type is the name of the provider.
	Type string `json:"type" protobuf:"bytes,8,opt,name=type"`
	// VolumeTypes contains constraints regarding allowed values for volume types in the 'workers' block in the Shoot specification.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +optional
	VolumeTypes []VolumeType `json:"volumeTypes,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,9,rep,name=volumeTypes"`
}

// KubernetesSettings contains constraints regarding allowed values of the 'kubernetes' block in the Shoot specification.
type KubernetesSettings struct {
	// Versions is the list of allowed Kubernetes versions with optional expiration dates for Shoot clusters.
	// +patchMergeKey=version
	// +patchStrategy=merge
	// +optional
	Versions []ExpirableVersion `json:"versions,omitempty" patchStrategy:"merge" patchMergeKey:"version" protobuf:"bytes,1,rep,name=versions"`
}

// MachineImage defines the name and multiple versions of the machine image in any environment.
type MachineImage struct {
	// Name is the name of the image.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Versions contains versions and expiration dates of the machine image
	// +patchMergeKey=version
	// +patchStrategy=merge
	Versions []ExpirableVersion `json:"versions" patchStrategy:"merge" patchMergeKey:"version" protobuf:"bytes,2,rep,name=versions"`
}

// ExpirableVersion contains a version and an expiration date.
type ExpirableVersion struct {
	// Version is the version identifier.
	Version string `json:"version" protobuf:"bytes,1,opt,name=version"`
	// ExpirationDate defines the time at which this version expires.
	// +optional
	ExpirationDate *metav1.Time `json:"expirationDate,omitempty" protobuf:"bytes,2,opt,name=expirationDate"`
	// Classification defines the state of a version (preview, supported, deprecated)
	// +optional
	Classification *VersionClassification `json:"classification,omitempty" protobuf:"bytes,3,opt,name=classification,casttype=VersionClassification"`
}

// MachineType contains certain properties of a machine type.
type MachineType struct {
	// CPU is the number of CPUs for this machine type.
	CPU resource.Quantity `json:"cpu" protobuf:"bytes,1,opt,name=cpu"`
	// GPU is the number of GPUs for this machine type.
	GPU resource.Quantity `json:"gpu" protobuf:"bytes,2,opt,name=gpu"`
	// Memory is the amount of memory for this machine type.
	Memory resource.Quantity `json:"memory" protobuf:"bytes,3,opt,name=memory"`
	// Name is the name of the machine type.
	Name string `json:"name" protobuf:"bytes,4,opt,name=name"`
	// Storage is the amount of storage associated with the root volume of this machine type.
	// +optional
	Storage *MachineTypeStorage `json:"storage,omitempty" protobuf:"bytes,5,opt,name=storage"`
	// Usable defines if the machine type can be used for shoot clusters.
	// +optional
	Usable *bool `json:"usable,omitempty" protobuf:"varint,6,opt,name=usable"`
}

// MachineTypeStorage is the amount of storage associated with the root volume of this machine type.
type MachineTypeStorage struct {
	// Class is the class of the storage type.
	Class string `json:"class" protobuf:"bytes,1,opt,name=class"`
	// StorageSize is the storage size.
	StorageSize resource.Quantity `json:"size" protobuf:"bytes,2,opt,name=size"`
	// Type is the type of the storage.
	Type string `json:"type" protobuf:"bytes,3,opt,name=type"`
}

// Region contains certain properties of a region.
type Region struct {
	// Name is a region name.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Zones is a list of availability zones in this region.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +optional
	Zones []AvailabilityZone `json:"zones,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=zones"`
	// Labels is an optional set of key-value pairs that contain certain administrator-controlled labels for this region.
	// It can be used by Gardener administrators/operators to provide additional information about a region, e.g. wrt
	// quality, reliability, access restrictions, etc.
	// +optional
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,3,rep,name=labels"`
}

// AvailabilityZone is an availability zone.
type AvailabilityZone struct {
	// Name is an an availability zone name.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// UnavailableMachineTypes is a list of machine type names that are not availability in this zone.
	// +optional
	UnavailableMachineTypes []string `json:"unavailableMachineTypes,omitempty" protobuf:"bytes,2,rep,name=unavailableMachineTypes"`
	// UnavailableVolumeTypes is a list of volume type names that are not availability in this zone.
	// +optional
	UnavailableVolumeTypes []string `json:"unavailableVolumeTypes,omitempty" protobuf:"bytes,3,rep,name=unavailableVolumeTypes"`
}

// VolumeType contains certain properties of a volume type.
type VolumeType struct {
	// Class is the class of the volume type.
	Class string `json:"class" protobuf:"bytes,1,opt,name=class"`
	// Name is the name of the volume type.
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`
	// Usable defines if the volume type can be used for shoot clusters.
	// +optional
	Usable *bool `json:"usable,omitempty" protobuf:"varint,3,opt,name=usable"`
}

const (
	// VolumeClassStandard is a constant for the standard volume class.
	VolumeClassStandard string = "standard"
	// VolumeClassPremium is a constant for the premium volume class.
	VolumeClassPremium string = "premium"
)

// VersionClassification is the logical state of a version according to https://github.com/gardener/gardener/blob/master/docs/operations/versioning.md
type VersionClassification string

const (
	// ClassificationPreview indicates that a version has recently been added and not promoted to "Supported" yet.
	// ClassificationPreview versions will not be considered for automatic Kubernetes and Machine Image patch version updates.
	ClassificationPreview VersionClassification = "preview"
	// ClassificationSupported indicates that a patch version is the recommended version for a shoot.
	// Using VersionMaintenance (see: https://github.com/gardener/gardener/docs/operation/versioning.md) there is one supported version per maintained minor version.
	// Supported versions are eligible for the automated Kubernetes and Machine image patch version update for shoot clusters in Gardener.
	ClassificationSupported VersionClassification = "supported"
	// ClassificationDeprecated indicates that a patch version should not be used anymore, should be updated to a new version
	// and will eventually expire.
	ClassificationDeprecated VersionClassification = "deprecated"
)
