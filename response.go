package tinybird

import (
	"encoding/json"
	"io"

	"github.com/olivere/ndjson"
)

// Basic JSON struct for all response by Tinybird.
type Response struct {
	Body                   io.ReadCloser // Body with the original data.
	Data                   []Row         `json:"data,omitempty"`                       // Data is part a tinybird response.
	Documentation          string        `json:"documentation,omitempty"`              // Documentation is part a tinybird response.
	Error                  string        `json:"error,omitempty"`                      // Error is part a tinybird response.
	Meta                   []Meta        `json:"meta,omitempty"`                       // Meta is part a tinybird response.
	Rows                   uint          `json:"rows,omitempty"`                       // Rows is part a tinybird response.
	RowsBeforeLimitAtLeast uint          `json:"rows_before_limit_at_least,omitempty"` // RowsBeforeLimitAtLeast is part a tinybird response.
	Statistics             Statistics    `json:"statistics,omitempty"`                 // Statistics is part a tinybird response.
	Status                 int           // Status is a HTTP status code, ej: 200, 400, 500 etc...
	NewLineDelimitedJSON   bool          // Save setting for NewLine-Delimited JSON.
}

// Generic row structure to allow any field with any type.
type Row map[string]any

// Specific field with data type.
type Meta struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// Statistics information about the response.
type Statistics struct {
	Elapsed   float64 `json:"elapsed,omitempty"`
	RowsRead  uint64  `json:"rows_read,omitempty"`
	BytesRead uint64  `json:"bytes_read,omitempty"`
}

// Convert response body to struct.
func (r *Response) Decode() error {
	if r.NewLineDelimitedJSON {
		return r.NDJSON()
	}

	return r.JSON()
}

// Convert body to NDJSON.
func (r *Response) NDJSON() error {
	var count uint
	var rows []Row

	ndjsonReader := ndjson.NewReader(r.Body)
	for ndjsonReader.Next() {
		var row Row
		if err := ndjsonReader.Decode(&row); err != nil {
			return err
		}
		rows = append(rows, row)
		count++
	}

	r.Data = rows
	r.Rows = count

	return nil
}

// Convert body to JSON.
func (r *Response) JSON() error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(body), &r)
}
