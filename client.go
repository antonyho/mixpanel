package mixpanel

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin/json"
	"io/ioutil"
	"net/http"
)

const (
	TrackingURL = "https://api.mixpanel.com/track/"
	UpdateURL   = "https://api.mixpanel.com/engage/"
)

//type Client interface {
//	Track(e Event) (success bool, err error)
//	Update(id string, u *UpdateOperation)
//}

type Event struct {
	Title string `json:"event"`
	Properties map[string]string `json:"properties"`
	DistinctID string `json:"-"`
	Time string `json:"-"`
	IP string `json:"-"`
	InsertID string `json:"-"`
	GroupKey string `json:"-"`
	GroupID string `json:"-"`
}


func (e Event) JSON() string {
	// Set pre-defined properties into the map
	if e.DistinctID != "" {
		e.Properties["distinct_id"] = e.DistinctID
	}
	if e.Time != "" {
		e.Properties["time"] = e.Time
	}
	if e.IP != "" {
		e.Properties["ip"] = e.IP
	}
	if e.InsertID != "" {
		e.Properties["$insert_id"] = e.InsertID
	}
	if e.GroupKey != "" {
		e.Properties["$group_key"] = e.GroupKey
	}
	if e.GroupID != "" {
		e.Properties["$group_id"] = e.GroupID
	}
	j, _ := json.Marshal(e)
	return string(j)
}

func NewEvent(event string, props map[string]string) *Event {
	// TODO - Check props keys and find any keyword which violate with pre-defined keywards
	return &Event{
		Title: event,
		Properties: props,
	}
}

type UpdateOperation struct {
}

func NewSetOperation() *UpdateOperation{
	return &UpdateOperation{}
}

func NewAddOperation() *UpdateOperation{
	return &UpdateOperation{}
}

func NewDeleteOperation() *UpdateOperation{
	return &UpdateOperation{}
}


type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

func (c Client) Track(e *Event) (success bool, err error) {
	e.Properties["token"] = c.token
	req := fmt.Sprintf("%s?data=%s", TrackingURL, base64.StdEncoding.EncodeToString([]byte(e.JSON())))
	resp, err := http.Get(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if string(respBody) == "1" {
		return true, nil
	}
	return false, nil
}

func (c Client) Update(id string, u *UpdateOperation) {}
