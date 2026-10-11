package bookings

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/labstack/echo/v4"

	"github.com/thorstenkramm/sithub/internal/api"
	"github.com/thorstenkramm/sithub/internal/areas"
)

// escapeCSVCell is security-critical: prefix spreadsheet formula/control characters with an
// apostrophe so untrusted names, emails and item labels are interpreted as text.
func escapeCSVCell(value string) string {
	if value != "" {
		switch value[0] {
		case '=', '+', '-', '@', '\t', '\r':
			return "'" + value
		}
	}
	return value
}

func writeExportCSV(output io.Writer, cfg *areas.Config, rows []ExportRow) error {
	if _, err := io.WriteString(output, "\xef\xbb\xbf"); err != nil {
		return fmt.Errorf("write csv BOM: %w", err)
	}
	if err := writeCSVRecord(output, []string{"Name", "Email", "Booked item", "Booking date"}); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}
	warned := map[string]bool{}
	for _, row := range rows {
		attrs := exportAttributes(cfg, &row, warned)
		if err := writeCSVRecord(output, []string{escapeCSVCell(attrs.UserName), escapeCSVCell(attrs.UserEmail),
			escapeCSVCell(attrs.BookedItem), attrs.BookingDate}); err != nil {
			return fmt.Errorf("write csv row: %w", err)
		}
	}
	return nil
}

// writeCSVRecord lets encoding/csv quote fields without normalizing their CR/LF characters.
// UseCRLF would strip bare CRs inside fields, so change only the record's final separator.
func writeCSVRecord(output io.Writer, fields []string) error {
	var record bytes.Buffer
	w := csv.NewWriter(&record)
	if err := w.Write(fields); err != nil {
		return fmt.Errorf("encode csv record: %w", err)
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("flush csv record: %w", err)
	}
	record.Truncate(record.Len() - 1)
	record.WriteString("\r\n")
	if _, err := record.WriteTo(output); err != nil {
		return fmt.Errorf("write csv record: %w", err)
	}
	return nil
}

// ExportCSVHandler downloads all matching bookings. Authorization belongs to the shared export
// route middleware. Both dates are required; pagination parameters are ignored, even if malformed.
func ExportCSVHandler(getConfig areas.ConfigGetter, store *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := getConfig()
		params, err := parseExportRequest(c, cfg, false)
		if err != nil {
			return api.WriteBadRequest(c, err.Error())
		}
		rows, _, err := ListExportBookings(c.Request().Context(), store, params.filter, 0, 0)
		if err != nil {
			return fmt.Errorf("list csv bookings: %w", err)
		}
		h := c.Response().Header()
		h.Set("Content-Type", "text/csv; charset=utf-8")
		h.Set("Cache-Control", "no-store")
		h.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="sithub-bookings-%s_%s.csv"`,
			params.filter.From, params.filter.To))
		return writeExportCSV(c.Response(), cfg, rows)
	}
}
