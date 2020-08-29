package main

import (
	"fmt"
	"os"

	"github.com/aaronsky/applereleaser/cmd"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {
	cmd.Execute(
		buildVersion(version, commit, date, builtBy),
		os.Exit,
		os.Args[1:],
	)
}

func buildVersion(version, commit, date, builtBy string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	return result
}

/*
var (
	keyID          = flag.String("kid", "", "key ID")
	issuerID       = flag.String("iss", "", "issuer ID")
	privateKeyPath = flag.String("privatekeypath", "", "path to a private key used to sign authorization token. You can also directly pass private key content to ASC_PRIVATE_KEY instead of this flag")
)

func main() {
	flag.Parse()

	auth, err := tokenConfig()
	if err != nil {
		log.Fatal(err)
	}
	client := asc.NewClient(auth.Client())

	// 1. select app
	appsResponse, _, err := client.Apps.ListApps(&asc.ListAppsQuery{
		FilterBundleID: []string{"com.wayfair.WayfairApp"},
	})
	if err != nil {
		log.Fatal(err)
	}
	apps := appsResponse.Data
	if len(apps) == 0 {

	}
	app := apps[0]

	// 2. list latest versions for app and get the next version string
	versionsResponse, _, err := client.Apps.ListAppStoreVersionsForApp(app.ID, &asc.ListAppStoreVersionsQuery{
		FilterAppStoreState: []string{"READY_FOR_SALE"},
		Limit:               1,
	})
	if err != nil {
		log.Fatal(err)
	}
	versions := versionsResponse.Data
	if len(versions) == 0 {

	}
	lastVersion := versions[0]
	lastVersionParsed, err := semver.New(*lastVersion.Attributes.VersionString)
	if err != nil {
		log.Fatal(err)
	}
	err = lastVersionParsed.IncrementMinor()
	if err != nil {
		log.Fatal(err)
	}

	// 3. select latest build
	buildsResponse, _, err := client.Builds.ListBuilds(&asc.ListBuildsQuery{
		FilterApp:     []string{app.ID},
		FilterVersion: []string{lastVersionParsed.FinalizeVersion()},
		Limit:         1,
	})
	if err != nil {
		log.Fatal(err)
	}
	builds := buildsResponse.Data
	if len(builds) == 0 {

	}
	latestBuild := builds[0]

	// 3. create new app store version
	newVersionResponse, _, err := client.Apps.CreateAppStoreVersion(&asc.AppStoreVersionCreateRequest{
		Attributes: asc.AppStoreVersionCreateRequestAttributes{
			Copyright:     asc.String(fmt.Sprintf("%d Wayfair LLC", time.Now().Year())),
			Platform:      asc.PlatformIOS,
			ReleaseType:   asc.String("AFTER_APPROVAL"),
			UsesIDFA:      asc.Bool(true),
			VersionString: lastVersionParsed.FinalizeVersion(),
		},
		Relationships: asc.AppStoreVersionCreateRequestRelationships{
			App: struct {
				Data asc.RelationshipsData "json:\"data\""
			}{
				Data: asc.RelationshipsData{
					ID:   app.ID,
					Type: "apps",
				},
			},
			Build: &struct {
				Data *asc.RelationshipsData "json:\"data,omitempty\""
			}{
				Data: &asc.RelationshipsData{
					ID:   latestBuild.ID,
					Type: "builds",
				},
			},
		},
		Type: "appStoreVersions",
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = newVersionResponse.Data

	// 4. set promotional text

	// 5. set app review info
	// 6. set phased release info
	// 7. set idfa declaration
	// 8. submit for review
}

func tokenConfig() (auth *asc.AuthTransport, err error) {
	if *keyID == "" {
		return nil, errors.New("no key ID provided to the -kid flag")
	} else if *issuerID == "" {
		return nil, errors.New("no issuer ID provided to the -iss flag")
	}

	var secret []byte
	if *privateKeyPath != "" {
		// Read private key file as []byte
		secret, err = ioutil.ReadFile(*privateKeyPath)
		if err != nil {
			return nil, err
		}
	} else if content, ok := os.LookupEnv("ASC_PRIVATE_KEY"); ok {
		secret = []byte(content)
	} else {
		return nil, errors.New("no private key provided to either the -privatekeypath flag or ASC_PRIVATE_KEY environment variable")
	}
	expiryDuration := 20 * time.Minute
	return asc.NewTokenConfig(*keyID, *issuerID, expiryDuration, secret)
}
*/
