package odata

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Metadata represents the OData service metadata
type Metadata struct {
	XMLName xml.Name `xml:"Edmx"`
	Version string   `xml:"Version,attr"`
	DataServices []DataService `xml:"DataServices>DataService"`
}

// DataService represents a data service within the metadata
type DataService struct {
	XMLName xml.Name `xml:"DataService"`
	Version string   `xml:"Version,attr"`
	Schema  Schema   `xml:"Schema"`
}

// Schema represents the schema definition
type Schema struct {
	XMLName xml.Name `xml:"Schema"`
	Namespace string `xml:"Namespace,attr"`
	EntityTypes []EntityType `xml:"EntityType"`
	EntitySets []EntitySet `xml:"EntitySet"`
}

// EntityType represents an entity type definition
type EntityType struct {
	XMLName xml.Name `xml:"EntityType"`
	Name    string   `xml:"Name,attr"`
	Key     Key      `xml:"Key"`
	Properties []Property `xml:"Property"`
}

// EntitySet represents an entity set
type EntitySet struct {
	XMLName xml.Name `xml:"EntitySet"`
	Name    string   `xml:"Name,attr"`
	EntityType string `xml:"EntityType,attr"`
}

// Key represents the key definition for an entity
type Key struct {
	XMLName xml.Name `xml:"Key"`
	PropertyRefs []PropertyRef `xml:"PropertyRef"`
}

// PropertyRef represents a property reference in a key
type PropertyRef struct {
	XMLName xml.Name `xml:"PropertyRef"`
	Name    string   `xml:"Name,attr"`
}

// Property represents a property definition
type Property struct {
	XMLName xml.Name `xml:"Property"`
	Name    string   `xml:"Name,attr"`
	Type    string   `xml:"Type,attr"`
	Nullable string  `xml:"Nullable,attr"`
}

// Client represents an OData metadata client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Headers    map[string]string
}

// NewClient creates a new OData metadata client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Follow redirects
				return nil
			},
		},
		Headers: make(map[string]string),
	}
}

// SetHeader sets a custom header for requests
func (c *Client) SetHeader(key, value string) {
	c.Headers[key] = value
}

// SetAuthToken sets an authorization token
func (c *Client) SetAuthToken(token string) {
	c.SetHeader("Authorization", "Bearer "+token)
}

// SetAPIKey sets an API key header
func (c *Client) SetAPIKey(apiKey string) {
	c.SetHeader("X-API-Key", apiKey)
}

// FetchMetadata retrieves the OData service metadata
func (c *Client) FetchMetadata() (*Metadata, error) {
	metadataURL := c.BaseURL + "/$metadata"
	
	req, err := http.NewRequest("GET", metadataURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add custom headers
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("metadata request failed with status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Preprocess the XML to handle namespace issues
	bodyStr := string(body)
	// Replace edmx: prefixes with empty strings for easier parsing
	bodyStr = strings.ReplaceAll(bodyStr, "edmx:", "")
	bodyStr = strings.ReplaceAll(bodyStr, "xmlns:edmx=\"http://docs.oasis-open.org/odata/ns/edmx\"", "")
	
	var metadata Metadata
	if err := xml.Unmarshal([]byte(bodyStr), &metadata); err != nil {
		// Log the first 500 characters of the response for debugging
		if len(bodyStr) > 500 {
			bodyStr = bodyStr[:500] + "..."
		}
		return nil, fmt.Errorf("failed to parse metadata XML: %w\nResponse preview: %s", err, bodyStr)
	}

	return &metadata, nil
}

// GetEntitySets returns all entity sets from the metadata
func (m *Metadata) GetEntitySets() []EntitySet {
	if len(m.DataServices) == 0 {
		return nil
	}
	return m.DataServices[0].Schema.EntitySets
}

// GetEntityTypes returns all entity types from the metadata
func (m *Metadata) GetEntityTypes() []EntityType {
	if len(m.DataServices) == 0 {
		return nil
	}
	return m.DataServices[0].Schema.EntityTypes
}

// FindEntityTypeByName finds an entity type by name
func (m *Metadata) FindEntityTypeByName(name string) *EntityType {
	entityTypes := m.GetEntityTypes()
	for _, et := range entityTypes {
		if et.Name == name {
			return &et
		}
	}
	return nil
}