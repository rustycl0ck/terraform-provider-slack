package slack

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/slack-go/slack"
)

const (
	conversationActionOnDestroyNone    = "none"
	conversationActionOnDestroyArchive = "archive"
)

var validateConversationActionOnDestroyValue = validation.StringInSlice([]string{
	conversationActionOnDestroyNone,
	conversationActionOnDestroyArchive,
}, false)

func resourceSlackConversation() *schema.Resource {
	return &schema.Resource{
		Description: "Channel-like things encountered in Slack",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"purpose": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_archived": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"action_on_destroy": {
				Type:         schema.TypeString,
				Description:  "Either of none or archive",
				Required:     true,
				ValidateFunc: validateConversationActionOnDestroyValue,
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
		ReadContext:   resourceSlackConversationRead,
		CreateContext: resourceSlackConversationCreate,
		UpdateContext: resourceSlackConversationUpdate,
		DeleteContext: resourceSlackConversationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func configureSlackConversation(d *schema.ResourceData, channel *slack.Channel) error {
	if err := d.Set("name", channel.Name); err != nil {
		return err
	}
	if err := d.Set("topic", channel.Topic.Value); err != nil {
		return err
	}
	if err := d.Set("purpose", channel.Purpose.Value); err != nil {
		return err
	}
	if err := d.Set("is_archived", channel.IsArchived); err != nil {
		return err
	}
	if err := d.Set("is_shared", channel.IsShared); err != nil {
		return err
	}
	if err := d.Set("is_ext_shared", channel.IsExtShared); err != nil {
		return err
	}
	if err := d.Set("is_org_shared", channel.IsOrgShared); err != nil {
		return err
	}
	if err := d.Set("created", channel.Created); err != nil {
		return err
	}
	if err := d.Set("creator", channel.Creator); err != nil {
		return err
	}

	// Required
	if err := d.Set("is_private", channel.IsPrivate); err != nil {
		return err
	}

	// Never support
	//_ = d.Set("members", channel.Members)
	//_ = d.Set("num_members", channel.NumMembers)
	//_ = d.Set("unread_count", channel.UnreadCount)
	//_ = d.Set("unread_count_display", channel.UnreadCountDisplay)
	//_ = d.Set("last_read", channel.Name)
	//_ = d.Set("latest", channel.Name)
	return nil
}

func resourceSlackConversationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*slack.Client)

	name := d.Get("name").(string)
	isPrivate := d.Get("is_private").(bool)

	log.Printf("[DEBUG] Creating Conversation: %s", name)
	channel, err := client.CreateConversationContext(ctx, name, isPrivate)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(channel.ID)
	return resourceSlackConversationRead(ctx, d, meta)
}

func resourceSlackConversationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*slack.Client)

	id := d.Id()

	log.Printf("[DEBUG] Reading Conversation: %s", d.Id())
	channel, err := client.GetConversationInfoContext(ctx, id, false)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := configureSlackConversation(d, channel); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSlackConversationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*slack.Client)
	id := d.Id()

	if _, err := client.RenameConversationContext(ctx, id, d.Get("name").(string)); err != nil {
		return diag.FromErr(err)
	}

	if topic, ok := d.GetOk("topic"); ok {
		if _, err := client.SetTopicOfConversationContext(ctx, id, topic.(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	if purpose, ok := d.GetOk("purpose"); ok {
		if _, err := client.SetPurposeOfConversationContext(ctx, id, purpose.(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	if isArchived, ok := d.GetOkExists("is_archived"); ok {
		if isArchived.(bool) {
			if err := client.ArchiveConversationContext(ctx, id); err != nil {
				if err.Error() != "already_archived" {
					return diag.FromErr(err)
				}
			}
		} else {
			if err := client.UnArchiveConversationContext(ctx, id); err != nil {
				if err.Error() != "not_archived" {
					return diag.FromErr(err)
				}
			}
		}
	}

	return resourceSlackConversationRead(ctx, d, meta)
}

func resourceSlackConversationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*slack.Client)
	id := d.Id()

	action := d.Get("action_on_destroy").(string)

	switch action {
	case conversationActionOnDestroyNone:
		log.Printf("[DEBUG] Do nothing on Conversation: %s (%s)", id, d.Get("name"))
	case conversationActionOnDestroyArchive:
		log.Printf("[DEBUG] Deleting(archive) Conversation: %s (%s)", id, d.Get("name"))
		if err := client.ArchiveConversationContext(ctx, id); err != nil {
			if err.Error() != "already_archived" {
				return diag.FromErr(err)
			}
		}
	default:
		return diag.FromErr(fmt.Errorf("unknown action was provided. (%s)", action))
	}

	d.SetId("")

	return nil
}
