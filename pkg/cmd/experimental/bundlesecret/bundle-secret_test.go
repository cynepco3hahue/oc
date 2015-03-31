package bundlesecret

import (
	"testing"
)

func TestValidate(t *testing.T) {

	tests := []struct {
		testName string
		args     []string
		expErr   bool
	}{
		{
			testName: "validArgs",
			args:     []string{"testSecret", "./bsFixtures/www.google.com"},
		},
		{
			testName: "noName",
			args:     []string{"./bsFixtures/www.google.com"},
			expErr:   true, //"Secret name is required"
		},
		{
			testName: "noFilesPassed",
			args:     []string{"testSecret"},
			expErr:   true, //"At least one source file or directory must be specified"
		},
	}

	for _, test := range tests {
		options := NewDefaultOptions()
		options.Complete(test.args)
		err := options.Validate()
		if err != nil && !test.expErr {
			t.Errorf("%s: unexpected error: %v", test.testName, err)
		}
	}
}

func TestCreateSecret(t *testing.T) {

	tests := []struct {
		testName string
		args     []string
		expErr   bool
		quiet    bool
	}{
		{
			testName: "validSources",
			args:     []string{"testSecret", "./bsFixtures/www.google.com", "./bsFixtures/dirNoSubdir"},
		},
		{
			testName: "invalidDNS",
			args:     []string{"testSecret", "./bsFixtures/invalid/invalid-DNS"},
			expErr:   true, // "/bsFixtures/invalid-DNS cannot be used as a key in a secret"
		},
		{
			testName: "filesSameName",
			args:     []string{"testSecret", "./bsFixtures/www.google.com", "./bsFixtures/multiple/www.google.com"},
			expErr:   true, // "Multiple files with the same name (www.google.com) cannot be included a secret"
		},
		{
			testName: "testQuietTrue",
			args:     []string{"testSecret", "./bsFixtures/dir"},
			quiet:    true,
		},
		{
			testName: "testQuietFalse",
			args:     []string{"testSecret", "./bsFixtures/dir"},
			expErr:   true, // "Skipping resource <resource path>"
		},
	}
	for _, test := range tests {
		options := NewDefaultOptions()
		options.Complete(test.args)
		options.Quiet = test.quiet

		err := options.Validate()
		if err != nil {
			t.Errorf("unexpected error")
		}
		_, err = options.CreateSecret()
		if err != nil && !test.expErr {
			t.Errorf("%s: unexpected error: %s", test.testName, err)
		}
	}
}
