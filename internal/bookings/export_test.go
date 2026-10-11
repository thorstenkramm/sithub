package bookings

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/areas"
)

func exportConfig() *areas.Config {
	return &areas.Config{Areas: []areas.Area{
		{
			ID: "office", Name: "Office", ItemGroups: []areas.ItemGroup{
				{ID: "cube-1", Name: "Cube 1", Items: []areas.Item{
					{ID: "desk-1", Name: "Tisch 3"}, {ID: "desk-2", Name: "Tisch 4"},
				}},
			},
		},
		{
			ID: "lab", Name: "Lab", ItemGroups: []areas.ItemGroup{
				{ID: "bench", Name: "Bench", Items: []areas.Item{{ID: "lab-1", Name: "Seat 1"}}},
			},
		},
		{ID: "empty", Name: "Empty Area"},
	}}
}

type exportResponse struct {
	Data []struct {
		ID         string                  `json:"id"`
		Type       string                  `json:"type"`
		Attributes BookingExportAttributes `json:"attributes"`
	} `json:"data"`
	Meta  api.PaginationMeta  `json:"meta"`
	Links api.PaginationLinks `json:"links"`
}

func runExport(t *testing.T, h echo.HandlerFunc, query url.Values) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings/export?"+query.Encode(), http.NoBody)
	rec := httptest.NewRecorder()
	require.NoError(t, h(echo.New().NewContext(req, rec)))
	return rec
}

func TestExportHandler_ResponseContent(t *testing.T) {
	t.Parallel()
	store := setupTestStore(t)
	seedExportUser(t, store, "u-booker", "Booker", "booker@example.com", false)
	seedExportUser(t, store, "u-colleague", "Colleague", "colleague@example.com", false)
	seedExportUser(t, store, "u-gone", "Gone Person", "gone@example.com", true)

	seedTestBookingFull(t, store, "on-behalf", "desk-1", "u-colleague", "u-booker", "2026-10-05")
	seedTestBookingWithGuest(t, store, "guest", "desk-2", "guest-1a2b3c4d", "u-booker", "2026-10-06",
		true, "Guest Person", "guest@example.com")
	seedTestBooking(t, store, "soft-deleted", "lab-1", "u-gone", "2026-10-07")
	seedTestBooking(t, store, "unknown-item", "removed-desk", "u-booker", "2026-10-08")

	cfg := exportConfig()
	h := ExportHandlerDynamic(func() *areas.Config { return cfg }, store)

	rec := runExport(t, h, url.Values{"filter[from]": {"2026-10-05"}, "filter[to]": {"2026-10-08"}})
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, api.JSONAPIContentType, rec.Header().Get(echo.HeaderContentType))

	var resp exportResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 4)
	assert.Equal(t, 4, resp.Meta.Total)
	assert.Equal(t, 20, resp.Meta.Limit)

	assert.Equal(t, "bookings", resp.Data[0].Type)
	assert.Equal(t, BookingExportAttributes{
		BookingDate: "2026-10-05", UserID: "u-colleague", UserName: "Colleague", UserEmail: "colleague@example.com",
		ItemID: "desk-1", ItemName: "Tisch 3", ItemGroupID: "cube-1", ItemGroupName: "Cube 1",
		AreaID: "office", AreaName: "Office", BookedItem: "Cube 1 / Tisch 3",
	}, resp.Data[0].Attributes, "on-behalf booking shows the colleague, not the booker")

	guest := resp.Data[1].Attributes
	assert.True(t, guest.IsGuest)
	assert.Equal(t, "Guest Person", guest.UserName)
	assert.Equal(t, "guest@example.com", guest.UserEmail)

	assert.Empty(t, guest.UserID, "guests have no avatar, their user_id is synthetic")
	assert.Equal(t, "Gone Person", resp.Data[2].Attributes.UserName, "soft-deleted user is included")
	assert.Equal(t, "u-gone", resp.Data[2].Attributes.UserID)

	unknown := resp.Data[3].Attributes
	assert.Equal(t, "removed-desk", unknown.BookedItem, "an item missing from the config keeps its row")
	assert.Equal(t, "removed-desk", unknown.ItemName)
	assert.Empty(t, unknown.AreaID)
	assert.Empty(t, unknown.ItemGroupName)
}

func TestExportHandler_AreaFilter(t *testing.T) {
	t.Parallel()
	store := setupTestStore(t)
	seedExportUser(t, store, "u1", "Ann", "ann@example.com", false)
	seedTestBooking(t, store, "b-office", "desk-1", "u1", "2026-10-05")
	seedTestBooking(t, store, "b-lab", "lab-1", "u1", "2026-10-05")
	seedTestBooking(t, store, "b-unknown", "removed-desk", "u1", "2026-10-05")

	cfg := exportConfig()
	h := ExportHandlerDynamic(func() *areas.Config { return cfg }, store)

	cases := []struct {
		name      string
		area      string
		wantTotal int
	}{
		{name: "all areas by omission", area: "", wantTotal: 3},
		{name: "office", area: "office", wantTotal: 1},
		{name: "lab", area: "lab", wantTotal: 1},
		{name: "area without items", area: "empty", wantTotal: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			query := url.Values{"filter[from]": {"2026-10-01"}, "filter[to]": {"2026-10-31"}}
			if tc.area != "" {
				query.Set("filter[area]", tc.area)
			}
			rec := runExport(t, h, query)
			require.Equal(t, http.StatusOK, rec.Code)

			var resp exportResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, tc.wantTotal, resp.Meta.Total)
			assert.Len(t, resp.Data, tc.wantTotal)
		})
	}
}

