package slack

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/slack-go/slack"
)

func dataSourceConversation() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"channel_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"purpose": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_archived": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_ext_shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_org_shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: dataSlackConversationRead,
	}
}

func dataSlackConversationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*slack.Client)

	channelId := d.Get("channel_id").(string)
	channel, err := client.GetConversationInfoContext(ctx, channelId, false)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(channel.ID)
	if err := d.Set("name", channel.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("topic", channel.Topic.Value); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("purpose", channel.Purpose.Value); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_private", channel.IsPrivate); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_archived", channel.IsArchived); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_shared", channel.IsShared); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_ext_shared", channel.IsExtShared); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_org_shared", channel.IsOrgShared); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created", channel.Created); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("creator", channel.Creator); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
