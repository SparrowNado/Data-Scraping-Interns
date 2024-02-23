// TODO: fix API fetch error handling

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	wg   sync.WaitGroup
	base = "https://access.redhat.com/hydra/rest/securitydata/cve.json"
)

type CVEData []struct {
	CVE                 string      `json:"CVE"`
	Severity            string      `json:"severity"`
	PublicDate          string      `json:"public_date"`
	Advisories          interface{} `json:"advisories"`
	Bugzilla            string      `json:"bugzilla"`
	BugzillaDescription string      `json:"bugzilla_description"`
	CvssScore           interface{} `json:"cvss_score"`
	CvssScoringVector   interface{} `json:"cvss_scoring_vector"`
	CWE                 string      `json:"CWE"`
	AffectedPackages    interface{} `json:"affected_packages"`
	PackageState        interface{} `json:"package_state"`
	ResourceURL         string      `json:"resource_url"`
	Cvss3ScoringVector  string      `json:"cvss3_scoring_vector,omitempty"`
	Cvss3Score          string      `json:"cvss3_score,omitempty"`
}

func main() {
	data, err := fetchJSON(base)
	if err != nil {
		panic(err)
	}

	downloadFiles(data)
}

func fetchJSON(url string) (CVEData, error) {
	var data CVEData

	resp, err := http.Get(url)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("failed to fetch JSON: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func downloadFiles(data CVEData) {
	wg.Add(len(data))

	for _, entry := range data {
		if entry.ResourceURL != "" {
			go downloadFile(entry.ResourceURL)
		}
	}

	wg.Wait()
}

func downloadFile(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Suppress the error message for 403 Forbidden
		if resp.StatusCode != http.StatusForbidden {
			return fmt.Errorf("failed to download file: %s", resp.Status)
		}
	}

	fileName := getFileNameFromURL(url)
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func getFileNameFromURL(url string) string {
	// Extract the last segment of the URL as the file name
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}
