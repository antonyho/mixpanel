package mixpanel

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient_Track(t *testing.T) {
	viper.SetEnvPrefix("mixpanel")
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	token := viper.Get("token")
	if token == nil {
		t.Fatalf("Mixpanel Token is not provided for the test. You can add MIXPANEL_TOKEN to your environment variable for the test.")
	}
	client := NewClient(token.(string))

	props := map[string]interface{}{"test": "testing"}
	event := NewEvent("go-test", props)
	event.DistinctID = "1"
	_, err := client.Track(event)
	assert.NoError(t, err)

	richEvent := NewEvent("go-test", props)
	richEvent.DistinctID = "2"
	richEvent.Time = uint(time.Now().Unix())
	richEvent.IP = "8.8.8.8"
	richEvent.GroupKey = "MPGO"
	richEvent.GroupID = "MPGOTEST"
	_, err = client.Track(richEvent)
	assert.NoError(t, err)
}

func TestUpdateOperation_JSON(t *testing.T) {
	operation := NewSetOperation(map[string]interface{}{"test": "testing"})
	operation.IgnoreAlias = true
	operation.IgnoreTime = false
	assert.EqualValues(t, "{\"$token\":\"\",\"$distinct_id\":\"\",\"$ignore_alias\":true,\"$set\":{\"test\":\"testing\"}}", operation.JSON())
}

func TestNewAddOperation(t *testing.T) {
	operation := NewAddOperation(map[string]interface{}{"test": "testing"})
	operation.DistinctID = "1"
	assert.EqualValues(t, "{\"$token\":\"\",\"$distinct_id\":\"1\",\"$add\":{\"test\":\"testing\"}}", operation.JSON())
}

func TestNewSetOperation(t *testing.T) {
	operation := NewSetOperation(map[string]interface{}{"test": "testing"})
	operation.DistinctID = "2"
	assert.EqualValues(t, "{\"$token\":\"\",\"$distinct_id\":\"2\",\"$set\":{\"test\":\"testing\"}}", operation.JSON())
}