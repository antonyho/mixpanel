// Package mixpanel provides the API client for making Mixpanel API calls.
package mixpanel

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	// Protocol is the protocol being used for Mixpanel API request.
	Protocol = "https"
	// Host is the hostname of Mixpanel API.
	Host = "api.mixpanel.com"
	// TrackingPath is the API URL path for Mixpanel event tracking.
	// https://developer.mixpanel.com/docs/http#section-events
	TrackingPath = "track/"
	// UpdatePath is the API URL path for Mixpanel customer profile updates.
	// https://developer.mixpanel.com/docs/http#section-profile-updates
	UpdatePath = "engage/"
)

// Event is the structure of Mixpanel event.
// https://developer.mixpanel.com/docs/http#section-events
type Event struct {
	Title      string                 `json:"event"`
	Token      string                 `json:"-"`
	Properties map[string]interface{} `json:"properties"`
	DistinctID string                 `json:"-"`
	Time       uint                   `json:"-"`
	IP         string                 `json:"-"`
	InsertID   string                 `json:"-"`
	GroupKey   string                 `json:"-"`
	GroupID    string                 `json:"-"`
}

// JSON returns the JSON format string of this event struct.
func (e Event) JSON() string {
	// Set pre-defined properties into the map
	if e.Token != "" {
		e.Properties["token"] = e.Token
	}
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

// SetToken is use by Mixpanel client to set the Mixpanel API token before making API request.
func (e *Event) SetToken(token string) {
	e.Token = token
}

func NewEvent(event string, props map[string]interface{}) *Event {
	// Is it necessary to check props keys and find any keyword
	// which violate with pre-defined keywords?
	return &Event{
		Title:      event,
		Properties: props,
	}
}

type UpdateOperation interface {
	JSON() string
	SetToken(token string)
}

// BasicUpdateOperation is the basic structure of a update operation.
// It has attributes which every kind of update operation consists.
type BasicUpdateOperation struct {
	Token       string `json:"$token"`
	DistinctID  string `json:"$distinct_id"`
	IP          string `json:"$ip,omitempty"`
	Time        uint   `json:"$time,omitempty"`
	IgnoreTime  bool   `json:"$ignore_time,omitempty"`
	IgnoreAlias bool   `json:"$ignore_alias,omitempty"`
}

// JSON returns the JSON format string of this update operation struct.
func (u BasicUpdateOperation) JSON() string {
	j, _ := json.Marshal(u)
	return string(j)
}

// SetToken is use by Mixpanel client to set the Mixpanel API token before making API request.
func (u *BasicUpdateOperation) SetToken(token string) {
	u.Token = token
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

func NewSetOperation(distinctID string, properties map[string]interface{}) UpdateOperation {
	return &SetOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		SetProperties:        properties,
	}
}

func NewSetOnceOperation(distinctID string, properties map[string]interface{}) UpdateOperation {
	return &SetOnceOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		SetOnceProperties:    properties,
	}
}

func NewAddOperation(distinctID string, properties map[string]interface{}) UpdateOperation {
	return &AddOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		AddProperties:        properties,
	}
}

func NewAppendOperation(distinctID string, properties map[string]interface{}) UpdateOperation {
	return &AppendOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		AppendProperties:     properties,
	}
}

func NewUnsetOperation(distinctID string, propertyNames []string) UpdateOperation {
	return &UnsetOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		UnsetProperties:      propertyNames,
	}
}

func NewUnionOperation(distinctID string, properties map[string][]interface{}) UpdateOperation {
	return &UnionOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		UnionProperties:      properties,
	}
}

func NewRemovalOperation(distinctID string, properties map[string]interface{}) UpdateOperation {
	return &RemoveOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
		RemoveProperties:     properties,
	}
}

func NewDeleteOperation(distinctID string) UpdateOperation {
	return &DeleteOperation{
		BasicUpdateOperation: BasicUpdateOperation{DistinctID: distinctID},
	}
}

// Client is the client for making Mixpanel API.
type Client struct {
	token string
}

// NewClient is the constructor for Mixpanel API client.
// It takes the Mixpanel API token string.
func NewClient(token string) *Client {
	return &Client{token: token}
}

// Track makes event tracking API call to Mixpanel.
// It returns the simple response from Mixpanel with a boolean to indicate the API call was successful or not.
// It returns error when the RESTful request to Mixpanel has error.
func (c Client) Track(e *Event) (success bool, err error) {
	e.SetToken(c.token)
	values := url.Values{}
	values.Set("data", base64.StdEncoding.EncodeToString([]byte(e.JSON())))
	req := url.URL{
		Scheme:   Protocol,
		Host:     Host,
		Path:     TrackingPath,
		RawQuery: values.Encode(),
	}
	resp, err := http.Get(req.String())
	if success = err == nil; !success {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()
	respBody, err := ioutil.ReadAll(resp.Body)
	if success = err == nil; !success {
		return
	}
	success = string(respBody) == "1"
	return
}

// Update makes profile update API call to Mixpanel.
// It returns the simple response from Mixpanel with a boolean to indicate the API call was successful or not.
// It returns error when the RESTful request to Mixpanel has error.
func (c Client) Update(u UpdateOperation) (success bool, err error) {
	u.SetToken(c.token)
	values := url.Values{}
	values.Set("data", base64.StdEncoding.EncodeToString([]byte(u.JSON())))
	req := url.URL{
		Scheme:   Protocol,
		Host:     Host,
		Path:     UpdatePath,
		RawQuery: values.Encode(),
	}
	resp, err := http.Get(req.String())
	if success = err == nil; !success {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()
	respBody, err := ioutil.ReadAll(resp.Body)
	if success = err == nil; !success {
		return
	}
	success = string(respBody) == "1"
	return
}
