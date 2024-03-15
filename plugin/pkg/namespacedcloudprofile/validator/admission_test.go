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

package validator_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/admission"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	gardencoreinformers "github.com/gardener/gardener/pkg/client/core/informers/externalversions"
	. "github.com/gardener/gardener/plugin/pkg/namespacedcloudprofile/validator"
)

var _ = Describe("Admission", func() {
	Describe("#Validate", func() {
		var (
			ctx                 context.Context
			admissionHandler    *ValidateNamespacedCloudProfile
			coreInformerFactory gardencoreinformers.SharedInformerFactory

			namespacedCloudProfile       gardencorev1beta1.NamespacedCloudProfile
			namespacedCloudProfileParent gardencorev1beta1.CloudProfileReference
			parentCloudProfile           gardencorev1beta1.CloudProfile
			machineType                  gardencorev1beta1.MachineType

			namespacedCloudProfileBase = gardencorev1beta1.NamespacedCloudProfile{
				ObjectMeta: metav1.ObjectMeta{
					Name: "profile",
				},
			}
			parentCloudProfileBase = gardencorev1beta1.CloudProfile{
				ObjectMeta: metav1.ObjectMeta{
					Name: "parent-profile",
				},
			}
			machineTypeBase = gardencorev1beta1.MachineType{
				Name: "my-machine",
			}
		)

		BeforeEach(func() {
			ctx = context.TODO()

			namespacedCloudProfile = *namespacedCloudProfileBase.DeepCopy()
			namespacedCloudProfileParent = gardencorev1beta1.CloudProfileReference{
				Kind: "CloudProfile",
				Name: parentCloudProfileBase.Name,
			}
			parentCloudProfile = *parentCloudProfileBase.DeepCopy()
			machineType = machineTypeBase

			admissionHandler, _ = New()
			admissionHandler.AssignReadyFunc(func() bool { return true })
			coreInformerFactory = gardencoreinformers.NewSharedInformerFactory(nil, 0)
			admissionHandler.SetCoreInformerFactory(coreInformerFactory)
		})

		It("should not allow creating a NamespacedCloudProfile with an invalid parent reference", func() {
			namespacedCloudProfile.Spec.Parent = gardencorev1beta1.CloudProfileReference{Kind: "CloudProfile", Name: "idontexist"}

			attrs := admission.NewAttributesRecord(&namespacedCloudProfile, nil, gardencorev1beta1.Kind("NamespacedCloudProfile").WithVersion("version"), "", namespacedCloudProfile.Name, gardencorev1beta1.Resource("namespacedcloudprofile").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(MatchError(ContainSubstring("parent CloudProfile could not be found")))
		})

		It("should allow creating a NamespacedCloudProfile with a valid parent reference", func() {
			Expect(coreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(&parentCloudProfile)).To(Succeed())

			namespacedCloudProfile.Spec.Parent = namespacedCloudProfileParent

			attrs := admission.NewAttributesRecord(&namespacedCloudProfile, nil, gardencorev1beta1.Kind("NamespacedCloudProfile").WithVersion("version"), "", namespacedCloudProfile.Name, gardencorev1beta1.Resource("namespacedcloudprofile").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(Succeed())
		})

		It("should not allow creating a NamespacedCloudProfile that defines a machineType of the parent CloudProfile", func() {
			parentCloudProfile.Spec.MachineTypes = []gardencorev1beta1.MachineType{machineType}
			Expect(coreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(&parentCloudProfile)).To(Succeed())

			namespacedCloudProfile.Spec.Parent = namespacedCloudProfileParent
			namespacedCloudProfile.Spec.MachineTypes = []gardencorev1beta1.MachineType{machineType}

			attrs := admission.NewAttributesRecord(&namespacedCloudProfile, nil, gardencorev1beta1.Kind("NamespacedCloudProfile").WithVersion("version"), "", namespacedCloudProfile.Name, gardencorev1beta1.Resource("namespacedcloudprofile").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(MatchError(ContainSubstring("NamespacedCloudProfile attempts to rewrite MachineType of parent CloudProfile with machineType")))
			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(MatchError(ContainSubstring("my-machine")))
		})

		It("should allow creating a NamespacedCloudProfile that defines a machineType of the parent CloudProfile if it was added to the NamespacedCloudProfile first", func() {
			Expect(coreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(&parentCloudProfile)).To(Succeed())

			namespacedCloudProfile.Spec.Parent = namespacedCloudProfileParent
			namespacedCloudProfile.Spec.MachineTypes = []gardencorev1beta1.MachineType{machineType}

			attrs := admission.NewAttributesRecord(&namespacedCloudProfile, nil, gardencorev1beta1.Kind("NamespacedCloudProfile").WithVersion("version"), "", namespacedCloudProfile.Name, gardencorev1beta1.Resource("namespacedcloudprofile").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(Succeed())

			parentCloudProfile.Spec.MachineTypes = []gardencorev1beta1.MachineType{machineType}

			attrs = admission.NewAttributesRecord(&namespacedCloudProfile, &namespacedCloudProfile, gardencorev1beta1.Kind("NamespacedCloudProfile").WithVersion("version"), "", namespacedCloudProfile.Name, gardencorev1beta1.Resource("namespacedcloudprofile").WithVersion("version"), "", admission.Update, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(Succeed())
		})

		It("should allow creating a NamespacedCloudProfile that defines a machineType of the parent CloudProfile if it was added to the NamespacedCloudProfile first but is changed", func() {
			Expect(coreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(&parentCloudProfile)).To(Succeed())

			namespacedCloudProfile.Spec.Parent = namespacedCloudProfileParent
			namespacedCloudProfile.Spec.MachineTypes = []gardencorev1beta1.MachineType{machineType}

			attrs := admission.NewAttributesRecord(&namespacedCloudProfile, nil, gardencorev1beta1.Kind("NamespacedCloudProfile").WithVersion("version"), "", namespacedCloudProfile.Name, gardencorev1beta1.Resource("namespacedcloudprofile").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(Succeed())

			oldNamespacedCloudProfile := *namespacedCloudProfile.DeepCopy()
			namespacedCloudProfile.Spec.MachineImages = []gardencorev1beta1.MachineImage{{Name: "my-image"}}
			parentCloudProfile.Spec.MachineTypes = []gardencorev1beta1.MachineType{machineType}

			attrs = admission.NewAttributesRecord(&namespacedCloudProfile, &oldNamespacedCloudProfile, gardencorev1beta1.Kind("NamespacedCloudProfile").WithVersion("version"), "", namespacedCloudProfile.Name, gardencorev1beta1.Resource("namespacedcloudprofile").WithVersion("version"), "", admission.Update, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(Succeed())
		})

		It("should allow creating a NamespacedCloudProfile that defines a different machineType than the parent CloudProfile", func() {
			Expect(coreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(&parentCloudProfile)).To(Succeed())

			namespacedCloudProfile.Spec.Parent = namespacedCloudProfileParent
			namespacedCloudProfile.Spec.MachineTypes = []gardencorev1beta1.MachineType{{Name: "my-other-machine"}}

			parentCloudProfile.Spec.MachineTypes = []gardencorev1beta1.MachineType{machineType}

			attrs := admission.NewAttributesRecord(&namespacedCloudProfile, nil, gardencorev1beta1.Kind("NamespacedCloudProfile").WithVersion("version"), "", namespacedCloudProfile.Name, gardencorev1beta1.Resource("namespacedcloudprofile").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(Succeed())
		})
	})

	Describe("#Register", func() {
		It("should register the plugin", func() {
			plugins := admission.NewPlugins()
			Register(plugins)

			registered := plugins.Registered()
			Expect(registered).To(HaveLen(1))
			Expect(registered).To(ContainElement("NamespacedCloudProfileValidator"))
		})
	})

	Describe("#New", func() {
		It("should only handle CREATE and UPDATE operations", func() {
			dr, err := New()
			Expect(err).ToNot(HaveOccurred())
			Expect(dr.Handles(admission.Create)).To(BeTrue())
			Expect(dr.Handles(admission.Update)).To(BeTrue())
			Expect(dr.Handles(admission.Connect)).To(BeFalse())
			Expect(dr.Handles(admission.Delete)).To(BeFalse())
		})
	})
})
