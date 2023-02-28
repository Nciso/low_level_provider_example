package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/Nciso/low_level_provider_example/internal/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	debugFlag := flag.Bool("debug", false, "Start provider in debug mode.")
	flag.Parse()

	ctx := context.Background()
	providers := []func() tfprotov5.ProviderServer{
		// Example terraform-plugin-sdk/v2 providers
		provider.PluginProviderServer,
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var serveOpts []tf5server.ServeOpt

	if *debugFlag {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve(
		"company.io/namespace/provider",
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
