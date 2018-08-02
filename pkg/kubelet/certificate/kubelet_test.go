/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package certificate

import (
	"net"
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	//v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kubernetes/pkg/kubelet/apis/kubeletconfig"
)

func TestAddressesToHostnamesAndIPs(t *testing.T) {
	tests := []struct {
		name         string
		addresses    []v1.NodeAddress
		wantDNSNames []string
		wantIPs      []net.IP
	}{
		{
			name:         "empty",
			addresses:    nil,
			wantDNSNames: nil,
			wantIPs:      nil,
		},
		{
			name:         "ignore empty values",
			addresses:    []v1.NodeAddress{{Type: v1.NodeHostName, Address: ""}},
			wantDNSNames: nil,
			wantIPs:      nil,
		},
		{
			name: "ignore invalid IPs",
			addresses: []v1.NodeAddress{
				{Type: v1.NodeInternalIP, Address: "1.2"},
				{Type: v1.NodeExternalIP, Address: "3.4"},
			},
			wantDNSNames: nil,
			wantIPs:      nil,
		},
		{
			name: "dedupe values",
			addresses: []v1.NodeAddress{
				{Type: v1.NodeHostName, Address: "hostname"},
				{Type: v1.NodeExternalDNS, Address: "hostname"},
				{Type: v1.NodeInternalDNS, Address: "hostname"},
				{Type: v1.NodeInternalIP, Address: "1.1.1.1"},
				{Type: v1.NodeExternalIP, Address: "1.1.1.1"},
			},
			wantDNSNames: []string{"hostname"},
			wantIPs:      []net.IP{net.ParseIP("1.1.1.1")},
		},
		{
			name: "order values",
			addresses: []v1.NodeAddress{
				{Type: v1.NodeHostName, Address: "hostname-2"},
				{Type: v1.NodeExternalDNS, Address: "hostname-1"},
				{Type: v1.NodeInternalDNS, Address: "hostname-3"},
				{Type: v1.NodeInternalIP, Address: "2.2.2.2"},
				{Type: v1.NodeExternalIP, Address: "1.1.1.1"},
				{Type: v1.NodeInternalIP, Address: "3.3.3.3"},
			},
			wantDNSNames: []string{"hostname-1", "hostname-2", "hostname-3"},
			wantIPs:      []net.IP{net.ParseIP("1.1.1.1"), net.ParseIP("2.2.2.2"), net.ParseIP("3.3.3.3")},
		},
		{
			name: "handle IP and DNS hostnames",
			addresses: []v1.NodeAddress{
				{Type: v1.NodeHostName, Address: "hostname"},
				{Type: v1.NodeHostName, Address: "1.1.1.1"},
			},
			wantDNSNames: []string{"hostname"},
			wantIPs:      []net.IP{net.ParseIP("1.1.1.1")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDNSNames, gotIPs := addressesToHostnamesAndIPs(tt.addresses)
			if !reflect.DeepEqual(gotDNSNames, tt.wantDNSNames) {
				t.Errorf("addressesToHostnamesAndIPs() gotDNSNames = %v, want %v", gotDNSNames, tt.wantDNSNames)
			}
			if !reflect.DeepEqual(gotIPs, tt.wantIPs) {
				t.Errorf("addressesToHostnamesAndIPs() gotIPs = %v, want %v", gotIPs, tt.wantIPs)
			}
		})
	}
}

func TestKubeletCertificateExecManager(t *testing.T) {
	tests := []struct {
		name       string
		addresses  []v1.NodeAddress
		certDir    string
		nodeName   string
		pluginPath string
	}{
		{
			name:       "env and output",
			pluginPath: "testdata/plugin.sh",
			certDir:    "testdata/certs",
			nodeName:   "localhost",
			addresses: []v1.NodeAddress{
				{Type: v1.NodeHostName, Address: "localhost"},
				{Type: v1.NodeInternalIP, Address: "127.0.0.1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAddresses := func() []v1.NodeAddress {
				return tt.addresses
			}
			cfg := &kubeletconfig.KubeletConfiguration{
				ServerCertExecEnvSubject:  "test_subject",
				ServerCertExecEnvDNSnames: "test_dns",
				ServerCertExecEnvIPnames:  "test_ip",
				ServerCertExecEnvOther: map[string]string{
					"test_other": "foo",
				},
			}
			m, err := NewKubeletServerCertificateExecManager(tt.pluginPath, cfg, types.NodeName(tt.nodeName), getAddresses, tt.certDir)
			if err != nil {
				t.Errorf("NewKubeletServerCertificateExecManager returned err: %v", err)
			}
			m.Start()
			cert := m.Current()
			if cert == nil {
				t.Errorf("no certificate")
			}
		})
	}
}

