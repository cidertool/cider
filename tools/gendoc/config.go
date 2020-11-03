/**
Copyright (C) 2020 Aaron Sky.

This file is part of Cider, a tool for automating submission
of apps to Apple's App Stores.

Cider is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Cider is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Cider.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/internal/closer"
	"github.com/cidertool/cider/pkg/config"
	"github.com/spf13/cobra"
)

const (
	docsConfigFrontmatterTemplate = `---
layout: page
nav_order: %d
---

# Configuration
{: .no_toc }

`
	docsConfigTableOfContents = `
<details open markdown="block">
  <summary>
    Table of Contents
  </summary>
  {: .text-delta }
- TOC
{:toc}
</details>

`
	docsConfigTerminologyDisclaimer = `
- [x] An X here means the field is required.
- [ ] This field is optional and can be omitted.

`
)

// nolint: gochecknoglobals
var docsConfigExampleProject = config.Project{
	"My App": {
		BundleID:              "com.myproject.MyApp",
		PrimaryLocale:         "en-US",
		UsesThirdPartyContent: asc.Bool(false),
		Availability: &config.Availability{
			AvailableInNewTerritories: asc.Bool(false),
			Pricing: []config.PriceSchedule{
				{Tier: "0"},
			},
			Territories: []string{"USA"},
		},
		Categories: &config.Categories{
			Primary:   "SOCIAL_NETWORKING",
			Secondary: "GAMES",
			SecondarySubcategories: [2]string{
				"GAMES_SIMULATION",
				"GAMES_RACING",
			},
		},
		Localizations: config.AppLocalizations{
			"en-US": {
				Name:     "My App",
				Subtitle: "Not Your App",
			},
		},
		Versions: config.Version{
			Platform:             config.PlatformiOS,
			Copyright:            "2020 Me",
			EarliestReleaseDate:  nil,
			ReleaseType:          config.ReleaseTypeAfterApproval,
			PhasedReleaseEnabled: true,
			IDFADeclaration:      nil,
			Localizations: config.VersionLocalizations{
				"en-US": {
					Description:  "My App for cool people",
					Keywords:     "Apps, Cool, Mine",
					WhatsNewText: `Thank you for using My App! I bring you updates every week so this continues to be my app.`,
					PreviewSets: config.PreviewSets{
						config.PreviewTypeiPhone65: []config.Preview{
							{
								File: config.File{
									Path: "assets/store/iphone65/preview.mp4",
								},
							},
						},
					},
					ScreenshotSets: config.ScreenshotSets{
						config.ScreenshotTypeiPhone65: []config.File{
							{Path: "assets/store/iphone65/app.jpg"},
						},
					},
				},
			},
		},
		Testflight: config.Testflight{
			EnableAutoNotify: true,
			Localizations: config.TestflightLocalizations{
				"en-US": {
					Description: "My App for cool people using the beta",
				},
			},
		},
	},
}

func runDocsConfigCmd(cmd *cobra.Command, args []string) error {
	var path string
	if len(args) == 0 {
		path = defaultDocsPath
	} else {
		path = args[0]
	}

	path = filepath.Join(path, "configuration.md")
	log.WithField("path", path).Info("generating configuration documentation")

	err := genConfigMarkdown(path)
	if err != nil {
		log.Error("generation failed")
	} else {
		log.Info("generation completed successfully")
	}

	return err
}

func genConfigMarkdown(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer closer.Close(f)

	r, err := newRenderer()
	if err != nil {
		return err
	}

	r.Header = func() string {
		return fmt.Sprintf(docsConfigFrontmatterTemplate, 4)
	}

	r.Footer = func() string {
		contents, err := ioutil.ReadFile(filepath.Join(filepath.Dir(path), "configuration-footer.md"))
		if err != nil {
			log.Error(err.Error())
			return ""
		}

		return string(contents)
	}

	return r.Render(f)
}
