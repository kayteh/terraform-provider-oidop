# Oidop Terraform Provider

This is a not the official Podio provider. This provider follows a very small set of goals, please do not use it in production.

It is tied to an unofficial podio-go SDK as well, maintained by only myself.

## API Key "Trust" Levels

A major caveat to this whole system is the use of "trust" levels to moderate API access. This terraform provider attempts to make certain operations that require higher trust break out gracefully. The default behavior is to not do this.

You need a trust level of 2 (defaulting to 0) to use this to it's full potential. Contact Podio support and let them know you are using this tool and you can link them to this documentation as to why.

This Terraform provider needs to:

- Create, modify, and delete spaces
- Invite and remove members from spaces

## Current Target: MMF1: Kanban

- [x] Workspace creation
- [x] App creation
- [x] Template/Field creation (Text + Category)
- [ ] Icon search picker
- [ ] Item creation

## Building and testing locally

You need to have the Podio-Go SDK (https://github.com/kayteh/podio-go) cloned in the same directory as the directory you are housing this repo in. So your directory structure should look like:
```
some_folder/
|- terraform-provider-podio/
|- podio-go/
| ...
```

### Build the Provider:

On Windows:
```
go mod tidy
go build -o podio-terraform-provider.exe
```

On Mac/Linux:
```
go mod tidy
go build -o podio-terraform-provider
```

Please refer the section titled "Build the provider" on the following page for details about how to override Terraform's plugin lookup mechanism to make it load your provider binary upon `terraform init`: https://learn.hashicorp.com/tutorials/terraform/plugin-framework-create?in=terraform/providers#build-the-provider

## Usage

Please refer the provider docs at in the [docs/](./docs/) directory

There are also some good examples in the [examples/](./examples/) directory