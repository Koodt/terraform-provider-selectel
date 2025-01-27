package selectel

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/go-selvpcclient/v2/selvpcclient/resell/v2/licenses"
)

func resourceVPCLicenseV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCLicenseV2Create,
		ReadContext:   resourceVPCLicenseV2Read,
		DeleteContext: resourceVPCLicenseV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      hashServers,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceVPCLicenseV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	projectID := d.Get("project_id").(string)
	region := d.Get("region").(string)
	licenseType := d.Get("type").(string)
	opts := licenses.LicenseOpts{
		Licenses: []licenses.LicenseOpt{
			{
				Region:   region,
				Type:     licenseType,
				Quantity: 1,
			},
		},
	}

	log.Print(msgCreate(objectLicense, opts))
	newLicenses, _, err := licenses.Create(ctx, resellV2Client, projectID, opts)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectLicense, err))
	}

	if len(newLicenses) != 1 {
		return diag.FromErr(errReadFromResponse(objectLicense))
	}

	d.SetId(strconv.Itoa(newLicenses[0].ID))

	return resourceVPCLicenseV2Read(ctx, d, meta)
}

func resourceVPCLicenseV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	log.Print(msgGet(objectLicense, d.Id()))
	license, response, err := licenses.Get(ctx, resellV2Client, d.Id())
	if err != nil {
		if response != nil {
			if response.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}

		return diag.FromErr(errGettingObject(objectLicense, d.Id(), err))
	}

	d.Set("project_id", license.ProjectID)
	d.Set("region", license.Region)
	d.Set("status", license.Status)
	d.Set("network_id", license.NetworkID)
	d.Set("subnet_id", license.SubnetID)
	d.Set("port_id", license.PortID)
	d.Set("type", license.Type)
	associatedServers := serversMapsFromStructs(license.Servers)
	if err := d.Set("servers", associatedServers); err != nil {
		log.Print(errSettingComplexAttr("servers", err))
	}

	return nil
}

func resourceVPCLicenseV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	log.Print(msgDelete(objectLicense, d.Id()))
	response, err := licenses.Delete(ctx, resellV2Client, d.Id())
	if err != nil {
		if response != nil {
			if response.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}

		return diag.FromErr(errDeletingObject(objectLicense, d.Id(), err))
	}

	return nil
}
