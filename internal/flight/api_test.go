package flight

import (
	"net/http"
	"testing"
	"time"

	"github.com/nvnoskov/dynamo-backend/internal/auth"
	"github.com/nvnoskov/dynamo-backend/internal/entity"
	"github.com/nvnoskov/dynamo-backend/internal/test"
	"github.com/nvnoskov/dynamo-backend/pkg/log"
)

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	repo := &mockRepository{items: []entity.Flight{
		{
			ID:            "123",
			Name:          "flight123",
			Number:        "123",
			Departure:     "Minsk",
			DepartureTime: time.Now(),
			Destination:   "Stockholm",
			ArrivalTime:   time.Now(),
			Fare:          "100EUR",
			Duration:      "3 hours",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}}
	RegisterHandlers(router.Group(""),
		NewService(repo,
			logger),
		auth.MockAuthHandler, logger)
	header := auth.MockAuthHeader()

	tests := []test.APITestCase{
		{"get all", "GET", "/flights", "", header, http.StatusOK, `*"total_count":1*`},
		{"get 123", "GET", "/flights/123", "", header, http.StatusOK, `*flight123*`},
		{"get unknown", "GET", "/flights/1234", "", header, http.StatusNotFound, ""},
		{"create ok", "POST", "/flights", `{"name": "BOEING 737-400","number": "UR-CSV","departure": "MALMÖ, SWEDEN2","departure_time": "2020-10-01T14:36:38Z","destination": "MERZIFON, TURKEY","arrival_time": "2020-10-01T17:36:38Z","fare": "100EUR"}`, header, http.StatusCreated, "*BOEING 737-400*"},
		{"create ok count", "GET", "/flights", "", header, http.StatusOK, `*"total_count":2*`},
		{"create auth error", "POST", "/flights", `{"name": "BOEING 737-400","number": "UR-CSV","departure": "MALMÖ, SWEDEN2","departure_time": "2020-10-01T14:36:38Z","destination": "MERZIFON, TURKEY","arrival_time": "2020-10-01T17:36:38Z","fare": "100EUR"}`, nil, http.StatusUnauthorized, ""},
		{"create input error", "POST", "/flights", `"name":"test"}`, header, http.StatusBadRequest, ""},
		{"update ok", "PUT", "/flights/123", `{"name": "flightxyz","number": "UR-CSV","departure": "MALMÖ, SWEDEN2","departure_time": "2020-10-01T14:36:38Z","destination": "MERZIFON, TURKEY","arrival_time": "2020-10-01T17:36:38Z","fare": "100EUR"}`, header, http.StatusOK, "*flightxyz*"},
		{"update verify", "GET", "/flights/123", "", header, http.StatusOK, `*flightxyz*`},
		{"update auth error", "PUT", "/flights/123", `{"name":"flightxyz"}`, nil, http.StatusUnauthorized, ""},
		{"update input error", "PUT", "/flights/123", `"name":"flightxyz"}`, header, http.StatusBadRequest, ""},
		{"delete ok", "DELETE", "/flights/123", ``, header, http.StatusOK, "*flightxyz*"},
		{"delete verify", "DELETE", "/flights/123", ``, header, http.StatusNotFound, ""},
		{"delete auth error", "DELETE", "/flights/123", ``, nil, http.StatusUnauthorized, ""},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
