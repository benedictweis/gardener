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

package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"

	"github.com/gardener/gardener/pkg/apis/core"
	namespacedcloudprofileregistry "github.com/gardener/gardener/pkg/registry/core/namespacedcloudprofile"
)

// REST implements a RESTStorage for NamespacedCloudProfile
type REST struct {
	*genericregistry.Store
}

// NamespacedCloudProfileStorage implements the storage for NamespacedCloudProfiles.
type NamespacedCloudProfileStorage struct {
	NamespacedCloudProfile *REST
}

// NewStorage creates a new NamespacedCloudProfileStorage object.
func NewStorage(optsGetter generic.RESTOptionsGetter) NamespacedCloudProfileStorage {
	namespacedCloudProfileRest := NewREST(optsGetter)

	return NamespacedCloudProfileStorage{
		NamespacedCloudProfile: namespacedCloudProfileRest,
	}
}

// NewREST returns a RESTStorage object that will work with NamespacedCloudProfile objects.
func NewREST(optsGetter generic.RESTOptionsGetter) *REST {
	store := &genericregistry.Store{
		NewFunc:                   func() runtime.Object { return &core.NamespacedCloudProfile{} },
		NewListFunc:               func() runtime.Object { return &core.NamespacedCloudProfileList{} },
		DefaultQualifiedResource:  core.Resource("namespacedcloudprofiles"),
		SingularQualifiedResource: core.Resource("namespacedcloudprofile"),
		EnableGarbageCollection:   true,

		CreateStrategy: namespacedcloudprofileregistry.Strategy,
		UpdateStrategy: namespacedcloudprofileregistry.Strategy,
		DeleteStrategy: namespacedcloudprofileregistry.Strategy,

		TableConvertor: newTableConvertor(),
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		panic(err)
	}
	return &REST{store}
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"nsprofile", "nscpfl"}
}
