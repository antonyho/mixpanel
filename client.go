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

type UpdateOperation interface {
	JSON() string
}

type BasicUpdateOperation struct {
	Token string `json:"$token"`
	DistinctID string `json:"$distinct_id"`
	IP string `json:"$ip,omitempty"`
	Time uint `json:"$time,omitempty"`
	IgnoreTime bool `json:"$ignore_time,omitempty"`
	IgnoreAlias bool `json:"$ignore_alias,omitempty"`
}


func (u BasicUpdateOperation) JSON() string {
	j, _ := json.Marshal(u)
	return string(j)
}

type SetOperation struct {
	BasicUpdateOperation
	SetProperties map[string]interface{} `json:"$set"`
}

func (u SetOperation) JSON() string {
	j, _ := json.Marshal(u)
	return string(j)
}

type SetOnceOperation struct {
	BasicUpdateOperation
	SetOnceProperties map[string]interface{} `json:"$set_once"`
}

func (u SetOnceOperation) JSON() string {
	j, _ := json.Marshal(u)
	return string(j)
}

type AddOperation struct {
	BasicUpdateOperation
	AddProperties map[string]interface{} `json:"$add"`
}

func (u AddOperation) JSON() string {
	j, _ := json.Marshal(u)
	return string(j)
}

type AppendOperation struct {
	BasicUpdateOperation
	AppendProperties map[string]interface{} `json:"$append"`
}

func (u AppendOperation) JSON() string {
	j, _ := json.Marshal(u)
	return string(j)
}

type UnsetOperation struct {
	BasicUpdateOperation
	UnsetProperties []string `json:"$unset"`
}

func (u UnsetOperation) JSON() string {
	j, _ := json.Marshal(u)
	return string(j)
}

type RemoveOperation struct {
	BasicUpdateOperation
	RemoveProperties map[string]interface{} `json:"$remove"`
}

func (u RemoveOperation) JSON() string {
	j, _ := json.Marshal(u)
	return string(j)
}

type UnionOperation struct {
	BasicUpdateOperation
	UnionProperties map[string][]interface{} `json:"$union"`
}

func (u UnionOperation) JSON() string {
	j, _ := json.Marshal(u)
	return string(j)
}

type DeleteOperation struct {
	BasicUpdateOperation
	Delete string `json:"$delete"`
}

func (u DeleteOperation) JSON() string {
	j, _ := json.Marshal(u)
	return string(j)
}

func NewSetOperation(distinctID string, properties map[string]interface{}) UpdateOperation{
	return &SetOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		SetProperties: properties,
	}
}

func NewSetOnceOperation(distinctID string, properties map[string]interface{}) UpdateOperation{
	return &SetOnceOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		SetOnceProperties: properties,
	}
}

func NewAddOperation(distinctID string, properties map[string]interface{}) UpdateOperation{
	return &AddOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		AddProperties: properties,
	}
}

func NewAppendOperation(distinctID string, properties map[string]interface{}) UpdateOperation{
	return &AppendOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		AppendProperties: properties,
	}
}

func NewUnsetOperation(distinctID string, propertyNames []string) UpdateOperation{
	return &UnsetOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		UnsetProperties: propertyNames,
	}
}

func NewUnionOperation(distinctID string, properties map[string][]interface{}) UpdateOperation{
	return &UnionOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		UnionProperties: properties,
	}
}

func NewRemovalOperation(distinctID string, properties map[string]interface{}) UpdateOperation{
	return &RemoveOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		RemoveProperties: properties,
	}
}

func NewDeleteOperation(distinctID string) UpdateOperation{
	return &DeleteOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
	}
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
