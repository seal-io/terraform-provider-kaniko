package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/seal-io/terraform-provider-kaniko/kaniko"
)

var version = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/seal-io/kaniko",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), kaniko.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