func TestExportHandler_PaginationKeepsFilters(t *testing.T) {
	t.Parallel()
	store := setupTestStore(t)
	seedExportUser(t, store, "u1", "Ann", "ann@example.com", false)
	seedTestBooking(t, store, "b1", "desk-1", "u1", "2026-10-05")
	seedTestBooking(t, store, "b2", "desk-2", "u1", "2026-10-05")
	seedTestBooking(t, store, "b3", "desk-1", "u1", "2026-10-06")
	cfg := exportConfig()
	h := ExportHandlerDynamic(func() *areas.Config { return cfg }, store)

	rec := runExport(t, h, url.Values{
		"filter[from]": {"2026-10-01"}, "filter[to]": {"2026-10-31"}, "filter[area]": {"office"},
		"page[limit]": {"2"}, "page[offset]": {"2"},
	})
	require.Equal(t, http.StatusOK, rec.Code)

	var resp exportResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, api.PaginationMeta{Total: 3, Limit: 2, Offset: 2}, resp.Meta)
	require.NotNil(t, resp.Links.Prev)
	assert.Nil(t, resp.Links.Next)
	for _, link := range []string{resp.Links.Self, resp.Links.First, *resp.Links.Prev, resp.Links.Last} {
		parsed, err := url.Parse(link)
		require.NoError(t, err)
		assert.Equal(t, "/api/v1/bookings/export", parsed.Path)
		assert.Equal(t, "2026-10-01", parsed.Query().Get("filter[from]"), "from is preserved in %s", link)
		assert.Equal(t, "2026-10-31", parsed.Query().Get("filter[to]"), "to is preserved in %s", link)
		assert.Equal(t, "office", parsed.Query().Get("filter[area]"), "area is preserved in %s", link)
	}
}

func TestExportHandler_InvalidInput(t *testing.T) {
	t.Parallel()
	store := setupTestStore(t)
	cfg := exportConfig()
	h := ExportHandlerDynamic(func() *areas.Config { return cfg }, store)

	cases := []struct {
		name       string
		query      url.Values
		wantDetail string
	}{
		{
			name: "missing from", query: url.Values{"filter[to]": {"2026-10-08"}},
			wantDetail: `"filter[from]" is required`,
		},
		{
			// Plain "from" is stripped by privacy extensions, so the API only accepts filter[...] names.
			name: "plain from and to are not filters", query: url.Values{"from": {"2026-10-01"}, "to": {"2026-10-08"}},
			wantDetail: `"filter[from]" is required`,
		},
		{name: "missing to", query: url.Values{"filter[from]": {"2026-10-08"}}, wantDetail: `"filter[to]" is required`},
		{
			name: "unparseable from", query: url.Values{"filter[from]": {"08.10.2026"}, "filter[to]": {"2026-10-08"}},
			wantDetail: `invalid "filter[from]" date, use YYYY-MM-DD`,
		},
		{
			name: "unparseable to", query: url.Values{"filter[from]": {"2026-10-01"}, "filter[to]": {"tomorrow"}},
			wantDetail: `invalid "filter[to]" date, use YYYY-MM-DD`,
		},
		{
			name: "from after to", query: url.Values{"filter[from]": {"2026-10-09"}, "filter[to]": {"2026-10-08"}},
			wantDetail: `"filter[from]" must not be after "filter[to]"`,
		},
		{
			name: "unknown area",
			query: url.Values{
				"filter[from]": {"2026-10-01"}, "filter[to]": {"2026-10-08"}, "filter[area]": {"nowhere"},
			},
			wantDetail: "unknown area: nowhere",
		},
		{
			name: "invalid limit",
			query: url.Values{
				"filter[from]": {"2026-10-01"}, "filter[to]": {"2026-10-08"}, "page[limit]": {"abc"},
			},
			wantDetail: "page[limit]",
		},
		{
			name: "negative offset",
			query: url.Values{
				"filter[from]": {"2026-10-01"}, "filter[to]": {"2026-10-08"}, "page[offset]": {"-1"},
			},
			wantDetail: "page[offset]",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := runExport(t, h, tc.query)
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var resp api.ErrorResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			require.Len(t, resp.Errors, 1)
			assert.Equal(t, "400", resp.Errors[0].Status)
			assert.Contains(t, resp.Errors[0].Detail, tc.wantDetail)
		})
	}
}

func TestExportHandler_SameDayRangeIsValid(t *testing.T) {
	t.Parallel()
	store := setupTestStore(t)
	seedExportUser(t, store, "u1", "Ann", "ann@example.com", false)
	seedTestBooking(t, store, "b1", "desk-1", "u1", "2026-10-05")
	cfg := exportConfig()
	h := ExportHandlerDynamic(func() *areas.Config { return cfg }, store)

	rec := runExport(t, h, url.Values{"filter[from]": {"2026-10-05"}, "filter[to]": {"2026-10-05"}})
	require.Equal(t, http.StatusOK, rec.Code)

	var resp exportResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, 1, resp.Meta.Total)
}

func TestFormatBookedItem(t *testing.T) {
	t.Parallel()
	cfg := exportConfig()
	loc, ok := cfg.FindItemLocation("desk-2")
	require.True(t, ok)

	assert.Equal(t, "Cube 1 / Tisch 4", FormatBookedItem(loc, "desk-2"))
	assert.Equal(t, "gone", FormatBookedItem(nil, "gone"))
}
