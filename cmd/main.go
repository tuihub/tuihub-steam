package main

import (
	"context"
	"os"

	porter "github.com/tuihub/protos/pkg/librarian/porter/v1"
	librarian "github.com/tuihub/protos/pkg/librarian/v1"
	"github.com/tuihub/tuihub-go"
	"github.com/tuihub/tuihub-go/logger"
	"github.com/tuihub/tuihub-steam/internal"
)

// go build -ldflags "-X main.version=x.y.z".
var (
	// version is the version of the compiled software.
	version string
)

func main() {
	apiKey, exist := os.LookupEnv("STEAM_API_KEY")
	if !exist || apiKey == "" {
		panic("STEAM_API_KEY is required")
	}
	config := tuihub.PorterConfig{
		Name:       "tuihub-steam",
		Version:    version,
		GlobalName: "github.com/tuihub/tuihub-steam",
		FeatureSummary: &porter.PorterFeatureSummary{
			SupportedAccounts: []*porter.PorterFeatureSummary_Account{
				{
					Platform: tuihub.WellKnownToString(
						librarian.WellKnownAccountPlatform_WELL_KNOWN_ACCOUNT_PLATFORM_STEAM,
					),
					AppRelationTypes: []librarian.AccountAppRelationType{
						librarian.AccountAppRelationType_ACCOUNT_APP_RELATION_TYPE_OWN,
					},
				},
			},
			SupportedAppInfoSources: []string{
				tuihub.WellKnownToString(librarian.WellKnownAppInfoSource_WELL_KNOWN_APP_INFO_SOURCE_STEAM),
			},
			SupportedFeedSources:        nil,
			SupportedNotifyDestinations: nil,
		},
	}
	server, err := tuihub.NewPorter(
		context.Background(),
		config,
		internal.NewHandler(apiKey),
	)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	if err = server.Run(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
