package mixpanel

import (
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestClient_Track(t *testing.T) {
	os.Setenv("MIXPANEL_TOKEN", "ac347c9880fa89b7e51a7136b83c81e")
	viper.SetEnvPrefix("mixpanel")
	viper.AutomaticEnv()
	c := NewClient(viper.Get("token").(string))
	props := map[string]string{"test": "testing"}
	event := NewEvent("go-test", props)
	event.DistinctID = "1"
	result, err := c.Track(event)
	t.Logf("%v %v", result, err)
}
