package acctest

import (
	"context"
	"os"
	"testing"

	"github.com/Nciso/low_level_provider_example/internal/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func CreateMuxFactories(experimental bool) map[string]func() (tfprotov5.ProviderServer, error) {
	prov := map[string]func() (tfprotov5.ProviderServer, error){
		"provider": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov5.ProviderServer{
				// Example terraform-plugin-sdk/v2 providers
				provider.PluginProviderServer,
			}

			muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
	return prov
}

func TestAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.

	os.Setenv("endpoint", "someendpoint.com")
	os.Setenv("token", "sometoken")
	os.Setenv("TF_LOG", "DEBUG")
}
