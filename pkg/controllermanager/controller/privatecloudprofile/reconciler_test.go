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

package privatecloudprofile

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	mockclient "github.com/gardener/gardener/pkg/mock/controller-runtime/client"
	kubernetesutils "github.com/gardener/gardener/pkg/utils/kubernetes"
)

var _ = Describe("Reconciler", func() {
	var (
		ctx  = context.TODO()
		ctrl *gomock.Controller
		c    *mockclient.MockClient

		privateCloudProfileName string
		parentCloudProfileName  string
		fakeErr                 error
		reconciler              reconcile.Reconciler
		_                       *gardencorev1beta1.PrivateCloudProfile
		_                       *gardencorev1beta1.CloudProfile
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		c = mockclient.NewMockClient(ctrl)

		privateCloudProfileName = "test-privatecloudprofile"
		parentCloudProfileName = "test-cloudprofile"
		fakeErr = fmt.Errorf("fake err")
		reconciler = &Reconciler{Client: c, Recorder: &record.FakeRecorder{}}
		_ = &gardencorev1beta1.PrivateCloudProfile{
			ObjectMeta: metav1.ObjectMeta{
				Name: privateCloudProfileName,
			},
			Spec: gardencorev1beta1.PrivateCloudProfileSpec{
				Parent: parentCloudProfileName,
			},
		}
		_ = &gardencorev1beta1.CloudProfile{
			ObjectMeta: metav1.ObjectMeta{
				Name: parentCloudProfileName,
			},
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should return nil because object not found", func() {
		c.EXPECT().Get(gomock.Any(), kubernetesutils.Key(privateCloudProfileName), gomock.AssignableToTypeOf(&gardencorev1beta1.PrivateCloudProfile{})).Return(apierrors.NewNotFound(schema.GroupResource{}, ""))

		result, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: privateCloudProfileName}})
		Expect(result).To(Equal(reconcile.Result{}))
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return err because object reading failed", func() {
		c.EXPECT().Get(gomock.Any(), kubernetesutils.Key(privateCloudProfileName), gomock.AssignableToTypeOf(&gardencorev1beta1.PrivateCloudProfile{})).Return(fakeErr)

		result, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: privateCloudProfileName}})
		Expect(result).To(Equal(reconcile.Result{}))
		Expect(err).To(MatchError(fakeErr))
	})

	//It("should return an err because object reading failed on parent cloud profile", func() {
	//	c.EXPECT().Get(gomock.Any(), kubernetesutils.Key(privateCloudProfileName), gomock.AssignableToTypeOf(&gardencorev1beta1.PrivateCloudProfile{})).DoAndReturn(func(_ context.Context, _ client.ObjectKey, obj *gardencorev1beta1.PrivateCloudProfile, _ ...client.GetOption) error {
	//		*obj = *privateCloudProfile
	//		return nil
	//	})
	//
	//	c.EXPECT().Get(gomock.Any(), kubernetesutils.Key(parentCloudProfileName), gomock.AssignableToTypeOf(&gardencorev1beta1.CloudProfile{})).DoAndReturn(func(_ context.Context, _ client.ObjectKey, obj *gardencorev1beta1.CloudProfile, _ ...client.GetOption) error {
	//		*obj = *parentCloudProfile
	//		return nil
	//	})
	//
	//	c.EXPECT().Patch(gomock.Any(), gomock.AssignableToTypeOf(&gardencorev1beta1.PrivateCloudProfile{}), gomock.Any()).DoAndReturn(func(_ context.Context, o client.Object, patch client.Patch, opts ...client.PatchOption) error {
	//		Expect(patch.Data(o)).To(BeEquivalentTo(fmt.Sprintf(`{"status":{"cloudProfile":{"metadata":{"name":"test-cloudprofile"}}}}`)))
	//		return nil
	//	})
	//
	//	result, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: privateCloudProfileName}})
	//	Expect(result).To(Equal(reconcile.Result{}))
	//	Expect(err).To(MatchError(ContainSubstring("hi")))
	//})
})
