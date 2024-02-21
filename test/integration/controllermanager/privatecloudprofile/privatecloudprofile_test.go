// Copyright 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package privatecloudprofile_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	. "github.com/gardener/gardener/pkg/utils/test/matchers"
)

var _ = Describe("PrivateCloudProfile controller tests", func() {
	const parentCloudProfileName = testID + "-my-profile"

	var (
		parentCloudProfile  *gardencorev1beta1.CloudProfile
		privateCloudProfile *gardencorev1beta1.PrivateCloudProfile
		mergedCloudProfile  *gardencorev1beta1.CloudProfile
		shoot               *gardencorev1beta1.Shoot
	)

	BeforeEach(func() {
		parentCloudProfile = &gardencorev1beta1.CloudProfile{
			ObjectMeta: metav1.ObjectMeta{
				Name: parentCloudProfileName,
			},
			Spec: gardencorev1beta1.CloudProfileSpec{
				Type: "some-type",
				Kubernetes: gardencorev1beta1.KubernetesSettings{
					Versions: []gardencorev1beta1.ExpirableVersion{{Version: "1.2.3"}},
				},
				MachineImages: []gardencorev1beta1.MachineImage{
					{
						Name: "some-image",
						Versions: []gardencorev1beta1.MachineImageVersion{
							{ExpirableVersion: gardencorev1beta1.ExpirableVersion{Version: "4.5.6"}},
						},
					},
				},
				MachineTypes: []gardencorev1beta1.MachineType{{
					Name:   "some-type",
					CPU:    resource.MustParse("1"),
					GPU:    resource.MustParse("0"),
					Memory: resource.MustParse("1Gi"),
				}},
				Regions: []gardencorev1beta1.Region{
					{Name: "some-region"},
				},
			},
		}

		privateCloudProfile = &gardencorev1beta1.PrivateCloudProfile{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: testID + "-",
				Namespace:    testNamespace.Name,
				Labels:       map[string]string{testID: testRunID},
			},
			Spec: gardencorev1beta1.PrivateCloudProfileSpec{
				Parent: parentCloudProfileName,
				Kubernetes: &gardencorev1beta1.KubernetesSettings{
					Versions: []gardencorev1beta1.ExpirableVersion{{Version: "1.2.4"}},
				},
				MachineImages: []gardencorev1beta1.MachineImage{
					{
						Name: "some-other-image",
						Versions: []gardencorev1beta1.MachineImageVersion{
							{ExpirableVersion: gardencorev1beta1.ExpirableVersion{Version: "4.5.7"}},
						},
					},
				},
				MachineTypes: []gardencorev1beta1.MachineType{{
					Name:   "some-other-type",
					CPU:    resource.MustParse("2"),
					GPU:    resource.MustParse("0"),
					Memory: resource.MustParse("2Gi"),
				}},
				Regions: []gardencorev1beta1.Region{
					{Name: "some-other-region"},
				},
			},
		}

		mergedCloudProfile = &gardencorev1beta1.CloudProfile{
			ObjectMeta: metav1.ObjectMeta{
				Name: parentCloudProfileName,
			},
			Spec: gardencorev1beta1.CloudProfileSpec{
				Type: "some-type",
				Kubernetes: gardencorev1beta1.KubernetesSettings{
					Versions: []gardencorev1beta1.ExpirableVersion{{Version: "1.2.3"}, {Version: "1.2.4"}},
				},
				MachineImages: []gardencorev1beta1.MachineImage{
					{
						Name: "some-image",
						Versions: []gardencorev1beta1.MachineImageVersion{
							{ExpirableVersion: gardencorev1beta1.ExpirableVersion{Version: "4.5.6"}},
						},
					},
					{
						Name: "some-other-image",
						Versions: []gardencorev1beta1.MachineImageVersion{
							{ExpirableVersion: gardencorev1beta1.ExpirableVersion{Version: "4.5.7"}},
						},
					},
				},
				MachineTypes: []gardencorev1beta1.MachineType{
					{
						Name:   "some-type",
						CPU:    resource.MustParse("1"),
						GPU:    resource.MustParse("0"),
						Memory: resource.MustParse("1Gi"),
					},
					{
						Name:   "some-other-type",
						CPU:    resource.MustParse("2"),
						GPU:    resource.MustParse("0"),
						Memory: resource.MustParse("2 i"),
					}},
				Regions: []gardencorev1beta1.Region{
					{Name: "some-region"},
					{Name: "some-other-region"},
				},
			},
		}

		shoot = &gardencorev1beta1.Shoot{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: testID + "-",
				Namespace:    testNamespace.Name,
			},
			Spec: gardencorev1beta1.ShootSpec{
				SecretBindingName: ptr.To("my-provider-account"),
				Region:            "foo-region",
				Provider: gardencorev1beta1.Provider{
					Type: privateCloudProfile.Spec.Parent,
					Workers: []gardencorev1beta1.Worker{
						{
							Name:    "cpu-worker",
							Minimum: 2,
							Maximum: 2,
							Machine: gardencorev1beta1.Machine{Type: "large"},
						},
					},
				},
				Kubernetes: gardencorev1beta1.Kubernetes{Version: "1.26.1"},
				Networking: &gardencorev1beta1.Networking{Type: ptr.To("foo-networking")},
			},
		}
	})

	JustBeforeEach(func() {
		By("Create parent CloudProfile")
		Expect(testClient.Create(ctx, parentCloudProfile)).To(Succeed())
		log.Info("Created parent CloudProfile for test", "parentCloudProfile", client.ObjectKeyFromObject(parentCloudProfile))

		By("Create PrivateCloudProfile")
		Expect(testClient.Create(ctx, privateCloudProfile)).To(Succeed())
		log.Info("Created PrivateCloudProfile for test", "privateCloudProfile", client.ObjectKeyFromObject(privateCloudProfile))

		DeferCleanup(func() {
			By("Delete ParentCloudProfile")
			Expect(client.IgnoreNotFound(testClient.Delete(ctx, parentCloudProfile))).To(Succeed())

			By("Delete PrivateCloudProfile")
			Expect(client.IgnoreNotFound(testClient.Delete(ctx, privateCloudProfile))).To(Succeed())
		})

		if shoot != nil {
			By("Create Shoot")
			shoot.Spec.CloudProfileName = privateCloudProfile.Name
			Expect(testClient.Create(ctx, shoot)).To(Succeed())
			log.Info("Created shoot for test", "shoot", client.ObjectKeyFromObject(shoot))

			By("Wait until manager has observed Shoot")
			// Use the manager's cache to ensure it has observed the Shoot.
			// Otherwise, the controller might clean up the PrivateCloudProfile too early because it thinks all referencing
			// Shoots are gone. Similar to https://github.com/gardener/gardener/issues/6486 and
			// https://github.com/gardener/gardener/issues/6607.
			Eventually(func() error {
				return mgrClient.Get(ctx, client.ObjectKeyFromObject(shoot), &gardencorev1beta1.Shoot{})
			}).Should(Succeed())

			DeferCleanup(func() {
				By("Delete Shoot")
				Expect(client.IgnoreNotFound(testClient.Delete(ctx, shoot))).To(Succeed())
			})
		}
	})

	Context("no shoot referencing the PrivateCloudProfile", func() {
		BeforeEach(func() {
			shoot = nil
		})

		It("should add the finalizer and release it on deletion", func() {
			By("Ensure finalizer got added")
			Eventually(func(g Gomega) {
				g.Expect(testClient.Get(ctx, client.ObjectKeyFromObject(privateCloudProfile), privateCloudProfile)).To(Succeed())
				g.Expect(privateCloudProfile.Finalizers).To(ConsistOf("gardener"))
			}).Should(Succeed())

			By("Delete PrivateCloudProfile")
			Expect(testClient.Delete(ctx, privateCloudProfile)).To(Succeed())

			By("Ensure PrivateCloudProfile is released")
			Eventually(func() error {
				return testClient.Get(ctx, client.ObjectKeyFromObject(privateCloudProfile), privateCloudProfile)
			}).Should(BeNotFoundError())
		})
	})

	Context("shoots referencing the PrivateCloudProfile", func() {
		JustBeforeEach(func() {
			By("Ensure finalizer got added")
			Eventually(func(g Gomega) {
				g.Expect(testClient.Get(ctx, client.ObjectKeyFromObject(privateCloudProfile), privateCloudProfile)).To(Succeed())
				g.Expect(privateCloudProfile.Finalizers).To(ConsistOf("gardener"))
			}).Should(Succeed())

			By("Delete PrivateCloudProfile")
			Expect(testClient.Delete(ctx, privateCloudProfile)).To(Succeed())
		})

		It("should add the finalizer and not release it on deletion since there still is a referencing shoot", func() {
			By("Ensure PrivateCloudProfile is not released")
			Consistently(func() error {
				return testClient.Get(ctx, client.ObjectKeyFromObject(privateCloudProfile), privateCloudProfile)
			}).Should(Succeed())
		})

		It("should add the finalizer and release it on deletion after the shoot got deleted", func() {
			By("Delete Shoot")
			Expect(testClient.Delete(ctx, shoot)).To(Succeed())

			By("Ensure PrivateCloudProfile is released")
			Eventually(func() error {
				return testClient.Get(ctx, client.ObjectKeyFromObject(privateCloudProfile), privateCloudProfile)
			}).Should(BeNotFoundError())
		})
	})

	Context("merging the PrivateCloudProfile with the parent Cloud Profile", func() {
		FIt("should merge the PrivateCloudProfile with the parent Cloud profile correctly", func() {
			By("Ensure PrivateCloudProfile is present")
			Eventually(func() error {
				return testClient.Get(ctx, client.ObjectKeyFromObject(privateCloudProfile), privateCloudProfile)
			}).Should(Succeed())

			By("Ensure PrivateCloudProfile rendered CloudProfile was correctly merged")
			Expect(privateCloudProfile.Status.CloudProfile.Name).To(Equal(mergedCloudProfile.Name))
			Expect(privateCloudProfile.Status.CloudProfile.CreationTimestamp).To(Not(BeNil()))
			Expect(privateCloudProfile.Status.CloudProfile.Spec).To(Equal(mergedCloudProfile.Spec))
		})
	})
})
