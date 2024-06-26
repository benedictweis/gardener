//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	authentication "github.com/gardener/gardener/pkg/apis/authentication"
	v1 "k8s.io/api/core/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*CredentialsBinding)(nil), (*authentication.CredentialsBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_CredentialsBinding_To_authentication_CredentialsBinding(a.(*CredentialsBinding), b.(*authentication.CredentialsBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.CredentialsBinding)(nil), (*CredentialsBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_CredentialsBinding_To_v1alpha1_CredentialsBinding(a.(*authentication.CredentialsBinding), b.(*CredentialsBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CredentialsBindingList)(nil), (*authentication.CredentialsBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_CredentialsBindingList_To_authentication_CredentialsBindingList(a.(*CredentialsBindingList), b.(*authentication.CredentialsBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.CredentialsBindingList)(nil), (*CredentialsBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_CredentialsBindingList_To_v1alpha1_CredentialsBindingList(a.(*authentication.CredentialsBindingList), b.(*CredentialsBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CredentialsBindingProvider)(nil), (*authentication.CredentialsBindingProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_CredentialsBindingProvider_To_authentication_CredentialsBindingProvider(a.(*CredentialsBindingProvider), b.(*authentication.CredentialsBindingProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authentication.CredentialsBindingProvider)(nil), (*CredentialsBindingProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_CredentialsBindingProvider_To_v1alpha1_CredentialsBindingProvider(a.(*authentication.CredentialsBindingProvider), b.(*CredentialsBindingProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*authentication.KubeconfigRequest)(nil), (*AdminKubeconfigRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_KubeconfigRequest_To_v1alpha1_AdminKubeconfigRequest(a.(*authentication.KubeconfigRequest), b.(*AdminKubeconfigRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*authentication.KubeconfigRequest)(nil), (*ViewerKubeconfigRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authentication_KubeconfigRequest_To_v1alpha1_ViewerKubeconfigRequest(a.(*authentication.KubeconfigRequest), b.(*ViewerKubeconfigRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*AdminKubeconfigRequest)(nil), (*authentication.KubeconfigRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_AdminKubeconfigRequest_To_authentication_KubeconfigRequest(a.(*AdminKubeconfigRequest), b.(*authentication.KubeconfigRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*ViewerKubeconfigRequest)(nil), (*authentication.KubeconfigRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ViewerKubeconfigRequest_To_authentication_KubeconfigRequest(a.(*ViewerKubeconfigRequest), b.(*authentication.KubeconfigRequest), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_CredentialsBinding_To_authentication_CredentialsBinding(in *CredentialsBinding, out *authentication.CredentialsBinding, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_CredentialsBindingProvider_To_authentication_CredentialsBindingProvider(&in.Provider, &out.Provider, s); err != nil {
		return err
	}
	out.CredentialsRef = in.CredentialsRef
	out.Quotas = *(*[]v1.ObjectReference)(unsafe.Pointer(&in.Quotas))
	return nil
}

// Convert_v1alpha1_CredentialsBinding_To_authentication_CredentialsBinding is an autogenerated conversion function.
func Convert_v1alpha1_CredentialsBinding_To_authentication_CredentialsBinding(in *CredentialsBinding, out *authentication.CredentialsBinding, s conversion.Scope) error {
	return autoConvert_v1alpha1_CredentialsBinding_To_authentication_CredentialsBinding(in, out, s)
}

func autoConvert_authentication_CredentialsBinding_To_v1alpha1_CredentialsBinding(in *authentication.CredentialsBinding, out *CredentialsBinding, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_authentication_CredentialsBindingProvider_To_v1alpha1_CredentialsBindingProvider(&in.Provider, &out.Provider, s); err != nil {
		return err
	}
	out.CredentialsRef = in.CredentialsRef
	out.Quotas = *(*[]v1.ObjectReference)(unsafe.Pointer(&in.Quotas))
	return nil
}

// Convert_authentication_CredentialsBinding_To_v1alpha1_CredentialsBinding is an autogenerated conversion function.
func Convert_authentication_CredentialsBinding_To_v1alpha1_CredentialsBinding(in *authentication.CredentialsBinding, out *CredentialsBinding, s conversion.Scope) error {
	return autoConvert_authentication_CredentialsBinding_To_v1alpha1_CredentialsBinding(in, out, s)
}

func autoConvert_v1alpha1_CredentialsBindingList_To_authentication_CredentialsBindingList(in *CredentialsBindingList, out *authentication.CredentialsBindingList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]authentication.CredentialsBinding)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_CredentialsBindingList_To_authentication_CredentialsBindingList is an autogenerated conversion function.
func Convert_v1alpha1_CredentialsBindingList_To_authentication_CredentialsBindingList(in *CredentialsBindingList, out *authentication.CredentialsBindingList, s conversion.Scope) error {
	return autoConvert_v1alpha1_CredentialsBindingList_To_authentication_CredentialsBindingList(in, out, s)
}

func autoConvert_authentication_CredentialsBindingList_To_v1alpha1_CredentialsBindingList(in *authentication.CredentialsBindingList, out *CredentialsBindingList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]CredentialsBinding)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_authentication_CredentialsBindingList_To_v1alpha1_CredentialsBindingList is an autogenerated conversion function.
func Convert_authentication_CredentialsBindingList_To_v1alpha1_CredentialsBindingList(in *authentication.CredentialsBindingList, out *CredentialsBindingList, s conversion.Scope) error {
	return autoConvert_authentication_CredentialsBindingList_To_v1alpha1_CredentialsBindingList(in, out, s)
}

func autoConvert_v1alpha1_CredentialsBindingProvider_To_authentication_CredentialsBindingProvider(in *CredentialsBindingProvider, out *authentication.CredentialsBindingProvider, s conversion.Scope) error {
	out.Type = in.Type
	return nil
}

// Convert_v1alpha1_CredentialsBindingProvider_To_authentication_CredentialsBindingProvider is an autogenerated conversion function.
func Convert_v1alpha1_CredentialsBindingProvider_To_authentication_CredentialsBindingProvider(in *CredentialsBindingProvider, out *authentication.CredentialsBindingProvider, s conversion.Scope) error {
	return autoConvert_v1alpha1_CredentialsBindingProvider_To_authentication_CredentialsBindingProvider(in, out, s)
}

func autoConvert_authentication_CredentialsBindingProvider_To_v1alpha1_CredentialsBindingProvider(in *authentication.CredentialsBindingProvider, out *CredentialsBindingProvider, s conversion.Scope) error {
	out.Type = in.Type
	return nil
}

// Convert_authentication_CredentialsBindingProvider_To_v1alpha1_CredentialsBindingProvider is an autogenerated conversion function.
func Convert_authentication_CredentialsBindingProvider_To_v1alpha1_CredentialsBindingProvider(in *authentication.CredentialsBindingProvider, out *CredentialsBindingProvider, s conversion.Scope) error {
	return autoConvert_authentication_CredentialsBindingProvider_To_v1alpha1_CredentialsBindingProvider(in, out, s)
}
