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

package validation_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	gomegatypes "github.com/onsi/gomega/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"

	"github.com/gardener/gardener/pkg/apis/core"
	. "github.com/gardener/gardener/pkg/apis/core/validation"
)

var _ = Describe("NamespacedCloudProfile Validation Tests ", func() {
	Describe("#ValidateNamespacedCloudProfile", func() {
		It("should forbid empty NamespacedCloudProfile resources", func() {
			namespacedCloudProfile := &core.NamespacedCloudProfile{
				ObjectMeta: metav1.ObjectMeta{},
				Spec:       core.NamespacedCloudProfileSpec{},
			}

			errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

			Expect(errorList).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("metadata.name"),
				})), PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("metadata.namespace"),
				})), PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeInvalid),
					"Field": Equal("spec.parent.kind"),
				})), PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("spec.parent.name"),
				}))))
		})

		It("should allow NamespacedCloudProfile resource with only name and parent field", func() {
			namespacedCloudProfile := &core.NamespacedCloudProfile{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "profile",
					Namespace: "default",
				},
				Spec: core.NamespacedCloudProfileSpec{
					Parent: core.CloudProfileReference{
						Kind: "CloudProfile",
						Name: "profile-parent",
					},
				},
			}

			errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

			Expect(errorList).To(BeEmpty())
		})

		Context("tests for unknown cloud profiles", func() {
			var (
				namespacedCloudProfile *core.NamespacedCloudProfile

				duplicatedKubernetes = core.KubernetesSettings{
					Versions: []core.ExpirableVersion{{Version: "1.11.4"}, {Version: "1.11.4", Classification: &previewClassification}},
				}
				duplicatedRegions = []core.Region{
					{
						Name: regionName,
						Zones: []core.AvailabilityZone{
							{Name: zoneName},
						},
					},
					{
						Name: regionName,
						Zones: []core.AvailabilityZone{
							{Name: zoneName},
						},
					},
				}
				duplicatedZones = []core.Region{
					{
						Name: regionName,
						Zones: []core.AvailabilityZone{
							{Name: zoneName},
							{Name: zoneName},
						},
					},
				}
			)

			BeforeEach(func() {
				namespacedCloudProfile = &core.NamespacedCloudProfile{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "profile",
						Namespace: "default",
					},
					Spec: core.NamespacedCloudProfileSpec{
						Parent: core.CloudProfileReference{
							Kind: "CloudProfile",
							Name: "unknown",
						},
						SeedSelector: &core.SeedSelector{
							LabelSelector: metav1.LabelSelector{
								MatchLabels: map[string]string{"foo": "bar"},
							},
						},
						Kubernetes: &core.KubernetesSettings{
							Versions: []core.ExpirableVersion{{
								Version: "1.11.4",
							}},
						},
						MachineImages: []core.MachineImage{
							{
								Name: "some-machineimage",
								Versions: []core.MachineImageVersion{
									{
										ExpirableVersion: core.ExpirableVersion{
											Version: "1.2.3",
										},
										CRI: []core.CRI{{Name: "containerd"}},
									},
								},
							},
						},
						Regions: []core.Region{
							{
								Name: regionName,
								Zones: []core.AvailabilityZone{
									{Name: zoneName},
								},
							},
						},
						MachineTypes: machineTypesConstraint,
						VolumeTypes:  volumeTypesConstraint,
					},
				}
			})

			It("should not return any errors", func() {
				errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

				Expect(errorList).To(BeEmpty())
			})

			DescribeTable("namespacedCloudProfile metadata",
				func(objectMeta metav1.ObjectMeta, matcher gomegatypes.GomegaMatcher) {
					namespacedCloudProfile.ObjectMeta = objectMeta
					namespacedCloudProfile.ObjectMeta.Namespace = "default"

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(matcher)
				},

				Entry("should forbid namespacedCloudProfile with empty metadata",
					metav1.ObjectMeta{},
					ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("metadata.name"),
					}))),
				),
				Entry("should forbid namespacedCloudProfile with empty name",
					metav1.ObjectMeta{Name: ""},
					ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("metadata.name"),
					}))),
				),
				Entry("should allow namespacedCloudProfile with '.' in the name",
					metav1.ObjectMeta{Name: "profile.test"},
					BeEmpty(),
				),
				Entry("should forbid namespacedCloudProfile with '_' in the name (not a DNS-1123 subdomain)",
					metav1.ObjectMeta{Name: "profile_test"},
					ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("metadata.name"),
					}))),
				),
			)

			It("should forbid not specifying a parent", func() {
				namespacedCloudProfile.Spec.Parent = core.CloudProfileReference{Kind: "", Name: ""}

				errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

				Expect(errorList).To(ConsistOf(
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("spec.parent.name"),
					})),
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.parent.kind"),
					}))))
			})

			It("should forbid not specifying an invalid parent kind", func() {
				namespacedCloudProfile.Spec.Parent = core.CloudProfileReference{Kind: "SomeOtherCloudPofile", Name: "my-profile"}

				errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

				Expect(errorList).To(ConsistOf(
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.parent.kind"),
					}))))
			})

			It("should forbid ca bundles with unsupported format", func() {
				namespacedCloudProfile.Spec.CABundle = ptr.To("unsupported")

				errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeInvalid),
					"Field": Equal("spec.caBundle"),
				}))))
			})

			Context("kubernetes version constraints", func() {
				It("should enforce that at least one version has been defined", func() {
					namespacedCloudProfile.Spec.Kubernetes.Versions = []core.ExpirableVersion{}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("spec.kubernetes.versions"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.kubernetes.versions"),
					}))))
				})

				It("should forbid versions of a not allowed pattern", func() {
					namespacedCloudProfile.Spec.Kubernetes.Versions = []core.ExpirableVersion{{Version: "1.11"}}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.kubernetes.versions[0]"),
					}))))
				})

				It("should forbid expiration date on latest kubernetes version", func() {
					expirationDate := &metav1.Time{Time: time.Now().AddDate(0, 0, 1)}
					namespacedCloudProfile.Spec.Kubernetes.Versions = []core.ExpirableVersion{
						{
							Version:        "1.1.0",
							Classification: &supportedClassification,
						},
						{
							Version:        "1.2.0",
							Classification: &deprecatedClassification,
							ExpirationDate: expirationDate,
						},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.kubernetes.versions[].expirationDate"),
					}))))
				})

				It("should forbid duplicated kubernetes versions", func() {
					namespacedCloudProfile.Spec.Kubernetes = &duplicatedKubernetes

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(
						PointTo(MatchFields(IgnoreExtras, Fields{
							"Type":  Equal(field.ErrorTypeDuplicate),
							"Field": Equal(fmt.Sprintf("spec.kubernetes.versions[%d].version", len(duplicatedKubernetes.Versions)-1)),
						}))))
				})

				It("should forbid invalid classification for kubernetes versions", func() {
					classification := core.VersionClassification("dummy")
					namespacedCloudProfile.Spec.Kubernetes.Versions = []core.ExpirableVersion{
						{
							Version:        "1.1.0",
							Classification: &classification,
						},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":     Equal(field.ErrorTypeNotSupported),
						"Field":    Equal("spec.kubernetes.versions[0].classification"),
						"BadValue": Equal(classification),
					}))))
				})

				It("only allow one supported version per minor version", func() {
					namespacedCloudProfile.Spec.Kubernetes.Versions = []core.ExpirableVersion{
						{
							Version:        "1.1.0",
							Classification: &supportedClassification,
						},
						{
							Version:        "1.1.1",
							Classification: &supportedClassification,
						},
					}
					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeForbidden),
						"Field": Equal("spec.kubernetes.versions[1]"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeForbidden),
						"Field": Equal("spec.kubernetes.versions[0]"),
					}))))
				})
			})

			Context("machine image validation", func() {
				It("should forbid an empty list of machine images", func() {
					namespacedCloudProfile.Spec.MachineImages = []core.MachineImage{}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("spec.machineImages"),
					}))))
				})

				It("should forbid duplicate names in list of machine images", func() {
					namespacedCloudProfile.Spec.MachineImages = []core.MachineImage{
						{
							Name: "some-machineimage",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version:        "3.4.6",
										Classification: &supportedClassification,
									},
									CRI: []core.CRI{{Name: "containerd"}},
								},
							},
						},
						{
							Name: "some-machineimage",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version:        "3.4.5",
										Classification: &previewClassification,
									},
									CRI: []core.CRI{{Name: "containerd"}},
								},
							},
						},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeDuplicate),
						"Field": Equal("spec.machineImages[1]"),
					}))))
				})

				It("should forbid machine images with no version", func() {
					namespacedCloudProfile.Spec.MachineImages = []core.MachineImage{
						{Name: "some-machineimage"},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("spec.machineImages[0].versions"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.machineImages"),
					}))))
				})

				It("should forbid machine images with an invalid machine image update strategy", func() {
					updateStrategy := core.MachineImageUpdateStrategy("dummy")
					namespacedCloudProfile.Spec.MachineImages = []core.MachineImage{
						{
							Name: "some-machineimage",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version:        "3.4.6",
										Classification: &supportedClassification,
									},
									CRI: []core.CRI{{Name: "containerd"}},
								},
							},
							UpdateStrategy: &updateStrategy,
						},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeNotSupported),
						"Field": Equal("spec.machineImages[0].updateStrategy"),
					}))))
				})

				It("should allow machine images with a valid machine image update strategy", func() {
					updateStrategy := core.UpdateStrategyMinor
					namespacedCloudProfile.Spec.MachineImages = []core.MachineImage{
						{
							Name: "some-machineimage",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version:        "3.4.6",
										Classification: &supportedClassification,
									},
									CRI: []core.CRI{{Name: "containerd"}},
								},
							},
							UpdateStrategy: &updateStrategy,
						},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(BeEmpty())
				})

				It("should forbid nonSemVer machine image versions", func() {
					namespacedCloudProfile.Spec.MachineImages = []core.MachineImage{
						{
							Name: "some-machineimage",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version:        "0.1.2",
										Classification: &supportedClassification,
									},
									CRI: []core.CRI{{Name: "containerd"}},
								},
							},
						},
						{
							Name: "xy",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version:        "a.b.c",
										Classification: &supportedClassification,
									},
									CRI: []core.CRI{{Name: "containerd"}},
								},
							},
						},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.machineImages"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.machineImages[1].versions[0].version"),
					}))))
				})

				It("should allow expiration date on latest machine image version", func() {
					expirationDate := &metav1.Time{Time: time.Now().AddDate(0, 0, 1)}
					namespacedCloudProfile.Spec.MachineImages = []core.MachineImage{
						{
							Name: "some-machineimage",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version:        "0.1.2",
										ExpirationDate: expirationDate,
										Classification: &previewClassification,
									},
									CRI: []core.CRI{{Name: "containerd"}},
								},
								{
									ExpirableVersion: core.ExpirableVersion{
										Version:        "0.1.1",
										Classification: &supportedClassification,
									},
									CRI: []core.CRI{{Name: "containerd"}},
								},
							},
						},
						{
							Name: "xy",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version:        "0.1.1",
										ExpirationDate: expirationDate,
										Classification: &supportedClassification,
									},
									CRI: []core.CRI{{Name: "containerd"}},
								},
							},
						},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)
					Expect(errorList).To(BeEmpty())
				})

				It("should forbid invalid classification for machine image versions", func() {
					classification := core.VersionClassification("dummy")
					namespacedCloudProfile.Spec.MachineImages = []core.MachineImage{
						{
							Name: "some-machineimage",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version:        "0.1.2",
										Classification: &classification,
									},
									CRI: []core.CRI{{Name: "containerd"}},
								},
							},
						},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)
					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":     Equal(field.ErrorTypeNotSupported),
						"Field":    Equal("spec.machineImages[0].versions[0].classification"),
						"BadValue": Equal(classification),
					}))))
				})

				It("should allow valid CPU architecture for machine image versions", func() {
					namespacedCloudProfile.Spec.MachineImages = []core.MachineImage{
						{
							Name: "some-machineimage",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version: "0.1.2",
									},
									CRI:           []core.CRI{{Name: "containerd"}},
									Architectures: []string{"amd64", "arm64"},
								},
								{
									ExpirableVersion: core.ExpirableVersion{
										Version: "0.1.3",
									},
									CRI:           []core.CRI{{Name: "containerd"}},
									Architectures: []string{"amd64"},
								},
								{
									ExpirableVersion: core.ExpirableVersion{
										Version: "0.1.4",
									},
									CRI:           []core.CRI{{Name: "containerd"}},
									Architectures: []string{"arm64"},
								},
							},
						},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)
					Expect(errorList).To(BeEmpty())
				})

				It("should forbid invalid CPU architecture for machine image versions", func() {
					namespacedCloudProfile.Spec.MachineImages = []core.MachineImage{
						{
							Name: "some-machineimage",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version: "0.1.2",
									},
									CRI:           []core.CRI{{Name: "containerd"}},
									Architectures: []string{"foo"},
								},
							},
						},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)
					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeNotSupported),
						"Field": Equal("spec.machineImages[0].versions[0].architecture"),
					}))))
				})

				It("should allow valid kubeletVersionConstraint for machine image versions", func() {
					namespacedCloudProfile.Spec.MachineImages = []core.MachineImage{
						{
							Name: "some-machineimage",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version: "0.1.2",
									},
									CRI:                      []core.CRI{{Name: "containerd"}},
									KubeletVersionConstraint: ptr.To("< 1.26"),
								},
								{
									ExpirableVersion: core.ExpirableVersion{
										Version: "0.1.3",
									},
									CRI:                      []core.CRI{{Name: "containerd"}},
									KubeletVersionConstraint: ptr.To(">= 1.26"),
								},
							},
						},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)
					Expect(errorList).To(BeEmpty())
				})

				It("should forbid invalid kubeletVersionConstraint for machine image versions", func() {
					namespacedCloudProfile.Spec.MachineImages = []core.MachineImage{
						{
							Name: "some-machineimage",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version: "0.1.2",
									},
									CRI:                      []core.CRI{{Name: "containerd"}},
									KubeletVersionConstraint: ptr.To(""),
								},
								{
									ExpirableVersion: core.ExpirableVersion{
										Version: "0.1.3",
									},
									CRI:                      []core.CRI{{Name: "containerd"}},
									KubeletVersionConstraint: ptr.To("invalid-version"),
								},
							},
						},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)
					Expect(errorList).To(ConsistOf(
						PointTo(MatchFields(IgnoreExtras, Fields{
							"Type":  Equal(field.ErrorTypeInvalid),
							"Field": Equal("spec.machineImages[0].versions[0].kubeletVersionConstraint"),
						})),
						PointTo(MatchFields(IgnoreExtras, Fields{
							"Type":  Equal(field.ErrorTypeInvalid),
							"Field": Equal("spec.machineImages[0].versions[1].kubeletVersionConstraint"),
						})),
					))
				})
			})

			It("should forbid if no cri is present", func() {
				namespacedCloudProfile.Spec.MachineImages[0].Versions[0].CRI = nil

				errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("spec.machineImages[0].versions[0].cri"),
				}))))
			})

			It("should forbid non-supported container runtime interface names", func() {
				namespacedCloudProfile.Spec.MachineImages = []core.MachineImage{
					{
						Name: "invalid-cri-name",
						Versions: []core.MachineImageVersion{
							{
								ExpirableVersion: core.ExpirableVersion{
									Version: "0.1.2",
								},
								CRI: []core.CRI{
									{
										Name: "invalid-cri-name",
									},
									{
										Name: "docker",
									},
								},
							},
						},
					},
					{
						Name: "valid-cri-name",
						Versions: []core.MachineImageVersion{
							{
								ExpirableVersion: core.ExpirableVersion{
									Version: "0.1.2",
								},
								CRI: []core.CRI{
									{
										Name: core.CRINameContainerD,
									},
								},
							},
						},
					},
				}

				errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeNotSupported),
					"Field":  Equal("spec.machineImages[0].versions[0].cri[0].name"),
					"Detail": Equal("supported values: \"containerd\""),
				})), PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeNotSupported),
					"Field":  Equal("spec.machineImages[0].versions[0].cri[1].name"),
					"Detail": Equal("supported values: \"containerd\""),
				})),
				))
			})

			It("should forbid duplicated container runtime interface names", func() {
				namespacedCloudProfile.Spec.MachineImages[0].Versions[0].CRI = []core.CRI{
					{
						Name: core.CRINameContainerD,
					},
					{
						Name: core.CRINameContainerD,
					},
				}

				errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeDuplicate),
					"Field": Equal("spec.machineImages[0].versions[0].cri[1]"),
				}))))
			})

			It("should forbid duplicated container runtime names", func() {
				namespacedCloudProfile.Spec.MachineImages[0].Versions[0].CRI = []core.CRI{
					{
						Name: core.CRINameContainerD,
						ContainerRuntimes: []core.ContainerRuntime{
							{
								Type: "cr1",
							},
							{
								Type: "cr1",
							},
						},
					},
				}

				errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeDuplicate),
					"Field": Equal("spec.machineImages[0].versions[0].cri[0].containerRuntimes[1].type"),
				}))))
			})

			Context("machine types validation", func() {
				It("should enforce that at least one machine type has been defined", func() {
					namespacedCloudProfile.Spec.MachineTypes = []core.MachineType{}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("spec.machineTypes"),
					}))))
				})

				It("should enforce uniqueness of machine type names", func() {
					namespacedCloudProfile.Spec.MachineTypes = []core.MachineType{
						namespacedCloudProfile.Spec.MachineTypes[0],
						namespacedCloudProfile.Spec.MachineTypes[0],
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeDuplicate),
						"Field": Equal("spec.machineTypes[1].name"),
					}))))
				})

				It("should forbid machine types with unsupported property values", func() {
					namespacedCloudProfile.Spec.MachineTypes = invalidMachineTypes

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("spec.machineTypes[0].name"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.machineTypes[0].cpu"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.machineTypes[0].gpu"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.machineTypes[0].memory"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.machineTypes[0].storage.minSize"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.machineTypes[1].storage.size"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.machineTypes[2].storage"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.machineTypes[3].storage"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeNotSupported),
						"Field": Equal("spec.machineTypes[3].architecture"),
					})),
					))
				})

				It("should allow machine types with valid values", func() {
					namespacedCloudProfile.Spec.MachineTypes = machineTypesConstraint

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)
					Expect(errorList).To(BeEmpty())
				})
			})

			Context("regions validation", func() {
				It("should forbid regions with unsupported name values", func() {
					namespacedCloudProfile.Spec.Regions = []core.Region{
						{
							Name:  "",
							Zones: []core.AvailabilityZone{{Name: ""}},
						},
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("spec.regions[0].name"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("spec.regions[0].zones[0].name"),
					}))))
				})

				It("should forbid duplicated region names", func() {
					namespacedCloudProfile.Spec.Regions = duplicatedRegions

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(
						PointTo(MatchFields(IgnoreExtras, Fields{
							"Type":  Equal(field.ErrorTypeDuplicate),
							"Field": Equal(fmt.Sprintf("spec.regions[%d].name", len(duplicatedRegions)-1)),
						}))))
				})

				It("should forbid duplicated zone names", func() {
					namespacedCloudProfile.Spec.Regions = duplicatedZones

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(
						PointTo(MatchFields(IgnoreExtras, Fields{
							"Type":  Equal(field.ErrorTypeDuplicate),
							"Field": Equal(fmt.Sprintf("spec.regions[0].zones[%d].name", len(duplicatedZones[0].Zones)-1)),
						}))))
				})

				It("should forbid invalid label specifications", func() {
					namespacedCloudProfile.Spec.Regions[0].Labels = map[string]string{
						"this-is-not-allowed": "?*&!@",
						"toolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolong": "toolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolongtoolong",
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(
						PointTo(MatchFields(IgnoreExtras, Fields{
							"Type":  Equal(field.ErrorTypeInvalid),
							"Field": Equal("spec.regions[0].labels"),
						})),
						PointTo(MatchFields(IgnoreExtras, Fields{
							"Type":  Equal(field.ErrorTypeInvalid),
							"Field": Equal("spec.regions[0].labels"),
						})),
						PointTo(MatchFields(IgnoreExtras, Fields{
							"Type":  Equal(field.ErrorTypeInvalid),
							"Field": Equal("spec.regions[0].labels"),
						})),
					))
				})
			})

			Context("volume types validation", func() {
				It("should enforce uniqueness of volume type names", func() {
					namespacedCloudProfile.Spec.VolumeTypes = []core.VolumeType{
						namespacedCloudProfile.Spec.VolumeTypes[0],
						namespacedCloudProfile.Spec.VolumeTypes[0],
					}

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeDuplicate),
						"Field": Equal("spec.volumeTypes[1].name"),
					}))))
				})

				It("should forbid volume types with unsupported property values", func() {
					namespacedCloudProfile.Spec.VolumeTypes = invalidVolumeTypes

					errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

					Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("spec.volumeTypes[0].name"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeRequired),
						"Field": Equal("spec.volumeTypes[0].class"),
					})), PointTo(MatchFields(IgnoreExtras, Fields{
						"Type":  Equal(field.ErrorTypeInvalid),
						"Field": Equal("spec.volumeTypes[0].minSize"),
					})),
					))
				})
			})

			It("should forbid unsupported seed selectors", func() {
				namespacedCloudProfile.Spec.SeedSelector.MatchLabels["foo"] = "no/slash/allowed"

				errorList := ValidateNamespacedCloudProfile(namespacedCloudProfile)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeInvalid),
					"Field": Equal("spec.seedSelector.matchLabels"),
				}))))
			})
		})
	})

	var (
		cloudProfileNew *core.NamespacedCloudProfile
		cloudProfileOld *core.NamespacedCloudProfile
		dateInThePast   = &metav1.Time{Time: time.Now().AddDate(-5, 0, 0)}
	)

	Describe("#ValidateNamespacedCloudProfileSpecUpdate", func() {
		BeforeEach(func() {
			cloudProfileNew = &core.NamespacedCloudProfile{
				ObjectMeta: metav1.ObjectMeta{
					ResourceVersion: "dummy",
					Name:            "dummy",
					Namespace:       "dummy",
				},
				Spec: core.NamespacedCloudProfileSpec{
					Parent: core.CloudProfileReference{
						Kind: "CloudProfile",
						Name: "aws-profile",
					},
					MachineImages: []core.MachineImage{
						{
							Name: "some-machineimage",
							Versions: []core.MachineImageVersion{
								{
									ExpirableVersion: core.ExpirableVersion{
										Version:        "1.2.3",
										Classification: &supportedClassification,
									},
									CRI: []core.CRI{{Name: "containerd"}},
								},
							},
						},
					},
					Kubernetes: &core.KubernetesSettings{
						Versions: []core.ExpirableVersion{
							{
								Version:        "1.17.2",
								Classification: &deprecatedClassification,
							},
						},
					},
					Regions: []core.Region{
						{
							Name: regionName,
							Zones: []core.AvailabilityZone{
								{Name: zoneName},
							},
						},
					},
					MachineTypes: machineTypesConstraint,
					VolumeTypes:  volumeTypesConstraint,
				},
			}
			cloudProfileOld = cloudProfileNew.DeepCopy()
		})

		Context("Removed Kubernetes versions", func() {
			It("deleting version - should not return any errors", func() {
				versions := []core.ExpirableVersion{
					{Version: "1.17.2", Classification: &deprecatedClassification},
					{Version: "1.17.1", Classification: &deprecatedClassification, ExpirationDate: dateInThePast},
					{Version: "1.17.0", Classification: &deprecatedClassification, ExpirationDate: dateInThePast},
				}
				cloudProfileNew.Spec.Kubernetes.Versions = versions[0:1]
				cloudProfileOld.Spec.Kubernetes.Versions = versions
				errorList := ValidateNamespacedCloudProfileUpdate(cloudProfileNew, cloudProfileOld)

				Expect(errorList).To(BeEmpty())
			})
		})

		Context("Removed MachineImage versions", func() {
			It("deleting version - should not return any errors", func() {
				versions := []core.MachineImageVersion{
					{
						ExpirableVersion: core.ExpirableVersion{
							Version:        "2135.6.2",
							Classification: &deprecatedClassification,
						},
						CRI: []core.CRI{{Name: "containerd"}},
					},
					{
						ExpirableVersion: core.ExpirableVersion{
							Version:        "2135.6.1",
							Classification: &deprecatedClassification,
							ExpirationDate: dateInThePast,
						},
						CRI: []core.CRI{{Name: "containerd"}},
					},
					{
						ExpirableVersion: core.ExpirableVersion{
							Version:        "2135.6.0",
							Classification: &deprecatedClassification,
							ExpirationDate: dateInThePast,
						},
						CRI: []core.CRI{{Name: "containerd"}},
					},
				}
				cloudProfileNew.Spec.MachineImages[0].Versions = versions[0:1]
				cloudProfileOld.Spec.MachineImages[0].Versions = versions
				errorList := ValidateNamespacedCloudProfileUpdate(cloudProfileNew, cloudProfileOld)

				Expect(errorList).To(BeEmpty())
			})
		})
	})
})
