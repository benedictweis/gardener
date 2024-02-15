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

package validation

import (
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1validation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/gardener/gardener/pkg/apis/core"
	"github.com/gardener/gardener/pkg/utils"
)

// ValidatePrivateCloudProfile validates a CloudProfile object.
func ValidatePrivateCloudProfile(cloudProfile *core.PrivateCloudProfile) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, apivalidation.ValidateObjectMeta(&cloudProfile.ObjectMeta, true, ValidateName, field.NewPath("metadata"))...)
	allErrs = append(allErrs, ValidatePrivateCloudProfileSpec(&cloudProfile.Spec, field.NewPath("spec"))...)

	return allErrs
}

// ValidatePrivateCloudProfileUpdate validates a CloudProfile object before an update.
func ValidatePrivateCloudProfileUpdate(newProfile, oldProfile *core.PrivateCloudProfile) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, apivalidation.ValidateObjectMetaUpdate(&newProfile.ObjectMeta, &oldProfile.ObjectMeta, field.NewPath("metadata"))...)
	allErrs = append(allErrs, ValidatePrivateCloudProfileSpecUpdate(&newProfile.Spec, &oldProfile.Spec, field.NewPath("spec"))...)
	allErrs = append(allErrs, ValidatePrivateCloudProfile(newProfile)...)

	return allErrs
}

// ValidatePrivateCloudProfileSpecUpdate validates the spec update of a CloudProfile
func ValidatePrivateCloudProfileSpecUpdate(_, _ *core.PrivateCloudProfileSpec, _ *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	return allErrs
}

// ValidatePrivateCloudProfileSpec validates the specification of a CloudProfile object.
func ValidatePrivateCloudProfileSpec(spec *core.PrivateCloudProfileSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(spec.Parent) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("parent"), "must provide a parent cloud profile"))
	}

	if spec.Kubernetes != nil {
		allErrs = append(allErrs, validateKubernetesSettings(*spec.Kubernetes, fldPath.Child("kubernetes"))...)
	}
	if spec.MachineTypes != nil {
		allErrs = append(allErrs, validateMachineImages(spec.MachineImages, fldPath.Child("machineImages"))...)
	}
	if spec.MachineTypes != nil {
		allErrs = append(allErrs, validateMachineTypes(spec.MachineTypes, fldPath.Child("machineTypes"))...)
	}
	if spec.VolumeTypes != nil {
		allErrs = append(allErrs, validateVolumeTypes(spec.VolumeTypes, fldPath.Child("volumeTypes"))...)
	}
	if spec.Regions != nil {
		allErrs = append(allErrs, validateRegions(spec.Regions, fldPath.Child("regions"))...)
	}
	if spec.SeedSelector != nil {
		allErrs = append(allErrs, metav1validation.ValidateLabelSelector(&spec.SeedSelector.LabelSelector, metav1validation.LabelSelectorValidationOptions{AllowInvalidLabelValueInSelector: true}, fldPath.Child("seedSelector"))...)
	}
	if spec.CABundle != nil {
		_, err := utils.DecodeCertificate([]byte(*(spec.CABundle)))
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("caBundle"), *(spec.CABundle), "caBundle is not a valid PEM-encoded certificate"))
		}
	}

	return allErrs
}
