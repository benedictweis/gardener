// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package initializer

import (
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	quotav1 "k8s.io/apiserver/pkg/quota/v1"
	"k8s.io/client-go/dynamic"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	authenticationclientset "github.com/gardener/gardener/pkg/client/authentication/clientset/versioned"
	authenticationinformers "github.com/gardener/gardener/pkg/client/authentication/informers/externalversions"
	gardencoreclientset "github.com/gardener/gardener/pkg/client/core/clientset/versioned"
	gardencoreinformers "github.com/gardener/gardener/pkg/client/core/informers/externalversions"
	seedmanagementclientset "github.com/gardener/gardener/pkg/client/seedmanagement/clientset/versioned"
	seedmanagementinformers "github.com/gardener/gardener/pkg/client/seedmanagement/informers/externalversions"
	settingsinformers "github.com/gardener/gardener/pkg/client/settings/informers/externalversions"
)

// New constructs new instance of PluginInitializer
func New(
	coreInformers gardencoreinformers.SharedInformerFactory,
	coreClient gardencoreclientset.Interface,
	seedManagementInformers seedmanagementinformers.SharedInformerFactory,
	seedManagementClient seedmanagementclientset.Interface,
	settingsInformers settingsinformers.SharedInformerFactory,
	authenticationInformers authenticationinformers.SharedInformerFactory,
	authenticationClient authenticationclientset.Interface,
	kubeInformers kubeinformers.SharedInformerFactory,
	kubeClient kubernetes.Interface,
	dynamicClient dynamic.Interface,
	authz authorizer.Authorizer,
	quotaConfiguration quotav1.Configuration,
) admission.PluginInitializer {
	return pluginInitializer{
		coreInformers: coreInformers,
		coreClient:    coreClient,

		seedManagementInformers: seedManagementInformers,
		seedManagementClient:    seedManagementClient,

		settingsInformers: settingsInformers,

		authenticationInformers: authenticationInformers,
		authenticationClient:    authenticationClient,

		kubeInformers: kubeInformers,
		kubeClient:    kubeClient,

		dynamicClient: dynamicClient,

		authorizer: authz,

		quotaConfiguration: quotaConfiguration,
	}
}

// Initialize checks the initialization interfaces implemented by each plugin
// and provide the appropriate initialization data
func (i pluginInitializer) Initialize(plugin admission.Interface) {
	if wants, ok := plugin.(WantsCoreInformerFactory); ok {
		wants.SetCoreInformerFactory(i.coreInformers)
	}
	if wants, ok := plugin.(WantsCoreClientSet); ok {
		wants.SetCoreClientSet(i.coreClient)
	}

	if wants, ok := plugin.(WantsSeedManagementInformerFactory); ok {
		wants.SetSeedManagementInformerFactory(i.seedManagementInformers)
	}
	if wants, ok := plugin.(WantsSeedManagementClientSet); ok {
		wants.SetSeedManagementClientSet(i.seedManagementClient)
	}

	if wants, ok := plugin.(WantsAuthenticationInformerFactory); ok {
		wants.SetAuthenticationInformerFactory(i.authenticationInformers)
	}
	if wants, ok := plugin.(WantsAuthenticationClientSet); ok {
		wants.SetAuthenticationClientSet(i.authenticationClient)
	}

	if wants, ok := plugin.(WantsSettingsInformerFactory); ok {
		wants.SetSettingsInformerFactory(i.settingsInformers)
	}

	if wants, ok := plugin.(WantsKubeInformerFactory); ok {
		wants.SetKubeInformerFactory(i.kubeInformers)
	}
	if wants, ok := plugin.(WantsKubeClientset); ok {
		wants.SetKubeClientset(i.kubeClient)
	}

	if wants, ok := plugin.(WantsDynamicClient); ok {
		wants.SetDynamicClient(i.dynamicClient)
	}

	if wants, ok := plugin.(WantsAuthorizer); ok {
		wants.SetAuthorizer(i.authorizer)
	}

	if wants, ok := plugin.(WantsQuotaConfiguration); ok {
		wants.SetQuotaConfiguration(i.quotaConfiguration)
	}
}
