// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/gardener/gardener/pkg/client/authentication/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// CredentialsBindings returns a CredentialsBindingInformer.
	CredentialsBindings() CredentialsBindingInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// CredentialsBindings returns a CredentialsBindingInformer.
func (v *version) CredentialsBindings() CredentialsBindingInformer {
	return &credentialsBindingInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
