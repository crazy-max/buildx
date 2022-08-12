// Copyright 2022 Docker Buildx authors
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

package context

import (
	"os"

	"github.com/docker/cli/cli/context"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// FromKubeConfig creates a Kubernetes endpoint from a Kubeconfig file
func FromKubeConfig(kubeconfig, kubeContext, namespaceOverride string) (Endpoint, error) {
	cfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: kubeContext, Context: clientcmdapi.Context{Namespace: namespaceOverride}})
	ns, _, err := cfg.Namespace()
	if err != nil {
		return Endpoint{}, err
	}
	clientcfg, err := cfg.ClientConfig()
	if err != nil {
		return Endpoint{}, err
	}
	var ca, key, cert []byte
	if ca, err = readFileOrDefault(clientcfg.CAFile, clientcfg.CAData); err != nil {
		return Endpoint{}, err
	}
	if key, err = readFileOrDefault(clientcfg.KeyFile, clientcfg.KeyData); err != nil {
		return Endpoint{}, err
	}
	if cert, err = readFileOrDefault(clientcfg.CertFile, clientcfg.CertData); err != nil {
		return Endpoint{}, err
	}
	var tlsData *context.TLSData
	if ca != nil || cert != nil || key != nil {
		tlsData = &context.TLSData{
			CA:   ca,
			Cert: cert,
			Key:  key,
		}
	}
	var usernamePassword *UsernamePassword
	if clientcfg.Username != "" || clientcfg.Password != "" {
		usernamePassword = &UsernamePassword{
			Username: clientcfg.Username,
			Password: clientcfg.Password,
		}
	}
	return Endpoint{
		EndpointMeta: EndpointMeta{
			EndpointMetaBase: context.EndpointMetaBase{
				Host:          clientcfg.Host,
				SkipTLSVerify: clientcfg.Insecure,
			},
			DefaultNamespace: ns,
			AuthProvider:     clientcfg.AuthProvider,
			Exec:             clientcfg.ExecProvider,
			UsernamePassword: usernamePassword,
		},
		TLSData: tlsData,
	}, nil
}

func readFileOrDefault(path string, defaultValue []byte) ([]byte, error) {
	if path != "" {
		return os.ReadFile(path)
	}
	return defaultValue, nil
}
