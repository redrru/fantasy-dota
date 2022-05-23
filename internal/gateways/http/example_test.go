//go:build unit
// +build unit

package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/redrru/fantasy-dota/pkg/server"
)

func TestGetExample(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	srv := server.ServerInterfaceWrapper{Handler: NewServer()}

	err := srv.GetExample(c)
	assert.Error(t, err)
}
