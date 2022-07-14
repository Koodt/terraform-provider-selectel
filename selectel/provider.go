package selectel

import (
	"context"

	"github.com/gophercloud/utils/terraform/auth"
	"github.com/gophercloud/utils/terraform/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
)

const (
	objectFloatingIP              = "floating IP"
	objectKeypair                 = "keypair"
	objectLicense                 = "license"
	objectProject                 = "project"
	objectProjectQuotas           = "quotas for project"
	objectRole                    = "role"
	objectSubnet                  = "subnet"
	objectToken                   = "token"
	objectUser                    = "user"
	objectCluster                 = "cluster"
	objectKubeConfig              = "kubeconfig"
	objectKubeVersions            = "kube-versions"
	objectNodegroup               = "nodegroup"
	objectDomain                  = "domain"
	objectRecord                  = "record"
	objectDatastore               = "datastore"
	objectDatabase                = "database"
	objectGrant                   = "grant"
	objectExtension               = "extension"
	objectDatastoreTypes          = "datastore-types"
	objectAvailableExtensions     = "available-extensions"
	objectFlavors                 = "flavors"
	objectConfigurationParameters = "configuration-parameters"
	objectPrometheusMetricToken   = "prometheus-metric-token"
	objectFeatureGates            = "feature-gates"
	objectAdmissionControllers    = "admission-controllers"
)

// This is a global MutexKV for use within this plugin.
var selMutexKV = mutexkv.NewMutexKV()

var descriptions = map[string]string{
	"auth_url": "The Identity authentication URL.",

	"cloud": "An entry in a `clouds.yaml` file to use.",

	"region": "The OpenStack region to connect to.",

	"user_name": "Username to login with.",

	"user_id": "User ID to login with.",

	"application_credential_id": "Application Credential ID to login with.",

	"application_credential_name": "Application Credential name to login with.",

	"application_credential_secret": "Application Credential secret to login with.",

	"tenant_id": "The ID of the Tenant (Identity v2) or Project (Identity v3)\n" +
		"to login with.",

	"tenant_name": "The name of the Tenant (Identity v2) or Project (Identity v3)\n" +
		"to login with.",

	"password": "Password to login with.",

	"token": "Authentication token to use as an alternative to username/password.",

	"user_domain_name": "The name of the domain where the user resides (Identity v3).",

	"user_domain_id": "The ID of the domain where the user resides (Identity v3).",

	"project_domain_name": "The name of the domain where the project resides (Identity v3).",

	"project_domain_id": "The ID of the domain where the proejct resides (Identity v3).",

	"domain_id": "The ID of the Domain to scope to (Identity v3).",

	"domain_name": "The name of the Domain to scope to (Identity v3).",

	"default_domain": "The name of the Domain ID to scope to if no other domain is specified. Defaults to `default` (Identity v3).",

	"insecure": "Trust self-signed certificates.",

	"cacert_file": "A Custom CA certificate.",

	"cert": "A client certificate to authenticate with.",

	"key": "A client private key to authenticate with.",

	"endpoint_type": "The catalog endpoint type to use.",

	"endpoint_overrides": "A map of services with an endpoint to override what was\n" +
		"from the Keystone catalog",

	"swauth": "Use Swift's authentication system instead of Keystone. Only used for\n" +
		"interaction with Swift.",

	"use_octavia": "If set to `true`, API requests will go the Load Balancer\n" +
		"service (Octavia) instead of the Networking service (Neutron).",

	"disable_no_cache_header": "If set to `true`, the HTTP `Cache-Control: no-cache` header will not be added by default to all API requests.",

	"delayed_auth": "If set to `false`, OpenStack authorization will be perfomed,\n" +
		"every time the service provider client is called. Defaults to `true`.",

	"allow_reauth": "If set to `false`, OpenStack authorization won't be perfomed\n" +
		"automatically, if the initial auth token get expired. Defaults to `true`",

	"max_retries": "How many times HTTP connection should be retried until giving up.",
}

