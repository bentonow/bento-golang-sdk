package bento

import (
	"encoding/json"
	"time"
)

// BroadcastType represents the type of broadcast
type BroadcastType string

const (
	BroadcastTypePlain BroadcastType = "plain"
	BroadcastTypeRaw   BroadcastType = "raw"
)

// CommandType represents subscriber command types
type CommandType string

const (
	CommandAddTag         CommandType = "add_tag"
	CommandAddTagViaEvent CommandType = "add_tag_via_event"
	CommandRemoveTag      CommandType = "remove_tag"
	CommandAddField       CommandType = "add_field"
	CommandRemoveField    CommandType = "remove_field"
	CommandSubscribe      CommandType = "subscribe"
	CommandUnsubscribe    CommandType = "unsubscribe"
	CommandChangeEmail    CommandType = "change_email"
)

// EventData represents a tracking event
type EventData struct {
	Type    string                 `json:"type"`
	Email   string                 `json:"email"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// SubscriberData represents subscriber information from the API
type SubscriberData struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		UUID           string                 `json:"uuid"`
		Email          string                 `json:"email"`
		Fields         map[string]interface{} `json:"fields"`
		CachedTagIDs   []string               `json:"cached_tag_ids"`
		UnsubscribedAt *string                `json:"unsubscribed_at"`
		NavigationURL  string                 `json:"navigation_url"`
	} `json:"attributes"`
}

// BroadcastData represents a broadcast message
type BroadcastData struct {
	Name             string        `json:"name"`
	Subject          string        `json:"subject"`
	Content          string        `json:"content"`
	Type             BroadcastType `json:"type"`
	From             ContactData   `json:"from"`
	InclusiveTags    string        `json:"inclusive_tags,omitempty"`
	ExclusiveTags    string        `json:"exclusive_tags,omitempty"`
	SegmentID        string        `json:"segment_id,omitempty"`
	BatchSizePerHour int           `json:"batch_size_per_hour"`
}

// ContactData represents contact information
type ContactData struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email"`
}

// CommandData represents a subscriber command
type CommandData struct {
	Command CommandType `json:"command"`
	Email   string      `json:"email"`
	Query   string      `json:"query"`
}

// TagData represents tag information from the API
type TagData struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name        string  `json:"name"`
		CreatedAt   string  `json:"created_at"`
		DiscardedAt *string `json:"discarded_at"`
		SiteID      int     `json:"site_id"`
	} `json:"attributes"`
}

type FieldAttributes struct {
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	Whitelisted *bool     `json:"whitelisted"`
	CreatedAt   time.Time `json:"created_at"`
}

type FieldData struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes FieldAttributes `json:"attributes"`
}

type FieldsResponse struct {
	Data []FieldData `json:"data"`
}

// BlacklistData represents blacklist check parameters
type BlacklistData struct {
	Domain    string `json:"domain,omitempty"`
	IPAddress string `json:"ip,omitempty"`
}

// ValidationData represents email validation parameters
type ValidationData struct {
	EmailAddress string `json:"email"`
	FullName     string `json:"name,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	IPAddress    string `json:"ip,omitempty"`
}

type ValidationResponse struct {
	Valid bool `json:"valid"`
}

// GenderData represents gender prediction parameters
type GenderData struct {
	FullName string `json:"name"`
}

// GeoLocationData represents IP geolocation parameters
type GeoLocationData struct {
	IPAddress string `json:"ip"`
}

// APIResponse represents the standard API response wrapper
type APIResponse struct {
	Data struct {
		ID         string          `json:"id"`
		Type       string          `json:"type"`
		Attributes json.RawMessage `json:"attributes"`
	} `json:"data"`
}

// ChartType Bento Reports
type ChartType string

const (
	ChartTypeCounter   ChartType = "counter"
	ChartTypeColumn    ChartType = "column_chart"
	ChartTypeArea      ChartType = "area_chart"
	ChartTypeLineChart ChartType = "line_chart"
)

func (c ChartType) IsValid() bool {
	switch c {
	case ChartTypeCounter, ChartTypeColumn, ChartTypeArea, ChartTypeLineChart:
		return true
	default:
		return false
	}
}

type ReportDataPoint struct {
	Group string `json:"g"`
	Date  string `json:"x"`
	Value int    `json:"y"`
}

type ReportResponse struct {
	ChartStyle ChartType         `json:"chart_style"`
	Data       []ReportDataPoint `json:"data"`
	ReportName string            `json:"report_name"`
	ReportType string            `json:"report_type"`
}

// EmailData represents the structure for creating an email
type EmailData struct {
	To               string                 `json:"to"`
	From             string                 `json:"from"`
	Subject          string                 `json:"subject"`
	HTMLBody         string                 `json:"html_body"`
	Transactional    bool                   `json:"transactional"`
	Personalizations map[string]interface{} `json:"personalizations,omitempty"`
}
