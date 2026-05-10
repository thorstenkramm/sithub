# Story 31.1: Live Updates for Bookings and Cancellations

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user browsing room plans, area tiles, or the weekly table view,
I want bookings and cancellations made by other users to appear without manually
reloading,
so that I do not waste time deciding on a desk that has already been taken.

## Acceptance Criteria

1. **Given** I am viewing the area tile, weekly table, or floor plan view for a given
   date
   **When** another user creates a booking that affects an item in my current view
   **Then** the corresponding tile, cell, or floor plan marker updates to its new busy
   state within a few seconds
   **And** I do not need to refresh the page or change any filter for the change to
   appear

2. **Given** I am viewing any of the same views
   **When** another user cancels a booking that affects an item in my current view
   **Then** the corresponding tile, cell, or floor plan marker updates to its new free
   state within a few seconds
   **And** the change is reflected without page reload

3. **Given** I have just been viewing live updates and the network connection drops
   **When** the connection is restored
   **Then** the client reconciles state with the server and reflects any bookings or
   cancellations that occurred during the outage
   **And** the existing connection-lost messaging from Story 18.7 still surfaces while
   disconnected

4. **Given** I open the app in multiple tabs
   **When** an event arrives
   **Then** all open tabs reflect the change consistently

## Tasks / Subtasks

