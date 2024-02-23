package main

import "encoding/xml"

type OvalDefinitions struct {
	XMLName        xml.Name `xml:"oval_definitions" json:"oval_definitions,omitempty"`
	SchemaLocation string   `xml:"schemaLocation,attr" json:"schemalocation,omitempty"`
	Xmlns          string   `xml:"xmlns,attr" json:"xmlns,omitempty"`
	Xsi            string   `xml:"xsi,attr" json:"xsi,omitempty"`
	Oval           string   `xml:"oval,attr" json:"oval,omitempty"`
	OvalDef        string   `xml:"oval-def,attr" json:"oval-def,omitempty"`
	Generator      struct {
		ProductName   string `xml:"product_name"`
		SchemaVersion string `xml:"schema_version"`
		Timestamp     string `xml:"timestamp"`
	} `xml:"generator" json:"generator,omitempty"`
	Definitions struct {
		Definition []struct {
			ID       string `xml:"id,attr" json:"id,omitempty"`
			Version  string `xml:"version,attr" json:"version,omitempty"`
			Class    string `xml:"class,attr" json:"class,omitempty"`
			Metadata struct {
				Title    string `xml:"title"`
				Affected struct {
					Family   string `xml:"family,attr" json:"family,omitempty"`
					Platform string `xml:"platform"`
				} `xml:"affected" json:"affected,omitempty"`
				Reference struct {
					RefID  string `xml:"ref_id,attr" json:"ref_id,omitempty"`
					RefURL string `xml:"ref_url,attr" json:"ref_url,omitempty"`
					Source string `xml:"source,attr" json:"source,omitempty"`
				} `xml:"reference" json:"reference,omitempty"`
				Description string `xml:"description"`
			} `xml:"metadata" json:"metadata,omitempty"`
			Criteria struct {
				Operator  string `xml:"operator,attr" json:"operator,omitempty"`
				Criterion []struct {
					TestRef string `xml:"test_ref,attr" json:"test_ref,omitempty"`
					Comment string `xml:"comment,attr" json:"comment,omitempty"`
				} `xml:"criterion" json:"criterion,omitempty"`
				Criteria []struct {
					Operator  string `xml:"operator,attr" json:"operator,omitempty"`
					Criterion []struct {
						TestRef string `xml:"test_ref,attr" json:"test_ref,omitempty"`
						Comment string `xml:"comment,attr" json:"comment,omitempty"`
					} `xml:"criterion" json:"criterion,omitempty"`
					Criteria struct {
						Operator  string `xml:"operator,attr" json:"operator,omitempty"`
						Criterion []struct {
							TestRef string `xml:"test_ref,attr" json:"test_ref,omitempty"`
							Comment string `xml:"comment,attr" json:"comment,omitempty"`
						} `xml:"criterion" json:"criterion,omitempty"`
					} `xml:"criteria" json:"criteria,omitempty"`
				} `xml:"criteria" json:"criteria,omitempty"`
			} `xml:"criteria" json:"criteria,omitempty"`
		} `xml:"definition" json:"definition,omitempty"`
	} `xml:"definitions" json:"definitions,omitempty"`
	Tests struct {
		RpminfoTest []struct {
			ID      string `xml:"id,attr" json:"id,omitempty"`
			Version string `xml:"version,attr" json:"version,omitempty"`
			Comment string `xml:"comment,attr" json:"comment,omitempty"`
			Check   string `xml:"check,attr" json:"check,omitempty"`
			Xmlns   string `xml:"xmlns,attr" json:"xmlns,omitempty"`
			Object  struct {
				ObjectRef string `xml:"object_ref,attr" json:"object_ref,omitempty"`
			} `xml:"object" json:"object,omitempty"`
			State struct {
				StateRef string `xml:"state_ref,attr" json:"state_ref,omitempty"`
			} `xml:"state" json:"state,omitempty"`
		} `xml:"rpminfo_test" json:"rpminfo_test,omitempty"`
	} `xml:"tests" json:"tests,omitempty"`
	Objects struct {
		RpminfoObject []struct {
			ID      string `xml:"id,attr" json:"id,omitempty"`
			Version string `xml:"version,attr" json:"version,omitempty"`
			Xmlns   string `xml:"xmlns,attr" json:"xmlns,omitempty"`
			Name    string `xml:"name"`
		} `xml:"rpminfo_object" json:"rpminfo_object,omitempty"`
	} `xml:"objects" json:"objects,omitempty"`
	States struct {
		RpminfoState []struct {
			ID          string `xml:"id,attr" json:"id,omitempty"`
			AttrVersion string `xml:"version,attr" json:"version,omitempty"`
			Xmlns       string `xml:"xmlns,attr" json:"xmlns,omitempty"`
			Version     struct {
				Operation string `xml:"operation,attr" json:"operation,omitempty"`
			} `xml:"version" json:"version,omitempty"`
			Evr struct {
				Datatype  string `xml:"datatype,attr" json:"datatype,omitempty"`
				Operation string `xml:"operation,attr" json:"operation,omitempty"`
			} `xml:"evr" json:"evr,omitempty"`
		} `xml:"rpminfo_state" json:"rpminfo_state,omitempty"`
	} `xml:"states" json:"states,omitempty"`
}
