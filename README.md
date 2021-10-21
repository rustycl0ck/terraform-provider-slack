# Terraform Provider for Slack

This is a [Terraform](https://www.terraform.io/) provider for [Slack](https://slack.com)

## Maintainers

This provider plugin is maintained by:

* [@KeisukeYamashita](https;//github.com/KeisukeYamashita)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= v0.12.0 (v0.11.x may work but not supported actively)

## Limitations

**I do not have any Plus or Enterprise Grid workspace which I'm free to use, unfortunately.**

That's why several resources, e.g. a slack user, have not been supported yet. 

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/KeisukeYamashita/terraform-provider-slack`

```console
$ mkdir -p $GOPATH/src/github.com/KeisukeYamashita; cd $GOPATH/src/github.com/KeisukeYamashita
$ git clone git@github.com:KeisukeYamashita/terraform-provider-slack
Enter the provider directory and build the provider

$ cd $GOPATH/src/github.com/KeisukeYamashita/terraform-provider-slack
$ make build
```

## Credits

This was forked from [jmatsu/terraform-provider-slack](https://github.com/jmatsu/terraform-provider-slack), thank you very much for open sourcing this project!
