//Oracle-Oval

package main

import (
	"bytes"
	"compress/bzip2"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type OvalDefinitions struct {
	XMLName        xml.Name `xml:"oval_definitions"`
	Xmlns          string   `xml:"xmlns,attr"`
	Oval           string   `xml:"oval,attr"`
	OvalDef        string   `xml:"oval-def,attr"`
	UnixDef        string   `xml:"unix-def,attr"`
	RedDef         string   `xml:"red-def,attr"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Generator      struct {
		ProductName    string `xml:"product_name"`
		ProductVersion string `xml:"product_version"`
		SchemaVersion  string `xml:"schema_version"`
		Timestamp      string `xml:"timestamp"`
	} `xml:"generator"`
	Definitions struct {
		Definition struct {
			ID       string `xml:"id,attr"`
			Version  string `xml:"version,attr"`
			Class    string `xml:"class,attr"`
			Metadata struct {
				Title    string `xml:"title"`
				Affected struct {
					Family   string `xml:"family,attr"`
					Platform string `xml:"platform"`
				} `xml:"affected"`
				Reference struct {
					Source string `xml:"source,attr"`
					RefID  string `xml:"ref_id,attr"`
					RefURL string `xml:"ref_url,attr"`
				} `xml:"reference"`
				Description string `xml:"description"`
				Advisory    struct {
					Severity string `xml:"severity"`
					Rights   string `xml:"rights"`
					Issued   struct {
						Date string `xml:"date,attr"`
					} `xml:"issued"`
				} `xml:"advisory"`
			} `xml:"metadata"`
			Criteria struct {
				Operator  string `xml:"operator,attr"`
				Criterion struct {
					TestRef string `xml:"test_ref,attr"`
					Comment string `xml:"comment,attr"`
				} `xml:"criterion"`
				Criteria struct {
					Operator string `xml:"operator,attr"`
					Criteria []struct {
						Operator  string `xml:"operator,attr"`
						Criterion struct {
							TestRef string `xml:"test_ref,attr"`
							Comment string `xml:"comment,attr"`
						} `xml:"criterion"`
						Criteria struct {
							Operator string `xml:"operator,attr"`
							Criteria []struct {
								Operator  string `xml:"operator,attr"`
								Criterion []struct {
									TestRef string `xml:"test_ref,attr"`
									Comment string `xml:"comment,attr"`
								} `xml:"criterion"`
							} `xml:"criteria"`
						} `xml:"criteria"`
					} `xml:"criteria"`
				} `xml:"criteria"`
			} `xml:"criteria"`
		} `xml:"definition"`
	} `xml:"definitions"`
	Tests struct {
		RpminfoTest []struct {
			ID      string `xml:"id,attr"`
			Version string `xml:"version,attr"`
			Comment string `xml:"comment,attr"`
			Check   string `xml:"check,attr"`
			Xmlns   string `xml:"xmlns,attr"`
			Object  struct {
				ObjectRef string `xml:"object_ref,attr"`
			} `xml:"object"`
			State struct {
				StateRef string `xml:"state_ref,attr"`
			} `xml:"state"`
		} `xml:"rpminfo_test"`
	} `xml:"tests"`
	Objects struct {
		RpminfoObject []struct {
			Xmlns   string `xml:"xmlns,attr"`
			ID      string `xml:"id,attr"`
			Version string `xml:"version,attr"`
			Name    string `xml:"name"`
		} `xml:"rpminfo_object"`
	} `xml:"objects"`
	States struct {
		RpminfoState []struct {
			Xmlns          string `xml:"xmlns,attr"`
			ID             string `xml:"id,attr"`
			AttrVersion    string `xml:"version,attr"`
			SignatureKeyid struct {
				Operation string `xml:"operation,attr"`
			} `xml:"signature_keyid"`
			Version struct {
				Operation string `xml:"operation,attr"`
			} `xml:"version"`
			Arch struct {
				Operation string `xml:"operation,attr"`
			} `xml:"arch"`
			Evr struct {
				Datatype  string `xml:"datatype,attr"`
				Operation string `xml:"operation,attr"`
			} `xml:"evr"`
		} `xml:"rpminfo_state"`
	} `xml:"states"`
}

func main() {
	url := "https://linux.oracle.com/security/oval/"
	files, err := fetchFiles(url)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			err := processFile(file)
			if err != nil {
				fmt.Printf("Error processing %s: %s\n", file, err)
			}
		}(file)
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
		if strings.HasSuffix(link, ".xml.bz2") {
			files = append(files, url+link)
		}
	}

	return files, nil
}

func extractLinks(html string) []string {
	var links []string
	startIndex := 0
	for {
		linkStart := strings.Index(html[startIndex:], "href=\"")
		if linkStart == -1 {
			break
		}
		linkStart += startIndex + len("href=\"")
		linkEnd := strings.Index(html[linkStart:], "\"")
		if linkEnd == -1 {
			break
		}
		linkEnd += linkStart
		links = append(links, html[linkStart:linkEnd])
		startIndex = linkEnd
	}
	return links
}

func processFile(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Extract filename from URL
	filename := filepath.Base(url)

	var xmlContent []byte

	// Check if the file is compressed
	if strings.HasSuffix(filename, ".bz2") {
		// Decompress the bz2 file
		reader := bytes.NewReader(body)
		bzip2Reader := bzip2.NewReader(reader)
		xmlContent, err = io.ReadAll(bzip2Reader)
		if err != nil {
			return err
		}
	} else {
		// File is not compressed, use the body as is
		xmlContent = body
	}

	// Unmarshal XML content
	var ovalDefinitions OvalDefinitions
	err = xml.Unmarshal(xmlContent, &ovalDefinitions)
	if err != nil {
		return err
	}

	// Create JSON filename
	jsonFilename := strings.TrimSuffix(filename, ".xml.bz2") + ".json"

	// Marshal OvalDefinitions to JSON
	jsonData, err := json.MarshalIndent(ovalDefinitions, "", "    ")
	if err != nil {
		return err
	}

	// Write JSON content to a file
	err = os.WriteFile(jsonFilename, jsonData, 0644)
	if err != nil {
		return err
	}

	// Create XML filename
	xmlFilename := strings.TrimSuffix(filename, ".xml.bz2") + ".xml"

	// Marshal OvalDefinitions to XML
	xmlData, err := xml.MarshalIndent(ovalDefinitions, "", "    ")
	if err != nil {
		return err
	}

	// Write XML content to a file
	err = os.WriteFile(xmlFilename, xmlData, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Files %s and %s processed successfully.\n", jsonFilename, xmlFilename)
	return nil
}
