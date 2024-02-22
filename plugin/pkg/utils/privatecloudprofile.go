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

	"github.com/gardener/gardener/pkg/apis/core"
	"github.com/gardener/gardener/pkg/client/core/listers/core/internalversion"
)

// GetCloudProfile gets determine whether a given CloudProfile name is a PrivateCloudProfile or a CloudProfile and returns the appropriate object
func GetCloudProfile(name string, cloudProfileLister internalversion.CloudProfileLister, PrivateCloudProfileLister internalversion.PrivateCloudProfileLister, namespace string) (*core.CloudProfile, error) {
	cloudProfile, err := cloudProfileLister.Get(name)
	if err == nil {
		return cloudProfile, nil
	}
	// TODO how to compare errors correctly
	//if !errors.Is(err, apierrors.NewNotFound(core.Resource("cloudprofile"), name)) {
	//	return nil, err
	//}
	privateCloudProfile, perr := PrivateCloudProfileLister.PrivateCloudProfiles(namespace).Get(name)
	if perr == nil {
		return &privateCloudProfile.Status.CloudProfile, nil
	}
	//if !errors.Is(perr, apierrors.NewNotFound(core.Resource("privatecloudprofile"), name)) {
	//	return nil, err
	//}
	return nil, fmt.Errorf("could not get (private) cloud profile: %+v, %+v", err.Error(), perr.Error())
}
