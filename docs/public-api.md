# SitHub Public API

Use API keys to read booking data from finance tools and scripts. This API is intended for server-side
clients. Cross-origin browser access is not enabled.

## Create an API key

An administrator opens the user menu (or mobile drawer), selects **API Keys**, then **Add key**.
Enter a name, choose an expiry of 3, 12 or 24 months, and select **Read bookings**.

> [!WARNING]
> The full key is shown only once. Copy it into your secret store before selecting **Done**.
> Only the first 10 characters remain visible in the key list. A lost key must be revoked and replaced.

## Authenticate

Send `Authorization: Bearer <key>` over HTTPS. The bearer credential takes precedence over a session
cookie: an invalid key returns 401 even when a valid administrator cookie is present. Keys are distinct
machine identities, not user accounts. Send the exact key as issued; changing its case invalidates it.

Other Authorization schemes, such as Basic authentication from a reverse proxy, do not authenticate an
API key. Session-based access still requires a permitted user with the data exporter role.

> [!IMPORTANT]
> Every request below `/api/` that carries a Bearer header is checked as an API key, including the
> public ping and login endpoints. A reverse proxy must not inject its own `Authorization: Bearer`
> header into SitHub requests, otherwise those requests return 401.

## Scopes

| Scope | Grants |
| --- | --- |
| Read bookings (`read-bookings`) | GET booking export as JSON or CSV, GET user avatars |

Only the two export endpoints and the avatar endpoint below accept keys. Other protected endpoints return
403, including user management, key management, `/me`, and booking creation. Public ping and authentication endpoints
accept a valid key but grant it no user permissions. An invalid key returns 401 on every endpoint.

## Endpoints

### GET /api/v1/bookings/export

Returns JSON:API (`application/vnd.api+json`) with `data`, pagination `meta` and `links`.

| Parameter | Meaning |
| --- | --- |
| `filter[from]` | Required first booking date, inclusive, `YYYY-MM-DD` |
| `filter[to]` | Required last booking date, inclusive, `YYYY-MM-DD` |
| `filter[area]` | Optional top-level area ID; omit or leave empty for all areas |
| `page[limit]` | Positive page size, default 20, capped at 100 |
| `page[offset]` | Non-negative offset, default 0 |

The filters use the JSON:API `filter[...]` names. Plain `from`/`to` parameters are not accepted, because
privacy browser extensions remove a `from` parameter from URLs.

Results are ordered by booking date, case-insensitive person name, then booking ID. `user_id`, name and
email describe the person the booking is for, including guests and colleagues. `user_id` is empty for
guests and for users that no longer exist; both have no avatar. Deactivated users' bookings
are included. Locations removed from configuration retain their raw item IDs.

```json
{
  "data": [
    {
      "type": "bookings",
      "id": "booking-id",
      "attributes": {
        "booking_date": "2026-10-08",
        "user_id": "3edbc42b-6157-493d-aad6-bb55120f94dc",
        "user_name": "Alex",
        "user_email": "alex@example.com",
        "is_guest": false,
        "item_id": "desk_1",
        "item_name": "Desk 1",
        "item_group_id": "test_room",
        "item_group_name": "Test Room",
        "area_id": "test_area",
        "area_name": "Test Area",
        "booked_item": "Test Room / Desk 1"
      }
    }
  ],
  "meta": { "total": 1, "limit": 20, "offset": 0 },
  "links": {
    "self": "/api/v1/bookings/export?filter%5Bfrom%5D=2026-10-01&filter%5Bto%5D=2026-10-08&page%5Blimit%5D=20&page%5Boffset%5D=0",
    "first": "/api/v1/bookings/export?filter%5Bfrom%5D=2026-10-01&filter%5Bto%5D=2026-10-08&page%5Blimit%5D=20&page%5Boffset%5D=0",
    "prev": null,
    "next": null,
    "last": "/api/v1/bookings/export?filter%5Bfrom%5D=2026-10-01&filter%5Bto%5D=2026-10-08&page%5Blimit%5D=20&page%5Boffset%5D=0"
  }
}
```

Follow `links.next` until it is null. Pagination is not a snapshot: concurrent booking changes may
shift rows between pages. CSV returns all matching rows in a single response.

### GET /api/v1/avatars/{user_id}

Returns the profile photo of a user as `image/png`. Use the `user_id` of a booking from the JSON export, for
example to show today's presence on door signage. The response carries `Cache-Control: max-age=300`.
A user without a photo returns 404, so show initials or a placeholder instead. The CSV export has no
user ID. Keys can read avatars but cannot upload or delete them.

### GET /api/v1/bookings/export.csv

Uses the same `filter[from]`, `filter[to]` and `filter[area]` filters. Pagination parameters are ignored.
The response is `text/csv; charset=utf-8`, UTF-8 with a BOM, comma-separated with CRLF record endings and
RFC 4180 quoting.
Columns are **Name**, **Email**, **Booked item**, **Booking date**. Ordering matches the JSON export.
Free-text cells beginning with `=`, `+`, `-`, `@`, TAB or CR are prefixed with an apostrophe to prevent
spreadsheet formula execution.

The download filename is `sithub-bookings-<from>_<to>.csv`.

## Errors

| Status | Meaning |
| --- | --- |
| 400 | Bad or missing dates, reversed range, unknown area, bad pagination |
| 401 | No credentials, or a malformed, unknown, revoked or expired key |
| 403 | Key scope or endpoint not permitted; session lacks exporter role |

Rejected bearer credentials include `WWW-Authenticate: Bearer error="invalid_token"`. A key without the
required scope gets 403 with `WWW-Authenticate: Bearer error="insufficient_scope"`.
Error bodies use JSON:API:

```json
{
  "errors": [
    {
      "status": "401",
      "title": "Unauthorized",
      "detail": "API key has expired",
      "code": "auth_required"
    }
  ]
}
```

## Expiry and revocation

Expiry is computed in UTC using calendar months from creation. Month-end overflow normalizes into the
following month (for example, November 30 plus three months can be March 2). A key is invalid at its
exact expiry instant and thereafter; expired keys remain listed until revoked.

Administrators revoke a key on **API Keys** by selecting its revoke button and confirming. Revocation
deletes the key and takes effect on subsequent requests; it does not cancel an export already authorized.
For rotation, create a new key, update your integration, then revoke the old key. Deleting or demoting
the creating administrator does not revoke the key automatically.

## Examples

Populate `SITHUB_API_KEY` from your secret store or a secure environment-variable prompt rather than
typing the key into shell history. Replace the host and `test_area` with your deployment's values.
`-g` disables curl's bracket expansion for the `filter[...]` and `page[...]` parameters.

```bash
export SITHUB_API_KEY="<your-key>"
export SITHUB_URL="http://localhost:9900"

# Get Bookings as JSOAN
curl -g --fail-with-body -H "Authorization: Bearer ${SITHUB_API_KEY}" \
  "${SITHUB_URL}/api/v1/bookings/export?filter[from]=2026-10-01&filter[to]=2026-10-08&page[limit]=50&page[offset]=0"

# Get Bookings as CSV
curl -g --fail-with-body -H "Authorization: Bearer ${SITHUB_API_KEY}" -o bookings.csv \
  "${SITHUB_URL}/api/v1/bookings/export.csv?filter[from]=2026-10-01&filter[to]=2026-10-08&filter[area]=test_area"
  
# Get the avatar of a user
curl -g --fail-with-body -H "Authorization: Bearer ${SITHUB_API_KEY}"   \
  "${SITHUB_URL}/api/v1/avatars/<USER-ID-FROM-BOOKINGS>"
```
