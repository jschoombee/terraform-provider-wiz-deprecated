package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"shell.com/terraform-provider-wiz/provider"
)

//go:generate terraform fmt -recursive ./examples/

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"

	// goreleaser can also pass the specific commit if you want
	// commit  string = ""
)

func main() {
	opts := tfsdk.ServeOpts{
		// TODO: Update this string with the published name of your provider.
		Name: "registry.terraform.io/shell/wiz",
	}

	err := tfsdk.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
