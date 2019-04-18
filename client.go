package mixpanel

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	Properties map[string]interface{} `json:"properties"`
	DistinctID string `json:"-"`
	Time uint `json:"-"`
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
	if e.Time != 0 {
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

func NewEvent(event string, props map[string]interface{}) *Event {
	// Is it necessary to check props keys and find any keyword
	// which violate with pre-defined keywords?
	return &Event{
		Title: event,
		Properties: props,
	}
}

type UpdateOperation struct {
	Token string `json:"$token"`
	DistinctID string `json:"$distinct_id"`
	IP string `json:"$ip,omitempty"`
	Time uint `json:"$time,omitempty"`
	IgnoreTime bool `json:"$ignore_time,omitempty"`
	IgnoreAlias bool `json:"$ignore_alias,omitempty"`
	SetProperties map[string]interface{} `json:"$set,omitempty"`
	SetOnceProperties map[string]interface{} `json:"$set_once,omitempty"`
	AddProperties map[string]interface{} `json:"$add,omitempty"`
	AppendProperties map[string]interface{} `json:"$append,omitempty"`
	UnsetProperties []string `json:"$unset,omitempty"`
	RemoveProperties map[string]map[string]interface{} `json:"$remove,omitempty"`
	UnionProperties map[string]map[string]interface{} `json:"$union,omitempty"`
	Delete *string `json:"$delete,omitempty"`
}

func NewSetOperation(distinctID string, properties map[string]interface{}) *UpdateOperation{
	return &UpdateOperation{
		DistinctID: distinctID,
		SetProperties: properties,
	}
}

func NewSetOnceOperation(distinctID string, properties map[string]interface{}) *UpdateOperation{
	return &UpdateOperation{
		DistinctID: distinctID,
		SetOnceProperties: properties,
	}
}

func NewAddOperation(distinctID string, properties map[string]interface{}) *UpdateOperation{
	return &UpdateOperation{
		DistinctID: distinctID,
		AddProperties: properties,
	}
}

func NewAppendOperation(distinctID string, properties map[string]interface{}) *UpdateOperation{
	return &UpdateOperation{
		DistinctID: distinctID,
		AppendProperties: properties,
	}
}

func NewUnsetOperation(distinctID string, propertyNames []string) *UpdateOperation{
	return &UpdateOperation{
		DistinctID: distinctID,
		UnsetProperties: propertyNames,
	}
}

// TODO Constructor for Remove Operation
//func NewRemovalOperation(distinctID string, properties map[string]interface{}) *UpdateOperation{
//	return &UpdateOperation{
//		DistinctID: distinctID,
//		RemoveProperties: properties,
//	}
//}

// TODO Constructor for Union Updation Operation

func NewDeleteOperation(distinctID string) *UpdateOperation{
	return &UpdateOperation{
		DistinctID: distinctID,
		Delete: new(string),
	}
}

func (u *UpdateOperation) JSON() string {
	j, _ := json.Marshal(u)
	return string(j)
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
