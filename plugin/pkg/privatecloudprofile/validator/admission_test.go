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

	"github.com/gardener/gardener/pkg/apis/core"
	gardencoreinformers "github.com/gardener/gardener/pkg/client/core/informers/internalversion"
	. "github.com/gardener/gardener/plugin/pkg/privatecloudprofile/validator"
)

var _ = Describe("Admission", func() {
	Describe("#Validate", func() {
		var (
			ctx                 context.Context
			admissionHandler    *ValidatePrivateCloudProfile
			coreInformerFactory gardencoreinformers.SharedInformerFactory

			privateCloudProfile core.PrivateCloudProfile
			parentCloudProfile  core.CloudProfile
			machineType         core.MachineType

			parentCloudProfileName = "parent-profile"

			privateCloudProfileBase = core.PrivateCloudProfile{
				ObjectMeta: metav1.ObjectMeta{
					Name: "profile",
				},
			}
			parentCloudProfileBase = core.CloudProfile{
				ObjectMeta: metav1.ObjectMeta{
					Name: parentCloudProfileName,
				},
			}
			machineTypeBase = core.MachineType{
				Name: "my-machine",
			}
		)

		BeforeEach(func() {
			ctx = context.TODO()

			privateCloudProfile = *privateCloudProfileBase.DeepCopy()
			parentCloudProfile = *parentCloudProfileBase.DeepCopy()
			machineType = machineTypeBase

			admissionHandler, _ = New()
			admissionHandler.AssignReadyFunc(func() bool { return true })
			coreInformerFactory = gardencoreinformers.NewSharedInformerFactory(nil, 0)
			admissionHandler.SetInternalCoreInformerFactory(coreInformerFactory)
		})

		It("should not allow creating a PrivateCloudProfile with an invalid parent reference", func() {
			privateCloudProfile.Spec.Parent = "idontexist"

			attrs := admission.NewAttributesRecord(&privateCloudProfile, nil, core.Kind("PrivateCloudProfile").WithVersion("version"), "", privateCloudProfile.Name, core.Resource("privatecloudprofile").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(MatchError(ContainSubstring("parent CloudProfile could not be found")))
		})

		It("should allow creating a PrivateCloudProfile with a valid parent reference", func() {
			Expect(coreInformerFactory.Core().InternalVersion().CloudProfiles().Informer().GetStore().Add(&parentCloudProfile)).To(Succeed())

			privateCloudProfile.Spec.Parent = parentCloudProfileName

			attrs := admission.NewAttributesRecord(&privateCloudProfile, nil, core.Kind("PrivateCloudProfile").WithVersion("version"), "", privateCloudProfile.Name, core.Resource("privatecloudprofile").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(Succeed())
		})

		It("should not allow creating a PrivateCloudProfile that defines a machineType of the parent CloudProfile", func() {
			parentCloudProfile.Spec.MachineTypes = []core.MachineType{machineType}
			Expect(coreInformerFactory.Core().InternalVersion().CloudProfiles().Informer().GetStore().Add(&parentCloudProfile)).To(Succeed())

			privateCloudProfile.Spec.Parent = parentCloudProfileName
			privateCloudProfile.Spec.MachineTypes = []core.MachineType{machineType}

			attrs := admission.NewAttributesRecord(&privateCloudProfile, nil, core.Kind("PrivateCloudProfile").WithVersion("version"), "", privateCloudProfile.Name, core.Resource("privatecloudprofile").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(MatchError(ContainSubstring("PrivateCloudProfile attempts to rewrite MachineType of parent CloudProfile with machineType")))
			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(MatchError(ContainSubstring("my-machine")))
		})

		It("should allow creating a PrivateCloudProfile that defines a machineType of the parent CloudProfile if it was added to the PrivateCloudProfile first", func() {
			Expect(coreInformerFactory.Core().InternalVersion().CloudProfiles().Informer().GetStore().Add(&parentCloudProfile)).To(Succeed())

			privateCloudProfile.Spec.Parent = parentCloudProfileName
			privateCloudProfile.Spec.MachineTypes = []core.MachineType{machineType}

			attrs := admission.NewAttributesRecord(&privateCloudProfile, nil, core.Kind("PrivateCloudProfile").WithVersion("version"), "", privateCloudProfile.Name, core.Resource("privatecloudprofile").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(Succeed())

			parentCloudProfile.Spec.MachineTypes = []core.MachineType{machineType}

			attrs = admission.NewAttributesRecord(&privateCloudProfile, &privateCloudProfile, core.Kind("PrivateCloudProfile").WithVersion("version"), "", privateCloudProfile.Name, core.Resource("privatecloudprofile").WithVersion("version"), "", admission.Update, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(Succeed())
		})

		It("should allow creating a PrivateCloudProfile that defines a machineType of the parent CloudProfile if it was added to the PrivateCloudProfile first but is changed", func() {
			Expect(coreInformerFactory.Core().InternalVersion().CloudProfiles().Informer().GetStore().Add(&parentCloudProfile)).To(Succeed())

			privateCloudProfile.Spec.Parent = parentCloudProfileName
			privateCloudProfile.Spec.MachineTypes = []core.MachineType{machineType}

			attrs := admission.NewAttributesRecord(&privateCloudProfile, nil, core.Kind("PrivateCloudProfile").WithVersion("version"), "", privateCloudProfile.Name, core.Resource("privatecloudprofile").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(Succeed())

			oldPrivateCloudProfile := *privateCloudProfile.DeepCopy()
			privateCloudProfile.Spec.MachineImages = []core.MachineImage{{Name: "my-image"}}
			parentCloudProfile.Spec.MachineTypes = []core.MachineType{machineType}

			attrs = admission.NewAttributesRecord(&privateCloudProfile, &oldPrivateCloudProfile, core.Kind("PrivateCloudProfile").WithVersion("version"), "", privateCloudProfile.Name, core.Resource("privatecloudprofile").WithVersion("version"), "", admission.Update, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(Succeed())
		})

		It("should allow creating a PrivateCloudProfile that defines a different machineType than the parent CloudProfile", func() {
			Expect(coreInformerFactory.Core().InternalVersion().CloudProfiles().Informer().GetStore().Add(&parentCloudProfile)).To(Succeed())

			privateCloudProfile.Spec.Parent = parentCloudProfileName
			privateCloudProfile.Spec.MachineTypes = []core.MachineType{{Name: "my-other-machine"}}

			parentCloudProfile.Spec.MachineTypes = []core.MachineType{machineType}

			attrs := admission.NewAttributesRecord(&privateCloudProfile, nil, core.Kind("PrivateCloudProfile").WithVersion("version"), "", privateCloudProfile.Name, core.Resource("privatecloudprofile").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(ctx, attrs, nil)).To(Succeed())
		})
	})

	Describe("#Register", func() {
		It("should register the plugin", func() {
			plugins := admission.NewPlugins()
			Register(plugins)

			registered := plugins.Registered()
			Expect(registered).To(HaveLen(1))
			Expect(registered).To(ContainElement("PrivateCloudProfileValidator"))
		})
	})

	Describe("#New", func() {
		It("should only handle CREATE operations", func() {
			dr, err := New()
			Expect(err).ToNot(HaveOccurred())
			Expect(dr.Handles(admission.Create)).To(BeTrue())
			Expect(dr.Handles(admission.Update)).To(BeTrue())
			Expect(dr.Handles(admission.Connect)).To(BeFalse())
			Expect(dr.Handles(admission.Delete)).To(BeFalse())
		})
	})
})
