// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package sigr

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/signalr/armsignalr"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/cmendible/azqr/internal/scanners"
)

func TestSignalRScanner_Rules(t *testing.T) {
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
			name: "SignalRScanner DiagnosticSettings",
			fields: fields{
				rule: "DiagnosticSettings",
				target: &armsignalr.ResourceInfo{
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
			name: "SignalRScanner Availability Zones",
			fields: fields{
				rule: "AvailabilityZones",
				target: &armsignalr.ResourceInfo{
					SKU: &armsignalr.ResourceSKU{
						Name: to.StringPtr("Premium"),
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
			name: "SignalRScanner SLA",
			fields: fields{
				rule:                "SLA",
				target:              &armsignalr.ResourceInfo{},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "99.9%",
			},
		},
		{
			name: "SignalRScanner Private Endpoint",
			fields: fields{
				rule: "Private",
				target: &armsignalr.ResourceInfo{
					Properties: &armsignalr.Properties{
						PrivateEndpointConnections: []*armsignalr.PrivateEndpointConnection{
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
			name: "SignalRScanner SKU",
			fields: fields{
				rule: "SKU",
				target: &armsignalr.ResourceInfo{
					SKU: &armsignalr.ResourceSKU{
						Name: to.StringPtr("Premium"),
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "Premium",
			},
		},
		{
			name: "SignalRScanner CAF",
			fields: fields{
				rule: "CAF",
				target: &armsignalr.ResourceInfo{
					Name: to.StringPtr("sigr-test"),
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
			s := &SignalRScanner{
				diagnosticsSettings: tt.fields.diagnosticsSettings,
			}
			rules := s.GetRules()
			b, w := rules[tt.fields.rule].Eval(tt.fields.target, tt.fields.scanContext)
			got := want{
				broken: b,
				result: w,
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SignalRScanner Rule.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
