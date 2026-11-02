package events_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/UTDNebula/nebula-api/rest/events"
	"github.com/UTDNebula/nebula-api/rest/schema"
	"github.com/UTDNebula/nebula-api/shared"
	"github.com/gin-gonic/gin"
)

type testRepository[T any] func(context.Context, string) (schema.MultiBuildingEvents[T], error)

func (r testRepository[T]) GetEventsByDate(ctx context.Context, date string) (schema.MultiBuildingEvents[T], error) {
	return r(ctx, date)
}

func TestAstraRoutes(t *testing.T) {
	testEventRoutes(t, "astra", events.AstraRoute, `{"activity_name":"Study session"}`, "activity_name")
}

func TestMazevoRoutes(t *testing.T) {
	testEventRoutes(t, "mazevo", events.MazevoRoute, `{"eventName":"Study session"}`, "eventName")
}

func TestCalendarRoutes(t *testing.T) {
	testEventRoutes(t, "calendar", events.CalendarRoute, `{"summary":"Study session"}`, "summary")
}

// Exercise real HTTP routes. Only the data returned by the repository is controlled.
func testEventRoutes[T any](t *testing.T, source string, register func(*gin.Engine, shared.EventRepository[T]), eventJSON, eventField string) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	const date = "2026-11-02"
	const roomJSON = `{"room":"1A","events":[]}`
	const buildingJSON = `{"building":"SCI","rooms":[{"room":"1A","events":[]},{"room":"2B","events":[]}]}`
	const documentJSON = `{"date":"2026-11-02","buildings":[{"building":"SCI","rooms":[{"room":"1A","events":[]},{"room":"2B","events":[]}]},{"building":"LIB","rooms":[]}]}`
	var document schema.MultiBuildingEvents[T]
	if err := json.Unmarshal([]byte(documentJSON), &document); err != nil {
		t.Fatal(err)
	}

	missingBuilding := "Building not found. Available: SCI, LIB"
	missingRoom := "Room not found. Available in this building: 1A, 2B"
	missingDate := "No events found for the specified date"
	if source == "mazevo" {
		missingBuilding = "No events found for the specified building"
		missingRoom = "No events found for that specific building and room"
		missingDate = missingBuilding
	}
	if source == "calendar" {
		missingRoom = "No events found for the specified room. Available rooms: 1A, 2B"
	}
	roomMissingDate := missingDate
	if source == "mazevo" {
		roomMissingDate = missingRoom
	}

	type routeCase struct {
		name     string
		suffix   string
		document schema.MultiBuildingEvents[T]
		err      error
		status   int
		data     string
	}
	cases := []routeCase{
		{"date events", date, document, nil, 200, documentJSON},
		{"building events", date + "/SCI", document, nil, 200, buildingJSON},
		{"room events", date + "/SCI/1A", document, nil, 200, roomJSON},
		{"date not found", date, schema.MultiBuildingEvents[T]{}, shared.ErrEventsNotFound, 200, `{"date":"2026-11-02","buildings":[]}`},
		{"building date not found", date + "/SCI", schema.MultiBuildingEvents[T]{}, shared.ErrEventsNotFound, 404, jsonString(missingDate)},
		{"room date not found", date + "/SCI/1A", schema.MultiBuildingEvents[T]{}, shared.ErrEventsNotFound, 404, jsonString(roomMissingDate)},
		{"wrapped missing date", date, schema.MultiBuildingEvents[T]{}, fmt.Errorf("lookup: %w", shared.ErrEventsNotFound), 200, `{"date":"2026-11-02","buildings":[]}`},
		{"building missing", date + "/UNKNOWN", document, nil, 404, jsonString(missingBuilding)},
		{"room missing", date + "/SCI/UNKNOWN", document, nil, 404, jsonString(missingRoom)},
		{"invalid date preserves empty response", "not-a-date", schema.MultiBuildingEvents[T]{}, shared.ErrEventsNotFound, 200, `{"date":"not-a-date","buildings":[]}`},
	}
	empty := schema.MultiBuildingEvents[T]{Date: date, Buildings: []schema.SingleBuildingEvents[T]{}}
	emptyBuildingMessage := "Building not found. Available: "
	emptyRoomMessage := "Building not found"
	if source == "mazevo" {
		emptyBuildingMessage = missingBuilding
		emptyRoomMessage = missingRoom
	}
	if source == "calendar" {
		emptyBuildingMessage = "No buildings found for the specified date"
		emptyRoomMessage = emptyBuildingMessage
	}
	cases = append(cases,
		routeCase{"empty date", date, empty, nil, 200, `{"date":"2026-11-02","buildings":[]}`},
		routeCase{"no buildings for building route", date + "/SCI", empty, nil, 404, jsonString(emptyBuildingMessage)},
		routeCase{"no buildings for room route", date + "/SCI/1A", empty, nil, 404, jsonString(emptyRoomMessage)},
	)
	roomBuildingMessage := "Building not found"
	if source == "mazevo" {
		roomBuildingMessage = missingRoom
	}
	noRoomsMessage := "Room not found. Available in this building: "
	if source == "mazevo" {
		noRoomsMessage = missingRoom
	}
	if source == "calendar" {
		noRoomsMessage = "No rooms found for the specified building"
	}
	cases = append(cases,
		routeCase{"building missing for room route", date + "/UNKNOWN/1A", document, nil, 404, jsonString(roomBuildingMessage)},
		routeCase{"building has no rooms", date + "/LIB/1A", document, nil, 404, jsonString(noRoomsMessage)},
	)
	for _, suffix := range []string{date, date + "/SCI", date + "/SCI/1A"} {
		cases = append(cases, routeCase{"database error " + suffix, suffix, schema.MultiBuildingEvents[T]{}, errors.New("database unavailable"), 500, `"database unavailable"`})
	}
	// Astra and Calendar trim input and match without case sensitivity. Mazevo does neither.
	if source == "mazevo" {
		cases = append(cases,
			routeCase{"building case mismatch", date + "/sci", document, nil, 404, jsonString(missingBuilding)},
			routeCase{"room case mismatch", date + "/SCI/1a", document, nil, 404, jsonString(missingRoom)},
		)
	} else {
		cases = append(cases,
			routeCase{"building case and spaces", date + "/%20sci%20", document, nil, 200, buildingJSON},
			routeCase{"room case and spaces", date + "/%20sci%20/%201a%20", document, nil, 200, roomJSON},
		)
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			register(router, testRepository[T](func(ctx context.Context, requestedDate string) (schema.MultiBuildingEvents[T], error) {
				expectedDate := date
				if tc.suffix == "not-a-date" {
					expectedDate = "not-a-date"
				}
				if requestedDate != expectedDate {
					return schema.MultiBuildingEvents[T]{}, errors.New("unexpected date")
				}
				return tc.document, tc.err
			}))
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/"+source+"/"+tc.suffix, nil))
			assertResponse(t, response, tc.status, tc.data)
		})
	}

	t.Run("event content survives each route", func(t *testing.T) {
		var event T
		if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
			t.Fatal(err)
		}
		data := schema.MultiBuildingEvents[T]{Date: date, Buildings: []schema.SingleBuildingEvents[T]{{Building: "SCI", Rooms: []schema.RoomEvents[T]{{Room: "1A", Events: []T{event}}}}}}
		router := gin.New()
		register(router, testRepository[T](func(context.Context, string) (schema.MultiBuildingEvents[T], error) { return data, nil }))
		for _, suffix := range []string{date, date + "/SCI", date + "/SCI/1A"} {
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/"+source+"/"+suffix, nil))
			if response.Code != 200 {
				t.Fatalf("status = %d", response.Code)
			}
			var body struct {
				Data json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			var result map[string]json.RawMessage
			if err := json.Unmarshal(body.Data, &result); err != nil {
				t.Fatal(err)
			}
			if suffix == date {
				var buildings []map[string]json.RawMessage
				if err := json.Unmarshal(result["buildings"], &buildings); err != nil {
					t.Fatal(err)
				}
				result = buildings[0]
			}
			if suffix != date+"/SCI/1A" {
				var rooms []map[string]json.RawMessage
				if err := json.Unmarshal(result["rooms"], &rooms); err != nil {
					t.Fatal(err)
				}
				result = rooms[0]
			}
			var eventData []map[string]json.RawMessage
			if err := json.Unmarshal(result["events"], &eventData); err != nil {
				t.Fatal(err)
			}
			if len(eventData) != 1 || string(eventData[0][eventField]) != `"Study session"` {
				t.Fatalf("unexpected event data: %s", body.Data)
			}
		}
	})
	t.Run("request cancellation reaches event lookup", func(t *testing.T) {
		router := gin.New()
		register(router, testRepository[T](func(ctx context.Context, date string) (schema.MultiBuildingEvents[T], error) {
			return schema.MultiBuildingEvents[T]{}, ctx.Err()
		}))
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		request := httptest.NewRequest(http.MethodGet, "/"+source+"/"+date, nil).WithContext(ctx)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assertResponse(t, response, 500, `"context canceled"`)
	})
	t.Run("preflight", func(t *testing.T) {
		router := gin.New()
		register(router, testRepository[T](func(context.Context, string) (schema.MultiBuildingEvents[T], error) {
			t.Fatal("preflight must not fetch events")
			return schema.MultiBuildingEvents[T]{}, nil
		}))
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodOptions, "/"+source, nil))
		if response.Code != 200 || response.Body.String() != "{}" {
			t.Fatalf("unexpected preflight: %d %s", response.Code, response.Body.String())
		}
	})
}

func jsonString(s string) string {
	data, _ := json.Marshal(s)
	return string(data)
}

func assertResponse(t *testing.T, response *httptest.ResponseRecorder, status int, data string) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, status, response.Body.String())
	}
	var body struct {
		Status  int             `json:"status"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	message := "success"
	if status >= 400 {
		message = "error"
	}
	if body.Status != status || body.Message != message {
		t.Fatalf("unexpected response envelope: %s", response.Body.String())
	}
	var got, want any
	if err := json.Unmarshal(body.Data, &got); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(data), &want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("data = %s, want %s", body.Data, data)
	}
}
