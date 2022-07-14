package selectel

import (
	"testing"

	"github.com/gophercloud/utils/terraform/auth"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	config := &Config{
		selectel_token: "secret",
		gophercloud_auth: auth.Config{
			Region: "ru-3",
		},
	}

	err := config.Validate()

	assert.NoError(t, err)
}

func TestValidateNoToken(t *testing.T) {
	config := &Config{}

	expected := "token must be specified"

	actual := config.Validate()

	assert.EqualError(t, actual, expected)
}

func TestValidateErrRegion(t *testing.T) {
	config := &Config{
		selectel_token: "secret",
		gophercloud_auth: auth.Config{
			Region: "unknown region",
		},
	}

	expected := "region is invalid: unknown region"

	actual := config.Validate()

	assert.EqualError(t, actual, expected)
}
