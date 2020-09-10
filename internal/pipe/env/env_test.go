package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	ctx := context.New(config.Project{})
	pipe := Pipe{}
	var err error

	assert.Equal(t, "loading environment variables", pipe.String())

	err = pipe.Run(ctx)
	assert.Error(t, err)

	// Pass with ASC_KEY_ID but fail on ASC_ISSUER_ID
	err = os.Setenv("ASC_KEY_ID", "TEST")
	assert.NoError(t, err)
	err = pipe.Run(ctx)
	assert.Error(t, err)

	// Pass with ASC_ISSUER_ID but fail on ASC_PRIVATE_KEY_PATH
	err = os.Setenv("ASC_ISSUER_ID", "TEST")
	assert.NoError(t, err)
	err = pipe.Run(ctx)
	assert.Error(t, err)

	// Check ASC_PRIVATE_KEY_PATH but fail because nothing exists at that path
	err = os.Setenv("ASC_PRIVATE_KEY_PATH", "TEST")
	assert.NoError(t, err)
	err = pipe.Run(ctx)
	assert.Error(t, err)

	// Check ASC_PRIVATE_KEY_PATH but silently fail because it isn't required
	err = os.Setenv("TEST", "path/to/no/file")
	assert.NoError(t, err)
	env, err := loadEnvFromPath("TEST", false)
	assert.NoError(t, err)
	assert.Empty(t, env)

	file, err := ioutil.TempFile("", "fake_key")
	if err != nil {
		assert.FailNow(t, "temp file creation produced an error", err)
	}
	defer rmFile(file)

	// Check ASC_PRIVATE_KEY_PATH but fail because the contents of the file do not contain a real key
	err = os.Setenv("ASC_PRIVATE_KEY_PATH", file.Name())
	assert.NoError(t, err)
	err = pipe.Run(ctx)
	assert.Error(t, err)

	// This key is a mock key generated by the following command:
	//
	//   openssl ecparam -name prime256v1 -genkey -noout | openssl pkcs8 -topk8 -nocrypt -out key.pem
	//
	// This will generate the ASN.1 PKCS#8 representation of the private key needed
	// to create a valid token. If you are looking at this test to see how to make a key,
	// reference Apple's documentation on this subject instead.
	//
	// https://developer.apple.com/documentation/appstoreconnectapi/creating_api_keys_for_app_store_connect_api
	key := `
-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgELiCwZa9oGedoUR7
8Vr36M6WOkEBGZh2YsUVL0kCIJ6hRANCAASQtP/ZdZBW6UdwJeyz09ws2nx5OOUA
tra43bY9mLeVK0zrTn/3jvjTHEdD3HcRJgau1jshXG4IHXSW9yXj9x3V
-----END PRIVATE KEY-----
`
	_, err = file.WriteString(key)
	assert.NoError(t, err)

	// Pass with ASC_PRIVATE_KEY_PATH, and create the credentials object
	err = pipe.Run(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, ctx.Credentials)
}

// rmFile closes an open descriptor.
func rmFile(f *os.File) {
	err := os.Remove(f.Name())
	if err != nil {
		fmt.Println(err)
	}
}
