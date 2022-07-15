package selectel

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell/v2/tokens"
)

func authBySelectelTokenSelectelToken() {
}

func authBySelectelTokenKeystoneToken(ctx context.Context, d *schema.ResourceData, meta interface{}) (*tokens.Token, error) {
	config := meta.(*Config)
	resellV2Client := config.resellV2Client()
	tokenOpts := tokens.TokenOpts{
		ProjectID: d.Get("project_id").(string),
	}

	log.Print(msgCreate(objectToken, tokenOpts))
	token, _, err := tokens.Create(ctx, resellV2Client, tokenOpts)
	if err != nil {
		return nil, err
	}
	return token, err
}

func authByCredentialsKeystoneToken() {
}
