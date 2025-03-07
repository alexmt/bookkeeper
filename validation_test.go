package bookkeeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateAndCanonicalizeRequest(t *testing.T) {
	testCases := []struct {
		name       string
		req        RenderRequest
		assertions func(RenderRequest, error)
	}{
		{
			name: "missing RepoURL",
			req:  RenderRequest{},
			assertions: func(req RenderRequest, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "RepoURL is a required field")
			},
		},
		{
			name: "invalid RepoURL",
			req: RenderRequest{
				RepoURL: "foobar",
			},
			assertions: func(req RenderRequest, err error) {
				require.Error(t, err)
				require.Contains(
					t,
					err.Error(),
					"does not appear to be a valid git repository URL",
				)
			},
		},
		{
			name: "missing Password",
			req: RenderRequest{
				RepoURL: "https://github.com/akuity/foobar",
			},
			assertions: func(req RenderRequest, err error) {
				require.Error(t, err)
				require.Contains(
					t,
					err.Error(),
					"RepoCreds.Password is a required field",
				)
			},
		},
		{
			name: "missing TargetBranch",
			req: RenderRequest{
				RepoURL: "https://github.com/akuity/foobar",
				RepoCreds: RepoCredentials{
					Password: "foobar",
				},
				Ref: "1abcdef2",
			},
			assertions: func(req RenderRequest, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "TargetBranch is a required field")
			},
		},
		{
			name: "invalid TargetBranch",
			req: RenderRequest{
				RepoURL: "https://github.com/akuity/foobar",
				RepoCreds: RepoCredentials{
					Password: "foobar",
				},
				Ref:          "1abcdef2",
				TargetBranch: "env/dev*", // * is an invalid character
			},
			assertions: func(req RenderRequest, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "is an invalid branch name")
			},
		},
		{
			name: "empty string image",
			req: RenderRequest{
				RepoURL: "https://github.com/akuity/foobar",
				RepoCreds: RepoCredentials{
					Password: "foobar",
				},
				Ref:          "1abcdef2",
				TargetBranch: "env/dev",
				Images:       []string{""}, // no good
			},
			assertions: func(req RenderRequest, err error) {
				require.Error(t, err)
				require.Contains(
					t,
					err.Error(),
					"Images must not contain any empty strings",
				)
			},
		},
		{
			name: "validation succeeds",
			req: RenderRequest{
				RepoURL: "  https://github.com/akuity/foobar  ",
				RepoCreds: RepoCredentials{
					Password: "  foobar  ",
				},
				Ref:          "  1abcdef2 ",
				TargetBranch: "  refs/heads/env/dev  ",
				Images:       []string{" akuity/some-image "}, // no good
			},
			assertions: func(req RenderRequest, err error) {
				require.NoError(t, err)
				require.Equal(t, "https://github.com/akuity/foobar", req.RepoURL)
				require.Equal(t, "foobar", req.RepoCreds.Password)
				require.Equal(t, "1abcdef2", req.Ref)
				require.Equal(t, "env/dev", req.TargetBranch)
				require.Equal(t, []string{"akuity/some-image"}, req.Images)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.assertions(validateAndCanonicalizeRequest(testCase.req))
		})
	}
}
