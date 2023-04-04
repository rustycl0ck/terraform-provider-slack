package slack

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConversation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSlackConversationRead,

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
	client := meta.(*Team).client
	conversationId := d.Get("channel_id").(string)

	logger := meta.(*Team).logger.withTags(map[string]interface{}{
		"data":            "conversation",
		"conversation_id": conversationId,
	})

	logger.trace(ctx, "Start reading a conversation")

	channel, err := client.GetConversationInfoContext(ctx, conversationId, false)

	if err != nil {
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Slack provider couldn't read conversation %s due to *%s*", conversationId, err.Error()),
				Detail:   fmt.Sprintf("Please refer to %s for the details.", "https://api.slack.com/methods/conversations.info"),
			},
		}
	} else {
		logger.trace(ctx, "Got a response from Slack api")
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

	logger.debug(ctx, "Conversation #%s (isArchived = %t)", d.Get("name").(string), d.Get("is_archived").(bool))

	return nil
}
