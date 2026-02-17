package bake

import (
	"testing"

	"github.com/docker/buildx/util/buildflags"
	"github.com/moby/buildkit/client/llb"
	"github.com/stretchr/testify/require"
)

func TestGitAuthSecretsFromEnviron(t *testing.T) {
	t.Run("base keys", func(t *testing.T) {
		secrets := gitAuthSecretsFromEnviron([]string{
			bakeGitAuthTokenEnv + "=token",
			bakeGitAuthHeaderEnv + "=basic",
		})
		require.Equal(t, []string{
			llb.GitAuthTokenKey + "|" + bakeGitAuthTokenEnv,
			llb.GitAuthHeaderKey + "|" + bakeGitAuthHeaderEnv,
		}, secretPairs(secrets))
	})
	t.Run("derives host from remote url", func(t *testing.T) {
		secrets := gitAuthSecretsFromEnviron([]string{
			bakeGitAuthTokenEnv + "=token",
			bakeGitAuthHeaderEnv + "=basic",
		}, "https://example.com/org/repo.git")
		require.Equal(t, []string{
			llb.GitAuthTokenKey + "|" + bakeGitAuthTokenEnv,
			llb.GitAuthTokenKey + ".example.com|" + bakeGitAuthTokenEnv,
			llb.GitAuthHeaderKey + "|" + bakeGitAuthHeaderEnv,
			llb.GitAuthHeaderKey + ".example.com|" + bakeGitAuthHeaderEnv,
		}, secretPairs(secrets))
	})
	t.Run("ignores host suffixed keys", func(t *testing.T) {
		secrets := gitAuthSecretsFromEnviron([]string{
			bakeGitAuthTokenEnv + ".example.com=token",
			bakeGitAuthHeaderEnv + ".example.com=basic",
		})
		require.Empty(t, secrets)
	})
}

func secretPairs(secrets buildflags.Secrets) []string {
	out := make([]string, 0, len(secrets))
	for _, s := range secrets {
		out = append(out, s.ID+"|"+s.Env)
	}
	return out
}