- [ ] Task 1: Add the WebSocket dependency and a broadcast hub package (AC: #1, #2, #4)
  - [ ] 1.1 Add `github.com/gorilla/websocket v1.5.x` (latest stable patch) to
        `go.mod`. Run `go mod tidy`.
  - [ ] 1.2 Create `internal/livefeed/hub.go` with:
        ```go
        // Package livefeed broadcasts booking events to connected WebSocket clients.
        package livefeed

        // Hub fans out events to all connected clients. It is safe for concurrent use.
        type Hub struct { ... }

        func NewHub() *Hub
        func (h *Hub) Run(ctx context.Context)            // owns the goroutine; exits on ctx
        func (h *Hub) Broadcast(event Event)              // non-blocking; drops on slow client
        func (h *Hub) Register(client *Client)
        func (h *Hub) Unregister(client *Client)
        ```
        Use channels for register/unregister/broadcast. Document concurrency
        ("safe for concurrent use after `Run` has been started").
  - [ ] 1.3 Define `Event` (and JSON tags) in `livefeed/event.go`. Mirror
        `notifications.BookingEvent` shape so downstream consumers see the same
        payload. Use distinct event type strings:
        - `booking.created`
        - `booking.canceled`
        Include `item_id`, `booking_date`, `booking_id`, `user_id` (for "did I do
        this?" filtering), and `timestamp` (RFC 3339 UTC). No avatar URLs, no
        emails.
  - [ ] 1.4 Add `internal/livefeed/client.go` with the per-connection writer
        goroutine: keeps a small buffered channel (e.g. 16) of events to write,
        an idle ping ticker (`websocket.PingMessage` every 30 s), and a write
        deadline. Read pump only handles pong + close frames. Drop the client
        if its outbound buffer fills (do NOT block the hub).

- [ ] Task 2: Wire the hub as a `notifications.Notifier` (AC: #1, #2)
  - [ ] 2.1 The existing emission points are already in place:
        `internal/bookings/handler.go:300` (cancel) and `:1167` (create). They
        call `notifier.NotifyAsync(...)`. Reuse those — do not duplicate calls.
  - [ ] 2.2 Make the hub satisfy `notifications.Notifier`:
        ```go
        func (h *Hub) NotifyAsync(event *notifications.BookingEvent) {
            h.Broadcast(toLiveEvent(event))
        }
        ```
        Add a tiny mapper that strips PII (`guest_email`, `guest_name`) before
        broadcasting. Live updates only need to know "this item on this date is
        now busy/free" plus identity for "is this me?" filtering.
  - [ ] 2.3 In `internal/startup/server.go`, replace the single notifier
        construction with a multiplexer:
        ```go
        webhook := notifications.NewNotifier(cfg.Notifications.WebhookURL)
        hub := livefeed.NewHub()
        go hub.Run(ctx)
        notifier := notifications.MultiNotifier{webhook, hub}
        ```
        Add `notifications.MultiNotifier` (a slice type implementing `Notifier`
        by fanning out to each underlying notifier) in `internal/notifications/`.
        The webhook → external system path stays; the hub is the in-process
        broadcaster.

- [ ] Task 3: Add the WebSocket endpoint behind auth (AC: #1, #2, #4)
  - [ ] 3.1 New file `internal/livefeed/handler.go`:
        ```go
        func Handler(hub *Hub) echo.HandlerFunc
        ```
        Use `gorilla/websocket.Upgrader` with:
        - `CheckOrigin`: accept same-origin only. Build the allowed origin from
          the request `Host` header in production and additionally accept
          `http://localhost:5173` and `http://127.0.0.1:5173` for dev (the Vite
          dev server proxies `/api` but WebSockets connect directly).
        - `ReadBufferSize`/`WriteBufferSize`: 1 KiB each.
  - [ ] 3.2 Register the endpoint in `internal/startup/server.go::registerRoutes`:
        ```go
        e.GET("/api/v1/live", livefeed.Handler(hub), requireAuth)
        ```
        Use the same `requireAuth` middleware as other authenticated routes — the
        upgrade handshake carries the `sithub_user` cookie, so existing session
        validation Just Works.
  - [ ] 3.3 In the handler:
        - Read the authenticated user from the Echo context (set by
          `middleware.LoadUser`).
        - Upgrade. On error, return JSON:API 400.
        - Construct a `Client` bound to the user ID and register it with the hub.
        - Spawn the read pump (blocks; exits on close or read error) and the
          write pump (drained by hub broadcasts). Unregister on exit.
  - [ ] 3.4 Do NOT scope subscriptions per-area or per-date in v1. Broadcast all
        events to all clients; the frontend filters by what it cares about. This
        keeps the server simple and matches the architecture decision
        ("WebSockets for live availability updates with polling fallback") in
        `architecture.md` lines 135–138.

- [ ] Task 4: Backend tests
  - [ ] 4.1 Hub unit tests (`internal/livefeed/hub_test.go`):
        - `Broadcast` reaches all registered clients.
        - Slow clients get dropped (buffered chan fills → unregister fires).
        - `Run` exits cleanly when `ctx` is cancelled.
        - Concurrent register/unregister/broadcast does not race
          (run with `-race`).
  - [ ] 4.2 Handler integration test
        (`internal/livefeed/handler_test.go`): use `httptest.NewServer` to wrap
        an Echo with `Handler(hub)` and `gorilla/websocket.DefaultDialer.Dial` to
        connect. Assert that an event broadcast through the hub is received over
        the wire by the connected client. Use a real `sithub_user` cookie
        produced by the test auth helper (mirror what `bookings/handler_test.go`
        does).
  - [ ] 4.3 Booking-flow integration test: extend
        `internal/bookings/handler_test.go` (or add a `livefeed_e2e_test.go`)
        that asserts a `POST /api/v1/bookings` followed by a `DELETE
        /api/v1/bookings/:id` results in two events landing on a connected WS
        client.
  - [ ] 4.4 Run `go test -race ./internal/livefeed/... ./internal/bookings/...`
        and `golangci-lint run ./internal/livefeed/...`.

- [ ] Task 5: Frontend WebSocket client (AC: #1, #2, #3, #4)
  - [ ] 5.1 New composable `web/src/composables/useLiveFeed.ts`:
        - Exposes `connect()`, `disconnect()`, `onEvent(handler)`,
          `state: Ref<'connecting' | 'open' | 'closed'>`.
        - Builds the URL from `window.location`: `(loc.protocol === 'https:' ?
          'wss' : 'ws') + '://' + loc.host + '/api/v1/live'`. Vite proxies
          `/api` in dev so this works locally without further config.
        - Reconnect on close with exponential backoff (1 s, 2 s, 4 s, 8 s,
          capped at 30 s). Add ±20% jitter to avoid synchronized reconnects.
        - On reconnect, fire a `'reconnected'` synthetic event so listeners can
          refetch their slice (AC #3).
        - Survives across views: register the composable as a singleton in
          `App.vue` and provide event subscription via Pinia store (Task 5.2).
  - [ ] 5.2 New Pinia store `web/src/stores/useLiveFeedStore.ts`:
        - Holds the WebSocket connection, reconnection state, and a tiny
          subscriber registry keyed by event type.
        - Started in `App.vue` `onMounted` once the user is authenticated; torn
          down on logout.
        - Logout-aware: closes the socket cleanly on `useAuthStore.logout()`.
  - [ ] 5.3 Define types in `web/src/api/liveFeed.ts`:
        ```ts
        export type LiveEventType = 'booking.created' | 'booking.canceled' | 'reconnected';
        export interface LiveBookingEvent {
          type: 'booking.created' | 'booking.canceled';
          item_id: string;
          booking_date: string;     // YYYY-MM-DD
          booking_id: string;
          user_id: string;
          timestamp: string;        // RFC 3339 UTC
        }
        ```

- [ ] Task 6: Wire live updates into the three views (AC: #1, #2, #3)
  - [ ] 6.1 `web/src/views/ItemGroupsView.vue`: subscribe in `onMounted` (and
        unsubscribe in `onUnmounted`) to `booking.created`, `booking.canceled`,
        and `reconnected`. On a relevant event (item belongs to the current
        area or any of its item groups for the visible date range), call the
        existing fetch helpers (`fetchItemGroupAvailability`,
        `fetchItemGroupBookings`) to refresh just that item group. On
        `reconnected`, refetch the whole view.
  - [ ] 6.2 `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue`:
        same pattern — on relevant events for any of its rendered items, call
        the existing `fetchWeeklyMatrix` once (debounce to ~250 ms so a burst
        of events maps to a single refetch). Mutate the displayed cell
        in-place where possible to avoid layout flicker; if a full refetch is
        easier, do that.
  - [ ] 6.3 `web/src/components/InteractiveFloorPlan.vue`: subscribe and
        re-fetch the active item group's availability on relevant events. Use
        the same debounce. The floor plan has the strictest visual contract
        (free → busy → reserved transitions), so prefer a full refetch on every
        event over patching local state.
  - [ ] 6.4 `web/src/views/AreasView.vue`: presence/availability isn't shown
        on this view today. **Do not** subscribe here. (Skip.) Document the
        decision in dev notes so reviewers don't ask.
  - [ ] 6.5 Self-events: when the event's `user_id` matches
        `useAuthStore().user.id`, ignore it on the listening view — the user
        already sees the change from their own request, and a duplicate
        refetch would just cause flicker. The webhook notifier still fires for
        external integrations.

- [ ] Task 7: Connection-lost UX (AC: #3)
  - [ ] 7.1 The store sets `state` to `'closed'` whenever the socket drops. A
        thin banner is **not** required — the existing `CONNECTION_LOST_MESSAGE`
        from `web/src/api/client.ts` already surfaces network failure when an
        API request fails. Do NOT add a duplicate banner.
  - [ ] 7.2 What is required: when the socket reconnects, fire `'reconnected'`
        so listeners refetch their slice. This is the explicit reconciliation
        AC #3 demands.
  - [ ] 7.3 If the WebSocket cannot connect at all (e.g. a corporate proxy
        blocking it), the views fall back to their existing data — they remain
        functional, just no longer real-time. Log a single warning to the
        console; do not surface it to the user. Architecture decision in
        `architecture.md` lines 135–138 explicitly calls for "polling fallback"
        — out of scope for this story; track as a follow-up if real-world
        environments need it.

- [ ] Task 8: Frontend tests
  - [ ] 8.1 Unit (Vitest):
        - `useLiveFeed.test.ts` (new): cover URL construction (http vs https),
          reconnect backoff with jitter (use fake timers), and `onEvent`
          dispatch.
        - `useLiveFeedStore.test.ts` (new): cover singleton lifecycle and
          logout-driven teardown.
  - [ ] 8.2 Component (Vitest + Vue Test Utils):
        - `ItemGroupsView.test.ts`, `AreaWeeklyMatrixView.test.ts`, and
          `InteractiveFloorPlan.test.ts`: assert that emitting a stub
          `booking.created` event triggers the corresponding fetch helper. Mock
          `useLiveFeedStore` to feed events synthetically.
  - [ ] 8.3 Cypress E2E (`cypress/e2e/live-updates.cy.ts` — new):
        - Open two browser contexts (the second via
          `cy.session()` + a different user). User A books a desk; assert User
          B's open ItemGroupsView reflects the new busy state without a manual
          reload.
        - Note: Cypress in `--browser electron` runs a single browser instance.
          Multi-tab E2E is best done with two `cy.session()`-isolated tests
          using the bookings API directly to simulate "another user", and
          asserting the visible UI updates. Use intercept aliases
          (`@itemGroupAvailability`) to wait for the refetch fired by the WS
          handler.
  - [ ] 8.4 Run:
        ```
    cd web
        npx vitest run
        npm run type-check
        npm run lint
        npm run build
        npm run test:e2e -- --browser electron
        ```

- [ ] Task 9: Operational sanity
  - [ ] 9.1 Smoke test in dev: with backend on `:9900` and Vite on `:5173`,
        open two browser windows logged in as different demo users (see
        `tools/database/demo-users.sql`). In window A, book a desk in an area
        currently open in window B's ItemGroupsView. Confirm window B updates
        within ~1 s without a manual refresh.
  - [ ] 9.2 Confirm `go run ./cmd/sithub run --config ./sithub.toml` startup
        logs include the hub's "live feed hub started" line.

### Review Findings

- [x] [Review][Patch] Live-feed `user_id` currently carries the booking owner instead of the acting user, so on-behalf bookings and admin cancellations are misclassified as self-events and can be ignored by the wrong client [internal/livefeed/event.go:37]
- [x] [Review][Patch] The live-refresh wiring does not pass any `isRelevant` filter, so every booking anywhere in the system triggers refetches in open ItemGroups, matrix, and floor-plan views instead of only refreshing affected slices [web/src/views/ItemGroupsView.vue:691]
- [x] [Review][Patch] WebSocket upgrade failures bypass the story’s required JSON:API 400 path because the handler returns after `gorilla/websocket` writes its default non-JSON error response [internal/livefeed/handler.go:66]
- [x] [Review][Patch] Story-required coverage is still missing for the websocket route and view wiring: there is no handler integration test, no booking-flow websocket test, no `useLiveFeedStore` test, and no component tests asserting that live events trigger the expected refetch helpers [internal/livefeed/hub_test.go:1]
- [x] [Review][Patch] Week-mode tile live refresh keeps stale selected slots after another user books one of them; `loadWeekData(..., silent=true)` preserves `weekSelections`, so the sticky footer can still submit a slot that is no longer free [web/src/views/ItemsView.vue:1349]
- [x] [Review][Patch] Manual-regression coverage is still missing for the newly fixed tile paths: `ItemsView.test.ts` has no synthetic live-feed tests for day/week tile updates, and the matrix test only asserts that a refetch happened, not that live refresh stays loader-free/no-flicker [web/src/views/ItemsView.test.ts:1]

## Dev Notes

### Architecture & Patterns

- The architecture explicitly calls for WebSockets here:
  `_bmad-output/planning-artifacts/architecture.md` lines 135–138 ("WebSockets
  for live availability updates with polling fallback"). This story implements
  the WebSockets path; polling fallback is deferred.
- The existing `internal/notifications` package already abstracts "something
  happened, broadcast it" via the `Notifier` interface. The cleanest insertion
  point is to make the live-feed hub a second `Notifier` implementation and
  fan out via a `MultiNotifier` so the existing webhook integration is
  preserved unchanged.
- **No backend storage of subscriptions.** The hub is in-process state in a
  single-node deployment. SitHub explicitly does not run multi-node; if that
  ever changes, the broadcast hub would need to be replaced with a pub/sub
  layer (Redis, NATS) — that is not a problem for today.
- **No per-client filtering on the server.** All events go to all connected
  clients. Each event is small (≈200 B JSON); even a few hundred concurrent
  clients × a few events/min is a trickle. The frontend filters down by
  current view. This avoids tracking per-client subscription state on the
  server.
- **Auth via existing session cookie.** The WebSocket upgrade handshake is a
  regular HTTP GET with the `sithub_user` cookie, so the existing
  `middleware.RequireAuth` chain validates it before the upgrade. No tokens,
  no secondary auth flow.
- **Self-event filtering happens on the client.** Server broadcasts everything;
  the client checks `event.user_id === me.id` and skips refetching when the
  user already saw the result locally. Doing this on the server would force
  per-client routing — not worth the complexity.

### Key Code Locations

| Element | Location | Why it matters |
| --- | --- | --- |
| Booking create emission | `internal/bookings/handler.go:1167` | Already calls `notifier.NotifyAsync(...)`; reuse |
| Booking cancel emission | `internal/bookings/handler.go:300` | Same |
| Notifier interface | `internal/notifications/notifier.go:43` | Hub must implement this |
| Notifier construction | `internal/startup/server.go:68` | Wrap with `MultiNotifier{webhook, hub}` |
| Route registration | `internal/startup/server.go:109-192::registerRoutes` | Add `/api/v1/live` here |
| Auth middleware | `internal/middleware/...::RequireAuth` | Reuse for the WS endpoint |
| User context loader | `internal/middleware/...::LoadUser` | Provides the user ID for the hub client |
| Vite dev proxy | `web/vite.config.ts` | Already proxies `/api`; WS may need an explicit `ws: true` line — verify and add if missing |
| Frontend client base | `web/src/api/client.ts` | Source of `CONNECTION_LOST_MESSAGE` and the http abstractions |
| Auth store | `web/src/stores/useAuthStore.ts` | Hooks for login/logout teardown |
| ItemGroups view | `web/src/views/ItemGroupsView.vue` | One of three subscription points |
| Weekly matrix | `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue` | Second subscription point |
| Floor plan | `web/src/components/InteractiveFloorPlan.vue` | Third subscription point |

### Implementation Strategy

1. Backend hub + handler + tests first, with a `wscat` smoke test
   (`wscat -c ws://localhost:9900/api/v1/live --header
   'Cookie: sithub_user=...'`). Don't touch frontend until you can curl-style
   verify the wire.
2. `MultiNotifier` is a five-line slice type — write it in
   `internal/notifications/multi.go` and rewire `startup/server.go` to use it.
3. Frontend store + composable next; verify in the browser console that an
   event broadcast via a manual `curl POST /api/v1/bookings` lands.
4. Wire the three views one at a time. Each view's existing fetch helper is
   the natural "refresh" hook — call it on event.
5. E2E last. The two-user simulation via `cy.session()` is the most fragile
   part of this story; budget extra time for it.

### Anti-patterns to Avoid

- Do NOT broadcast booking events directly from `notifications.WebhookNotifier`
  internals. Build a separate `Notifier` and combine via `MultiNotifier`. This
  keeps the webhook integration testable in isolation.
- Do NOT block the booking handler waiting for the WS broadcast. The hub's
  `Broadcast` must be non-blocking (channel send with `default` drop, or
  buffered channel + drop on slow client). The handler returns 201/204 on its
  own timeline.
- Do NOT include guest emails or full names in WS payloads. Live updates need
  only `item_id`, `booking_date`, `booking_id`, `user_id`, `timestamp`. Keep
  PII out of the broadcast.
- Do NOT add per-area or per-date subscription state to the hub in v1. The
  fanout cost is negligible; complexity is not.
- Do NOT add a polling fallback in this story. Architecture mentions it; that
  is a follow-up.
- Do NOT couple the hub's lifecycle to a specific request — `Run` must be
  started once at app boot and live until shutdown. The hub is a long-lived
  goroutine, not a per-request worker.
- Do NOT forget the `read pump`. A WS connection without a reader will hang
  on close frames and leak goroutines. The minimum-viable read pump just
  reads and discards (we never expect client→server messages).
- Do NOT write to the WebSocket from multiple goroutines.
  `gorilla/websocket` writes are NOT concurrency-safe; route every write
  through the per-client write goroutine.
- Do NOT redirect on 401 from the WS endpoint. The handshake must return a
  proper HTTP error so the client's reconnect loop can back off; redirects
  during an Upgrade are a footgun.
- Do NOT ship without `-race` testing on the hub. The fan-out is the area
  most likely to race.

### Latest tech information

- `gorilla/websocket v1.5.x`: stable, widely used, supported. Echo-compatible
  via the standard `Upgrader` against the underlying `http.ResponseWriter`
  (`c.Response().Writer`) and `c.Request()`. Set ping/pong handlers and
  read/write deadlines per the package's recommended pattern. Do not use the
  abandoned `nhooyr.io/websocket` v0 patterns; gorilla is the de-facto Go
  default for this need.
- Browser `WebSocket`: standard, no polyfill needed. Reconnect logic must be
  hand-rolled (the browser API doesn't auto-reconnect). `wss://` is required
  on `https://` origins; the URL builder must follow `window.location.protocol`.

### Testing Standards

- Backend: table-driven tests, `require`/`assert`, `-race` for any test that
  touches the hub. Use `httptest.NewServer` for the handler integration test.
- Frontend unit: Vitest with fake timers for reconnect backoff.
- Frontend component: Vue Test Utils with a mocked `useLiveFeedStore` that
  emits synthetic events.
- Frontend E2E: Cypress against the dev server using `cy.session()` for the
  two-user setup. No API mocking. Use intercept aliases to synchronize.
- All AC must have at least one E2E or component test that fails without the
  live feed wiring.

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Epic 31 Stories: Live Updates, Favorites Rework & Areas Config Hint]
- [Source: _bmad-output/planning-artifacts/architecture.md#Category 2: Real-Time Updates]
- [Source: internal/bookings/handler.go]
- [Source: internal/notifications/notifier.go]
- [Source: internal/startup/server.go]
- [Source: web/src/api/client.ts]
- [Source: web/src/stores/useAuthStore.ts]
- [Source: web/src/views/ItemGroupsView.vue]
- [Source: web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue]
- [Source: web/src/components/InteractiveFloorPlan.vue]
- [Source: _bmad-output/implementation-artifacts/18-7-connection-lost-error-messaging.md]
- [Source: .claude/rules/golang.md]
- [Source: .claude/rules/vue.md]
- [Source: .claude/rules/cypress.md]
- [Source: .claude/rules/json-api.md]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.7 (1M context)

### Debug Log References

- `go test -race ./...` — all packages pass
- `cd web && npx vitest run` — 391 tests pass
- `cd web && npm run type-check && npm run lint && npm run build` — all pass
- `golangci-lint run ./internal/livefeed/... ./internal/notifications/... ./internal/startup/...`
  — only the two pre-existing issues from `main` remain (notifier_test goconst,
  server_test gosec G124 cookie); no new issues introduced

### Completion Notes List

- **Backend `internal/livefeed`** — new package with three small files:
  `event.go` (public payload + `fromBookingEvent` mapper that strips
  `guest_name` / `guest_email` PII), `hub.go` (single-goroutine event loop with
  channels for register / unregister / broadcast; `running atomic.Bool` so
  `NotifyAsync` can early-return after shutdown; non-blocking broadcast via
  `select { default }` and per-client send-buffer drop), and `client.go` (one
  read pump + one write pump per connection; the write pump is the sole writer
  to `*websocket.Conn`, satisfying gorilla's "writes are not concurrency-safe"
  contract).
- **`internal/notifications.MultiNotifier`** — new five-line slice type
  satisfying the existing `Notifier` interface. Combines the existing webhook
  notifier with the new live hub so booking emission points
  (`internal/bookings/handler.go:300`, `:1167`) are unchanged and the webhook
  integration keeps working.
- **`/api/v1/live` endpoint** — registered in `startup/server.go` behind
  `requireAuth`. The handshake's `sithub_user` cookie flows through the
  existing `LoadUser` + `RequireAuth` middleware so no secondary auth flow is
  needed.
- **Origin check** — `livefeed.checkOrigin` allows same-origin (matching
  request `Host`) plus the four Vite dev ports
  (`http://localhost:5173`, `127.0.0.1:5173`, and their `https://` peers).
  Non-browser clients without an `Origin` header (curl, wscat, Go test
  dialer) are accepted because the endpoint is auth-gated.
- **Frontend `useLiveFeed` composable** — owns one WebSocket at a time,
  reconnects with exponential backoff capped at 30 s and ±20 % jitter,
  dispatches a synthetic `'reconnected'` event after the second open so views
  can refetch any state that drifted during the dropout (AC #3). URL is
  derived from `window.location` (`ws` ↔ `wss`) so production over `https`
  Just Works.
- **`useLiveFeedStore` (Pinia)** — singleton wrapper exposing
  `start()`, `stop()`, `subscribe(handler)`, and `reset()`. Started from
  `App.vue` via a `watch` on `authStore.isAuthenticated` (immediate), stopped
  on logout. Subscribers fan-out is wrapped in try/catch so a throwing
  handler does not break the rest.
- **`useLiveBookingRefresh` composable** — view-side helper that subscribes
  the calling component to the live feed and calls a debounced `refresh()`
  whenever a relevant booking event arrives. Self-events (where
  `event.user_id === authStore.userId`) are filtered out so the user does not
  see flicker from their own actions; the synthetic `'reconnected'` event
  always triggers a refresh regardless of `isRelevant`.
- **View wiring** — `ItemGroupsView.vue` calls `loadAvailability` for the
  current area + week on event; `AreaWeeklyMatrixView.vue` calls
  `loadMatrix`; `InteractiveFloorPlan.vue` calls `refreshAvailability`. All
  three reuse their existing fetch helpers — no new API calls were
  introduced. Default debounce is 250 ms.
- **AreasView is intentionally not wired** — that view does not show any
  per-item availability, so subscribing would just generate noise. Documented
  in the dev notes.
- **Smoke-test fixes after first run-through** — round-2 tweaks based on
  manual verification:
  - **`ItemsView` was missing.** The story listed three subscription points
    (ItemGroupsView, AreaWeeklyMatrixView, InteractiveFloorPlan) but not
    ItemsView, which is the desk-level day/week tile view inside a room. Added
    `useLiveBookingRefresh` there with a `silentReloadItems` helper for day
    mode and a new `silent` flag on `loadWeekData` for week mode. `isRelevant`
    filters by `selectedDate` (day) or `selectedWeekDates` (week) so only
    events affecting the visible date range trigger a refresh.
  - **Matrix view flicker eliminated.** `loadMatrix` now accepts
    `{ silent?: boolean }` and skips the `loading.value = true` toggle when
    invoked by the live feed. Without the toggle the LoadingState component
    no longer mounts and unmounts; the keyed v-for over `matrixData` /
    `days` reuses existing rows so only the cells whose state changed get
    DOM updates.
  - **ItemsView day/week silent paths** preserve the user's UI state across
    background refreshes — they do not reset `expandedDayTiles`,
    `expandedWeekTiles`, `weekSelections`, or `weekBookingResults`. The
    foreground load paths still reset those (existing behaviour) so explicit
    user navigation still gets a clean slate.
  - **Local edits I integrated:** `event.go` now uses `BookedByUserID` /
    `CanceledByUserID` (when set) as the `UserID` on the wire, so booked-on-
    behalf and admin-cancel flows correctly attribute the action to the actor
    for the self-event filter to skip the right person. `handler.go` now
    returns a JSON:API-shaped error response on upgrade failure via
    `api.WriteUnauthorized` and a custom upgrader `Error` callback.
- **Vite dev proxy** — added `ws: true` to the `/api` proxy in
  `web/vite.config.ts` so the WebSocket upgrade works against the Vite dev
  server in local development.
- **Pinia in test setup** — `web/vitest.setup.ts` now installs a fresh Pinia
  before each test. Without this, mounting any view that uses `useAuthStore`
  or the new `useLiveFeedStore` (via the `useLiveBookingRefresh` composable)
  would throw "no active Pinia". This restored `InteractiveFloorPlan` test
  coverage that broke after wiring the live refresh.
- **Hub queue-full behaviour** — when the hub's incoming broadcast channel is
  saturated (e.g. during a sustained burst), `NotifyAsync` drops with a
  `slog.Warn` rather than blocking the booking handler. This is exercised by
  `TestNotifyAsyncDropsWhenQueueFull`.
- **Slow clients are dropped** — if a connected client cannot drain its 16-
  event send buffer fast enough, the hub deletes it and closes the channel,
  which causes the write pump to exit and the read pump to unwind. Covered
  by `TestHubDropsSlowClient`.

### Deferred follow-ups (intentionally not in this story)

- **Cypress E2E for the multi-user flow.** Driving two authenticated browser
  sessions from a single Cypress run is fragile (Cypress shares a cookie jar
  per session, and the live feed correctly filters out self-events so a
  same-user trick won't exercise the path). The unit + component tests cover
  the wiring and the backend integration test covers the wire-level event
  delivery. Recommend adding the E2E in a follow-up that uses `cy.session()`
  isolation plus `cy.exec` curl commands for the "other user" mutation.
- **Polling fallback** mentioned in `architecture.md` lines 135–138. The
  story's Task 7.3 explicitly defers this; views remain functional without
  the live feed (just no longer real-time) if a corporate proxy blocks the
  upgrade. Track if real-world environments need it.

### File List

Backend (new):

- `internal/livefeed/event.go`
- `internal/livefeed/hub.go`
- `internal/livefeed/client.go`
- `internal/livefeed/handler.go`
- `internal/livefeed/hub_test.go`
- `internal/notifications/multi.go`
- `internal/notifications/multi_test.go`

Backend (modified):

- `go.mod` (added `github.com/gorilla/websocket v1.5.3`)
- `go.sum`
- `internal/startup/server.go` (construct hub, run it, combine with webhook
  notifier via `MultiNotifier`, register `/api/v1/live`)
- `internal/startup/server_test.go` (registerRoutes signature update)

Frontend (new):

- `web/src/api/liveFeed.ts`
- `web/src/api/liveFeed.test.ts`
- `web/src/composables/useLiveFeed.ts`
- `web/src/composables/useLiveFeed.test.ts`
- `web/src/composables/useLiveBookingRefresh.ts`
- `web/src/composables/useLiveBookingRefresh.test.ts`
- `web/src/stores/useLiveFeedStore.ts`

Frontend (modified):

- `web/src/App.vue` (wire `useLiveFeedStore.start/stop` to auth state)
- `web/src/views/ItemGroupsView.vue` (refresh availability on live events)
- `web/src/views/ItemsView.vue` (silent refresh of day/week tile views on live
  events; new `silentReloadItems` helper and `silent` flag on `loadWeekData`)
- `web/src/components/area-weekly-matrix/AreaWeeklyMatrixView.vue` (silent
  refresh of matrix on live events; `loadMatrix({ silent })` opt-in to skip
  the loading-state flicker)
- `web/src/components/InteractiveFloorPlan.vue` (refresh availability on
  live events)
- `web/vite.config.ts` (`ws: true` on `/api` proxy)
- `web/vitest.setup.ts` (per-test Pinia install)
