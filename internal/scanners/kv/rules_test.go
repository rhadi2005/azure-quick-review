// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package kv

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/cmendible/azqr/internal/scanners"
)

func TestKeyVaultScanner_Rules(t *testing.T) {
	type fields struct {
		rule                string
		target              interface{}
		scanContext         *scanners.ScanContext
		diagnosticsSettings scanners.DiagnosticsSettings
	}
	type want struct {
		broken bool
		result string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "KeyVaultScanner DiagnosticSettings",
			fields: fields{
				rule: "DiagnosticSettings",
				target: &armkeyvault.Vault{
					ID: to.StringPtr("test"),
				},
				scanContext: &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{
					HasDiagnosticsFunc: func(resourceId string) (bool, error) {
						return true, nil
					},
				},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "KeyVaultScanner SLA",
			fields: fields{
				rule:                "SLA",
				target:              &armkeyvault.Vault{},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "99.99%",
			},
		},
		{
			name: "KeyVaultScanner Private Endpoint",
			fields: fields{
				rule: "Private",
				target: &armkeyvault.Vault{
					Properties: &armkeyvault.VaultProperties{
						PrivateEndpointConnections: []*armkeyvault.PrivateEndpointConnectionItem{
							{
								ID: to.StringPtr("test"),
							},
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "KeyVaultScanner SKU",
			fields: fields{
				rule: "SKU",
				target: &armkeyvault.Vault{
					Properties: &armkeyvault.VaultProperties{
						SKU: &armkeyvault.SKU{
							Name: getSKUName(),
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "standard",
			},
		},
		{
			name: "KeyVaultScanner CAF",
			fields: fields{
				rule: "CAF",
				target: &armkeyvault.Vault{
					Name: to.StringPtr("kv-test"),
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "KeyVaultScanner soft delete enabled",
			fields: fields{
				rule: "kv-008",
				target: &armkeyvault.Vault{
					Properties: &armkeyvault.VaultProperties{
						EnableSoftDelete: to.BoolPtr(true),
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "KeyVaultScanner purge protection enabled",
			fields: fields{
				rule: "kv-009",
				target: &armkeyvault.Vault{
					Properties: &armkeyvault.VaultProperties{
						EnablePurgeProtection: to.BoolPtr(true),
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KeyVaultScanner{
				diagnosticsSettings: tt.fields.diagnosticsSettings,
			}
			rules := s.GetRules()
			b, w := rules[tt.fields.rule].Eval(tt.fields.target, tt.fields.scanContext)
			got := want{
				broken: b,
				result: w,
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyVaultScanner Rule.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getSKUName() *armkeyvault.SKUName {
	s := armkeyvault.SKUNameStandard
	return &s
}
