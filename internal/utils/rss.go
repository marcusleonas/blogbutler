package utils

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// --- RSS Structure Definitions ---

type Post struct {
	PostTitle   string
	PostPath    string
	Author      string
	Date        string
	Description string
	Content     string
}

// RSS is the root element
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

// Channel contains feed metadata and items
type Channel struct {
	XMLName       xml.Name `xml:"channel"`
	Title         string   `xml:"title"`         // Feed Title (e.g., "My Blog Feed")
	Link          string   `xml:"link"`          // URL to the website (e.g., "https://example.com/")
	Description   string   `xml:"description"`   // Feed Description
	Language      string   `xml:"language"`      // e.g., "en-us"
	PubDate       string   `xml:"pubDate"`       // Feed publication date (RFC1123Z format)
	LastBuildDate string   `xml:"lastBuildDate"` // Feed last build date (RFC1123Z format)
	Generator     string   `xml:"generator"`     // Optional: Name of generator
	Items         []Item   `xml:"item"`          // Slice of feed items
}

// Item represents a single entry in the feed
type Item struct {
	XMLName        xml.Name `xml:"item"`
	Title          string   `xml:"title"`           // Post Title
	Link           string   `xml:"link"`            // Full URL to the post
	Description    string   `xml:"description"`     // Post Description (Placeholder)
	Author         string   `xml:"author"`          // Optional: email (author@example.com (Name)) (Placeholder)
	PubDate        string   `xml:"pubDate"`         // Post publication date (RFC1123Z format) (Placeholder)
	Guid           Guid     `xml:"guid"`            // Unique identifier for the item
	ContentEncoded string   `xml:"content:encoded"` // The actual Guid value (often the link)
}

// Guid represents the globally unique identifier for an item
// Often the Link, marked as a permalink
type Guid struct {
	XMLName     xml.Name `xml:"guid"`
	IsPermaLink string   `xml:"isPermaLink,attr"`
	Value       string   `xml:",chardata"` // The actual Guid value (often the link)
}

// --- RSS Generation Function ---

// GenerateRSS creates an RSS 2.0 feed string from a slice of Posts.
// baseURL is the root URL of your site (e.g., "https://example.com")
// channelTitle, channelDescription are metadata for the feed itself.
func GenerateRSS(posts []Post, baseURL, channelTitle, channelDescription string) (string, error) {
	now := time.Now().Format(time.RFC1123Z) // Standard RSS date format

	// Ensure baseURL ends correctly for joining with PostPath
	baseURL = strings.TrimSuffix(baseURL, "/")

	channel := Channel{
		Title:         channelTitle,
		Link:          baseURL + "/", // Link to the main site
		Description:   channelDescription,
		Language:      "en-gb", // Placeholder: Set your language
		PubDate:       now,     // Placeholder: Use feed generation time
		LastBuildDate: now,     // Placeholder: Use feed generation time
		Generator:     "BlogButler CLI",
		Items:         make([]Item, 0, len(posts)),
	}

	for _, post := range posts {
		// Ensure post path starts with a slash if needed
		fullPostURL := baseURL + post.PostPath // Assumes PostPath starts with "/"

		item := Item{
			Title: post.PostTitle,
			Link:  fullPostURL,
			// --- Placeholders for missing data ---
			Description: post.Description,
			Author:      post.Author,
			PubDate:     now, // Placeholder: Use feed generation time, ideally use actual post time
			// --- End Placeholders ---
			Guid: Guid{
				IsPermaLink: "true", // Assume the link is a permalink
				Value:       fullPostURL,
			},
			ContentEncoded: post.Content,
		}
		channel.Items = append(channel.Items, item)
	}

	rss := RSS{
		Version: "2.0",
		Channel: channel,
	}

	// Marshal the RSS struct into XML bytes with indentation
	xmlBytes, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal RSS to XML: %w", err)
	}

	// Prepend the standard XML header
	xmlOutput := xml.Header + string(xmlBytes)

	return xmlOutput, nil
}
