// Copyright 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	gardencoreinformers "github.com/gardener/gardener/pkg/client/core/informers/externalversions"
	gardencorev1beta1listers "github.com/gardener/gardener/pkg/client/core/listers/core/v1beta1"
	admissionutils "github.com/gardener/gardener/plugin/pkg/utils"
)

var _ = Describe("NamespacedCloudProfile", func() {
	Describe("#GetCloudProfile", func() {
		var (
			coreInformerFactory          gardencoreinformers.SharedInformerFactory
			cloudProfileLister           gardencorev1beta1listers.CloudProfileLister
			namespacedCloudProfileLister gardencorev1beta1listers.NamespacedCloudProfileLister

			namespaceName    = "foo"
			cloudProfileName = "profile-1"

			cloudProfile           *gardencorev1beta1.CloudProfile
			namespacedCloudProfile *gardencorev1beta1.NamespacedCloudProfile
		)

		BeforeEach(func() {
			coreInformerFactory = gardencoreinformers.NewSharedInformerFactory(nil, 0)
			cloudProfileLister = coreInformerFactory.Core().V1beta1().CloudProfiles().Lister()
			namespacedCloudProfileLister = coreInformerFactory.Core().V1beta1().NamespacedCloudProfiles().Lister()

			cloudProfile = &gardencorev1beta1.CloudProfile{
				ObjectMeta: metav1.ObjectMeta{
					Name: cloudProfileName,
				},
			}

			namespacedCloudProfile = &gardencorev1beta1.NamespacedCloudProfile{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cloudProfileName,
					Namespace: namespaceName,
				},
			}
		})

		It("returns an error if neither a CloudProfile nor a NamespacedCloudProfile could be found", func() {
			res, err := admissionutils.GetCloudProfile(cloudProfileLister, namespacedCloudProfileLister, cloudProfileName, namespaceName)
			Expect(res).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("returns CloudProfile if present", func() {
			Expect(coreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(cloudProfile)).To(Succeed())

			res, err := admissionutils.GetCloudProfile(cloudProfileLister, namespacedCloudProfileLister, cloudProfileName, namespaceName)
			Expect(res).To(Equal(cloudProfile))
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns NamespacedCloudProfile if present", func() {
			Expect(coreInformerFactory.Core().V1beta1().NamespacedCloudProfiles().Informer().GetStore().Add(namespacedCloudProfile)).To(Succeed())

			res, err := admissionutils.GetCloudProfile(cloudProfileLister, namespacedCloudProfileLister, cloudProfileName, namespaceName)
			Expect(res.Spec).To(Equal(namespacedCloudProfile.Status.CloudProfileSpec))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
