package main

import (
	bento "github.com/bentonow/bento-golang-sdk"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func main() {
	config := &bento.Config{
		PublishableKey: "your publishable key",
		SecretKey:      "your secret key",
		SiteUUID:       "your site uuid",
		Timeout:        10 * time.Second,
	}

	client, err := bento.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Find a subscriber
	fmt.Println("\n=== Finding Subscriber ===")
	subscriber, err := client.FindSubscriber(ctx, "example@test.com")
	if err != nil {
		log.Printf("Find subscriber error: %v", err)
	} else {
		prettyPrint(subscriber)
	}

	// Create a subscriber
	fmt.Println("\n=== Creating Subscriber ===")
	createInput := &bento.SubscriberInput{
		Email:     "example@test.com",
		FirstName: "Jesse",
		LastName:  "Test",
		Tags:      []string{"test-tag"},
		Fields: map[string]interface{}{
			"company": "Test Company",
			"role":    "Developer",
		},
	}

	newSubscriber, err := client.CreateSubscriber(ctx, createInput)
	if err != nil {
		log.Printf("Create subscriber error: %v", err)
	} else {
		prettyPrint(newSubscriber)
	}

	//  Import subscribers
	fmt.Println("\n=== Importing Subscribers ===")
	importSubscribers := []*bento.SubscriberInput{
		{
			Email:     "example@test.com",
			FirstName: "Jesse",
			LastName:  "Import1",
			Tags:      []string{"imported"},
			Fields: map[string]interface{}{
				"imported_at": time.Now().Format(time.RFC3339),
			},
		},
		{
			Email:     "example+test@test.com",
			FirstName: "Jesse",
			LastName:  "Import2",
			Tags:      []string{"imported"},
			Fields: map[string]interface{}{
				"imported_at": time.Now().Format(time.RFC3339),
			},
		},
	}

	err = client.ImportSubscribers(ctx, importSubscribers)
	if err != nil {
		log.Printf("Import subscribers error: %v", err)
	} else {
		fmt.Println("Successfully imported subscribers")
	}

	// Track Events
	fmt.Println("\n=== Tracking Events ===")
	events := []bento.EventData{
		{
			Type:  "$completed_onboarding",
			Email: "example@test.com",
			Fields: map[string]interface{}{
				"onboarding_type": "api_test",
				"timestamp":       time.Now().Format(time.RFC3339),
			},
		},
	}

	err = client.TrackEvent(ctx, events)
	if err != nil {
		log.Printf("Track event error: %v", err)
	} else {
		fmt.Println("Successfully tracked event")
	}

	// Broadcast Creation
	fmt.Println("\n=== Creating Broadcast ===")
	broadcasts := []bento.BroadcastData{
		{
			Name:    "API Test Broadcast",
			Subject: "Test Email",
			Content: "<p>Hello {{ visitor.first_name }}</p>",
			Type:    bento.BroadcastTypePlain,
			From: bento.ContactData{
				Name:  "Example Test",
				Email: "example@test.com",
			},
			InclusiveTags:    "test-tag",
			BatchSizePerHour: 1000,
		},
	}

	err = client.CreateBroadcast(ctx, broadcasts)
	if err != nil {
		log.Printf("Create broadcast error: %v", err)
	} else {
		fmt.Println("Successfully created broadcast")
	}

	// Get Site Stats
	fmt.Println("\n=== Getting Site Stats ===")
	stats, err := client.GetSiteStats(ctx)
	if err != nil {
		log.Printf("Get site stats error: %v", err)
	} else {
		prettyPrint(stats)
	}

	// Get all tags
	fmt.Println("=== Getting All Tags ===")
	tags, err := client.GetTags(ctx)
	if err != nil {
		log.Printf("Get tags error: %v", err)
	} else {
		for _, tag := range tags {
			fmt.Printf("Tag: %s (ID: %s, Created: %s)\n",
				tag.Attributes.Name,
				tag.ID,
				tag.Attributes.CreatedAt)
		}
	}

	// Create a new tag
	fmt.Println("\n=== Creating New Tag ===")
	newTag, err := client.CreateTag(ctx, "go-sdk-test-tag")
	if err != nil {
		log.Printf("Create tag error: %v", err)
	} else {
		fmt.Println("Created new tag:")
		prettyPrint(newTag)
	}

	// Validate email
	fmt.Println("=== Validating Email ===")
	validationData := &bento.ValidationData{
		EmailAddress: "example@test.com",
		FullName:     "John Snow",
		UserAgent:    "Go-SDK-Test",
		IPAddress:    "1.1.1.1",
	}

	result, err := client.ValidateEmail(ctx, validationData)
	if err != nil {
		log.Printf("Email validation error: %v", err)
	} else {
		fmt.Printf("Email validation result: %v\n", result.Valid)
	}

	// Subscriber Command
	fmt.Println("\n=== Executing Subscriber Command ===")
	commands := []bento.CommandData{
		{
			Command: bento.CommandAddTag,
			Email:   "example@test.com",
			Query:   "api-test-tag",
		},
	}

	err = client.SubscriberCommand(ctx, commands)
	if err != nil {
		log.Printf("Subscriber command error: %v", err)
	} else {
		fmt.Println("Successfully executed subscriber command")
	}

	// Get all fields
	fmt.Println("\n=== Getting All Fields ===")
	fields, err := client.GetFields(ctx)
	if err != nil {
		log.Printf("Get fields error: %v", err)
	} else {
		for _, field := range fields {
			fmt.Printf("Field: %s (ID: %s)\n",
				field.Attributes.Name,
				field.ID)
			fmt.Printf("  Key: %s\n", field.Attributes.Key)
			fmt.Printf("  Created: %s\n", field.Attributes.CreatedAt.Format(time.RFC3339))
			fmt.Printf("  Type: %s\n\n", field.Type)
		}
	}

	// Create a new field
	fmt.Println("\n=== Creating New Field ===")
	newField, err := client.CreateField(ctx, "purchase_amount")
	if err != nil {
		log.Printf("Create field error: %v", err)
	} else {
		fmt.Println("Created new field:")
		prettyPrint(newField)
	}

	// Example of validation error
	fmt.Println("\n=== Testing Field Creation Validation ===")
	invalidField, err := client.CreateField(ctx, "")
	if err != nil {
		fmt.Printf("Expected validation error occurred: %v\n", err)
	} else {
		prettyPrint(invalidField)
	}

	// Check Blacklist Status
	fmt.Println("=== Checking Blacklist Status ===")
	blacklistData := &bento.BlacklistData{
		Domain:    "example.com",
		IPAddress: "1.1.1.1",
	}

	blacklistResult, err := client.GetBlacklistStatus(ctx, blacklistData)
	if err != nil {
		log.Printf("Blacklist check error: %v", err)
	} else {
		fmt.Println("Blacklist check result:")
		prettyPrint(blacklistResult)
	}

	// Content Moderation
	fmt.Println("\n=== Content Moderation ===")
	content := "Hello world! This is a test message."
	moderationResult, err := client.GetContentModeration(ctx, content)
	if err != nil {
		log.Printf("Content moderation error: %v", err)
	} else {
		fmt.Println("Content moderation result:")
		prettyPrint(moderationResult)
	}

	// EGender Prediction
	fmt.Println("\n=== Gender Prediction ===")
	fullName := "John Smith"
	genderResult, err := client.GetGender(ctx, fullName)
	if err != nil {
		log.Printf("Gender prediction error: %v", err)
	} else {
		fmt.Println("Gender prediction result:")
		prettyPrint(genderResult)
	}

	// IP Geolocation
	fmt.Println("\n=== IP Geolocation ===")
	ipAddress := "8.8.8.8"
	geoResult, err := client.GeoLocateIP(ctx, ipAddress)
	if err != nil {
		log.Printf("Geolocation error: %v", err)
	} else {
		fmt.Println("Geolocation result:")
		prettyPrint(geoResult)
	}

	// Example of error handling with invalid IP
	fmt.Println("\n=== Testing Invalid IP ===")
	invalidIP := "invalid.ip"
	_, err = client.GeoLocateIP(ctx, invalidIP)
	if err != nil {
		fmt.Printf("Expected validation error occurred: %v\n", err)
	}

	// Example of multiple blacklist checks
	fmt.Println("\n=== Multiple Blacklist Checks ===")
	domains := []string{"example.com", "test.com", "sample.org"}
	for _, domain := range domains {
		blacklistData := &bento.BlacklistData{
			Domain: domain,
		}
		result, err := client.GetBlacklistStatus(ctx, blacklistData)
		if err != nil {
			log.Printf("Blacklist check error for %s: %v", domain, err)
		} else {
			fmt.Printf("Blacklist result for %s:\n", domain)
			prettyPrint(result)
		}
	}

	// Example of content moderation with different types of content
	fmt.Println("\n=== Multiple Content Moderation Tests ===")
	contentSamples := []string{
		"This is a normal message.",
		"This message contains some keywords that might need moderation.",
		"A simple test of the moderation system.",
	}

	for _, content := range contentSamples {
		result, err := client.GetContentModeration(ctx, content)
		if err != nil {
			log.Printf("Content moderation error for '%s': %v", content, err)
		} else {
			fmt.Printf("Moderation result for: '%s'\n", content)
			prettyPrint(result)
		}
	}

	// Example of gender prediction with multiple names
	fmt.Println("\n=== Multiple Gender Predictions ===")
	names := []string{"John Smith", "Mary Johnson", "Pat Taylor"}
	for _, name := range names {
		result, err := client.GetGender(ctx, name)
		if err != nil {
			log.Printf("Gender prediction error for %s: %v", name, err)
		} else {
			fmt.Printf("Gender prediction for %s:\n", name)
			prettyPrint(result)
		}
	}

	// Example of IP geolocation with multiple IPs
	fmt.Println("\n=== Multiple IP Geolocations ===")
	ips := []string{"8.8.8.8", "1.1.1.1", "208.67.222.222"}
	for _, ip := range ips {
		result, err := client.GeoLocateIP(ctx, ip)
		if err != nil {
			log.Printf("Geolocation error for %s: %v", ip, err)
		} else {
			fmt.Printf("Geolocation result for %s:\n", ip)
			prettyPrint(result)
		}
	}

	// Get Site Stats
	fmt.Println("=== Getting Site Stats ===")
	siteStats, err := client.GetSiteStats(ctx)
	if err != nil {
		log.Printf("Site stats error: %v", err)
	} else {
		fmt.Println("Site statistics:")
		prettyPrint(siteStats)
	}

	// Get Segment Stats
	fmt.Println("\n=== Getting Segment Stats ===")
	// Example segment IDs - replace with actual segment IDs
	segmentIDs := []string{"segment_123", "segment_456", "segment_789"}

	for _, segmentID := range segmentIDs {
		fmt.Printf("\nGetting stats for segment: %s\n", segmentID)
		segmentStats, err := client.GetSegmentStats(ctx, segmentID)
		if err != nil {
			log.Printf("Segment stats error for %s: %v", segmentID, err)
			continue
		}
		prettyPrint(segmentStats)
	}

	// Get Report Stats
	fmt.Println("\n=== Getting Report Stats ===")
	// Example report IDs - replace with actual report IDs
	reportIDs := []string{"report_612roaJQBO4pUqjzd3RjW8Vn", "report_a0qX4grDmZVQsj3MmK6VY179", "report_a0qX4grDmZVQsj3MmK6VY179"}

	for _, reportID := range reportIDs {
		fmt.Printf("\nGetting stats for report: %s\n", reportID)
		reportStats, err := client.GetReportStats(ctx, reportID)
		if err != nil {
			log.Printf("Report stats error for %s: %v", reportID, err)
			continue
		}
		prettyPrint(reportStats)
	}

	// Error Handling for Invalid Segment ID
	fmt.Println("\n=== Testing Invalid Segment ID ===")
	invalidSegmentStats, err := client.GetSegmentStats(ctx, "")
	if err != nil {
		fmt.Printf("Expected error for invalid segment ID: %v\n", err)
	} else {
		prettyPrint(invalidSegmentStats)
	}

	// Error Handling for Invalid Report ID
	fmt.Println("\n=== Testing Invalid Report ID ===")
	invalidReportStats, err := client.GetReportStats(ctx, "")
	if err != nil {
		fmt.Printf("Expected error for invalid report ID: %v\n", err)
	} else {
		prettyPrint(invalidReportStats)
	}

	emails := []bento.EmailData{
		{
			To:            "example@test.com",
			From:          "example@test.com",
			Subject:       "Reset Password",
			HTMLBody:      "<p>Here is a link to reset your password ... {{ link }}</p>",
			Transactional: true,
			Personalizations: map[string]interface{}{
				"link": "https://example.com/test",
			},
		},
	}

	// Send the email
	results, err := client.CreateEmails(ctx, emails)
	if err != nil {
		log.Fatalf("Failed to send emails: %v", err)
	}

	fmt.Printf("Successfully queued %d emails for delivery\n", results)

	// Example of sending multiple emails
	multipleEmails := []bento.EmailData{
		{
			To:            "example@test.com",
			From:          "example@test.com",
			Subject:       "Welcome!",
			HTMLBody:      "<p>Welcome to our service 1, {{ name }}!</p>",
			Transactional: true,
			Personalizations: map[string]interface{}{
				"name": "User 1",
			},
		},
		{
			To:            "example@test.com",
			From:          "example@test.com",
			Subject:       "Welcome!",
			HTMLBody:      "<p>Welcome to our service 2, {{ name }}!</p>",
			Transactional: true,
			Personalizations: map[string]interface{}{
				"name": "User 2",
			},
		},
	}

	results, err = client.CreateEmails(ctx, multipleEmails)
	if err != nil {
		log.Fatalf("Failed to send multiple emails: %v", err)
	}

	fmt.Printf("Successfully queued %d emails for delivery\n", results)
}
