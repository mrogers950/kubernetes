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
	"k8s.io/apimachinery/pkg/types"
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

var nodeCert = `-----BEGIN CERTIFICATE REQUEST-----
MIIDBDCCAewCAQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0B
AQEFAAOCAQ8AMIIBCgKCAQEAypwEc8kgzCWH9R3HoC6tdUdbalF5FrWGyv84Dplb
/qdyO3++rqh3fnqGUWaErAy0TaC1v5lYihX4ASlCPbWvTTZBgzaFlxBdAm3NPYr+
NRv4iBvMuPDOA0s9jVOK3vudqHbl7mYX46eHbyogX8Iu1lgzS97KRLje20lGWrEh
iMDjjOcBijlk7t7ugW3id58RO52GXQrliDRVmjiWX8XBU7TpbSa/meWut3nrbZId
6QYyzWbSVW74rEcjRxD8vhQ8wVTFhsSJ2Gp/IbRToYeIozIwm4p0SGBxqz3KwQmj
IzxFcgqG/Y+yQLn5IeoDvXXrOVbG+nBtFo1Tjpw9NOy1xwIDAQABoIGqMIGnBgkq
hkiG9w0BCQ4xgZkwgZYwHQYDVR0OBBYEFF7nH5h5PHo09n6jB9xUKqHdLRJ/MAkG
A1UdEwQCMAAwCwYDVR0PBAQDAgWgMBoGA1UdEQQTMBGCCWxvY2FsaG9zdIcEfwAA
ATAsBglghkgBhvhCAQ0EHxYdT3BlblNTTCBHZW5lcmF0ZWQgQ2VydGlmaWNhdGUw
EwYDVR0lBAwwCgYIKwYBBQUHAwEwDQYJKoZIhvcNAQELBQADggEBADo7mX9uSdwG
Dga6raavcHDGOsMNJ3N5Q6pE05s5+yPHBsHQlbWjhWRpuxNBAoLCR0sF92H97BFl
NZuxBAIWRetNsaWffg99Zgp/cKeLxdH96jhX5L+KTc7/KmB7H0aDj2flB85vdCOR
bfs0jXhfM3FBxFwQw73FSQORwhq4zujmSS8vKELxp/H9XR0f7vKlz4JE6S+Yfq2b
On9kpogRVsff1Ic+b8PRFc+B0a6DYUw+LR/SsJaa/4mBOKzjlMWYiM3fSP9hVZkh
b+F+u00qh2g+67m9D+lN/CHJluybdFhxBVszSZqB/MK1H77YxItLJC1GRUdx9Zx2
dPH+u69r214=
-----END CERTIFICATE REQUEST-----`

var nodeKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAypwEc8kgzCWH9R3HoC6tdUdbalF5FrWGyv84Dplb/qdyO3++
rqh3fnqGUWaErAy0TaC1v5lYihX4ASlCPbWvTTZBgzaFlxBdAm3NPYr+NRv4iBvM
uPDOA0s9jVOK3vudqHbl7mYX46eHbyogX8Iu1lgzS97KRLje20lGWrEhiMDjjOcB
ijlk7t7ugW3id58RO52GXQrliDRVmjiWX8XBU7TpbSa/meWut3nrbZId6QYyzWbS
VW74rEcjRxD8vhQ8wVTFhsSJ2Gp/IbRToYeIozIwm4p0SGBxqz3KwQmjIzxFcgqG
/Y+yQLn5IeoDvXXrOVbG+nBtFo1Tjpw9NOy1xwIDAQABAoIBACiFSXKqr++ENgu0
t/72NuS0r7i0sKX1Cg9BOcHZtAdbD8KMiuM9eCCIeqJ/AVuzcr/vu0mlboq3WBFL
Yh8bXgLwLewDFHag5CkfMqPzT2HpxSvbe3clWd5Yxuej5Ksx4VcW6GdkbbSvBawa
3bypBlsB6shqt0NFQfTTU8nBkTZbGk3pv1VBKKEQnrYsV03tSO5Ph5jppdUusYFm
2R8mgDwOR5PxcpN0uzA9f2XbKIPIUnhiNNcIss5jASTi5IPojOPgQurTq9aDR2XA
nsJB+h+AqZbwQUpzrCMAlK8Oz/xwUs6BqSptVFDY0yVNS1qVfap4xd5oNWxGksQV
o3P0kyECgYEA5WNy7gqrZ+3Kuk3vqWghDWmZ+jhKjeuGy45YnMYs7WRFK5ckBo6r
iucbj/1GwqIfmFAXUIDTlG5tJHL3s0Nmv5Q0Tjq5umjXpf/iiiQMLK/E8FgMpH0a
dMGqcFW49VBYa/CD2/HT6YjyDqy13Dj9ptYnMbMe3z/BG8puG3hoPJECgYEA4h1D
xRQ9HEUVOVKMFLoCvAzBeVh3UZmA5W7Bl904+qLN3wkGy3bIBiVlQkZtaY5AspPg
NkccJYsuPoDkx0Qp8qt0M6hWOmgSPBQMpSnzWoQF42YfuUemqFjfM0j4EEIEvGhD
O2vbryDMn4NHAk/u1pnxmlVwOM8LNM3hJb+QWNcCgYEAvAC7BHAINcDF49XWdCDc
3hJL2bFjIVgE/TZoV+1wiwwgSO6x3xH1dH2fsG6kHQclH/+cbCV5w3CR0UrMysaW
IrRD/k3RRP+CpxHGyPNsav+QSG/RxMqn8UN8/l6znZNBNQ5F8/EKfp/3y6Ev2BN5
iNCCBRDKX6zwB2fswGT6AZECgYAoB2pB721KHei99yEZYjytscxmgQTOi1BITa00
B1PY+w1bGKv9RQ/wFpqweutProFBm/Ara7dN5i/PnN3jcOvELBosMvbg7B+eRyZd
7ulH8utf8GpZUJfuYZ1R5O8VYbqY6BRO5q9Dd5kB/CmL/T6Y+zPMUKfHRtADDxd2
qU0SjQKBgDrOCBhxDn61cRyMyUMo7zXxsC+YgZYwtDyYtI/s37Px2nWO+lfGVX/q
ZhM/GfLa20Uqe9a9Khw7tJ1VSDvFayyv/c2WuszIjLXgE1wDLnOVcetQSn3eDag1
IFBW9JlHHZ5RgEDdy3/lmWMi4nLDsa+/UNXs++AXOSAB0o+yFUQP
-----END RSA PRIVATE KEY-----
`

func TestKubeletCertificateExecManager(t *testing.T) {
	tests := []struct {
		name       string
		addresses  []v1.NodeAddress
		keyFile    string
		certFile   string
		certDir    string
		nodeName   string
		pluginPath string
	}{
		{
			name:       "first",
			pluginPath: "testdata/plugin.sh",
			certDir:    "testdata/certs",
			keyFile:    "testdata/certs/node.key",
			certFile:   "testdata/certs/node.crt",
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
			m, err := NewKubeletServerCertificateExecManager(tt.pluginPath, &kubeletconfig.KubeletConfiguration{TLSCertFile: tt.certFile, TLSPrivateKeyFile: tt.keyFile}, types.NodeName(tt.nodeName), getAddresses, tt.certDir)
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
