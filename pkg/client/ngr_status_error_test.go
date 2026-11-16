package client

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNgrClient_DeleteRecordStatusErrors(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		wantNotFound bool
	}{
		{name: "not found", status: http.StatusNotFound, wantNotFound: true},
		{name: "server error", status: http.StatusInternalServerError, wantNotFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					if strings.Contains(req.URL.String(), "info?type=me") {
						writeForbiddenResponse(rw, ContentTypeJSON)

						return
					}

					rw.WriteHeader(tt.status)
				}),
			)
			defer server.Close()

			err := getNgrClient(server).DeleteRecord("689c413e-a057-11f0-8de9-0242ac120002")

			assert.Error(t, err)
			assert.Equal(t, tt.wantNotFound, errors.Is(err, ErrNgrNotFound))

			var statusErr *NgrStatusError

			assert.ErrorAs(t, err, &statusErr)
			assert.Equal(t, tt.status, statusErr.StatusCode)
		})
	}
}
