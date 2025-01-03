# Bento Go SDK
<img align="right" src="https://app.bentonow.com/brand/logoanim.gif">

> [!TIP]
> Need help? Join our [Discord](https://discord.gg/ssXXFRmt5F) or email jesse@bentonow.com for personalized support.

The Bento Go SDK makes it quick and easy to send emails and track events in your Go applications. We provide powerful and customizable APIs that can be used out-of-the-box to manage subscribers, track events, and send transactional emails.

Get started with our [📚 integration guides](https://docs.bentonow.com), or [📘 browse the SDK reference](https://docs.bentonow.com/subscribers).

## Features

* **Event Tracking**: Easily track custom events and user behavior in your Go applications.
* **Subscriber Management**: Import and manage subscribers directly using type-safe structures.
* **Idiomatic Go**: Designed using Go best practices and patterns.
* **Context Support**: All operations support context for cancellation and timeouts.
* **Strong Types**: Type-safe request and response handling.

## Requirements

- Go 1.18 or higher
- Bento API Keys

## Installation

Install the package via Go modules:

```bash
go get github.com/bentonow/bento-golang-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/bentonow/bento-go-sdk"
)

func main() {
    config := &bento.Config{
        PublishableKey: "your-key",
        SecretKey:     "your-secret",
        SiteUUID:      "your-uuid",
        Timeout:       10 * time.Second,
    }

    client, err := bento.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    subscriber, err := client.FindSubscriber(ctx, "test@example.com")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Found subscriber: %+v", subscriber)
}
```

## Core APIs

### Subscriber Management

#### Find Subscriber
Retrieves a subscriber by their email address:

```go
subscriber, err := client.FindSubscriber(ctx, "test@example.com")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Subscriber details: %+v\n", subscriber)
```

#### Create Subscriber
Creates a new subscriber in your account:

```go
input := &bento.SubscriberInput{
    Email:     "test@example.com",
    FirstName: "John",
    LastName:  "Doe",
    Tags:      []string{"new-user"},
    Fields: map[string]interface{}{
        "company": "Acme Inc",
        "role":    "Developer",
    },
}

newSubscriber, err := client.CreateSubscriber(ctx, input)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created subscriber: %+v\n", newSubscriber)
```

#### Import Subscribers
Batch import multiple subscribers:

```go
subscribers := []*bento.SubscriberInput{
    {
        Email:     "user1@example.com",
        FirstName: "John",
        LastName:  "Doe",
        Tags:      []string{"imported", "customer"},
        Fields: map[string]interface{}{
            "imported_at": time.Now().Format(time.RFC3339),
        },
    },
    {
        Email:     "user2@example.com",
        FirstName: "Jane",
        LastName:  "Smith",
        Tags:      []string{"imported", "prospect"},
        Fields: map[string]interface{}{
            "imported_at": time.Now().Format(time.RFC3339),
        },
    },
}

err = client.ImportSubscribers(ctx, subscribers)
if err != nil {
    log.Fatal(err)
}
```

### Event Tracking

#### Track Events
Send custom events to track user behavior:

```go
events := []bento.EventData{
    {
        Type:  "$completed_onboarding",
        Email: "user@example.com",
        Fields: map[string]interface{}{
            "onboarding_type": "api_test",
            "timestamp":       time.Now().Format(time.RFC3339),
        },
        Details: map[string]interface{}{
            "source": "api",
            "version": "1.0",
        },
    },
}

err = client.TrackEvent(ctx, events)
if err != nil {
    log.Fatal(err)
}
```

### Email Management

#### Send Transactional Emails
Send personalized transactional emails:

```go
emails := []bento.EmailData{
    {
        To:            "recipient@example.com",
        From:          "sender@yourdomain.com",
        Subject:       "Welcome to Our Service",
        HTMLBody:      "<p>Hello {{ name }}, welcome aboard!</p>",
        Transactional: true,
        Personalizations: map[string]interface{}{
            "name": "John Doe",
            "account_type": "premium",
        },
    },
}

results, err := client.CreateEmails(ctx, emails)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Successfully queued %d emails for delivery\n", results)
```

#### Send Multiple Emails
Send multiple emails in a single request:

```go
multipleEmails := []bento.EmailData{
    {
        To:            "user1@example.com",
        From:          "notifications@yourdomain.com",
        Subject:       "Your Daily Update",
        HTMLBody:      "<p>Hi {{ name }}, here's your daily summary...</p>",
        Transactional: true,
        Personalizations: map[string]interface{}{
            "name": "User 1",
            "summary_items": []string{"item1", "item2"},
        },
    },
    {
        To:            "user2@example.com",
        From:          "notifications@yourdomain.com",
        Subject:       "Your Daily Update",
        HTMLBody:      "<p>Hi {{ name }}, here's your daily summary...</p>",
        Transactional: true,
        Personalizations: map[string]interface{}{
            "name": "User 2",
            "summary_items": []string{"item3", "item4"},
        },
    },
}

results, err = client.CreateEmails(ctx, multipleEmails)
if err != nil {
    log.Fatal(err)
}
```

### Broadcast Management

#### Get Broadcasts
Retrieve a list of all broadcasts in your account:

```go
broadcasts, err := client.GetBroadcasts(ctx)
if err != nil {
    log.Fatal(err)
}
for _, broadcast := range broadcasts {
    fmt.Printf("Broadcast: %s (Type: %s)\n", broadcast.Name, broadcast.Type)
}
```

#### Create Broadcasts
Create new broadcast campaigns:

```go
broadcasts := []bento.BroadcastData{
    {
        Name:    "Campaign #1 Example",
        Subject: "Hello world Plain Text",
        Content: "<p>Hi {{ visitor.first_name }}</p>",
        Type:    bento.BroadcastTypePlain,
        From: bento.ContactData{
            Name:  "John Doe",
            Email: "sender@yourdomain.com",
        },
        InclusiveTags:    "lead,mql",
        ExclusiveTags:    "customers",
        SegmentID:        "segment_123456789",
        BatchSizePerHour: 1500,
    },
}

err = client.CreateBroadcast(ctx, broadcasts)
if err != nil {
    log.Fatal(err)
}
```

### Tag Management

#### Get Tags
Retrieve all tags in your account:

```go
tags, err := client.GetTags(ctx)
if err != nil {
    log.Fatal(err)
}
for _, tag := range tags {
    fmt.Printf("Tag: %s (ID: %s, Created: %s)\n",
        tag.Attributes.Name,
        tag.ID,
        tag.Attributes.CreatedAt)
}
```

#### Create Tag
Create a new tag:

```go
newTag, err := client.CreateTag(ctx, "go-sdk-test-tag")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created new tag: %s\n", newTag.Attributes.Name)
```

### Field Management

#### Get Fields
Retrieve all custom fields:

```go
fields, err := client.GetFields(ctx)
if err != nil {
    log.Fatal(err)
}
for _, field := range fields {
    fmt.Printf("Field: %s\n", field.Attributes.Key)
    fmt.Printf("  Name: %s\n", field.Attributes.Name)
    fmt.Printf("  Created: %s\n", field.Attributes.CreatedAt.Format(time.RFC3339))
}
```

#### Create Field
Create a new custom field:

```go
newField, err := client.CreateField(ctx, "purchase_amount")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created new field: %s\n", newField.Attributes.Key)
```

### Subscriber Commands

#### Execute Commands
Perform operations on subscribers:

```go
commands := []bento.CommandData{
    {
        Command: bento.CommandAddTag,
        Email:   "user@example.com",
        Query:   "new-tag",
    },
    {
        Command: bento.CommandRemoveTag,
        Email:   "user@example.com",
        Query:   "old-tag",
    },
}

err = client.SubscriberCommand(ctx, commands)
if err != nil {
    log.Fatal(err)
}
```

Available command types:
- `CommandAddTag`: Add a tag to a subscriber
- `CommandAddTagViaEvent`: Add a tag through an event
- `CommandRemoveTag`: Remove a tag from a subscriber
- `CommandAddField`: Add a field to a subscriber
- `CommandRemoveField`: Remove a field from a subscriber
- `CommandSubscribe`: Subscribe a user
- `CommandUnsubscribe`: Unsubscribe a user
- `CommandChangeEmail`: Change a user's email address

### Statistics APIs

#### Get Site Stats
Retrieve overall statistics for your site:

```go
stats, err := client.GetSiteStats(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Site statistics: %+v\n", stats)
```

#### Get Segment Stats
Retrieve statistics for a specific segment:

```go
segmentStats, err := client.GetSegmentStats(ctx, "segment_123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Segment stats: %+v\n", segmentStats)
```

#### Get Report Stats
Retrieve statistics for a specific report:

```go
reportStats, err := client.GetReportStats(ctx, "report_456")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Report stats: %+v\n", reportStats)
```

### Experimental APIs

#### Blacklist Status Check
Check if an IP address or domain is blacklisted:

```go
blacklistData := &bento.BlacklistData{
    Domain:    "example.com",
    IPAddress: "1.1.1.1",
}

result, err := client.GetBlacklistStatus(ctx, blacklistData)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Blacklist check result: %+v\n", result)
```

#### Email Validation
Validate email addresses with additional context:

```go
validationData := &bento.ValidationData{
    EmailAddress: "test@example.com",
    FullName:     "John Snow",
    UserAgent:    "Go-SDK-Test",
    IPAddress:    "1.1.1.1",
}

result, err := client.ValidateEmail(ctx, validationData)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Email validation result: valid=%v\n", result.Valid)
```

#### Content Moderation
Perform content moderation on text:

```go
content := "Hello world! This is a test message."
moderationResult, err := client.GetContentModeration(ctx, content)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Content moderation result: %+v\n", moderationResult)
```

#### Gender Prediction
Predict gender from a full name:

```go
genderResult, err := client.GetGender(ctx, "John Smith")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Gender prediction result: %+v\n", genderResult)
```

#### IP Geolocation
Geolocate an IP address:

```go
geoResult, err := client.GeoLocateIP(ctx, "8.8.8.8")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Geolocation result: %+v\n", geoResult)
```

### Error Handling

The SDK provides several predefined error types for better error handling:

```go
switch {
case errors.Is(err, bento.ErrInvalidConfig):
    log.Print("Invalid configuration")
case errors.Is(err, bento.ErrInvalidEmail):
    log.Print("Invalid email address")
case errors.Is(err, bento.ErrInvalidIPAddress):
    log.Print("Invalid IP address")
case errors.Is(err, bento.ErrAPIResponse):
    log.Print("Unexpected API response")
case errors.Is(err, bento.ErrInvalidRequest):
    log.Print("Invalid request parameters")
}
```

Available error types:
- `ErrInvalidConfig`: Configuration error
- `ErrInvalidEmail`: Invalid email format
- `ErrInvalidIPAddress`: Invalid IP address format
- `ErrInvalidRequest`: Invalid request parameters
- `ErrAPIResponse`: Unexpected API response
- `ErrInvalidName`: Invalid name format
- `ErrInvalidSegmentID`: Invalid segment ID
- `ErrInvalidContent`: Invalid content
- `ErrInvalidTags`: Invalid tags format
- `ErrInvalidBatchSize`: Invalid batch size

## Data Types

### Core Types

#### Config
Configuration for the Bento client:
```go
type Config struct {
    PublishableKey string
    SecretKey      string
    SiteUUID       string
    Timeout        time.Duration
}
```

#### EventData
Structure for tracking events:
```go
type EventData struct {
    Type    string                 `json:"type"`
    Email   string                 `json:"email"`
    Fields  map[string]interface{} `json:"fields,omitempty"`
    Details map[string]interface{} `json:"details,omitempty"`
}
```

#### SubscriberInput
Structure for creating/importing subscribers:
```go
type SubscriberInput struct {
    Email      string                 `json:"email"`
    FirstName  string                 `json:"first_name,omitempty"`
    LastName   string                 `json:"last_name,omitempty"`
    Tags       []string               `json:"tags,omitempty"`
    RemoveTags []string              `json:"remove_tags,omitempty"`
    Fields     map[string]interface{} `json:"fields,omitempty"`
}
```

#### EmailData
Structure for sending emails:
```go
type EmailData struct {
    To               string                 `json:"to"`
    From             string                 `json:"from"`
    Subject          string                 `json:"subject"`
    HTMLBody         string                 `json:"html_body"`
    Transactional    bool                   `json:"transactional"`
    Personalizations map[string]interface{} `json:"personalizations,omitempty"`
}
```

#### BroadcastData
Structure for creating broadcasts:
```go
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
```

## Best Practices

### Context Usage
Always use context for proper timeout and cancellation handling:
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

result, err := client.SomeOperation(ctx, params)
```

### Error Handling
Always check errors and handle them appropriately:
```go
result, err := client.SomeOperation(ctx, params)
if err != nil {
    switch {
    case errors.Is(err, bento.ErrInvalidConfig):
        // Handle configuration error
    case errors.Is(err, bento.ErrAPIResponse):
        // Handle API error
    default:
        // Handle unexpected error
    }
}
```

### Batch Operations
When performing batch operations, respect the limits:
- Maximum 60 emails per request
- Use reasonable batch sizes for subscriber imports
```go
// Split large imports into chunks
for _, chunk := range subscribers.Chunk(500) {
    err := client.ImportSubscribers(ctx, chunk)
    if err != nil {
        log.Printf("Failed to import chunk: %v", err)
    }
}
```

## Things to Know

1. All API methods support context for cancellation and timeouts
2. All methods perform input validation before making API calls
3. Strong types ensure type safety for requests and responses
4. Error handling follows Go best practices
5. The SDK uses the standard `net/http` client with configurable timeouts
6. All responses are properly typed for better type safety
7. The SDK supports concurrent usage and is safe for concurrent access

## Contributing

We welcome contributions! Please feel free to submit a Pull Request. Here are some ways you can help:

- Report bugs and issues
- Add new features
- Improve documentation
- Add tests
- Provide feedback

## License

The Bento SDK for Go is available as open source under the terms of the [MIT License](LICENSE.md).

## Support

Need help? Here are some ways to get support:

- Join our [Discord](https://discord.gg/ssXXFRmt5F)
- Email support at jesse@bentonow.com
- Check out our [documentation](https://docs.bentonow.com)
- Open an issue on GitHub