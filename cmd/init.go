package cmd

import (
	"os"

	"github.com/aaronsky/applereleaser/internal/closer"
	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/asc-go/asc"
	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const configDocString = "# This is an example applereleaser.yaml file with some sane defaults.\n"

type initCmd struct {
	cmd    *cobra.Command
	config string
}

func newInitCmd() *initCmd {
	var root = &initCmd{}
	var cmd = &cobra.Command{
		Use:           "init",
		Short:         "Generates an .applereleaser.yml file",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.OpenFile(root.config, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_EXCL, 0600)
			if err != nil {
				return err
			}
			defer closer.Close(f)

			log.Infof(color.New(color.Bold).Sprintf("Generating %s file", root.config))

			return initProject(f)
		},
	}

	cmd.Flags().StringVarP(&root.config, "config", "f", ".applereleaser.yml", "Path of configuration file to create")

	root.cmd = cmd
	return root
}

func initProject(f *os.File) error {
	var project = config.Project{
		Apps: make(map[string]config.App),
	}

	projectNamePrompt := promptui.Prompt{Label: "Project Name"}
	projectName, err := projectNamePrompt.Run()
	if err != nil {
		return err
	}
	project.Name = projectName

	var continueAppsSetup = true
	for continueAppsSetup {
		name, app, err := setupApp()
		if err != nil {
			return err
		}
		project.Apps[name] = app

		continuePrompt := promptui.Prompt{
			Label:     "Add more apps?",
			IsConfirm: true,
		}
		_, err = continuePrompt.Run()
		continueAppsSetup = err == nil
	}

	contents, err := project.String()
	if err != nil {
		return err
	}

	if _, err := f.WriteString(configDocString + contents); err != nil {
		return err
	}

	log.
		WithField("file", f.Name()).
		Info("config created; please edit accordingly to your needs")

	return nil
}

func setupApp() (string, config.App, error) {
	app := config.App{}
	var prompt promptui.Prompt
	var selec promptui.Select

	log.Info("let's set up an app in your project!")

	prompt = promptui.Prompt{Label: "App Name"}
	name, err := prompt.Run()
	if err != nil {
		return name, app, err
	}

	prompt = promptui.Prompt{Label: "Bundle ID"}
	bundleID, err := prompt.Run()
	if err != nil {
		return name, app, err
	}

	selec = promptui.Select{
		Label: "Platform",
		Items: []config.Platform{
			config.PlatformiOS,
			config.PlatformTvOS,
			config.PlatformMacOS,
		},
	}
	_, platform, err := selec.Run()
	if err != nil {
		return name, app, err
	}

	prompt = promptui.Prompt{
		Label:   "Primary Locale",
		Default: "en-US",
	}
	primaryLocale, err := prompt.Run()
	if err != nil {
		return name, app, err
	}

	prompt = promptui.Prompt{
		Label: "Price Tier",
	}
	tier, err := prompt.Run()
	if err != nil {
		return name, app, err
	}

	app.BundleID = bundleID
	app.PrimaryLocale = primaryLocale
	app.Availability = &config.Availability{
		AvailableInNewTerritories: asc.Bool(false),
		Pricing: []config.PriceSchedule{
			{Tier: tier},
		},
		Territories: []string{"US"},
	}
	app.Localizations = config.AppLocalizations{}
	app.Localizations[primaryLocale] = config.AppLocalization{
		Name: name,
	}

	app.Testflight.EnableAutoNotify = true
	app.Testflight.Localizations = config.TestflightLocalizations{}
	app.Testflight.Localizations[primaryLocale] = config.TestflightLocalization{}
	app.Testflight.ReviewDetails = &config.ReviewDetails{
		Contact: &config.ContactPerson{},
		DemoAccount: &config.DemoAccount{
			Required: false,
		},
	}

	app.Versions.Platform = config.Platform(platform)
	app.Versions.Localizations = config.VersionLocalizations{}
	app.Versions.Localizations[primaryLocale] = config.VersionLocalization{}
	app.Versions.PhasedReleaseEnabled = true
	app.Versions.ReleaseType = config.ReleaseTypeAfterApproval
	app.Versions.ReviewDetails = &config.ReviewDetails{
		Contact: &config.ContactPerson{},
		DemoAccount: &config.DemoAccount{
			Required: false,
		},
	}

	return name, app, nil
}
