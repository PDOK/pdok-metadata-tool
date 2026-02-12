package client

import "time"

// ValidationResult struct for unmarshalling the NGR validation response
type ValidationResult struct {
	Errors                     []any                      `json:"errors"`
	Infos                      []any                      `json:"infos"`
	UUID                       string                     `json:"uuid"`
	Metadata                   []int64                    `json:"metadata"`
	MetadataErrors             map[string][]MetadataError `json:"metadataErrors"`
	MetadataInfos              map[string][]any           `json:"metadataInfos"`
	NumberOfNullRecords        int                        `json:"numberOfNullRecords"`
	NumberOfRecordsProcessed   int                        `json:"numberOfRecordsProcessed"`
	NumberOfRecordsUnchanged   int                        `json:"numberOfRecordsUnchanged"`
	NumberOfRecordsWithErrors  int                        `json:"numberOfRecordsWithErrors"`
	NumberOfRecordNotFound     int                        `json:"numberOfRecordNotFound"`
	NumberOfRecordsNotEditable int                        `json:"numberOfRecordsNotEditable"`
	NumberOfRecords            int                        `json:"numberOfRecords"`
	Type                       string                     `json:"type"`
	StartIsoDateTime           time.Time                  `json:"startIsoDateTime"`
	EndIsoDateTime             time.Time                  `json:"endIsoDateTime"`
	EllapsedTimeInSeconds      int                        `json:"ellapsedTimeInSeconds"` // Typo in NGR API
	TotalTimeInSeconds         int                        `json:"totalTimeInSeconds"`
	Running                    bool                       `json:"running"`
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
