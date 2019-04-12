package mixpanel

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Track(t *testing.T) {
	viper.SetEnvPrefix("mixpanel")
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	token := viper.Get("token").(string)
	if token == "" {
		t.Fatalf("Mixpanel Token is not provided for the test. You can add MIXPANEL_TOKEN to your environment variable for the test.")
	}
	client := NewClient(token)

	props := map[string]string{"test": "testing"}
	event := NewEvent("go-test", props)
	event.DistinctID = "1"
	_, err := client.Track(event)
	assert.NoError(t, err)
}
