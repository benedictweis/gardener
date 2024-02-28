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

package utils

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/gardener/gardener/pkg/apis/core"
	"github.com/gardener/gardener/pkg/client/core/listers/core/internalversion"
)

// GetCloudProfile gets determine whether a given CloudProfile name is a NamespacedCloudProfile or a CloudProfile and returns the appropriate object
func GetCloudProfile(name string, cloudProfileLister internalversion.CloudProfileLister, NamespacedCloudProfileLister internalversion.NamespacedCloudProfileLister, namespace string) (*core.CloudProfile, error) {
	cloudProfile, err := cloudProfileLister.Get(name)
	if err == nil {
		return cloudProfile, nil
	}
	if !apierrors.IsNotFound(err) {
		return nil, err
	}
	namespacedCloudProfile, perr := NamespacedCloudProfileLister.NamespacedCloudProfiles(namespace).Get(name)
	if perr == nil {
		return &namespacedCloudProfile.Status.CloudProfile, nil
	}
	if !apierrors.IsNotFound(perr) {
		return nil, err
	}
	return nil, fmt.Errorf("could not get (private) cloud profile: %+v, %+v", err.Error(), perr.Error())
}
