// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package features

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cilium/cilium/pkg/policy/api"
)

func Test_ruleType(t *testing.T) {
	type args struct {
		r api.Rule
	}
	type metrics struct {
		npL3Ingested                float64
		npL3Present                 float64
		npHostNPIngested            float64
		npHostNPPresent             float64
		npDNSIngested               float64
		npDNSPresent                float64
		npHTTPIngested              float64
		npHTTPPresent               float64
		npHTTPHeaderMatchesIngested float64
		npHTTPHeaderMatchesPresent  float64
	}
	type wanted struct {
		wantRF      RuleFeatures
		wantMetrics metrics
	}
	tests := []struct {
		name string
		args args
		want wanted
	}{
		{
			name: "L3 from FromEndpoints",
			args: args{
				r: api.Rule{
					Ingress: []api.IngressRule{
						{
							IngressCommonRule: api.IngressCommonRule{
								FromEndpoints: []api.EndpointSelector{
									{},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					L3: true,
				},
				wantMetrics: metrics{
					npL3Ingested: 1,
					npL3Present:  1,
				},
			},
		},
		{
			name: "L3 from FromCIDRSet with CIDRGroupRef",
			args: args{
				r: api.Rule{
					Ingress: []api.IngressRule{
						{
							IngressCommonRule: api.IngressCommonRule{
								FromCIDRSet: []api.CIDRRule{
									{CIDRGroupRef: "some-group-ref"},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					L3: true,
				},
				wantMetrics: metrics{
					npL3Ingested: 1,
					npL3Present:  1,
				},
			},
		},
		{
			name: "L3 IngressDeny from FromCIDRSet with CIDRGroupRef",
			args: args{
				r: api.Rule{
					IngressDeny: []api.IngressDenyRule{
						{
							IngressCommonRule: api.IngressCommonRule{
								FromCIDRSet: []api.CIDRRule{
									{CIDRGroupRef: "some-group-ref"},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					L3: true,
				},
				wantMetrics: metrics{
					npL3Ingested: 1,
					npL3Present:  1,
				},
			},
		},
		{
			name: "L3 from Ingress ToNodes",
			args: args{
				r: api.Rule{
					Ingress: []api.IngressRule{
						{
							IngressCommonRule: api.IngressCommonRule{
								FromNodes: []api.EndpointSelector{
									{},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					L3:   true,
					Host: true,
				},
				wantMetrics: metrics{
					npL3Ingested:     1,
					npL3Present:      1,
					npHostNPIngested: 1,
					npHostNPPresent:  1,
				},
			},
		},
		{
			name: "L3 from Egress ToNodes",
			args: args{
				r: api.Rule{
					Egress: []api.EgressRule{
						{
							EgressCommonRule: api.EgressCommonRule{
								ToNodes: []api.EndpointSelector{
									{},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					L3:   true,
					Host: true,
				},
				wantMetrics: metrics{
					npL3Ingested:     1,
					npL3Present:      1,
					npHostNPIngested: 1,
					npHostNPPresent:  1,
				},
			},
		},
		{
			name: "No L3 rule present",
			args: args{
				r: api.Rule{
					Ingress: []api.IngressRule{
						{
							IngressCommonRule: api.IngressCommonRule{},
						},
					},
					Egress: []api.EgressRule{
						{
							EgressCommonRule: api.EgressCommonRule{},
						},
					},
				},
			},
		},
		{
			name: "L3 from IngressDeny FromNodes",
			args: args{
				r: api.Rule{
					IngressDeny: []api.IngressDenyRule{
						{
							IngressCommonRule: api.IngressCommonRule{
								FromNodes: []api.EndpointSelector{
									{},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					L3:   true,
					Host: true,
				},
				wantMetrics: metrics{
					npL3Ingested:     1,
					npL3Present:      1,
					npHostNPIngested: 1,
					npHostNPPresent:  1,
				},
			},
		},
		{
			name: "L3 from IngressDeny IsL3",
			args: args{
				r: api.Rule{
					IngressDeny: []api.IngressDenyRule{
						{
							IngressCommonRule: api.IngressCommonRule{
								FromCIDR: []api.CIDR{"192.168.0.0/24"},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					L3: true,
				},
				wantMetrics: metrics{
					npL3Ingested: 1,
					npL3Present:  1,
				},
			},
		},
		{
			name: "L3 from EgressDeny IsL3",
			args: args{
				r: api.Rule{
					EgressDeny: []api.EgressDenyRule{
						{
							EgressCommonRule: api.EgressCommonRule{
								ToCIDR: []api.CIDR{"192.168.0.0/24"},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					L3: true,
				},
				wantMetrics: metrics{
					npL3Ingested: 1,
					npL3Present:  1,
				},
			},
		},
		{
			name: "L3 from EgressDeny ToNodes",
			args: args{
				r: api.Rule{
					EgressDeny: []api.EgressDenyRule{
						{
							EgressCommonRule: api.EgressCommonRule{
								ToNodes: []api.EndpointSelector{
									{},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					L3:   true,
					Host: true,
				},
				wantMetrics: metrics{
					npL3Ingested:     1,
					npL3Present:      1,
					npHostNPIngested: 1,
					npHostNPPresent:  1,
				},
			},
		},
		{
			name: "Host from EgressDeny ToNodes",
			args: args{
				r: api.Rule{
					EgressDeny: []api.EgressDenyRule{
						{
							EgressCommonRule: api.EgressCommonRule{
								ToNodes: []api.EndpointSelector{
									{},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					L3:   true,
					Host: true,
				},
				wantMetrics: metrics{
					npL3Ingested:     1,
					npL3Present:      1,
					npHostNPIngested: 1,
					npHostNPPresent:  1,
				},
			},
		},
		{
			name: "Host from Egress ToNodes",
			args: args{
				r: api.Rule{
					Egress: []api.EgressRule{
						{
							EgressCommonRule: api.EgressCommonRule{
								ToNodes: []api.EndpointSelector{
									{},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					L3:   true,
					Host: true,
				},
				wantMetrics: metrics{
					npL3Ingested:     1,
					npL3Present:      1,
					npHostNPIngested: 1,
					npHostNPPresent:  1,
				},
			},
		},
		{
			name: "DNS rules",
			args: args{
				r: api.Rule{
					Egress: []api.EgressRule{
						{
							ToPorts: api.PortRules{
								{
									Rules: &api.L7Rules{
										DNS: []api.PortRuleDNS{
											{
												MatchName: "cilium.io",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					DNS: true,
				},
				wantMetrics: metrics{
					npDNSIngested: 1,
					npDNSPresent:  1,
				},
			},
		},
		{
			name: "FQDN rules",
			args: args{
				r: api.Rule{
					Egress: []api.EgressRule{
						{
							ToFQDNs: api.FQDNSelectorSlice{
								{
									MatchName:    "cilium.io",
									MatchPattern: "",
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					DNS: true,
				},
				wantMetrics: metrics{
					npDNSIngested: 1,
					npDNSPresent:  1,
				},
			},
		},
		{
			name: "HTTP ingress rules",
			args: args{
				r: api.Rule{
					Ingress: []api.IngressRule{
						{
							ToPorts: api.PortRules{
								{
									Rules: &api.L7Rules{
										HTTP: []api.PortRuleHTTP{
											{},
										},
									},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					HTTP: true,
				},
				wantMetrics: metrics{
					npHTTPIngested: 1,
					npHTTPPresent:  1,
				},
			},
		},
		{
			name: "HTTP egress rules",
			args: args{
				r: api.Rule{
					Egress: []api.EgressRule{
						{
							ToPorts: api.PortRules{
								{
									Rules: &api.L7Rules{
										HTTP: []api.PortRuleHTTP{
											{},
										},
									},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					HTTP: true,
				},
				wantMetrics: metrics{
					npHTTPIngested: 1,
					npHTTPPresent:  1,
				},
			},
		},
		{
			name: "HTTP matches ingress rules",
			args: args{
				r: api.Rule{
					Ingress: []api.IngressRule{
						{
							ToPorts: api.PortRules{
								{
									Rules: &api.L7Rules{
										HTTP: []api.PortRuleHTTP{
											{
												HeaderMatches: []*api.HeaderMatch{
													{},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					HTTP:              true,
					HTTPHeaderMatches: true,
				},
				wantMetrics: metrics{
					npHTTPIngested:              1,
					npHTTPPresent:               1,
					npHTTPHeaderMatchesIngested: 1,
					npHTTPHeaderMatchesPresent:  1,
				},
			},
		},
		{
			name: "HTTP matches egress rules",
			args: args{
				r: api.Rule{
					Egress: []api.EgressRule{
						{
							ToPorts: api.PortRules{
								{
									Rules: &api.L7Rules{
										HTTP: []api.PortRuleHTTP{
											{
												HeaderMatches: []*api.HeaderMatch{
													{},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: wanted{
				wantRF: RuleFeatures{
					HTTP:              true,
					HTTPHeaderMatches: true,
				},
				wantMetrics: metrics{
					npHTTPIngested:              1,
					npHTTPPresent:               1,
					npHTTPHeaderMatchesIngested: 1,
					npHTTPHeaderMatchesPresent:  1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := ruleType(tt.args.r)
			assert.Equalf(t, tt.want.wantRF, rt, "ruleType(%v)", tt.args.r)

			metrics := NewMetrics(true)
			metrics.AddRule(tt.args.r)

			assert.Equalf(t, tt.want.wantMetrics.npL3Ingested, metrics.NPL3Ingested.Get(), "NPL3Ingested different")
			assert.Equalf(t, tt.want.wantMetrics.npL3Present, metrics.NPL3Present.Get(), "NPL3Present different")
			assert.Equalf(t, tt.want.wantMetrics.npHostNPIngested, metrics.NPHostNPIngested.Get(), "NPHostNPIngested different")
			assert.Equalf(t, tt.want.wantMetrics.npHostNPPresent, metrics.NPHostNPPresent.Get(), "NPHostNPPresent different")
			assert.Equalf(t, tt.want.wantMetrics.npDNSIngested, metrics.NPDNSIngested.Get(), "NPDNSIngested different")
			assert.Equalf(t, tt.want.wantMetrics.npDNSPresent, metrics.NPDNSPresent.Get(), "NPDNSPresent different")
			assert.Equalf(t, tt.want.wantMetrics.npHTTPIngested, metrics.NPHTTPIngested.Get(), "NPDNSIngested different")
			assert.Equalf(t, tt.want.wantMetrics.npHTTPPresent, metrics.NPHTTPPresent.Get(), "NPDNSPresent different")
			assert.Equalf(t, tt.want.wantMetrics.npHTTPHeaderMatchesIngested, metrics.NPHTTPHeaderMatchesIngested.Get(), "NPHTTPHeaderMatchesIngested different")
			assert.Equalf(t, tt.want.wantMetrics.npHTTPHeaderMatchesPresent, metrics.NPHTTPHeaderMatchesPresent.Get(), "NPHTTPHeaderMatchesPresent different")

			metrics.DelRule(tt.args.r)

			assert.Equalf(t, tt.want.wantMetrics.npL3Ingested, metrics.NPL3Ingested.Get(), "NPL3Ingested different")
			assert.Equalf(t, float64(0), metrics.NPL3Present.Get(), "NPL3Present different")
			assert.Equalf(t, tt.want.wantMetrics.npHostNPIngested, metrics.NPHostNPIngested.Get(), "NPHostNPIngested different")
			assert.Equalf(t, float64(0), metrics.NPHostNPPresent.Get(), "NPHostNPPresent different")
			assert.Equalf(t, tt.want.wantMetrics.npDNSIngested, metrics.NPDNSIngested.Get(), "NPHostNPIngested different")
			assert.Equalf(t, float64(0), metrics.NPDNSPresent.Get(), "NPDNSPresent different")
			assert.Equalf(t, tt.want.wantMetrics.npHTTPIngested, metrics.NPHTTPIngested.Get(), "NPHTTPIngested different")
			assert.Equalf(t, float64(0), metrics.NPHTTPPresent.Get(), "NPHTTPPresent different")
			assert.Equalf(t, tt.want.wantMetrics.npHTTPHeaderMatchesIngested, metrics.NPHTTPHeaderMatchesIngested.Get(), "NPHTTPHeaderMatchesIngested different")
			assert.Equalf(t, float64(0), metrics.NPHTTPHeaderMatchesPresent.Get(), "NPHTTPHeaderMatchesPresent different")

		})
	}
}
