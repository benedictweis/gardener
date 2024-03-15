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

package namespacedcloudprofile

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"

	"github.com/gardener/gardener/pkg/api"
	"github.com/gardener/gardener/pkg/apis/core"
	"github.com/gardener/gardener/pkg/apis/core/validation"
)

type namespacedCloudProfileStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy defines the storage strategy for NamespacedCloudProfiles.
var Strategy = namespacedCloudProfileStrategy{api.Scheme, names.SimpleNameGenerator}

func (namespacedCloudProfileStrategy) NamespaceScoped() bool {
	return true
}

func (namespacedCloudProfileStrategy) PrepareForCreate(_ context.Context, obj runtime.Object) {
	namespacedCloudProfile := obj.(*core.NamespacedCloudProfile)

	dropExpiredVersions(namespacedCloudProfile)
	namespacedCloudProfile.Status = core.NamespacedCloudProfileStatus{}
}

func (namespacedCloudProfileStrategy) Validate(_ context.Context, obj runtime.Object) field.ErrorList {
	namespacedCloudProfile := obj.(*core.NamespacedCloudProfile)
	return validation.ValidateNamespacedCloudProfile(namespacedCloudProfile)
}

func (namespacedCloudProfileStrategy) Canonicalize(_ runtime.Object) {
}

func (namespacedCloudProfileStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (namespacedCloudProfileStrategy) PrepareForUpdate(_ context.Context, newObj, oldObj runtime.Object) {
	_ = oldObj.(*core.NamespacedCloudProfile)
	_ = newObj.(*core.NamespacedCloudProfile)
}

func (namespacedCloudProfileStrategy) AllowUnconditionalUpdate() bool {
	return true
}

func (namespacedCloudProfileStrategy) ValidateUpdate(_ context.Context, newObj, oldObj runtime.Object) field.ErrorList {
	oldProfile, newProfile := oldObj.(*core.NamespacedCloudProfile), newObj.(*core.NamespacedCloudProfile)
	return validation.ValidateNamespacedCloudProfileUpdate(newProfile, oldProfile)
}

// WarningsOnCreate returns warnings to the client performing the create.
func (namespacedCloudProfileStrategy) WarningsOnCreate(_ context.Context, _ runtime.Object) []string {
	return nil
}

// WarningsOnUpdate returns warnings to the client performing the update.
func (namespacedCloudProfileStrategy) WarningsOnUpdate(_ context.Context, _, _ runtime.Object) []string {
	return nil
}

func dropExpiredVersions(namespacedCloudProfile *core.NamespacedCloudProfile) {
	if namespacedCloudProfile.Spec.Kubernetes == nil {
		return
	}

	var validKubernetesVersions []core.ExpirableVersion

	for _, version := range namespacedCloudProfile.Spec.Kubernetes.Versions {
		if version.ExpirationDate != nil && version.ExpirationDate.Time.Before(time.Now()) {
			continue
		}
		validKubernetesVersions = append(validKubernetesVersions, version)
	}

	namespacedCloudProfile.Spec.Kubernetes.Versions = validKubernetesVersions

	for i, machineImage := range namespacedCloudProfile.Spec.MachineImages {
		var validMachineImageVersions []core.MachineImageVersion

		for _, version := range machineImage.Versions {
			if version.ExpirationDate != nil && version.ExpirationDate.Time.Before(time.Now()) {
				continue
			}
			validMachineImageVersions = append(validMachineImageVersions, version)
		}

		namespacedCloudProfile.Spec.MachineImages[i].Versions = validMachineImageVersions
	}
}

type namespacedCloudProfileStatusStrategy struct {
	namespacedCloudProfileStrategy
}

// StatusStrategy defines the storage strategy for the status subresource of BackupBuckets.
var StatusStrategy = namespacedCloudProfileStatusStrategy{Strategy}

func (namespacedCloudProfileStatusStrategy) PrepareForUpdate(_ context.Context, obj, old runtime.Object) {
	newBackupBucket := obj.(*core.BackupBucket)
	oldBackupBucket := old.(*core.BackupBucket)
	newBackupBucket.Spec = oldBackupBucket.Spec
}

func (namespacedCloudProfileStatusStrategy) ValidateUpdate(_ context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
