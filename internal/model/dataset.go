package model

type TestCase struct {
	Query        string   `json:"query"`
	RelevantDocs []string `json:"relevant_docs"`
}
