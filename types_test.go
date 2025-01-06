package bento_test

import (
	"encoding/json"
	"testing"
	"time"

	bento "bento-golang-sdk"
)

func TestChartTypeIsValid(t *testing.T) {
	tests := []struct {
		name      string
		chartType bento.ChartType
		want      bool
	}{
		{
			name:      "counter type",
			chartType: bento.ChartTypeCounter,
			want:      true,
		},
		{
			name:      "column type",
			chartType: bento.ChartTypeColumn,
			want:      true,
		},
		{
			name:      "area type",
			chartType: bento.ChartTypeArea,
			want:      true,
		},
		{
			name:      "line chart type",
			chartType: bento.ChartTypeLineChart,
			want:      true,
		},
		{
			name:      "invalid type",
			chartType: "invalid_chart_type",
			want:      false,
		},
		{
			name:      "empty type",
			chartType: "",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.chartType.IsValid(); got != tt.want {
				t.Errorf("ChartType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscriberDataJSONMarshaling(t *testing.T) {
	subscriber := bento.SubscriberData{
		ID:   "test_id",
		Type: "subscriber",
		Attributes: struct {
			UUID           string                 `json:"uuid"`
			Email          string                 `json:"email"`
			Fields         map[string]interface{} `json:"fields"`
			CachedTagIDs   []string              `json:"cached_tag_ids"`
			UnsubscribedAt *string               `json:"unsubscribed_at"`
			NavigationURL  string                `json:"navigation_url"`
		}{
			UUID:  "test_uuid",
			Email: "test@example.com",
			Fields: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
			},
			CachedTagIDs:  []string{"tag1", "tag2"},
			NavigationURL: "https://example.com",
		},
	}

	// Test marshaling
	data, err := json.Marshal(subscriber)
	if err != nil {
		t.Fatalf("Failed to marshal SubscriberData: %v", err)
	}

	// Test unmarshaling
	var unmarshaledSubscriber bento.SubscriberData
	if err := json.Unmarshal(data, &unmarshaledSubscriber); err != nil {
		t.Fatalf("Failed to unmarshal SubscriberData: %v", err)
	}

	// Verify key fields
	if unmarshaledSubscriber.ID != subscriber.ID {
		t.Errorf("ID mismatch: got %v, want %v", unmarshaledSubscriber.ID, subscriber.ID)
	}
	if unmarshaledSubscriber.Attributes.Email != subscriber.Attributes.Email {
		t.Errorf("Email mismatch: got %v, want %v", unmarshaledSubscriber.Attributes.Email, subscriber.Attributes.Email)
	}
}

func TestReportDataPointJSONMarshaling(t *testing.T) {
	dataPoint := bento.ReportDataPoint{
		Group: "test_group",
		Date:  "2024-01-01",
		Value: 100,
	}

	// Test marshaling
	data, err := json.Marshal(dataPoint)
	if err != nil {
		t.Fatalf("Failed to marshal ReportDataPoint: %v", err)
	}

	// Verify JSON structure
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if jsonMap["g"] != dataPoint.Group {
		t.Errorf("Group field incorrect: got %v, want %v", jsonMap["g"], dataPoint.Group)
	}
	if jsonMap["x"] != dataPoint.Date {
		t.Errorf("Date field incorrect: got %v, want %v", jsonMap["x"], dataPoint.Date)
	}
	if int(jsonMap["y"].(float64)) != dataPoint.Value {
		t.Errorf("Value field incorrect: got %v, want %v", jsonMap["y"], dataPoint.Value)
	}

	// Test unmarshaling
	var unmarshaledDataPoint bento.ReportDataPoint
	if err := json.Unmarshal(data, &unmarshaledDataPoint); err != nil {
		t.Fatalf("Failed to unmarshal ReportDataPoint: %v", err)
	}

	if unmarshaledDataPoint != dataPoint {
		t.Errorf("Data points don't match: got %+v, want %+v", unmarshaledDataPoint, dataPoint)
	}
}

func TestFieldAttributesJSONMarshaling(t *testing.T) {
	now := time.Now().UTC()
	whitelisted := true
	attrs := bento.FieldAttributes{
		Name:        "Test Field",
		Key:         "test_field",
		Whitelisted: &whitelisted,
		CreatedAt:   now,
	}

	// Test marshaling
	data, err := json.Marshal(attrs)
	if err != nil {
		t.Fatalf("Failed to marshal FieldAttributes: %v", err)
	}

	// Test unmarshaling
	var unmarshaledAttrs bento.FieldAttributes
	if err := json.Unmarshal(data, &unmarshaledAttrs); err != nil {
		t.Fatalf("Failed to unmarshal FieldAttributes: %v", err)
	}

	// Verify fields
	if unmarshaledAttrs.Name != attrs.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaledAttrs.Name, attrs.Name)
	}
	if unmarshaledAttrs.Key != attrs.Key {
		t.Errorf("Key mismatch: got %v, want %v", unmarshaledAttrs.Key, attrs.Key)
	}
	if *unmarshaledAttrs.Whitelisted != *attrs.Whitelisted {
		t.Errorf("Whitelisted mismatch: got %v, want %v", *unmarshaledAttrs.Whitelisted, *attrs.Whitelisted)
	}
	if !unmarshaledAttrs.CreatedAt.Equal(attrs.CreatedAt) {
		t.Errorf("CreatedAt mismatch: got %v, want %v", unmarshaledAttrs.CreatedAt, attrs.CreatedAt)
	}
}