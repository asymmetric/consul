package responses

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseArn(t *testing.T) {
	cases := map[string]struct {
		arn    string
		expArn *ParsedArn
	}{
		"assumed-role": {
			arn: "arn:aws:sts::000000000000:assumed-role/my-role/session-name",
			expArn: &ParsedArn{
				Partition:     "aws",
				AccountNumber: "000000000000",
				Type:          "assumed-role",
				Path:          "",
				FriendlyName:  "my-role",
				SessionInfo:   "session-name",
			},
		},
		"role": {
			arn: "arn:aws:iam::000000000000:role/my-role",
			expArn: &ParsedArn{
				Partition:     "aws",
				AccountNumber: "000000000000",
				Type:          "role",
				Path:          "",
				FriendlyName:  "my-role",
				SessionInfo:   "",
			},
		},
		"user": {
			arn: "arn:aws:iam::000000000000:user/my-user",
			expArn: &ParsedArn{
				Partition:     "aws",
				AccountNumber: "000000000000",
				Type:          "user",
				Path:          "",
				FriendlyName:  "my-user",
				SessionInfo:   "",
			},
		},
		"role with path": {
			arn: "arn:aws:iam::000000000000:role/path/my-role",
			expArn: &ParsedArn{
				Partition:     "aws",
				AccountNumber: "000000000000",
				Type:          "role",
				Path:          "path",
				FriendlyName:  "my-role",
				SessionInfo:   "",
			},
		},
		"role with path 2": {
			arn: "arn:aws:iam::000000000000:role/path/to/my-role",
			expArn: &ParsedArn{
				Partition:     "aws",
				AccountNumber: "000000000000",
				Type:          "role",
				Path:          "path/to",
				FriendlyName:  "my-role",
				SessionInfo:   "",
			},
		},
		"role with path 3": {
			arn: "arn:aws:iam::000000000000:role/some/path/to/my-role",
			expArn: &ParsedArn{
				Partition:     "aws",
				AccountNumber: "000000000000",
				Type:          "role",
				Path:          "some/path/to",
				FriendlyName:  "my-role",
				SessionInfo:   "",
			},
		},
		"user with path": {
			arn: "arn:aws:iam::000000000000:user/path/my-user",
			expArn: &ParsedArn{
				Partition:     "aws",
				AccountNumber: "000000000000",
				Type:          "user",
				Path:          "path",
				FriendlyName:  "my-user",
				SessionInfo:   "",
			},
		},

		// Invalid cases
		"empty string":               {arn: ""},
		"wildcard":                   {arn: "*"},
		"missing prefix":             {arn: ":aws:sts::000000000000:assumed-role/my-role/session-name"},
		"missing partition":          {arn: "arn::sts::000000000000:assumed-role/my-role/session-name"},
		"missing service":            {arn: "arn:aws:::000000000000:assumed-role/my-role/session-name"},
		"missing separator":          {arn: "arn:aws:sts:000000000000:assumed-role/my-role/session-name"},
		"missing account id":         {arn: "arn:aws:sts:::assumed-role/my-role/session-name"},
		"missing resource":           {arn: "arn:aws:sts::000000000000:"},
		"assumed-role missing parts": {arn: "arn:aws:sts::000000000000:assumed-role/my-role"},
		"role missing parts":         {arn: "arn:aws:sts::000000000000:role"},
		"role missing parts 2":       {arn: "arn:aws:sts::000000000000:role/"},
		"user missing parts":         {arn: "arn:aws:sts::000000000000:user"},
		"user missing parts 2":       {arn: "arn:aws:sts::000000000000:user/"},
		"unsupported service":        {arn: "arn:aws:ecs:us-east-1:000000000000:task/my-task/00000000000000000000000000000000"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			parsed, err := ParseArn(c.arn)
			if c.expArn != nil {
				require.NoError(t, err)
				require.Equal(t, c.expArn, parsed)
			} else {
				require.Error(t, err)
				require.Nil(t, parsed)
			}
		})
	}
}

func TestCanonicalArn(t *testing.T) {
	cases := map[string]struct {
		arn    string
		expArn string
	}{
		"assumed-role arn": {
			arn:    "arn:aws:sts::000000000000:assumed-role/my-role/session-name",
			expArn: "arn:aws:iam::000000000000:role/my-role",
		},
		"role arn": {
			arn:    "arn:aws:iam::000000000000:role/my-role",
			expArn: "arn:aws:iam::000000000000:role/my-role",
		},
		"role arn with path": {
			arn:    "arn:aws:iam::000000000000:role/path/to/my-role",
			expArn: "arn:aws:iam::000000000000:role/my-role",
		},
		"user arn": {
			arn:    "arn:aws:iam::000000000000:user/my-user",
			expArn: "arn:aws:iam::000000000000:user/my-user",
		},
		"user arn with path": {
			arn:    "arn:aws:iam::000000000000:user/path/to/my-user",
			expArn: "arn:aws:iam::000000000000:user/my-user",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			parsed, err := ParseArn(c.arn)
			require.NoError(t, err)
			require.Equal(t, c.expArn, parsed.CanonicalArn())
		})
	}
}
