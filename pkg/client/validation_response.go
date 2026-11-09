package client

import "time"

// ValidationResult struct for unmarshalling the NGR validation response
type ValidationResult struct {
	Errors                     []any                         `json:"errors"`
	Infos                      []any                         `json:"infos"`
	UUID                       string                        `json:"uuid"`
	Metadata                   []int64                       `json:"metadata"`
	MetadataErrors             map[string][]MetadataError    `json:"metadataErrors"`
	MetadataInfos              map[string][]any              `json:"metadataInfos"`
	ValidMetadata              map[string]MetadataError      `json:"validMetadata"`
	InvalidMetadata            map[string]MetadataError      `json:"invalidMetadata"`
	MetadataWithWarnings       map[string]MetadataError      `json:"metadataWithWarnings"`
	ValidationErrors           map[string][]SchematronReport `json:"validationErrors"`
	ValidationWarnings         map[string][]SchematronReport `json:"validationWarnings"`
	NumberOfNullRecords        int                           `json:"numberOfNullRecords"`
	NumberOfRecordsProcessed   int                           `json:"numberOfRecordsProcessed"`
	NumberOfRecordsUnchanged   int                           `json:"numberOfRecordsUnchanged"`
	NumberOfRecordsWithErrors  int                           `json:"numberOfRecordsWithErrors"`
	NumberOfRecordNotFound     int                           `json:"numberOfRecordNotFound"`
	NumberOfRecordsNotEditable int                           `json:"numberOfRecordsNotEditable"`
	NumberOfRecords            int                           `json:"numberOfRecords"`
	Type                       string                        `json:"type"`
	StartIsoDateTime           time.Time                     `json:"startIsoDateTime"`
	EndIsoDateTime             time.Time                     `json:"endIsoDateTime"`
	EllapsedTimeInSeconds      int                           `json:"ellapsedTimeInSeconds"` // Typo in NGR API
	TotalTimeInSeconds         int                           `json:"totalTimeInSeconds"`
	Running                    bool                          `json:"running"`
}

// MetadataError struct for unmarshalling the NGR validation response
type MetadataError struct {
	Message  string    `json:"message"`
	UUID     string    `json:"uuid"`
	Draft    bool      `json:"draft"`
	Approved bool      `json:"approved"`
	Date     time.Time `json:"date"`
	Stack    string    `json:"stack"`
}

// SchematronReport holds the messages of one schematron rule set, grouped by the title of the rule set.
type SchematronReport struct {
	Schematron string              `json:"schematron"`
	Messages   map[string][]string `json:"messages"`
}

// ValidationMessage is a single flattened validation message.
type ValidationMessage struct {
	Schematron string `json:"schematron"`
	Ruleset    string `json:"ruleset"`
	Message    string `json:"message"`
}

// Validation is the interpreted outcome of the validation of one record.
// Warnings never make a record invalid.
type Validation struct {
	Errors   []ValidationMessage
	Warnings []ValidationMessage
	// Validated is false when NGR did not validate the record, e.g. because it is not found or not editable.
	Validated bool
}

// Validation interprets the validation result for the record with the given NGR internal id.
// Both the current NGR output (validationErrors/validationWarnings) and the legacy output (metadataErrors) are read.
func (r ValidationResult) Validation(internalID string) Validation {
	v := Validation{
		Errors:    flattenSchematronReports(r.ValidationErrors[internalID]),
		Warnings:  flattenSchematronReports(r.ValidationWarnings[internalID]),
		Validated: len(r.Metadata) > 0,
	}

	for _, legacy := range r.MetadataErrors[internalID] {
		v.Errors = append(v.Errors, ValidationMessage{Message: legacy.Message})
	}

	// An invalid record without detailed messages still has to be reported as invalid.
	if invalid, ok := r.InvalidMetadata[internalID]; ok && len(v.Errors) == 0 {
		v.Errors = append(v.Errors, ValidationMessage{Message: invalid.Message})
	}

	return v
}

func flattenSchematronReports(reports []SchematronReport) []ValidationMessage {
	messages := []ValidationMessage{}

	for _, report := range reports {
		for ruleset, texts := range report.Messages {
			for _, text := range texts {
				messages = append(messages, ValidationMessage{
					Schematron: report.Schematron,
					Ruleset:    ruleset,
					Message:    text,
				})
			}
		}
	}

	return messages
}
