package bookings

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thorstenkramm/sithub/internal/areas"
)

func TestEscapeCSVCell(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"=SUM(A1)", "+1", "-1", "@cmd", "\tx", "\rx", "plain", "", "Jürgen Müßig", "a=b"} {
		t.Run(value, func(t *testing.T) {
			t.Parallel()
			want := value
			if value != "" && strings.ContainsRune("=+-@\t\r", rune(value[0])) {
				want = "'" + value
			}
			assert.Equal(t, want, escapeCSVCell(value))
		})
	}
}

func TestExportCSV(t *testing.T) {
	t.Parallel()
	store := setupTestStore(t)
	seedExportUser(t, store, "u", "Jürgen, \"Müßig\"", "=email", true)
	for i := 1; i <= 25; i++ {
		seedTestBooking(t, store, fmt.Sprintf("b%d", i), "desk-1", "u", fmt.Sprintf("2026-10-%02d", i))
	}
	cfg := exportConfig()
	h := ExportCSVHandler(func() *areas.Config { return cfg }, store)
	rec := runExport(t, h, url.Values{"filter[from]": {"2026-10-01"}, "filter[to]": {"2026-10-31"},
		"page[limit]": {"garbage"}, "page[offset]": {"999"}, "filter[area]": {"office"}})
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/csv; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Equal(t, `attachment; filename="sithub-bookings-2026-10-01_2026-10-31.csv"`,
		rec.Header().Get("Content-Disposition"))
	assert.Equal(t, "no-store", rec.Header().Get("Cache-Control"))
	require.True(t, strings.HasPrefix(rec.Body.String(), "\xef\xbb\xbf"))
	assert.Contains(t, rec.Body.String(), "\r\n")
	rows, err := csv.NewReader(strings.NewReader(strings.TrimPrefix(rec.Body.String(), "\xef\xbb\xbf"))).ReadAll()
	require.NoError(t, err)
	require.Len(t, rows, 26)
	assert.Equal(t, []string{"Name", "Email", "Booked item", "Booking date"}, rows[0])
	assert.Equal(t, []string{"Jürgen, \"Müßig\"", "'=email", "Cube 1 / Tisch 3", "2026-10-01"}, rows[1])
	query := url.Values{"filter[from]": {"2026-10-01"}, "filter[to]": {"2026-10-31"}, "page[limit]": {"20"}}
	assert.Equal(t, rec.Body.String(), runExport(t, h, query).Body.String(), "CSV ignores valid pagination too")
	query.Set("page[limit]", "100")
	jsonRec := runExport(t, ExportHandlerDynamic(func() *areas.Config { return cfg }, store), query)
	var jsonRows exportResponse
	require.NoError(t, json.Unmarshal(jsonRec.Body.Bytes(), &jsonRows))
	require.Len(t, jsonRows.Data, len(rows)-1)
	for i, resource := range jsonRows.Data {
		assert.Equal(t, resource.Attributes.BookingDate, rows[i+1][3])
		assert.Equal(t, resource.Attributes.UserName, rows[i+1][0])
	}
}

func TestExportCSVValidation(t *testing.T) {
	t.Parallel()
	for _, query := range []string{
		"", "from=bad&to=2026-10-08", "from=2026-10-08", "from=2026-10-08&to=2026-10-01",
		"from=2026-10-01&to=2026-10-08&area=missing",
	} {
		t.Run(query, func(t *testing.T) {
			t.Parallel()
			cfg := exportConfig()
			values, err := url.ParseQuery(query)
			require.NoError(t, err)
			rec := runExport(t, ExportCSVHandler(func() *areas.Config { return cfg }, setupTestStore(t)), values)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "application/vnd.api+json", rec.Header().Get("Content-Type"))
			assert.Empty(t, rec.Header().Get("Content-Disposition"))
		})
	}
}

func TestExportCSVPersonsAndEmptyArea(t *testing.T) {
	t.Parallel()
	store := setupTestStore(t)
	seedExportUser(t, store, "colleague", "Zoe", "zoe@example.com", false)
	seedTestBookingFull(t, store, "on-behalf", "desk-1", "colleague", "booker", "2026-10-08")
	seedTestBookingWithGuest(t, store, "guest", "desk-2", "guest-id", "booker", "2026-10-08",
		true, "=HYPERLINK(\"x\")", "+email")
	cfg := exportConfig()
	for _, area := range []string{"office", "empty"} {
		rec := runExport(t, ExportCSVHandler(func() *areas.Config { return cfg }, store),
			url.Values{"filter[from]": {"2026-10-08"}, "filter[to]": {"2026-10-08"}, "filter[area]": {area}})
		rows, err := csv.NewReader(strings.NewReader(strings.TrimPrefix(rec.Body.String(), "\xef\xbb\xbf"))).ReadAll()
		require.NoError(t, err)
		if area == "empty" {
			assert.Len(t, rows, 1)
			continue
		}
		require.Len(t, rows, 3)
		assert.Equal(t, []string{"'=HYPERLINK(\"x\")", "'+email", "Cube 1 / Tisch 4", "2026-10-08"}, rows[1])
		assert.Equal(t, "Zoe", rows[2][0])
	}
}

type csvFailWriter struct{ remaining int }

func (w *csvFailWriter) Write(p []byte) (int, error) {
	if w.remaining == 0 {
		return 0, errors.New("write failed")
	}
	w.remaining--
	return len(p), nil
}

func TestWriteExportCSVErrors(t *testing.T) {
	t.Parallel()
	for _, successfulWrites := range []int{0, 1} {
		err := writeExportCSV(&csvFailWriter{remaining: successfulWrites}, exportConfig(), nil)
		require.ErrorContains(t, err, "write failed")
	}
}

func TestWriteExportCSVPreservesFieldCharacters(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"Alice\rSmith", "\r=1+2", "Alice\nSmith", "Alice\r\nSmith"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var output strings.Builder
			err := writeExportCSV(&output, exportConfig(), []ExportRow{{
				PersonName: name, ItemID: "desk-1", BookingDate: "2026-10-08",
			}})
			require.NoError(t, err)
			assert.Equal(t, "\xef\xbb\xbfName,Email,Booked item,Booking date\r\n\""+
				escapeCSVCell(name)+"\",,Cube 1 / Tisch 3,2026-10-08\r\n", output.String())
		})
	}
}
