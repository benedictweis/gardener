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

	"github.com/go-logr/logr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/controllerutils/mapper"
)

// ControllerName is the name of this controller.
const ControllerName = "privatecloudprofile"

// AddToManager adds Reconciler to the given manager.
func (r *Reconciler) AddToManager(ctx context.Context, mgr manager.Manager) error {
	if r.Client == nil {
		r.Client = mgr.GetClient()
	}
	if r.Recorder == nil {
		r.Recorder = mgr.GetEventRecorderFor(ControllerName + "-controller")
	}

	c, err := builder.
		ControllerManagedBy(mgr).
		Named(ControllerName).
		For(&gardencorev1beta1.PrivateCloudProfile{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: ptr.Deref(r.Config.ConcurrentSyncs, 0),
		}).
		Build(r)
	if err != nil {
		return err
	}

	return c.Watch(
		source.Kind(mgr.GetCache(), &gardencorev1beta1.CloudProfile{}),
		mapper.EnqueueRequestsFrom(ctx, mgr.GetCache(), mapper.MapFunc(r.MapCloudProfileToPrivateCloudProfile), mapper.UpdateWithNew, c.GetLogger()),
	)
}

// MapCloudProfileToPrivateCloudProfile is a mapper.MapFunc for mapping a core.gardener.cloud/v1beta1.CloudProfile to core.gardener.cloud/v1beta1.PrivateCloudProfile
func (r *Reconciler) MapCloudProfileToPrivateCloudProfile(ctx context.Context, log logr.Logger, _ client.Reader, obj client.Object) []reconcile.Request {
	cloudProfile, ok := obj.(*gardencorev1beta1.CloudProfile)
	if !ok {
		return nil
	}

	privateCloudProfileList := &gardencorev1beta1.PrivateCloudProfileList{}
	// TODO make this work with the proper selector
	if err := r.Client.List(ctx, privateCloudProfileList /*, client.MatchingFields{"spec.parent": cloudProfile.Name}*/); err != nil {
		log.Error(err, "Failed to list privatecloudprofiles referencing this cloudprofile", "cloudProfileName", cloudProfile.Name)
		return nil
	}

	return mapper.ObjectListToRequests(privateCloudProfileList)
}
