package slack

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SLACK_TOKEN", nil),
				Description: "The OAuth token used to connect to Slack. This can also be set via the SLACK_TOKEN environment variable.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"slack_conversation": dataSourceConversation(),
			"slack_user":         dataSourceSlackUser(),
			"slack_usergroup":    dataSourceUserGroup(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"slack_conversation":       resourceSlackConversation(),
			"slack_usergroup":          resourceSlackUserGroup(),
			"slack_usergroup_channels": resourceSlackUserGroupChannels(),
			"slack_usergroup_members":  resourceSlackUserGroupMembers(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	if token == "" {
		return nil, diag.FromErr(errors.New("token must be configured"))
	}

	return slack.New(token), nil
}
