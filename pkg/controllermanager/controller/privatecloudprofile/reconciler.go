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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/pkg/controllermanager/apis/config"
	"github.com/gardener/gardener/pkg/controllerutils"
)

// Reconciler reconciles CloudProfiles.
type Reconciler struct {
	Client   client.Client
	Config   config.PrivateCloudProfileControllerConfiguration
	Recorder record.EventRecorder
}

// Reconcile performs the main reconciliation logic.
func (r *Reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := logf.FromContext(ctx)

	ctx, cancel := controllerutils.GetMainReconciliationContext(ctx, controllerutils.DefaultReconciliationTimeout)
	defer cancel()

	privateCloudProfile := &gardencorev1beta1.PrivateCloudProfile{}
	if err := r.Client.Get(ctx, request.NamespacedName, privateCloudProfile); err != nil {
		if apierrors.IsNotFound(err) {
			log.V(1).Info("Object is gone, stop reconciling")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, fmt.Errorf("error retrieving object from store: %w", err)
	}

	parentCloudProfile := &gardencorev1beta1.CloudProfile{}
	if err := r.Client.Get(ctx, client.ObjectKey{Name: privateCloudProfile.Spec.Parent}, parentCloudProfile); err != nil {
		if apierrors.IsNotFound(err) {
			log.V(1).Info("Parent object is gone, stop reconciling")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, fmt.Errorf("error retrieving object from store: %w", err)
	}

	if err := patchCloudProfileStatus(ctx, r.Client, privateCloudProfile, parentCloudProfile); err != nil {
		return reconcile.Result{}, err
	}

	// The deletionTimestamp labels the PrivateCloudProfile as intended to get deleted. Before deletion, it has to be ensured that
	// no Shoots and Seed are assigned to the PrivateCloudProfile anymore. If this is the case then the controller will remove
	// the finalizers from the PrivateCloudProfile so that it can be garbage collected.
	if privateCloudProfile.DeletionTimestamp != nil {
		if !sets.New(privateCloudProfile.Finalizers...).Has(gardencorev1beta1.GardenerName) {
			return reconcile.Result{}, nil
		}

		associatedShoots, err := controllerutils.DetermineShootsAssociatedTo(ctx, r.Client, privateCloudProfile)
		if err != nil {
			return reconcile.Result{}, err
		}

		if len(associatedShoots) == 0 {
			log.Info("No Shoots are referencing the PrivateCloudProfile, deletion accepted")

			if controllerutil.ContainsFinalizer(privateCloudProfile, gardencorev1beta1.GardenerName) {
				log.Info("Removing finalizer")
				if err := controllerutils.RemoveFinalizers(ctx, r.Client, privateCloudProfile, gardencorev1beta1.GardenerName); err != nil {
					return reconcile.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
				}
			}

			return reconcile.Result{}, nil
		}

		message := fmt.Sprintf("Cannot delete PrivateCloudProfile, because the following Shoots are still referencing it: %+v", associatedShoots)
		r.Recorder.Event(privateCloudProfile, corev1.EventTypeNormal, v1beta1constants.EventResourceReferenced, message)
		return reconcile.Result{}, fmt.Errorf(message)
	}

	if !controllerutil.ContainsFinalizer(privateCloudProfile, gardencorev1beta1.GardenerName) {
		log.Info("Adding finalizer")
		if err := controllerutils.AddFinalizers(ctx, r.Client, privateCloudProfile, gardencorev1beta1.GardenerName); err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
	}

	return reconcile.Result{}, nil
}

func patchCloudProfileStatus(ctx context.Context, c client.Client, privateCloudProfile *gardencorev1beta1.PrivateCloudProfile, parentCloudProfile *gardencorev1beta1.CloudProfile) error {
	patch := client.StrategicMergeFrom(privateCloudProfile.DeepCopy())
	mergeCloudProfiles(parentCloudProfile, privateCloudProfile)
	privateCloudProfile.Status.CloudProfile = *parentCloudProfile
	return c.Patch(ctx, privateCloudProfile, patch)
}

func mergeCloudProfiles(parentCloudProfile *gardencorev1beta1.CloudProfile, privateCloudProfile *gardencorev1beta1.PrivateCloudProfile) {
	if privateCloudProfile.Spec.Kubernetes != nil {
		parentCloudProfile.Spec.Kubernetes.Versions = append(parentCloudProfile.Spec.Kubernetes.Versions, privateCloudProfile.Spec.Kubernetes.Versions...)
	}
	parentCloudProfile.Spec.MachineImages = append(parentCloudProfile.Spec.MachineImages, privateCloudProfile.Spec.MachineImages...)
	parentCloudProfile.Spec.MachineTypes = append(parentCloudProfile.Spec.MachineTypes, privateCloudProfile.Spec.MachineTypes...)
	parentCloudProfile.Spec.Regions = append(parentCloudProfile.Spec.Regions, privateCloudProfile.Spec.Regions...)
	parentCloudProfile.Spec.VolumeTypes = append(parentCloudProfile.Spec.VolumeTypes, privateCloudProfile.Spec.VolumeTypes...)
	if privateCloudProfile.Spec.CABundle != nil {
		mergedCABundles := fmt.Sprintf("%s%s", *parentCloudProfile.Spec.CABundle, *privateCloudProfile.Spec.CABundle)
		parentCloudProfile.Spec.CABundle = &mergedCABundles
	}
	// TODO how to merge seedSelector
	// TODO Should metadata also be merged?
}
