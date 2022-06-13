package lib

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// GitTagFunc is a function that returns the git tag from CI env vars or the
// working tree if available.
var GitTagFunc = function.New(&function.Spec{
	Params: []function.Parameter{},
	Type:   function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		var tag string
		var err error
		for _, fn := range []func() (string, error){
			func() (string, error) {
				// https://buildkite.com/docs/pipelines/environment-variables
				return os.Getenv("BUILDKITE_TAG"), nil
			},
			func() (string, error) {
				// https://circleci.com/docs/2.0/env-vars/
				return os.Getenv("CIRCLE_TAG"), nil
			},
			func() (string, error) {
				// https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
				if ref := os.Getenv("GITHUB_REF"); strings.HasPrefix(ref, "refs/tags/") {
					return strings.TrimLeft(ref, "refs/tags/"), nil
				}
				return "", nil
			},
			func() (string, error) {
				// https://docs.semaphoreci.com/ci-cd-environment/environment-variables/
				if ref := os.Getenv("SEMAPHORE_GIT_REF"); strings.HasPrefix(ref, "refs/tags/") {
					return strings.TrimLeft(ref, "refs/tags/"), nil
				}
				return "", nil
			},
			func() (string, error) {
				// https://docs.travis-ci.com/user/environment-variables/#default-environment-variables
				return os.Getenv("TRAVIS_TAG"), nil
			},
			func() (string, error) {
				return gitRun("tag", "--points-at", "HEAD", "--sort", "-version:creatordate")
			},
			func() (string, error) {
				return gitRun("describe", "--tags", "--abbrev=0")
			},
		} {
			tag, err = fn()
			if tag != "" || err != nil {
				return cty.StringVal(tag), err
			}
		}
		return cty.StringVal(tag), err
	},
})

// GitTag returns the git tag from CI env vars or the working tree if available.
func GitTag() (cty.Value, error) {
	return GitTagFunc.Call(nil)
}

func gitRun(args ...string) (string, error) {
	var extraArgs = []string{
		"-c", "log.showSignature=false",
	}

	args = append(extraArgs, args...)
	var cmd = exec.Command("git", args...)

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", errors.New(strings.TrimSuffix(stderr.String(), "\n"))
	}

	return strings.ReplaceAll(strings.Split(stdout.String(), "\n")[0], "'", ""), nil
}
