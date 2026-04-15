# SitHub - A desk booking app for shared desk offices

SitHub is a user-friendly web application designed to facilitate desk bookings in shared office spaces.
It provides a seamless experience for both users and office administrators, ensuring efficient desk allocation and
management.

## Key Features

- **User-Friendly Interface**: Easy navigation and intuitive booking process on mobile and desktop.
- **Real-Time Availability**: View desk availability in real-time.
- **Notifications**: Receive alerts for upcoming bookings and changes.
- **Admin Dashboard**: Comprehensive tools for managing desk bookings and office operations.
- **Single Sign-On (SSO)**: Support for SSO integration with popular identity providers. Entra ID for now;
  other providers are coming soon.
- **Single Binary Distribution**: Deployable as a single binary for easy installation and management.
- **Built-in Database**: No external dependencies, ensuring minimal setup and maintenance overhead.

## Feature Breakdown

### Authentication

- Uses Entra ID for SSO integration, with plans to support additional providers in the future.
- Access to the app can be limited to a user group.
- Admin users are specified by Entra ID group membership.

### Test Authentication (Development Only)

For E2E and local development, you can bypass Entra ID with the `[test_auth]` section in `sithub.example.toml`.
This is only for local testing and must not be used in production.

- `test_auth.enabled = true` to enable test auth (or `SITHUB_TEST_AUTH_ENABLED=true`)
- `test_auth.user_id` to override the test user id (default: `test-user`)
- `test_auth.user_name` to override the display name (default: `Test User`)

### Desk Booking

- Locate available desks on an interactive floor plan.
- Users can book for a single day, an entire week, or a configurable number of days.
- Users can book for other users of the organization or guests not belonging to the organization without an account.
- Bookings can be made in advance or on the spot.
- Users can view and manage their bookings from the dashboard.
- Users can subscribe to notifications when someone books a desk in the same room.
- Rooms can be assigned to user groups for exclusive booking.

### User Interface

- Great user experience on mobile devices.
- Dark and light mode support.
- Multi-language support for English, German, Spanish, and French.

### Areas, Rooms, and Desk Management

- Areas are the physical locations where rooms are available.
- Rooms are the spaces within an area where desks are located.
- Desks are the individual workstations available for booking.
- Areas, rooms, desks, and desk equipment can be managed through a comprehensive
  [YAML configuration file](./sithub_areas.example.yaml).
- Point SitHub at the YAML file using `spaces.config_file` in `sithub.toml` or `--spaces-config-file`.
- Custom icons in the areas YAML file refer to [pictogrammers.com](https://pictogrammers.com/).
  If an item has no icon assigned, it inherits the icon from the higher-level area.

### Installation

- The application can be installed as a single binary, making deployment straightforward.
- No external dependencies are required, simplifying setup and maintenance.
- The built-in database ensures minimal setup and maintenance overhead.
- Download the app, create a configuration, define rooms and desks, start the server, set up a reverse proxy for SSL
  termination, and you are done.

## Tech Stack

- Go for the backend using the [Echo framework](https://echo.labstack.com/).
- Backend implements a clean REST API following the JSON:API specification.
- Vue 3 and Vuetify for the frontend.
- SQLite3 for the database.
- All frontend artifacts are embedded into the distributed binary.
- GitHub Actions for CI/CD.

### RESTful API

The backend implements a clean REST API following the JSON:API specification.
The API provides endpoints for managing areas, item groups, items, and bookings.
It supports CRUD operations and includes pagination and filtering capabilities.

You can view the API documentation by launching any OpenAPI viewer. Example:

```shell
npx redoc-cli serve ./docs/openapi.yaml
```

## FAQ

### Avatar sync from Entra ID fails

SitHub downloads profile photos from Microsoft Graph (`/me/photo/$value`) during
Entra ID login. Common reasons for sync failure:

- **Image format**: Microsoft Graph may return JPEG or PNG. SitHub decodes both formats
  and re-encodes as PNG. If the image uses an unsupported format (e.g., BMP, TIFF), the
  sync will fail and the fallback initials avatar is used instead.
- **Graph API permissions**: The Entra ID app registration needs the `User.Read` scope.
  Without it, the photo endpoint returns 403.
- **No profile photo**: If the user has no photo in Entra ID, Graph returns 404 and SitHub
  removes any stale avatar file. The initials avatar is shown.
- **Large photos**: Photos exceeding 512 KB are skipped. Ask users to resize their photo
  in Microsoft 365.

Check the server log for detailed diagnostics. Failed syncs log the user ID, HTTP status,
content-type, and response body size.
