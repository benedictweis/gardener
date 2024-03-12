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

package gardener_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	pkgclient "sigs.k8s.io/controller-runtime/pkg/client"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	gardenerutils "github.com/gardener/gardener/pkg/utils/gardener"
	kubernetesutils "github.com/gardener/gardener/pkg/utils/kubernetes"
	mockclient "github.com/gardener/gardener/third_party/mock/controller-runtime/client"
)

var _ = Describe("NamespacedCloudProfile", func() {
	Describe("#GetCloudProfile", func() {

		var (
			ctrl *gomock.Controller
			c    *mockclient.MockClient

			ctx     = context.TODO()
			fakeErr = fmt.Errorf("fake err")

			namespaceName    = "foo"
			cloudProfileName = "profile-1"

			cloudProfile           *gardencorev1beta1.CloudProfile
			namespacedCloudProfile *gardencorev1beta1.NamespacedCloudProfile
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			c = mockclient.NewMockClient(ctrl)

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

		AfterEach(func() {
			ctrl.Finish()
		})

		It("returns an error if neither a CloudProfile nor a NamespacedCloudProfile could be found", func() {
			c.EXPECT().Get(ctx, kubernetesutils.Key(cloudProfileName), gomock.AssignableToTypeOf(&gardencorev1beta1.CloudProfile{})).Return(fakeErr)

			res, err := gardenerutils.GetCloudProfile(ctx, c, cloudProfileName, namespaceName)
			Expect(res).To(BeNil())
			Expect(err).To(MatchError(fakeErr))
		})

		It("returns CloudProfile if present", func() {
			c.EXPECT().Get(ctx,
				kubernetesutils.Key(cloudProfileName),
				gomock.AssignableToTypeOf(&gardencorev1beta1.CloudProfile{}),
			).DoAndReturn(func(_ context.Context, _ client.ObjectKey, obj *gardencorev1beta1.CloudProfile, _ ...client.GetOption) error {
				cloudProfile.DeepCopyInto(obj)
				return nil
			})

			res, err := gardenerutils.GetCloudProfile(ctx, c, cloudProfileName, namespaceName)
			Expect(res).To(Equal(cloudProfile))
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns NamespacedCloudProfile if present", func() {
			c.EXPECT().Get(ctx, kubernetesutils.Key(cloudProfileName), gomock.AssignableToTypeOf(&gardencorev1beta1.CloudProfile{})).Return(apierrors.NewNotFound(schema.GroupResource{Group: "foo", Resource: "bar"}, cloudProfileName))

			c.EXPECT().Get(ctx,
				pkgclient.ObjectKey{Name: cloudProfileName, Namespace: namespaceName},
				gomock.AssignableToTypeOf(&gardencorev1beta1.NamespacedCloudProfile{}),
			).DoAndReturn(func(_ context.Context, _ client.ObjectKey, obj *gardencorev1beta1.NamespacedCloudProfile, _ ...client.GetOption) error {
				namespacedCloudProfile.DeepCopyInto(obj)
				return nil
			})

			res, err := gardenerutils.GetCloudProfile(ctx, c, cloudProfileName, namespaceName)
			Expect(res.Spec).To(Equal(namespacedCloudProfile.Status.CloudProfileSpec))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
