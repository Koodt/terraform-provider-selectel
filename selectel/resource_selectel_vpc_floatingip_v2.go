package selectel

import (
	"context"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/go-selvpcclient/v2/selvpcclient/resell/v2/floatingips"
)

func resourceVPCFloatingIPV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCFloatingIPV2Create,
		ReadContext:   resourceVPCFloatingIPV2Read,
		DeleteContext: resourceVPCFloatingIPV2Delete,
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
			"port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"floating_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
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

func resourceVPCFloatingIPV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	projectID := d.Get("project_id").(string)
	opts := floatingips.FloatingIPOpts{
		FloatingIPs: []floatingips.FloatingIPOpt{
			{
				Region:   d.Get("region").(string),
				Quantity: 1,
			},
		},
	}

	log.Print(msgCreate(objectFloatingIP, opts))
	floatingIPs, _, err := floatingips.Create(ctx, resellV2Client, projectID, opts)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectFloatingIP, err))
	}
	if len(floatingIPs) != 1 {
		return diag.FromErr(errReadFromResponse(objectFloatingIP))
	}

	d.SetId(floatingIPs[0].ID)

	return resourceVPCFloatingIPV2Read(ctx, d, meta)
}

func resourceVPCFloatingIPV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	log.Print(msgGet(objectFloatingIP, d.Id()))
	floatingIP, response, err := floatingips.Get(ctx, resellV2Client, d.Id())
	if err != nil {
		if response != nil {
			if response.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}

		return diag.FromErr(errGettingObject(objectFloatingIP, d.Id(), err))
	}

	d.Set("fixed_ip_address", floatingIP.FixedIPAddress)
	d.Set("floating_ip_address", floatingIP.FloatingIPAddress)
	d.Set("port_id", floatingIP.PortID)
	d.Set("project_id", floatingIP.ProjectID)
	d.Set("region", floatingIP.Region)
	d.Set("status", floatingIP.Status)

	associatedServers := serversMapsFromStructs(floatingIP.Servers)
	if err := d.Set("servers", associatedServers); err != nil {
		log.Print(errSettingComplexAttr("servers", err))
	}

	return nil
}

func resourceVPCFloatingIPV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()

	log.Print(msgDelete(objectFloatingIP, d.Id()))
	response, err := floatingips.Delete(ctx, resellV2Client, d.Id())
	if err != nil {
		if response != nil {
			if response.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}

		return diag.FromErr(errDeletingObject(objectFloatingIP, d.Id(), err))
	}

	return nil
}