// Provider returns the Selectel terraform provider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"selectel_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SEL_TOKEN", nil),
				Description: "Token to authorize with the Selectel API.",
			},
			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", ""),
				Description: descriptions["auth_url"],
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["region"],
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION_NAME", ""),
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USERNAME", ""),
				Description: descriptions["user_name"],
			},

			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_ID", ""),
				Description: descriptions["user_name"],
			},

			"application_credential_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_APPLICATION_CREDENTIAL_ID", ""),
				Description: descriptions["application_credential_id"],
			},

			"application_credential_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_APPLICATION_CREDENTIAL_NAME", ""),
				Description: descriptions["application_credential_name"],
			},

			"application_credential_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_APPLICATION_CREDENTIAL_SECRET", ""),
				Description: descriptions["application_credential_secret"],
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_ID",
					"OS_PROJECT_ID",
				}, ""),
				Description: descriptions["tenant_id"],
			},

			"tenant_name": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_NAME",
					"OS_PROJECT_NAME",
				}, ""),
				Description: descriptions["tenant_name"],
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSWORD", ""),
				Description: descriptions["password"],
			},

			"token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TOKEN",
					"OS_AUTH_TOKEN",
				}, ""),
				Description: descriptions["token"],
			},

			"user_domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_DOMAIN_NAME", ""),
				Description: descriptions["user_domain_name"],
			},

			"user_domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_DOMAIN_ID", ""),
				Description: descriptions["user_domain_id"],
			},

			"project_domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PROJECT_DOMAIN_NAME", ""),
				Description: descriptions["project_domain_name"],
			},

			"project_domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PROJECT_DOMAIN_ID", ""),
				Description: descriptions["project_domain_id"],
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DOMAIN_ID", ""),
				Description: descriptions["domain_id"],
			},

			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DOMAIN_NAME", ""),
				Description: descriptions["domain_name"],
			},

			"default_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DEFAULT_DOMAIN", "default"),
				Description: descriptions["default_domain"],
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_INSECURE", nil),
				Description: descriptions["insecure"],
			},

			"endpoint_type": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ENDPOINT_TYPE", ""),
			},

			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CACERT", ""),
				Description: descriptions["cacert_file"],
			},

			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: descriptions["cert"],
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: descriptions["key"],
			},

			"swauth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SWAUTH", false),
				Description: descriptions["swauth"],
			},

			"use_octavia": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USE_OCTAVIA", false),
				Description: descriptions["use_octavia"],
			},

			"delayed_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DELAYED_AUTH", true),
				Description: descriptions["delayed_auth"],
			},

			"allow_reauth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ALLOW_REAUTH", true),
				Description: descriptions["allow_reauth"],
			},

			"cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CLOUD", ""),
				Description: descriptions["cloud"],
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: descriptions["max_retries"],
			},

			"endpoint_overrides": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: descriptions["endpoint_overrides"],
			},

			"disable_no_cache_header": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["disable_no_cache_header"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"selectel_domains_domain_v1":                dataSourceDomainsDomainV1(),
			"selectel_dbaas_datastore_type_v1":          dataSourceDBaaSDatastoreTypeV1(),
			"selectel_dbaas_available_extension_v1":     dataSourceDBaaSAvailableExtensionV1(),
			"selectel_dbaas_flavor_v1":                  dataSourceDBaaSFlavorV1(),
			"selectel_dbaas_configuration_parameter_v1": dataSourceDBaaSConfigurationParameterV1(),
			"selectel_dbaas_prometheus_metric_token_v1": dataSourceDBaaSPrometheusMetricTokenV1(),
			"selectel_mks_kubeconfig_v1":                dataSourceMKSKubeconfigV1(),
			"selectel_mks_kube_versions_v1":             dataSourceMKSKubeVersionsV1(),
			"selectel_mks_feature_gates_v1":             dataSourceMKSFeatureGatesV1(),
			"selectel_mks_admission_controllers_v1":     dataSourceMKSAdmissionControllersV1(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"selectel_vpc_floatingip_v2":                resourceVPCFloatingIPV2(),
			"selectel_vpc_keypair_v2":                   resourceVPCKeypairV2(),
			"selectel_vpc_license_v2":                   resourceVPCLicenseV2(),
			"selectel_vpc_project_v2":                   resourceVPCProjectV2(),
			"selectel_vpc_role_v2":                      resourceVPCRoleV2(),
			"selectel_vpc_subnet_v2":                    resourceVPCSubnetV2(),
			"selectel_vpc_token_v2":                     resourceVPCTokenV2(),
			"selectel_vpc_user_v2":                      resourceVPCUserV2(),
			"selectel_vpc_vrrp_subnet_v2":               resourceVPCVRRPSubnetV2(),        // DEPRECATED
			"selectel_vpc_crossregion_subnet_v2":        resourceVPCCrossRegionSubnetV2(), // DEPRECATED
			"selectel_mks_cluster_v1":                   resourceMKSClusterV1(),
			"selectel_mks_nodegroup_v1":                 resourceMKSNodegroupV1(),
			"selectel_domains_domain_v1":                resourceDomainsDomainV1(),
			"selectel_domains_record_v1":                resourceDomainsRecordV1(),
			"selectel_dbaas_datastore_v1":               resourceDBaaSDatastoreV1(),
			"selectel_dbaas_user_v1":                    resourceDBaaSUserV1(),
			"selectel_dbaas_database_v1":                resourceDBaaSDatabaseV1(),
			"selectel_dbaas_grant_v1":                   resourceDBaaSGrantV1(),
			"selectel_dbaas_extension_v1":               resourceDBaaSExtensionV1(),
			"selectel_dbaas_prometheus_metric_token_v1": resourceDBaaSPrometheusMetricTokenV1(),
		},
	}
	provider.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return configureProvider(d, terraformVersion)
	}

	return provider
}

