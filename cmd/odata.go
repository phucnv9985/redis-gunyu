package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/mgtv-tech/redis-GunYu/config"
	"github.com/mgtv-tech/redis-GunYu/pkg/log"
	"github.com/mgtv-tech/redis-GunYu/pkg/odata"
)

// ODataCmd represents the OData metadata command
type ODataCmd struct {
	baseURL string
}

// NewODataCmd creates a new OData command
func NewODataCmd() *ODataCmd {
	return &ODataCmd{}
}

// Name returns the command name
func (c *ODataCmd) Name() string {
	return "odata"
}

// Run executes the OData metadata fetching command
func (c *ODataCmd) Run() error {
	// Get the OData URL from configuration or environment
	flags := config.GetFlag()
	odataURL := flags.ODataCmd.URL
	if odataURL == "" {
		odataURL = os.Getenv("ODATA_URL")
	}
	if odataURL == "" {
		return fmt.Errorf("OData URL not provided. Use -odata.url flag, set ODATA_URL environment variable, or provide as argument")
	}

	log.Infof("Fetching OData metadata from: %s", odataURL)

	client := odata.NewClient(odataURL)
	
	// Set authentication if provided
	if flags.ODataCmd.Token != "" {
		client.SetAuthToken(flags.ODataCmd.Token)
		log.Infof("Using OAuth token for authentication")
	} else if flags.ODataCmd.APIKey != "" {
		client.SetAPIKey(flags.ODataCmd.APIKey)
		log.Infof("Using API key for authentication")
	}
	
	metadata, err := client.FetchMetadata()
	if err != nil {
		return fmt.Errorf("failed to fetch OData metadata: %w", err)
	}

	log.Infof("Successfully fetched OData metadata")
	log.Infof("Found %d data services", len(metadata.DataServices))
	
	// Determine output destination
	output := flags.ODataCmd.Output
	if output == "" {
		output = "stdout"
	}
	
	// Format and output the metadata
	format := flags.ODataCmd.Format
	if err := c.outputMetadata(metadata, format, output); err != nil {
		return fmt.Errorf("failed to output metadata: %w", err)
	}

	return nil
}

// outputMetadata outputs the metadata in the specified format
func (c *ODataCmd) outputMetadata(metadata *odata.Metadata, format, output string) error {
	var data []byte
	var err error

	switch format {
	case "json":
		data, err = json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal metadata to JSON: %w", err)
		}
	case "xml":
		data, err = xml.MarshalIndent(metadata, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal metadata to XML: %w", err)
		}
	case "text":
		fallthrough
	default:
		data = []byte(c.formatMetadataAsText(metadata))
	}

	if output == "stdout" {
		fmt.Print(string(data))
	} else {
		if err := os.WriteFile(output, data, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		log.Infof("Metadata written to: %s", output)
	}

	return nil
}

// formatMetadataAsText formats the metadata as human-readable text
func (c *ODataCmd) formatMetadataAsText(metadata *odata.Metadata) string {
	result := "\n=== OData Service Metadata ===\n"
	result += fmt.Sprintf("Version: %s\n", metadata.Version)
	
	if len(metadata.DataServices) > 0 {
		dataService := metadata.DataServices[0]
		result += fmt.Sprintf("Data Service Version: %s\n", dataService.Version)
		result += fmt.Sprintf("Schema Namespace: %s\n", dataService.Schema.Namespace)
		
		// Print entity sets
		entitySets := metadata.GetEntitySets()
		if len(entitySets) > 0 {
			result += "\n=== Entity Sets ===\n"
			for _, es := range entitySets {
				result += fmt.Sprintf("- %s (EntityType: %s)\n", es.Name, es.EntityType)
			}
		}
		
		// Print entity types
		entityTypes := metadata.GetEntityTypes()
		if len(entityTypes) > 0 {
			result += "\n=== Entity Types ===\n"
			for _, et := range entityTypes {
				result += fmt.Sprintf("- %s\n", et.Name)
				if len(et.Properties) > 0 {
					result += "  Properties:\n"
					for _, prop := range et.Properties {
						nullable := ""
						if prop.Nullable == "false" {
							nullable = " (required)"
						}
						result += fmt.Sprintf("    - %s: %s%s\n", prop.Name, prop.Type, nullable)
					}
				}
				if len(et.Key.PropertyRefs) > 0 {
					result += "  Key: "
					for i, key := range et.Key.PropertyRefs {
						if i > 0 {
							result += ", "
						}
						result += key.Name
					}
					result += "\n"
				}
			}
		}
	}
	
	return result
}

// Stop gracefully stops the OData command
func (c *ODataCmd) Stop() error {
	log.Infof("Stopping OData command")
	return nil
}