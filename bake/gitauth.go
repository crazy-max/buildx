package bake

import (
	"os"
	"sort"
	"strings"

	"github.com/docker/buildx/util/buildflags"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/gitutil"
)

const (
	bakeGitAuthTokenEnv  = "BUILDX_BAKE_GIT_AUTH_TOKEN" // #nosec G101 -- environment variable key, not a credential
	bakeGitAuthHeaderEnv = "BUILDX_BAKE_GIT_AUTH_HEADER"
)

func gitAuthSecretsFromEnv(remoteURLs ...string) buildflags.Secrets {
	return gitAuthSecretsFromEnviron(os.Environ(), remoteURLs...)
}

func gitAuthSecretsFromEnviron(environ []string, remoteURLs ...string) buildflags.Secrets {
	hosts := gitAuthHostsFromURLs(remoteURLs)
	secrets := make(buildflags.Secrets, 0, 2)
	secrets = append(secrets, gitAuthSecretsForEnv(llb.GitAuthTokenKey, bakeGitAuthTokenEnv, environ, hosts)...)
	secrets = append(secrets, gitAuthSecretsForEnv(llb.GitAuthHeaderKey, bakeGitAuthHeaderEnv, environ, hosts)...)
	return secrets
}

func gitAuthSecretsForEnv(secretIDPrefix, envPrefix string, environ []string, hosts []string) buildflags.Secrets {
	envKey, ok := findGitAuthEnvKey(envPrefix, environ)
	if !ok {
		return nil
	}
	secrets := make(buildflags.Secrets, 0, len(hosts)+1)
	secrets = append(secrets, &buildflags.Secret{
		ID:  secretIDPrefix,
		Env: envKey,
	})
	for _, host := range hosts {
		secrets = append(secrets, &buildflags.Secret{
			ID:  secretIDPrefix + "." + host,
			Env: envKey,
		})
	}
	return secrets
}

func gitAuthHostsFromURLs(remoteURLs []string) []string {
	if len(remoteURLs) == 0 {
		return nil
	}
	hosts := make(map[string]struct{}, len(remoteURLs))
	for _, remoteURL := range remoteURLs {
		gitURL, err := gitutil.ParseURL(remoteURL)
		if err != nil || gitURL.Host == "" {
			continue
		}
		hosts[gitURL.Host] = struct{}{}
	}
	out := make([]string, 0, len(hosts))
	for host := range hosts {
		out = append(out, host)
	}
	sort.Strings(out)
	return out
}

func findGitAuthEnvKey(envKey string, environ []string) (string, bool) {
	for _, env := range environ {
		key, _, ok := strings.Cut(env, "=")
		if !ok {
			continue
		}
		if strings.EqualFold(key, envKey) {
			return key, true
		}
	}
	return "", false
}
