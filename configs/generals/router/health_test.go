package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samuskitchen/go-health-checker/configs/cache"
	events "github.com/samuskitchen/go-health-checker/configs/event"
	"github.com/samuskitchen/go-health-checker/configs/storage"
	"github.com/samuskitchen/go-health-checker/pkg/kit/enums"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type HTTPContextHealth struct {
	Req     *http.Request
	Res     *httptest.ResponseRecorder
	context echo.Context
}

func SetupHTTPContextHealth(method string, url string, body interface{}) HTTPContextHealth {
	path := fmt.Sprintf("%s%s", enums.BasePath, url)
	requestByte, _ := json.Marshal(body)
	requestReader := bytes.NewReader(requestByte)
	e := echo.New()
	req := httptest.NewRequest(method, path, requestReader)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	return HTTPContextHealth{req, res, c}
}

func TestHealthCheck(t *testing.T) {
	ctx := SetupHTTPContextHealth("GET", "/health", "")

	dbData := storage.Data{}
	cacheHazelcast := &cache.Cache{}
	rabbitClient := &events.RabbitEvent{}

	hHandler := NewHealthHandler(&dbData, cacheHazelcast, rabbitClient)

	err := hHandler.HealthChecker(ctx.context)

	assert.NoError(t, err)
}
