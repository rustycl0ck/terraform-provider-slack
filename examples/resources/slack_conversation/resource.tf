resource "slack_conversation" "conversation" {
  name    = "my-channel"
  topic   = "My personal channel"
  purpose = "Daily communication"

  is_archived = flase
  is_private  = flase

  action_on_destroy = "archive"
}
