package lib

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGit(t *testing.T) {
	out, err := gitRun("status")
	require.NoError(t, err)
	require.NotEmpty(t, out)

	out, err = gitRun("command-that-dont-exist")
	require.Error(t, err)
	require.Empty(t, out)
	require.Equal(t, "git: 'command-that-dont-exist' is not a git command. See 'git --help'.", err.Error())
}

func TestGitTagsPointsAt(t *testing.T) {
	mktmp(t)
	gitInit(t)
	gitCommit(t, "bar")
	gitTag(t, "v0.8.0")
	gitCommit(t, "foo")
	gitTag(t, "v0.9.0")

	out, err := gitRun("tag", "--points-at", "HEAD", "--sort", "-version:creatordate")
	require.NoError(t, err)
	require.Equal(t, "v0.9.0", out)
}

func TestGitDescribeTags(t *testing.T) {
	mktmp(t)
	gitInit(t)
	gitCommit(t, "bar")
	gitTag(t, "v0.8.0")
	gitCommit(t, "foo")
	gitTag(t, "v0.9.0")

	out, err := gitRun("describe", "--tags", "--abbrev=0")
	require.NoError(t, err)
	require.Equal(t, "v0.9.0", out)
}

func TestGitTag(t *testing.T) {
	tests := []struct {
		Envs     map[string]string
		Expected string
	}{
		{
			Envs: map[string]string{
				"BUILDKITE_TAG": "v0.9.0",
			},
			Expected: "v0.9.0",
		},
		{
			Envs: map[string]string{
				"CIRCLE_TAG": "v0.9.1",
			},
			Expected: "v0.9.1",
		},
		{
			Envs: map[string]string{
				"GITHUB_REF": "refs/tags/v0.9.2",
			},
			Expected: "v0.9.2",
		},
		{
			Envs: map[string]string{
				"SEMAPHORE_GIT_REF": "refs/tags/v0.9.3",
			},
			Expected: "v0.9.3",
		},
		{
			Envs: map[string]string{
				"TRAVIS_TAG": "v0.9.4",
			},
			Expected: "v0.9.4",
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%#v", test.Envs), func(t *testing.T) {
			for k, v := range test.Envs {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}
			result, err := GitTag()
			require.NoError(t, err)
			require.NotEmpty(t, result.AsString())
			require.Equal(t, test.Expected, result.AsString())
		})
	}
}

func gitInit(tb testing.TB) {
	tb.Helper()
	out, err := fakeGit("init")
	require.NoError(tb, err)
	require.Contains(tb, out, "Initialized empty Git repository")
	require.NoError(tb, err)
	gitCheckoutBranch(tb, "main")
	_, _ = fakeGit("branch", "-D", "master")
}

func gitCommit(tb testing.TB, msg string) {
	tb.Helper()
	out, err := fakeGit("commit", "--allow-empty", "-m", msg)
	require.NoError(tb, err)
	require.Contains(tb, out, "main", msg)
}

func gitTag(tb testing.TB, tag string) {
	tb.Helper()
	out, err := fakeGit("tag", tag)
	require.NoError(tb, err)
	require.Empty(tb, out)
}

func gitCheckoutBranch(tb testing.TB, name string) {
	tb.Helper()
	out, err := fakeGit("checkout", "-b", name)
	require.NoError(tb, err)
	require.Empty(tb, out)
}

func fakeGit(args ...string) (string, error) {
	allArgs := []string{
		"-c", "user.name=buildx",
		"-c", "user.email=buildx@docker.com",
		"-c", "commit.gpgSign=false",
		"-c", "tag.gpgSign=false",
		"-c", "log.showSignature=false",
	}
	allArgs = append(allArgs, args...)
	return gitRun(allArgs...)
}

func mktmp(tb testing.TB) string {
	tb.Helper()
	folder := tb.TempDir()
	current, err := os.Getwd()
	require.NoError(tb, err)
	require.NoError(tb, os.Chdir(folder))
	tb.Cleanup(func() {
		require.NoError(tb, os.Chdir(current))
	})
	return folder
}
