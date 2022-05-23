//go:build unit
// +build unit

package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestHttpClient(t *testing.T) {
	type args struct {
		status int
		result string
	}
	type want struct {
		result []byte
		err    bool
	}

	result := gofakeit.Word()

	testCases := []struct {
		name string
		args args
		want want
	}{
		{
			name: "StatusOk",
			args: args{
				status: http.StatusOK,
				result: result,
			},
			want: want{
				result: []byte(result),
			},
		},
		{
			name: "InternalServerError",
			args: args{
				status: http.StatusInternalServerError,
			},
			want: want{
				result: nil,
				err:    true,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.args.status)
				if tc.args.result != "" {
					_, err := fmt.Fprint(w, tc.args.result)
					assert.NoError(t, err)
				}
			}))
			defer ts.Close()

			httpClient := NewClient()
			resp, err := httpClient.Get(context.Background(), ts.URL)
			if tc.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want.result, resp)
		})
	}
}
