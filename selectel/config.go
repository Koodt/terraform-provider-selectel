package selectel

import (
	"errors"
	"strings"

	"github.com/gophercloud/utils/terraform/auth"
	"github.com/hashicorp/go-retryablehttp"
	domainsV1 "github.com/selectel/domains-go/pkg/v1"
	"github.com/selectel/go-selvpcclient/selvpcclient"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell"
	resellV2 "github.com/selectel/go-selvpcclient/selvpcclient/resell/v2"
)

// Config contains all available configuration options.
type Config struct {
	gophercloud_auth auth.Config
	selectel_token   string
}

// Validate performs config validation.
func (c *Config) Validate() error {
	if c.selectel_token == "" {
		return errors.New("token must be specified")
	}
	if c.gophercloud_auth.IdentityEndpoint == "" {
		c.gophercloud_auth.IdentityEndpoint = strings.Join([]string{resell.Endpoint, resellV2.APIVersion}, "/")
	}
	if c.gophercloud_auth.Region != "" {
		if err := validateRegion(c.gophercloud_auth.Region); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) resellV2Client() *selvpcclient.ServiceClient {
	return resellV2.NewV2ResellClientWithEndpoint(c.selectel_token, c.gophercloud_auth.IdentityEndpoint)
}

func (c *Config) domainsV1Client() *domainsV1.ServiceClient {
	domainsClient := domainsV1.NewDomainsClientV1WithDefaultEndpoint(c.selectel_token)
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil // Ignore retyablehttp client logs
	retryClient.RetryWaitMin = domainsV1DefaultRetryWaitMin
	retryClient.RetryWaitMax = domainsV1DefaultRetryWaitMax
	retryClient.RetryMax = domainsV1DefaultRetry
	domainsClient.HTTPClient = retryClient.StandardClient()

	return domainsClient
}