func ipNamesToIPs(ips []string) []net.IP {
	out := make([]net.IP, 0)
	for _, ip := range ips {
		if parsed := net.ParseIP(ip); parsed != nil {
			out = append(out, parsed)
		}
	}
	return out
}

func TestGetNodeExecEnv(t *testing.T) {
	tests := []struct {
		name         string
		nodeName     string
		kubeCfg      *kubeletconfig.KubeletConfiguration
		dnsNames     []string
		ipNames      []string
		expectedEnvs []api.ExecEnvVar
	}{
		{
			name:     "defaults no names",
			nodeName: "localhost",
			kubeCfg: &kubeletconfig.KubeletConfiguration{
				ServerCertExecEnvSubject:  "",
				ServerCertExecEnvDNSnames: "",
				ServerCertExecEnvIPnames:  "",
				ServerCertExecEnvOther:    nil,
			},
			dnsNames: nil,
			ipNames:  nil,
			expectedEnvs: []api.ExecEnvVar{
				{
					Name:  defaultSubjectEnv,
					Value: "CN=system:node:localhost,O=system:nodes",
				},
				{
					Name:  defaultDNSNamesEnv,
					Value: "",
				},
				{
					Name:  defaultIPNamesEnv,
					Value: "",
				},
			},
		},
		{
			name:     "defaults with names",
			nodeName: "localhost",
			kubeCfg: &kubeletconfig.KubeletConfiguration{
				ServerCertExecEnvSubject:  "",
				ServerCertExecEnvDNSnames: "",
				ServerCertExecEnvIPnames:  "",
				ServerCertExecEnvOther:    nil,
			},
			dnsNames: []string{"localhost"},
			ipNames:  []string{"127.0.0.1"},
			expectedEnvs: []api.ExecEnvVar{
				{
					Name:  defaultSubjectEnv,
					Value: "CN=system:node:localhost,O=system:nodes",
				},
				{
					Name:  defaultDNSNamesEnv,
					Value: "localhost",
				},
				{
					Name:  defaultIPNamesEnv,
					Value: "127.0.0.1",
				},
			},
		},
		{
			name:     "defaults multiple names",
			nodeName: "localhost",
			kubeCfg: &kubeletconfig.KubeletConfiguration{
				ServerCertExecEnvSubject:  "",
				ServerCertExecEnvDNSnames: "",
				ServerCertExecEnvIPnames:  "",
				ServerCertExecEnvOther:    nil,
			},
			dnsNames: []string{
				"localhost",
				"foo",
			},
			ipNames: []string{
				"127.0.0.1",
				"10.0.0.1",
			},
			expectedEnvs: []api.ExecEnvVar{
				{
					Name:  defaultSubjectEnv,
					Value: "CN=system:node:localhost,O=system:nodes",
				},
				{
					Name:  defaultDNSNamesEnv,
					Value: "localhost,foo",
				},
				{
					Name:  defaultIPNamesEnv,
					Value: "127.0.0.1,10.0.0.1",
				},
			},
		},
		{
			name:     "defaults multiple names",
			nodeName: "localhost",
			kubeCfg: &kubeletconfig.KubeletConfiguration{
				ServerCertExecEnvSubject:  "subject_names",
				ServerCertExecEnvDNSnames: "dns_names",
				ServerCertExecEnvIPnames:  "ip_names",
				ServerCertExecEnvOther: map[string]string{
					"other1": "foo",
					"other2": "bar",
				},
			},
			dnsNames: []string{
				"localhost",
				"foo",
			},
			ipNames: []string{
				"127.0.0.1",
				"10.0.0.1",
			},
			expectedEnvs: []api.ExecEnvVar{
				{
					Name:  "subject_names",
					Value: "CN=system:node:localhost,O=system:nodes",
				},
				{
					Name:  "dns_names",
					Value: "localhost,foo",
				},
				{
					Name:  "ip_names",
					Value: "127.0.0.1,10.0.0.1",
				},
				{
					Name:  "other1",
					Value: "foo",
				},
				{
					Name:  "other2",
					Value: "bar",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envs := getNodeExecEnv(tt.kubeCfg, tt.nodeName, tt.dnsNames, ipNamesToIPs(tt.ipNames))
			if !reflect.DeepEqual(tt.expectedEnvs, envs) {
				t.Errorf("expected output does not match: expected %v, got %v", tt.expectedEnvs, envs)
			}
		})
	}
}
