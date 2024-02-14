// Copyright 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package privatecloudprofile

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

type privateCloudProfileStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy defines the storage strategy for PrivateCloudProfiles.
var Strategy = privateCloudProfileStrategy{api.Scheme, names.SimpleNameGenerator}

func (privateCloudProfileStrategy) NamespaceScoped() bool {
	return true
}

func (privateCloudProfileStrategy) PrepareForCreate(_ context.Context, obj runtime.Object) {
	privateCloudprofile := obj.(*core.PrivateCloudProfile)

	dropExpiredVersions(privateCloudprofile)
}

func (privateCloudProfileStrategy) Validate(_ context.Context, obj runtime.Object) field.ErrorList {
	privateCloudProfile := obj.(*core.PrivateCloudProfile)
	return validation.ValidatePrivateCloudProfile(privateCloudProfile)
}

func (privateCloudProfileStrategy) Canonicalize(_ runtime.Object) {
}

func (privateCloudProfileStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (privateCloudProfileStrategy) PrepareForUpdate(_ context.Context, newObj, oldObj runtime.Object) {
	_ = oldObj.(*core.PrivateCloudProfile)
	_ = newObj.(*core.PrivateCloudProfile)
}

func (privateCloudProfileStrategy) AllowUnconditionalUpdate() bool {
	return true
}

func (privateCloudProfileStrategy) ValidateUpdate(_ context.Context, newObj, oldObj runtime.Object) field.ErrorList {
	oldProfile, newProfile := oldObj.(*core.PrivateCloudProfile), newObj.(*core.PrivateCloudProfile)
	return validation.ValidatePrivateCloudProfileUpdate(newProfile, oldProfile)
}

// WarningsOnCreate returns warnings to the client performing the create.
func (privateCloudProfileStrategy) WarningsOnCreate(_ context.Context, _ runtime.Object) []string {
	return nil
}

// WarningsOnUpdate returns warnings to the client performing the update.
func (privateCloudProfileStrategy) WarningsOnUpdate(_ context.Context, _, _ runtime.Object) []string {
	return nil
}

func dropExpiredVersions(cloudProfile *core.PrivateCloudProfile) {
	var validKubernetesVersions []core.ExpirableVersion

	for _, version := range cloudProfile.Spec.Kubernetes.Versions {
		if version.ExpirationDate != nil && version.ExpirationDate.Time.Before(time.Now()) {
			continue
		}
		validKubernetesVersions = append(validKubernetesVersions, version)
	}

	cloudProfile.Spec.Kubernetes.Versions = validKubernetesVersions

	for i, machineImage := range cloudProfile.Spec.MachineImages {
		var validMachineImageVersions []core.MachineImageVersion

		for _, version := range machineImage.Versions {
			if version.ExpirationDate != nil && version.ExpirationDate.Time.Before(time.Now()) {
				continue
			}
			validMachineImageVersions = append(validMachineImageVersions, version)
		}

		cloudProfile.Spec.MachineImages[i].Versions = validMachineImageVersions
	}
}
