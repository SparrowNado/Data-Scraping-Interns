package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	wg   sync.WaitGroup
	base = "https://access.redhat.com/hydra/rest/securitydata/cve"
)

type Welcome10 struct {
	Vulnerability Vulnerability `json:"Vulnerability"`
}

type Vulnerability struct {
	Link                 []Link               `json:"link"`
	Style                []Style              `json:"style"`
	DocumentDistribution DocumentDistribution `json:"DocumentDistribution"`
	ThreatSeverity       string               `json:"ThreatSeverity"`
	PublicDate           string               `json:"PublicDate"`
	Bugzilla             Bugzilla             `json:"Bugzilla"`
	Cvss3                Cvss3                `json:"CVSS3"`
	Cwe                  string               `json:"CWE"`
	Details              Details              `json:"Details"`
	Statement            DocumentDistribution `json:"Statement"`
	Mitigation           DocumentDistribution `json:"Mitigation"`
	PackageState         []PackageState       `json:"PackageState"`
	UpstreamFix          string               `json:"UpstreamFix"`
	References           DocumentDistribution `json:"References"`
	Name                 string               `json:"_name"`
}

type Bugzilla struct {
	ID      string `json:"_id"`
	URL     string `json:"_url"`
	XMLLang string `json:"_xml:lang"`
	Text    string `json:"__text"`
}

type Cvss3 struct {
	CVSS3BaseScore     string `json:"CVSS3BaseScore"`
	CVSS3ScoringVector string `json:"CVSS3ScoringVector"`
	Status             string `json:"_status"`
}

type Details struct {
	XMLLang string `json:"_xml:lang"`
	Source  string `json:"_source"`
	Text    string `json:"__text"`
}

type DocumentDistribution struct {
	XMLLang string `json:"_xml:lang"`
	Text    string `json:"__text"`
}

type Link struct {
	Type string `json:"_type"`
	Rel  string `json:"_rel"`
	ID   string `json:"_id"`
}

type PackageState struct {
	ProductName string   `json:"ProductName"`
	FixState    FixState `json:"FixState"`
	PackageName string   `json:"PackageName"`
	Cpe         string   `json:"_cpe"`
}

type Style struct {
	Lang string `json:"_lang"`
	Type string `json:"_type"`
	ID   string `json:"_id"`
}

type FixState string

const (
	Affected           FixState = "Affected"
	OutOfSupportScope  FixState = "Out of support scope"
	UnderInvestigation FixState = "Under investigation"
)

func main() {

	files, err := fetchFiles(base)
	if err != nil {
		panic(err)
	}

	wg.Add(len(files))
	for _, file := range files {
		go processFile(file)
	}
	wg.Wait()

	fmt.Println("All files processed successfully.")
}

func fetchFiles(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch files: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var files []string
	links := extractLinks(string(body))
	for _, link := range links {
		if strings.HasSuffix(link, ".xml") {
			files = append(files, link)
		}
	}

	return files, nil
}
func extractLinks(xmlContent string) []string {
	var links []string

	// Create a custom XML decoder that ignores errors
	decoder := xml.NewDecoder(strings.NewReader(xmlContent))
	decoder.Strict = false
	decoder.Entity = xml.HTMLEntity

	for {
		// Read tokens
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			// Ignore XML syntax errors
			continue
		}

		// Check if the token is a StartElement
		if se, ok := token.(xml.StartElement); ok {
			// Iterate over attributes of the StartElement
			for _, attr := range se.Attr {
				// If the attribute is a href, add its value to links
				if attr.Name.Local == "href" {
					links = append(links, attr.Value)
				}
			}
		}
	}

	return links
}
func processFile(file string) {
	defer wg.Done()
	fmt.Printf("Downloading file %s...\n", file)
	resp, err := http.Get(base + file)
	if err != nil {
		fmt.Printf("Error downloading file %s: %s\n", file, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to download file %s: %s\n", file, resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body for file %s: %s\n", file, err)
		return
	}

	// Write the XML content to a file
	err = os.WriteFile(file, body, 0644)
	if err != nil {
		fmt.Printf("Error writing XML content to file %s: %s\n", file, err)
		return
	}

	// Unmarshal XML content
	var welcome10 Welcome10 // No need for mypackage prefix here
	err = xml.Unmarshal(body, &welcome10)
	if err != nil {
		fmt.Printf("Error unmarshalling XML content for file %s: %s\n", file, err)
		return
	}

	// Create JSON filename
	jsonFilename := strings.TrimSuffix(file, ".xml") + ".json"

	// Marshal welcome10 to JSON
	jsonData, err := json.MarshalIndent(welcome10, "", "    ")
	if err != nil {
		fmt.Printf("Error marshalling JSON for file %s: %s\n", file, err)
		return
	}

	// Write JSON content to a file
	err = os.WriteFile(jsonFilename, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON content to file %s: %s\n", jsonFilename, err)
		return
	}

	fmt.Printf("Files %s and %s processed successfully.\n", file, jsonFilename)
}
