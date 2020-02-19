package models

import (
	"regexp"
	"testing"

	"github.com/monitoror/monitoror/pkg/validator"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestHTTPParams_IsValid(t *testing.T) {
	for _, testcase := range []struct {
		params   validator.SimpleValidator
		expected bool
	}{
		{&HTTPStatusParams{}, false},
		{&HTTPStatusParams{URL: "url"}, true},
		{&HTTPStatusParams{URL: "url", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HTTPStatusParams{URL: "url", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},

		{&HTTPRawParams{}, false},
		{&HTTPRawParams{URL: "url"}, true},
		{&HTTPRawParams{URL: "url", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HTTPRawParams{URL: "url", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},
		{&HTTPRawParams{URL: "url", Regex: "("}, false},
		{&HTTPRawParams{URL: "url", Regex: "(.*)"}, true},

		{&HTTPFormattedParams{}, false},
		{&HTTPFormattedParams{URL: "url"}, false},
		{&HTTPFormattedParams{URL: "url", Format: "unknown"}, false},
		{&HTTPFormattedParams{URL: "url", Format: "JSON", Key: ""}, false},
		{&HTTPFormattedParams{URL: "url", Format: "JSON", Key: "."}, false},
		{&HTTPFormattedParams{URL: "url", Format: "JSON", Key: "key"}, true},
		{&HTTPFormattedParams{URL: "url", Format: "JSON", Key: "key", StatusCodeMin: pointer.ToInt(300), StatusCodeMax: pointer.ToInt(299)}, false},
		{&HTTPFormattedParams{URL: "url", Format: "JSON", Key: "key", StatusCodeMin: pointer.ToInt(299), StatusCodeMax: pointer.ToInt(300)}, true},
		{&HTTPFormattedParams{URL: "url", Format: "JSON", Key: "key", Regex: "("}, false},
		{&HTTPFormattedParams{URL: "url", Format: "JSON", Key: "key", Regex: "(.*)"}, true},

		{&HTTPProxyParams{URL: ""}, false},
		{&HTTPProxyParams{URL: "url"}, true},
	} {
		assert.Equal(t, testcase.expected, testcase.params.IsValid())
	}
}

func TestHTTPParams_GetRegex(t *testing.T) {
	for _, testcase := range []struct {
		params         RegexProvider
		expectedRegex  string
		expectedRegexp *regexp.Regexp
	}{
		{&HTTPRawParams{}, "", nil},
		{&HTTPRawParams{Regex: ""}, "", nil},
		{&HTTPRawParams{Regex: "("}, "(", nil},
		{&HTTPRawParams{Regex: "(.*)"}, "(.*)", regexp.MustCompile("(.*)")},

		{&HTTPFormattedParams{}, "", nil},
		{&HTTPFormattedParams{Regex: ""}, "", nil},
		{&HTTPFormattedParams{Regex: "("}, "(", nil},
		{&HTTPFormattedParams{Regex: "(.*)"}, "(.*)", regexp.MustCompile("(.*)")},
	} {
		assert.Equal(t, testcase.expectedRegex, testcase.params.GetRegex())
		if isValidRegex(testcase.params) {
			assert.Equal(t, testcase.expectedRegexp, testcase.params.GetRegexp())
		}
	}
}

func TestHTTPFormattedParams_FormattedDataProvider(t *testing.T) {
	for _, testcase := range []struct {
		params         FormattedDataProvider
		expectedFormat string
		expectedKey    string
	}{
		{&HTTPFormattedParams{}, "", ""},
		{&HTTPFormattedParams{Format: JSONFormat}, JSONFormat, ""},
		{&HTTPFormattedParams{Format: YAMLFormat, Key: "key"}, YAMLFormat, "key"},
		{&HTTPFormattedParams{Format: XMLFormat, Key: "key"}, XMLFormat, "key"},
	} {
		assert.Equal(t, testcase.expectedFormat, testcase.params.GetFormat())
		assert.Equal(t, testcase.expectedKey, testcase.params.GetKey())
	}
}
