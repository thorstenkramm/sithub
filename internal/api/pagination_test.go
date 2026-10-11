//revive:disable-next-line var-naming
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newQueryContext(t *testing.T, rawQuery string) echo.Context {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/things?"+rawQuery, http.NoBody)
	return echo.New().NewContext(req, httptest.NewRecorder())
}

func TestParsePage(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		query   string
		want    Page
		wantErr string
	}{
		{name: "absent params use defaults", query: "", want: Page{Limit: DefaultPageLimit, Offset: 0}},
		{name: "valid values", query: "page[limit]=20&page[offset]=40", want: Page{Limit: 20, Offset: 40}},
		{name: "encoded brackets", query: "page%5Blimit%5D=5&page%5Boffset%5D=10", want: Page{Limit: 5, Offset: 10}},
		{name: "limit above max is capped", query: "page[limit]=101", want: Page{Limit: MaxPageLimit}},
		{name: "huge valid limit is capped", query: "page[limit]=9223372036854775807", want: Page{Limit: MaxPageLimit}},
		{name: "only offset", query: "page[offset]=3", want: Page{Limit: DefaultPageLimit, Offset: 3}},
		{name: "unknown page key ignored", query: "page[number]=2", want: Page{Limit: DefaultPageLimit}},
		{name: "limit zero", query: "page[limit]=0", wantErr: "page[limit]"},
		{name: "negative limit", query: "page[limit]=-1", wantErr: "page[limit]"},
		{name: "negative offset", query: "page[offset]=-1", wantErr: "page[offset]"},
		{name: "non-numeric limit", query: "page[limit]=abc", wantErr: "page[limit]"},
		{name: "non-numeric offset", query: "page[offset]=abc", wantErr: "page[offset]"},
		{name: "empty limit", query: "page[limit]=", wantErr: "page[limit]"},
		{name: "float limit", query: "page[limit]=1.5", wantErr: "page[limit]"},
		{name: "overflowing offset", query: "page[offset]=99999999999999999999", wantErr: "page[offset]"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			page, err := ParsePage(newQueryContext(t, tc.query))
			if tc.wantErr != "" {
				require.ErrorIs(t, err, ErrInvalidPagination)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, page)
		})
	}
}

const usersPath = "/api/v1/users"

func pageLink(t *testing.T, limit, offset int, extra url.Values) string {
	t.Helper()
	q := url.Values{}
	for k, v := range extra {
		q[k] = v
	}
	q.Set("page[limit]", strconv.Itoa(limit))
	q.Set("page[offset]", strconv.Itoa(offset))
	return usersPath + "?" + q.Encode()
}

func TestBuildPaginationLinks(t *testing.T) {
	t.Parallel()

	strPtr := func(s string) *string { return &s }

	cases := []struct {
		name  string
		query url.Values
		page  Page
		total int
		want  func(t *testing.T) PaginationLinks
	}{
		{
			name: "first page", page: Page{Limit: 20, Offset: 0}, total: 45,
			want: func(t *testing.T) PaginationLinks {
				return PaginationLinks{
					Self: pageLink(t, 20, 0, nil), First: pageLink(t, 20, 0, nil),
					Prev: nil, Next: strPtr(pageLink(t, 20, 20, nil)), Last: pageLink(t, 20, 40, nil),
				}
			},
		},
		{
			name: "last page", page: Page{Limit: 20, Offset: 40}, total: 45,
			want: func(t *testing.T) PaginationLinks {
				return PaginationLinks{
					Self: pageLink(t, 20, 40, nil), First: pageLink(t, 20, 0, nil),
					Prev: strPtr(pageLink(t, 20, 20, nil)), Next: nil, Last: pageLink(t, 20, 40, nil),
				}
			},
		},
		{
			name: "exact multiple of limit", page: Page{Limit: 20, Offset: 20}, total: 40,
			want: func(t *testing.T) PaginationLinks {
				return PaginationLinks{
					Self: pageLink(t, 20, 20, nil), First: pageLink(t, 20, 0, nil),
					Prev: strPtr(pageLink(t, 20, 0, nil)), Next: nil, Last: pageLink(t, 20, 20, nil),
				}
			},
		},
		{
			name: "empty collection", page: Page{Limit: 20, Offset: 0}, total: 0,
			want: func(t *testing.T) PaginationLinks {
				return PaginationLinks{
					Self: pageLink(t, 20, 0, nil), First: pageLink(t, 20, 0, nil),
					Prev: nil, Next: nil, Last: pageLink(t, 20, 0, nil),
				}
			},
		},
		{
			name: "offset not aligned to limit", page: Page{Limit: 20, Offset: 5}, total: 45,
			want: func(t *testing.T) PaginationLinks {
				return PaginationLinks{
					Self: pageLink(t, 20, 5, nil), First: pageLink(t, 20, 0, nil),
					Prev: strPtr(pageLink(t, 20, 0, nil)), Next: strPtr(pageLink(t, 20, 25, nil)),
					Last: pageLink(t, 20, 40, nil),
				}
			},
		},
		{
			name: "offset past the end", page: Page{Limit: 20, Offset: 100}, total: 45,
			want: func(t *testing.T) PaginationLinks {
				return PaginationLinks{
					Self: pageLink(t, 20, 100, nil), First: pageLink(t, 20, 0, nil),
					Prev: strPtr(pageLink(t, 20, 80, nil)), Next: nil, Last: pageLink(t, 20, 40, nil),
				}
			},
		},
		{
			name:  "extra query params preserved and self uses effective limit",
			query: url.Values{"from": {"2026-01-01"}, "page[limit]": {"500"}},
			page:  Page{Limit: MaxPageLimit, Offset: 0}, total: 150,
			want: func(t *testing.T) PaginationLinks {
				extra := url.Values{"from": {"2026-01-01"}}
				return PaginationLinks{
					Self: pageLink(t, 100, 0, extra), First: pageLink(t, 100, 0, extra),
					Prev: nil, Next: strPtr(pageLink(t, 100, 100, extra)),
					Last: pageLink(t, 100, 100, extra),
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			u := &url.URL{Path: usersPath, RawQuery: tc.query.Encode()}
			got := BuildPaginationLinks(u, tc.page, tc.total)
			assert.Equal(t, tc.want(t), got)
		})
	}
}

func TestWritePaginatedCollection(t *testing.T) {
	t.Parallel()

	c := newQueryContext(t, "page[limit]=2&page[offset]=2")
	rec, ok := c.Response().Writer.(*httptest.ResponseRecorder)
	require.True(t, ok)

	resources := []Resource{{ID: "c", Type: "things"}}
	require.NoError(t, WritePaginatedCollection(c, resources, Page{Limit: 2, Offset: 2}, 3, "write things"))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, JSONAPIContentType, rec.Header().Get(echo.HeaderContentType))

	var body struct {
		Data  []Resource      `json:"data"`
		Meta  PaginationMeta  `json:"meta"`
		Links PaginationLinks `json:"links"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Len(t, body.Data, 1)
	assert.Equal(t, PaginationMeta{Total: 3, Limit: 2, Offset: 2}, body.Meta)
	assert.NotNil(t, body.Links.Prev)
	assert.Nil(t, body.Links.Next)
}

func TestWriteCollectionOmitsPagination(t *testing.T) {
	t.Parallel()

	c := newQueryContext(t, "")
	rec, ok := c.Response().Writer.(*httptest.ResponseRecorder)
	require.True(t, ok)

	require.NoError(t, WriteCollection(c, []Resource{}, "write things"))
	assert.JSONEq(t, `{"data":[]}`, rec.Body.String())
}
