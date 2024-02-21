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

package validator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/admission"

	"github.com/gardener/gardener/pkg/apis/core"
	admissioninitializer "github.com/gardener/gardener/pkg/apiserver/admission/initializer"
	gardencoreinformers "github.com/gardener/gardener/pkg/client/core/informers/internalversion"
	gardencorelisters "github.com/gardener/gardener/pkg/client/core/listers/core/internalversion"
	plugin "github.com/gardener/gardener/plugin/pkg"
)

// Register registers a plugin.
func Register(plugins *admission.Plugins) {
	plugins.Register(plugin.PluginNamePrivateCloudProfileValidator, func(config io.Reader) (admission.Interface, error) {
		return New()
	})
}

// ValidatePrivateCloudProfile contains listers and admission handler.
type ValidatePrivateCloudProfile struct {
	*admission.Handler
	cloudProfileLister gardencorelisters.CloudProfileLister
	readyFunc          admission.ReadyFunc
}

var (
	_          = admissioninitializer.WantsInternalCoreInformerFactory(&ValidatePrivateCloudProfile{})
	readyFuncs []admission.ReadyFunc
)

// New creates a new ValidatePrivateCloudProfile admission plugin.
func New() (*ValidatePrivateCloudProfile, error) {
	return &ValidatePrivateCloudProfile{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}, nil
}

// AssignReadyFunc assigns the ready function to the admission handler.
func (v *ValidatePrivateCloudProfile) AssignReadyFunc(f admission.ReadyFunc) {
	v.readyFunc = f
	v.SetReadyFunc(f)
}

// SetInternalCoreInformerFactory gets Lister from SharedInformerFactory.
func (v *ValidatePrivateCloudProfile) SetInternalCoreInformerFactory(f gardencoreinformers.SharedInformerFactory) {
	cloudProfileInformer := f.Core().InternalVersion().CloudProfiles()
	v.cloudProfileLister = cloudProfileInformer.Lister()

	readyFuncs = append(readyFuncs, cloudProfileInformer.Informer().HasSynced)
}

// ValidateInitialization checks whether the plugin was correctly initialized.
func (v *ValidatePrivateCloudProfile) ValidateInitialization() error {
	if v.cloudProfileLister == nil {
		return errors.New("missing cloudProfile lister")
	}
	return nil
}

var _ admission.ValidationInterface = &ValidatePrivateCloudProfile{}

// Validate validates the PrivateCloudProfile
func (v *ValidatePrivateCloudProfile) Validate(_ context.Context, a admission.Attributes, _ admission.ObjectInterfaces) error {
	// Wait until the caches have been synced
	if v.readyFunc == nil {
		v.AssignReadyFunc(func() bool {
			for _, readyFunc := range readyFuncs {
				if !readyFunc() {
					return false
				}
			}
			return true
		})
	}
	if !v.WaitForReady() {
		return admission.NewForbidden(a, errors.New("not yet ready to handle request"))
	}

	if a.GetKind().GroupKind() != core.Kind("PrivateCloudProfile") {
		return nil
	}

	if a.GetSubresource() != "" {
		return nil
	}

	var oldPrivateCloudProfile = &core.PrivateCloudProfile{}

	privateCloudProfile, convertIsSuccessful := a.GetObject().(*core.PrivateCloudProfile)
	if !convertIsSuccessful {
		return apierrors.NewInternalError(errors.New("could not convert object to PrivateCloudProfile"))
	}

	// Exit early if shoot spec hasn't changed
	if a.GetOperation() == admission.Update {
		old, ok := a.GetOldObject().(*core.PrivateCloudProfile)
		if !ok {
			return apierrors.NewInternalError(errors.New("could not convert old resource into PrivateCloudProfile object"))
		}
		oldPrivateCloudProfile = old

		// do not ignore metadata updates to detect and prevent removal of the gardener finalizer or unwanted changes to annotations
		if reflect.DeepEqual(privateCloudProfile.Spec, oldPrivateCloudProfile.Spec) && reflect.DeepEqual(privateCloudProfile.ObjectMeta, oldPrivateCloudProfile.ObjectMeta) {
			return nil
		}
	}

	parentCloudProfileName := privateCloudProfile.Spec.Parent
	parentCloudProfile, err := v.cloudProfileLister.Get(parentCloudProfileName)
	if err != nil {
		return apierrors.NewBadRequest("parent CloudProfile could not be found")
	}

	validationContext := &validationContext{
		parentCloudProfile:     parentCloudProfile,
		privateCloudProfile:    privateCloudProfile,
		oldPrivateCloudProfile: oldPrivateCloudProfile,
	}

	if err := validationContext.validateMachineTypes(a); err != nil {
		return err
	}

	return nil
}

type validationContext struct {
	parentCloudProfile     *core.CloudProfile
	privateCloudProfile    *core.PrivateCloudProfile
	oldPrivateCloudProfile *core.PrivateCloudProfile
}

func (c *validationContext) validateMachineTypes(a admission.Attributes) error {
	if c.privateCloudProfile.Spec.MachineTypes == nil || c.parentCloudProfile.Spec.MachineTypes == nil {
		return nil
	}

	// TODO this does not feel very clean
	for _, machineType := range c.privateCloudProfile.Spec.MachineTypes {
	parentMachineTypesLoop:
		for _, parentMachineType := range c.parentCloudProfile.Spec.MachineTypes {
			if parentMachineType.Name != machineType.Name {
				continue
			}
			if a.GetOperation() == admission.Update {
				for _, oldMachineType := range c.oldPrivateCloudProfile.Spec.MachineTypes {
					if oldMachineType.Name == machineType.Name {
						continue parentMachineTypesLoop
					}
				}
			}
			return apierrors.NewBadRequest(fmt.Sprintf("PrivateCloudProfile attempts to rewrite MachineType of parent CloudProfile with machineType: %+v", machineType))
		}
	}

	return nil
}