func configureProvider(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	config := Config{
		gophercloud_auth: auth.Config{
			CACertFile:                  d.Get("cacert_file").(string),
			ClientCertFile:              d.Get("cert").(string),
			ClientKeyFile:               d.Get("key").(string),
			Cloud:                       d.Get("cloud").(string),
			DefaultDomain:               d.Get("default_domain").(string),
			DomainID:                    d.Get("domain_id").(string),
			DomainName:                  d.Get("domain_name").(string),
			EndpointOverrides:           d.Get("endpoint_overrides").(map[string]interface{}),
			EndpointType:                d.Get("endpoint_type").(string),
			IdentityEndpoint:            d.Get("auth_url").(string),
			Password:                    d.Get("password").(string),
			ProjectDomainID:             d.Get("project_domain_id").(string),
			ProjectDomainName:           d.Get("project_domain_name").(string),
			Region:                      d.Get("region").(string),
			Swauth:                      d.Get("swauth").(bool),
			Token:                       d.Get("token").(string),
			TenantID:                    d.Get("tenant_id").(string),
			TenantName:                  d.Get("tenant_name").(string),
			UserDomainID:                d.Get("user_domain_id").(string),
			UserDomainName:              d.Get("user_domain_name").(string),
			Username:                    d.Get("user_name").(string),
			UserID:                      d.Get("user_id").(string),
			ApplicationCredentialID:     d.Get("application_credential_id").(string),
			ApplicationCredentialName:   d.Get("application_credential_name").(string),
			ApplicationCredentialSecret: d.Get("application_credential_secret").(string),
			UseOctavia:                  d.Get("use_octavia").(bool),
			DelayedAuth:                 d.Get("delayed_auth").(bool),
			AllowReauth:                 d.Get("allow_reauth").(bool),
			MaxRetries:                  d.Get("max_retries").(int),
			DisableNoCacheHeader:        d.Get("disable_no_cache_header").(bool),
			TerraformVersion:            terraformVersion,
			SDKVersion:                  meta.SDKVersionString(),
			MutexKV:                     mutexkv.NewMutexKV(),
		},
		selectel_token: d.Get("selectel_token").(string),
	}
	if v, ok := d.GetOk("project_id"); ok {
		config.gophercloud_auth.TenantID = v.(string)
	}
	if v, ok := d.GetOk("selectel_region"); ok {
		config.gophercloud_auth.Region = v.(string)
	}
	if err := config.Validate(); err != nil {
		return nil, diag.FromErr(err)
	}

	return &config, nil
}
