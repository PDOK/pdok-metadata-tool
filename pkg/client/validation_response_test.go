package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const warningsOnlyResponse = `{"uuid":"x","metadata":[70043169],
"validMetadata":{"70043169":{"message":"Is valid","uuid":"u"}},
"invalidMetadata":{},
"metadataWithWarnings":{"70043169":{"message":"Has warnings","uuid":"u"}},
"validationErrors":{},
"validationWarnings":{"70043169":[{"schematron":"schematron-rules-warn-dcat-ap-nl-3","messages":{"DCAT-AP NL 3":["De periode wordt aanbevolen"]}}]},
"numberOfRecords":1,"numberOfRecordsProcessed":1}`

const errorsResponse = `{"uuid":"x","metadata":[7],
"validMetadata":{},
"invalidMetadata":{"7":{"message":"Is invalid","uuid":"u"}},
"metadataWithWarnings":{},
"validationErrors":{"7":[{"schematron":"schematron-rules-iso","messages":{"ISO":["E-mail ontbreekt","Titel ontbreekt"]}}]},
"validationWarnings":{},
"numberOfRecords":1,"numberOfRecordsProcessed":1}`

const notEditableResponse = `{"uuid":"x","metadata":[],"validMetadata":{},"invalidMetadata":{},
"validationErrors":{},"validationWarnings":{},"numberOfRecords":1,"numberOfRecordsProcessed":0,"numberOfRecordsNotEditable":1}`

func parse(t *testing.T, raw string) ValidationResult {
	t.Helper()

	var result ValidationResult
	assert.NoError(t, json.Unmarshal([]byte(raw), &result))

	return result
}

func TestValidationResult_Validation(t *testing.T) {
	t.Run("warnings only are no errors", func(t *testing.T) {
		v := parse(t, warningsOnlyResponse).Validation("70043169")
		assert.True(t, v.Validated)
		assert.Empty(t, v.Errors)
		assert.Len(t, v.Warnings, 1)
		assert.Equal(t, "De periode wordt aanbevolen", v.Warnings[0].Message)
		assert.Equal(t, "DCAT-AP NL 3", v.Warnings[0].Ruleset)
	})

	t.Run("errors", func(t *testing.T) {
		v := parse(t, errorsResponse).Validation("7")
		assert.True(t, v.Validated)
		assert.Len(t, v.Errors, 2)
		assert.Empty(t, v.Warnings)
	})

	t.Run("not validated", func(t *testing.T) {
		v := parse(t, notEditableResponse).Validation("1")
		assert.False(t, v.Validated)
		assert.Empty(t, v.Errors)
	})
}
