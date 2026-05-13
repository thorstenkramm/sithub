---
stepsCompleted: [step-01-validate-prerequisites, step-02-design-epics, step-03-create-stories, step-04-final-validation]
inputDocuments:
  - /Users/thorsten/projects/thorsten/sithub/_bmad-output/planning-artifacts/prd.md
  - /Users/thorsten/projects/thorsten/sithub/_bmad-output/planning-artifacts/architecture.md
  - /Users/thorsten/projects/thorsten/sithub/private/epic-25.md
  - /Users/thorsten/projects/thorsten/sithub/private/epic-26.md
  - /Users/thorsten/projects/thorsten/sithub/private/epic-27.md
  - /Users/thorsten/projects/thorsten/sithub/private/epic-28.md
  - /Users/thorsten/projects/thorsten/sithub/private/epic-30.md
  - /Users/thorsten/projects/thorsten/sithub/private/epic-31.md
  - /Users/thorsten/projects/thorsten/sithub/private/epic-32.md
  - /Users/thorsten/projects/thorsten/sithub/private/epic-33.md
lastEdited: '2026-05-11'
editHistory:
  - date: '2026-02-07'
    changes: "Updated Epic 1 for dual-source auth (Entra ID + local). Added FR28-FR35. Added Epic 11: User Management & Local Authentication with 8 stories. Updated NFR3, additional requirements, and coverage map."
  - date: '2026-02-08'
    changes: "Domain rename: reworded FR4-FR19 (rooms/desks to items). Added FR36-FR42 (weekly availability, booking notes, week booking, booker display, breadcrumbs, schema normalization, UI labels). Added Epic 12: Domain Rename & Data Normalization and Epic 13: Enhanced Booking Experience."
  - date: '2026-02-14'
    changes: "Added FR43-FR58 and Epics 14-17: UI Cleanup & Booking Simplification, Collapsible Item Tiles, User Preferences & Settings, Equipment Filter."
  - date: '2026-03-13'
    changes: "Added FR59-FR66 and Epic 18: Floor Plan Display & Config Consistency. Covers terminology rename, file location enforcement, floor plan serving/display, and connection error handling."
  - date: '2026-03-21'
    changes: "Added FR67-FR74 and Epic 19: User Feedback — Bug Fixes & Feature Requests. Covers cancel dialog bug, week selector/calendar fixes, equipment filter enhancements, week view cancellation, custom icons, and favorites."
  - date: '2026-03-22'
    changes: "Added FR75-FR82 and Epic 20: Interactive Floor Plans & UX Consistency. Covers favorite free-busy indicators, week/day memorization, consistent snackbar confirmations, floor plan positions in SQLite, admin floor plan editor, interactive floor plan overlay with free/busy, and first-level drill-down."
  - date: '2026-03-29'
    changes: "Added FR83-FR84 and Story 20.8: Floor Plan Booking UX Refinements. Covers multi-day booking dialog with weekday checkboxes, persistent overlay, mobile fullscreen layout, close/back navigation stack, drill-down safety enforcement, and precise booking error messages. Based on user testing feedback."
  - date: '2026-04-04'
    changes: "Added FR85-FR90 and Epic 21: i18n, UX Improvements & Booking Limits. Covers multilanguage UI with auto-detection, My Bookings layout reorder, visual fixes (equipment filter icon, floor plan button), and configurable booking limits (advance weeks, max per person with area overrides)."
  - date: '2026-04-05'
    changes: "Added FR91-FR100 and Epic 22: Bug Fixes, Avatars & Reserved Areas. Covers mobile UX audit findings (truncation, menu overflow, week mode readability, floor plan mobile), user avatar sync/upload, and reserved areas/items with YAML-based access control."
  - date: '2026-04-06'
    changes: "Added FR101-FR103 and Epic 23: UI Bug Fixes. Covers booking tile heart icon positioning, hidden booking limit error messages, and floor plan desktop width."
  - date: '2026-04-06'
    changes: "Added FR104-FR106 and Epic 24: Booking Warnings & Profile Consolidation. Covers item warning confirmation dialogs with don't-show-again, sequential warnings in week mode, and merging Settings into Profile."
  - date: '2026-04-09'
    changes: "Added FR107-FR117, UX-DR1-UX-DR14, and Epic 25: UX/UI Improvements — Floor Plan Editor, Booking & Avatar. Covers floor plan editor overhaul (sidebar removal, toolbar dropdowns, auto-save, undo removal, zoom redesign, canvas enlargement), subarea drill-down image enlargement, Entra ID avatar async sync, login spinner, and Profile Photo menu hiding for Entra ID users."
  - date: '2026-04-13'
    changes: "Added FR118-FR121 and Epic 26: Floor Plan Editor — Area Drawing Fixes. Fixes subarea selection tab switching, items dropdown visibility on Areas tab, draw mode for subareas, and rectangle locking during subarea editing."
  - date: '2026-04-15'
    changes: "Added FR122-FR125 and Epic 27: Avatar Sync Fix & Reserved Item Visibility. Fixes Entra ID avatar decoding for non-PNG formats, corrects reserved area free/busy display on floor plans, and shows occupancy on reserved items in list view."
  - date: '2026-04-15'
    changes: "Added FR126-FR128 and Epic 28: Date Selector Fix & Floor Plan Booker Names. Fixes date picker jumping to today after booking, adds booker name tooltips on floor plan avatars, and shows initials with name tooltip when avatars are disabled."
  - date: '2026-05-08'
    changes: "Added FR137-FR139 and Epic 30: Operator Validation, Editor Zoom Height & Optional Drill-Down. Adds startup detection of duplicate items in areas config, makes the floor plan editor canvas grow vertically on zoom-in, and adds an Area drill-down toggle for direct booking on large screens with per-device persistence."
  - date: '2026-05-11'
    changes: "Added FR143-FR145 and Epic 32: Booker Avatars on Item Tiles & Stable Colleague Booking Layout. Adds booker avatars (with initials fallback and full-name tooltip) to booked item tiles in day and week modes, and stabilizes the 'Book for a colleague' layout so the colleague dropdown renders inline next to the radio group on wide viewports without pushing subsequent containers down."
  - date: '2026-05-11'
    changes: "Back-filled FR140-FR142 for Epic 31 (Live Updates, Favorites Rework & Areas Config Hint) into the Functional Requirements list and the FR Coverage Map; added epic-31.md to inputDocuments. No story content changed."
  - date: '2026-05-11'
    changes: "Added FR146-FR152 and Epic 33: Equipment Filter Fixes, Compact Booking Controls & Login Page Rebranding. Fixes equipment-filter reset and table-view filter, widens table-view cancel popover, replaces booking-type radio with single checkbox + always-on dropdown (supersedes Story 32.3 layout), collapses booking controls to a single row, promotes Entra ID to primary login with official icon and 'more options' toggle, and embeds the new SitHub brand logo for the login page and header."
---

# sithub - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for sithub, decomposing the requirements from the PRD and
Architecture requirements into implementable stories.

## Requirements Inventory

### Functional Requirements

FR1: Users can authenticate via Entra ID SSO or local credentials (email and password).
Acceptance: the login page presents a username/password form and a "Login via Entra ID"
button; both methods result in an authenticated session with the user's name displayed.
FR2: The system determines user roles (regular vs admin) based on authentication source.
For Entra ID users, admin status is synced from group membership on every login and cannot
be changed locally. For local users, admin status is managed locally by administrators.
Acceptance: admins see admin-only controls; regular users do not; Entra ID admin status
reflects current group membership after each login.
FR3: Users can access the application only if they are authenticated. Acceptance:
unauthenticated users see only the login page and cannot view any booking data.
FR4: Users can view a list of available areas. Acceptance: after login, the UI lists all
configured areas.
FR5: Users can view a list of item groups within a selected area. Acceptance: selecting an
area shows only its item groups.
FR6: Users can view a list of items within a selected item group. Acceptance: selecting an
item group lists its items.
FR7: Users can view equipment details for each item. Acceptance: each item entry shows its
equipment list if configured.
FR8: Users can view current booking status for items. Acceptance: item entries show
available/occupied status for the selected date.
FR9: Users can book an item for a single day. Acceptance: selecting an item and date creates
a booking that appears in "My Bookings."
FR10: The system prevents double-booking of the same item for the same day. Acceptance: the
second attempt is rejected and no duplicate booking is created.
FR11: Users receive a message when a selected item becomes unavailable during booking.
Acceptance: the message states the item is no longer available for that date and prompts
the user to choose another item.
FR12: Users can view their upcoming bookings ("My Bookings"). Acceptance: the list includes
item, item group, area, and date for each future booking; booking notes are displayed if
present.
FR13: Users can cancel their own bookings. Acceptance: cancelling removes the booking from
all relevant lists and frees the item.
FR14: Admin users can cancel any booking. Acceptance: admins can cancel another user's
booking and the affected user sees the cancellation reflected in their list.
FR15: Users can view an item-group-level booking overview. Acceptance: for a selected item
group and date, the overview lists all booked items and associated users; booking notes are
displayed if present.
FR16: Users can view "Today's presence" for an area (who is in the office today).
Acceptance: the view lists all users with bookings in that area for today; booking notes
are displayed if present.
FR17: Operators can configure server settings via a configuration file. Acceptance: invalid
settings prevent startup with a descriptive error; valid settings allow startup.
FR18: Operators can configure areas, item groups, items, and equipment via a configuration
file. Acceptance: after restart, the UI reflects the updated space definitions.
FR19: The system can load and apply configuration changes on startup. Acceptance:
configuration changes take effect after restart without manual data migration steps.
FR20: Users can book on behalf of another user. Acceptance: the booking appears in both
users' booking lists and either can cancel.
FR21: Users can book items for guests outside the organization. Acceptance: a guest booking
stores guest name and contact and is visible as a guest booking in overviews.
FR22: Users can book multi-day or recurring reservations. Acceptance: the system creates
individual daily bookings and reports conflicts per day.
FR23: Users can view booking history. Acceptance: users can see past bookings with date
range filtering.
FR24: Users can receive notifications related to bookings. Acceptance: booking
creation/cancellation triggers a notification via the configured channel within 5 minutes.
FR25: Admins can manage item groups and items via an admin UI. Acceptance: admins can
add/edit/remove item groups and items; changes appear in discovery lists after save.
FR26: Users can book items using a graphical floor-map view. Acceptance: an item selected
on the map can be booked for a chosen date.
FR27: Admins can access advanced reporting and analytics. Acceptance: admins can view usage
summaries by area/item group and date range.
FR28: All users (Entra ID and local) are stored in a users table. Entra ID users are
inserted on first login and updated on subsequent logins. Acceptance: after login, the user
exists in the users table with correct source, name, and email.
FR29: Email addresses are unique across all authentication sources. Creating a local user
with an email that exists for an Entra ID user is blocked, and vice versa. Acceptance:
attempting to create a duplicate email returns an error regardless of authentication source.
FR30: Local users can log in with email and password. Acceptance: entering valid credentials
in the login form creates an authenticated session; invalid credentials show a descriptive
error message.
FR31: Local users can change their own password via the `/me` endpoint. Acceptance: after
changing the password, the old password no longer works and the new password (minimum 14
characters) is accepted.
FR32: Admin users can reset the password of any local user via the `/users/{id}` endpoint.
Acceptance: the affected user can log in with the new password; Entra ID user passwords
cannot be reset this way.
FR33: The system provides a `/me` endpoint returning the current user's profile information.
Acceptance: authenticated requests to `/me` return the user's id, email, name, role, and
authentication source.
FR34: The system provides a `/users` endpoint for user management. Acceptance: admins can
list, create, read, update, and delete local users; non-admin users can only read. Entra ID
users cannot be created or deleted via this endpoint.
FR35: A demo users SQL file is provided for development setup. Acceptance: running the SQL
file creates 15 users (2 admins, 13 regular users) with local credentials in the database.
FR36: Users can view a weekly availability preview for item groups. Acceptance: the item
group list view includes a calendar week selector (next 8 weeks, current week pre-selected);
each item group tile displays weekday indicators (MO-FR) colored green (at least one item
available) or red (fully booked) for the selected week; the week selector displays dates in
locale-aware format with week number (e.g., "2026-03-16 - Week 12").
FR37: Users can add, view, and edit free-text notes on their bookings. Acceptance: after
completing a booking, a confirmation message includes an "add note" action; notes are
visible in My Bookings, Today's Presence, and item detail views; notes longer than the
display width are truncated with an expand indicator that opens the full text; notes are
editable from the My Bookings view.
FR38: Users can toggle between day booking mode and week booking mode. Acceptance: the
selected mode is persisted in browser local storage and restored on next visit; in week
mode, the date selector becomes a calendar week selector (next 8 weeks); item tiles show
per-day checkboxes with booker names; existing bookings by other users cannot be unchecked;
a single "Confirm My Booking" button submits all selected days at once.
FR39: Users can see the display name of the person who booked an item. Acceptance: in the
item detail view, booked items show the booker's display name alongside the booking status.
FR40: Breadcrumbs in the navigation hierarchy are clickable and navigate to the
corresponding view. Acceptance: clicking any breadcrumb segment navigates to that level of
the hierarchy; the current page breadcrumb is not clickable.
FR41: The bookings table references users by `user_id` only; display names are resolved
via JOIN with the users table. Acceptance: the bookings table does not store denormalized
user name columns; all booking queries that require user names perform a JOIN; existing
bookings are migrated to remove redundant columns.
FR42: UI action labels use domain-neutral terminology. Acceptance: labels read "BOOK"
instead of "VIEW DESK" or "VIEW ROOM"; "BOOK THIS ITEM" instead of "BOOK THIS DESK";
booking confirmation messages reference the item name from the configuration (e.g.,
"Parking Lot 1 booked successfully") rather than the generic term "desk."
FR43: Navigation action labels are simplified. Acceptance: "VIEW ITEM GROUPS" becomes
"SELECT"; "VIEW ITEMS" becomes "SELECT"; "BOOK THIS ITEM" becomes "BOOK"; redundant page
titles and subtitles are removed from item group and item views.
FR44: Redundant "Not available" text is removed from booked items. Acceptance: booked items
show only the status chip, booker name, and notes; no "Not available for \<date\>" message.
FR45: Booker name and booking notes use readable font size. Acceptance: displayed at body-2
or larger, not caption size.
FR46: Booking result feedback uses icons. Acceptance: green checkmark for success, red
warning icon for failure, replacing text-based labels.
FR47: Guest booking option is removed. Acceptance: "Book for guest" radio and fields are no
longer shown.
FR48: Multi-day booking checkbox is removed from day mode. Acceptance: checkbox and
additional dates field no longer shown; week booking mode is the replacement.
FR49: Colleague booking uses a user dropdown. Acceptance: dropdown lists existing users by
display name; booking uses user ID; free-text fields removed.
FR50: Item tiles in week mode are collapsible. Acceptance: chevron toggles between folded
(compact M-F) and unfolded (one line per day, equipment, warnings).
FR51: Booked item tiles in day mode are collapsible. Acceptance: equipment and warnings
hidden by default; chevron reveals details; available items remain expanded.
FR52: Folded tiles with warnings show a warning icon. Acceptance: clicking the icon shows
the warning without unfolding.
FR53: Past date checkboxes disabled in week mode. Acceptance: dates before today are
grayed out and not interactive.
FR54: Full booker name on hover in week mode. Acceptance: tooltip shows full name for
truncated names in folded tiles.
FR55: Theme selector. Acceptance: auto/dark/light in user menu; stored in localStorage;
applied immediately.
FR56: Show weekends toggle. Acceptance: checkbox in user menu; adds Sat/Sun to booking
pages; stored in localStorage.
FR57: Change Password icon fix. Acceptance: icon visible in desktop and mobile menus.
FR58: Equipment filter. Acceptance: text input on booking page; non-matching items blurred;
info icon explains syntax; OR/AND/exact matching; case-insensitive; frontend-only.
FR59: Rename `[spaces]` config section to `[areas]`. Acceptance: sithub.toml uses `[areas]`
table; data models, CLI flags, and environment variables use "areas" terminology; the term
"space" is no longer used anywhere in the codebase.
FR60: Enforce `areas.config_file` inside `main.data_dir`. Acceptance: areas config file path
is resolved relative to `main.data_dir`; startup validation rejects paths outside
`main.data_dir`.
FR61: Floor plans image directory configuration. Acceptance: optional `floor_plans` key under
`[areas]` in sithub.toml specifies a directory inside `main.data_dir`; overridable via
`--areas-floor-plans` flag or `SITHUB_AREAS_FLOOR_PLANS` env var; if set, directory existence
is validated at startup; missing directory causes server exit with error.
FR62: Validate floor plan image references at startup. Acceptance: all images referenced in
`areas.config_file` are checked for existence and format (jpg, png, svg only); missing images
or unsupported formats cause server exit with descriptive error.
FR63: Authenticated floor plan image serving. Acceptance: floor plan images are served via an
authenticated API endpoint; unauthenticated requests are denied.
FR64: Area floor plan display. Acceptance: when an area has a floor plan, a "Floor plan"
button with icon appears next to the calendar week selector; clicking opens the floor plan
image in an overlay with the area name as heading; button is hidden when no floor plan exists.
FR65: Item group floor plan display. Acceptance: when an item group has a floor plan, a
"Floor plan" button with icon appears beneath the day/week selector; clicking opens the floor
plan image in an overlay with the item group name as heading; button is hidden when no floor
plan exists.
FR66: Connection lost error messaging. Acceptance: when the backend is unavailable, the
frontend shows a clear "Connection to server lost" error instead of misleading content like
"no areas available".
FR67: Cancel booking dialog not closing. Acceptance: after confirming a booking cancellation,
the confirmation dialog closes automatically.
FR68: Equipment filter on area/item-groups view. Acceptance: a text input on the item-groups
page filters item groups by equipment; non-matching groups are blurred and disabled; same
parsing rules as the existing equipment filter apply.
FR69: Equipment filter saving. Acceptance: saved filters are stored in browser local storage;
a save icon next to the input persists the current filter; a combobox allows selecting saved
filters; loading a saved filter turns the save icon into a delete icon; editing saved filters
is not supported.
FR70: Cancel bookings from week view. Acceptance: when a booking in the week view belongs to
the current user, a small red cancel icon appears next to the checkmark; clicking it cancels
the booking; other users' bookings show no cancel icon.
FR71: Week selector date range display. Acceptance: each option in the week selector shows
"DD.MM.-DD.MM.YYYY - Week N" (e.g. "23.03.-29.03.2026 - Week 13").
FR72: Calendar widget starts on Monday. Acceptance: the calendar date picker shows Monday as
the first column and Sunday as the last column.
FR73: Custom icons in areas YAML. Acceptance: an optional `icon` field at area, item group,
and item levels in the areas YAML specifies an MDI icon code; the frontend renders the
specified icon; missing or invalid icons fall back to the current defaults.
FR74: Favorites. Acceptance: heart outline icons on item group and item tiles allow marking
favorites; favorites are stored in browser local storage; clicking toggles the favorite state
with confirmation messages; item-groups view sorts: (1) third-level favorites A-Z,
(2) second-level favorites A-Z, (3) remaining items in YAML order; third-level favorites
appear as bookable tiles on the item-groups view.
FR75: Free-busy indicators on favorite tiles. Acceptance: promoted third-level favorite tiles
on the item-groups view show weekly availability dots matching regular item group tiles.
FR76: Memorize selected week. Acceptance: the selected calendar week persists across
navigation between areas and item groups; resets to the current week when the memorized
week is in the past.
FR77: Memorize selected day. Acceptance: the selected booking day persists across navigation;
resets to today after a successful booking.
FR78: Consistent snackbar confirmations. Acceptance: all success confirmations across the app
use the bottom snackbar style; no inline alerts for success feedback.
FR79: Interactive floor plan overlay. Acceptance: floor plan overlay shows item positions with
free/busy state per weekday; free items have green outlines, busy items have red
semi-transparent overlays; clicking a free item creates a booking.
FR80: First-level floor plan drill-down. Acceptance: clicking a sub-area on the first-level
floor plan opens its detail floor plan; fully booked sub-areas show a red overlay.
FR81: Floor plan editor. Acceptance: admin tool accessible from settings; displays floor plan
image with item list; admin draws rectangles to position items; positions are saved via API.
FR82: Floor plan positions in SQLite. Acceptance: item positions on floor plans are stored in
a `floor_plan_positions` table with floor plan filename, item ID, and rectangle coordinates;
CRUD API endpoints exist for reading, creating, updating, and deleting positions.
FR83: Multi-day floor plan booking with weekday checkboxes. Acceptance: clicking a free item
on the floor plan opens a dialog with weekday checkboxes; the current day is pre-checked;
past days and already-booked days are disabled; a summary and "Book now" button allow
multi-day booking in one action; error messages name the specific conflicting day.
FR84: Floor plan overlay UX polish. Acceptance: overlay is persistent (close button only);
mobile opens fullscreen with controls at the bottom; close/back navigates to previous screen
(higher-level floor plan or area page); drill-down is enforced when a detail floor plan
exists, preventing direct booking on first-level sub-areas.
FR85: Multilanguage UI. Acceptance: users can switch the UI language from settings; supported
languages are auto (browser detection), English, German, Spanish, French, and Ukrainian;
selection is stored in browser local storage and applied immediately; the language selector
displays a colored country flag for each language (UK flag for English).
FR86: My Bookings display reorder. Acceptance: each booking card shows the booked date on
the first line and the booked item (with area breadcrumb) on the second line, swapping the
current layout.
FR87: Equipment filter save icon. Acceptance: the save icon next to the equipment filter
input uses the `mdi-content-save` icon instead of the plus icon.
FR88: Floor plan button consistent height and position. Acceptance: the floor plan button
has the same height as the calendar week selector; when an area has a detail floor plan,
the button is positioned next to the calendar week selector (not below the booking mode
toggle).
FR89: Booking advance limit. Acceptance: an optional `weeks_in_advanced` integer under
`[bookings]` in sithub.toml limits how far ahead users can book; only the current plus the
next N weeks are shown and bookable; default is 5.
FR90: Maximum bookings per person. Acceptance: an optional `max_bookings_per_person` integer
under `[bookings]` in sithub.toml limits total active bookings per user across all areas;
default 0 means unlimited; the areas YAML supports the same key at area, item group, and
item levels to override the global limit; the most specific (deepest) matching limit applies;
exceeded limits produce a clear error naming the limit and scope (e.g., "You have exceeded
the maximum of 2 active bookings for 'Tiefgaragenstellplätze, Stellplatz 1'").
FR91: Translation bug fixes. Acceptance: booking limit error messages use frontend i18n
keys; weekday abbreviations in item group tiles use locale-aware labels (DE: MO, DI, MI,
DO, FR, SA, SO); "n/a" labels in week mode are translated or removed.
FR92: Language selector mobile layout. Acceptance: language and theme buttons in the
navigation drawer render without clipping on 390px-wide mobile screens.
FR93: Mobile text truncation. Acceptance: item names in card titles, week mode tile
headers, and My Bookings subtitles wrap to multiple lines instead of truncating with
ellipsis; history date filter fields stack vertically on narrow viewports.
FR94: Week mode mobile readability. Acceptance: booked user names in week mode columns
show initials or abbreviated form that does not overflow or collide with adjacent columns.
FR95: Floor plan mobile improvements. Acceptance: floor plan auto-zooms to fit viewport
width on mobile; floor plan images apply a dark-mode filter when dark theme is active;
floor plan editor shows a desktop-recommended banner on narrow viewports.
FR96: Favorites heart icon visibility. Acceptance: the favorite heart icon is visible on
all item tiles including those with warning badges.
FR97: User avatar sync from Entra ID. Acceptance: on each Entra ID login the user's
profile photo is downloaded from Microsoft Graph and stored locally; avatars are served
via an authenticated API endpoint; missing photos fall back to initials.
FR98: User avatar upload for local users. Acceptance: local users can upload, replace,
and delete a profile image from settings; images are stored under
`{data_dir}/avatars/{user_id}.png` with a reasonable size limit.
FR99: Avatar display integration. Acceptance: user avatars appear in the settings menu,
Today's Presence list, and optionally on floor plan overlays via a toggle.
FR100: Reserved areas and items. Acceptance: `reserved_for` in the areas YAML restricts
booking to listed user emails at area, item group, and item levels; hierarchical
validation rejects configs where a child reserves for users excluded by a parent; items
not bookable by the current user are disabled and blurred in the UI and floor plan.
FR104: Warning confirmation before booking. Acceptance: when a user attempts to book an
item that has a warning, a dialog appears showing the item name (truncated if necessary),
the warning text, CONFIRM and CANCEL buttons, and a "Don't show again" checkbox; confirming
proceeds with the booking; cancelling aborts; the don't-show-again status is stored per
item in browser localStorage and suppresses the dialog on future bookings of that item.
FR105: Sequential warning display in week booking mode. Acceptance: when booking multiple
items in week mode where different items have warnings, the warning dialogs are shown one
after another before the booking is submitted; each dialog identifies the item; the user
can cancel at any point which aborts the entire booking; items whose warnings were
previously dismissed via "Don't show again" are skipped.
FR106: Profile and Settings consolidation. Acceptance: the separate Settings menu option is
removed from the navigation; all settings (theme, language, show weekends, change password)
are accessible under the Profile menu; the Profile menu uses the current profile layout;
no settings functionality is lost.
FR107: Remove the floor plan editor sidebar and expand the canvas to full width. Acceptance:
the left-hand sidebar listing subareas and items is removed from the floor plan editor; the
canvas card expands to use the full 12-column width.
FR108: Replace the floor plan editor sidebar subarea list with a toolbar dropdown. Acceptance:
a v-select dropdown for subareas appears in the toolbar row (same area as the floor plan
selector); selecting a subarea in the dropdown has the same effect as the old sidebar click.
FR109: Replace the floor plan editor sidebar items list with a toolbar dropdown. Acceptance:
a v-select dropdown for items appears in the toolbar row; each option indicates positioned
or unpositioned status; selecting an unpositioned item enters draw mode; selecting a
positioned item selects its rectangle; a delete action is available for positioned items
from this dropdown.
FR110: Double the floor plan editor canvas height. Acceptance: the editor canvas area uses
approximately twice the vertical space compared to the current layout.
FR111: Auto-save the floor plan editor and remove the manual save button. Acceptance: after
a draw, move, or resize interaction completes (pointerup), changes are saved automatically
when unsaved changes exist; the manual Save button is removed; the unsaved changes chip is
replaced with a saving/saved indicator reflecting auto-save state.
FR112: Remove the undo function from the floor plan editor. Acceptance: the Undo button and
all undo-related logic are removed from the floor plan editor.
FR113: Reposition the zoom factor label in the floor plan editor toolbar. Acceptance: the
zoom percentage label appears between the minus and plus buttons (not next to them); the
layout is compact.
FR114: Enlarge subarea floor plan images when drilling down. Acceptance: when a user drills
into a subarea in the floor plan booking view, the floor plan image renders at an enlarged
size that fills the available viewport width; no horizontal scrollbars appear at the default
zoom level after drill-in.
FR115: Hide the "Profile Photo" menu item for Entra ID users. Acceptance: the "Profile
Photo" option is hidden in both desktop user menu and mobile navigation drawer when the
authenticated user's auth_source is not "internal"; Entra ID users cannot access the
avatar upload.
FR116: Make Entra ID avatar sync asynchronous. Acceptance: the avatar sync in the backend
CallbackHandler runs in a goroutine so the OAuth callback returns immediately; the avatar
downloads in the background; login is not blocked by avatar sync.
FR117: Show loading spinner on Entra ID login button. Acceptance: after clicking "Sign in
with Entra ID", the button shows a loading spinner and is disabled; visual feedback is
provided during the redirect to the OAuth flow.
FR118: Subarea selection from dropdown must respect the active tab (Areas/Items) and not
force-switch to Items mode.
FR119: Items dropdown must be hidden when the Areas tab is active.
FR120: Selecting an unpositioned subarea on the Areas tab must enter draw mode; selecting a
positioned subarea must select its rectangle.
FR121: When a subarea is selected for editing, all other subarea rectangles must be locked
(non-interactive) to prevent accidental modification.
FR122: Avatar sync from Entra ID must handle JPEG, PNG, and other common image formats;
failed syncs must log detailed diagnostics (user ID, HTTP status, content-type, body size)
and fall back to initials avatar.
FR123: Reserved areas on the interactive floor plan must show correct free/busy counts and
green/red indicators based on actual availability, not reservation status.
FR124: Users must be able to drill down into reserved areas on the floor plan and see
individual desk availability; free desks the user cannot book must show a "reserved" overlay.
FR125: Reserved items in the regular booking list view must show occupancy (booker names,
free/busy status) instead of an opaque veil; free reserved items show a lock badge and
block booking actions.
FR126: After completing a booking, the date selector must stay on the previously selected
date instead of jumping to today. Acceptance: when a user selects a future date (e.g.
30 April 2026), books an item, and the confirmation completes, the date picker still shows
30 April 2026 and the displayed bookings reflect that date.
FR127: On the interactive floor plan with "Show Avatar" enabled, hovering over (or tapping
on mobile) a booked item's avatar shows the full display name of the booker in a tooltip.
Acceptance: hovering over any avatar on the floor plan displays a tooltip with the person's
full name; tapping on mobile shows the same.
FR128: On the interactive floor plan with "Show Avatar" disabled, booked items show the
booker's initials (e.g. "AS" for Alexander Seidemann-Klamant) instead of an avatar image.
Hovering or tapping shows the full display name in a tooltip. Acceptance: initials are
derived from the first letters of the user's display name parts; the tooltip shows the
full name.
FR137: The server must detect duplicate item identifiers in the areas configuration at
startup and refuse to start with a clear error. Acceptance: if the same item identifier
(e.g. `desk29`) appears more than once across the loaded areas configuration, the server
logs an error naming the duplicate identifier and the offending locations and exits with a
non-zero status before opening the listening socket.
FR138: In the floor plan editor, the image container height must grow when the user zooms
in so that vertical scrolling is not required to see the full image height. Acceptance:
when the user zooms in, the box surrounding the floor plan image expands vertically to
match the scaled image height; only horizontal scrolling is needed to reach off-screen
content; on initial load the container height continues to be derived from the image.
FR139: On the room/floor plan booking view, an "Area drill-down" toggle controls whether
clicking an item drills into the detailed area/room view or starts booking directly. The
toggle is enabled by default on small screens and disabled by default on large screens.
The user's choice is persisted in browser local storage and applies to all subsequent
floor plan bookings on that device. Acceptance: a checkbox/toggle labelled "Area
drill-down" is visible beneath the room plan; toggling it changes click behavior on the
floor plan accordingly; the choice survives reloads on the same device and is independent
per device.
FR140: Bookings and cancellations made by other users must propagate to all open clients
in real time across the tile, weekly table, and floor plan views, without requiring an
explicit page reload. Acceptance: when another user creates a booking that affects an
item visible in the current view, the corresponding tile, table cell, or floor plan
marker updates to its new busy state within a few seconds; when another user cancels such
a booking, the same elements update to free within a few seconds; after a network outage
the client reconciles state with the server on reconnection; multiple open tabs reflect
the change consistently.
FR141: The favorites feature must be reworked into a virtual "Favorites" room that
contains only bookable items (desks) and behaves like any other room with respect to
drill-down, free/busy indicators, and interaction model. Acceptance: with at least one
favorite item, a "Favorites" tile appears as the first tile on the area/room overview;
with no favorites the tile is not shown; favoriting an entire area or room is not
supported; in the weekly table view, favorite item rows display a heart icon and no
sorting or filtering by favorites is offered; on the floor plan, a heart icon is rendered
with its center positioned exactly at the bottom-right corner of free favorite items, and
no heart is rendered for busy items regardless of favorite status; clicking any heart
icon (tile, table, floor plan) removes the item from favorites and the "Favorites" tile
disappears once no favorites remain.
FR142: The example configuration (`sithub.example.toml`) must clarify that the areas
configuration file has to live inside `data_dir`, using wording consistent with the
existing floor-plans-directory hint. Acceptance: the inline comment for the areas config
setting includes the sentence "Must be inside data_dir."; the wording matches the
existing hint used for the floor plans directory setting; the file remains valid TOML and
follows the comment style described in `.claude/rules/toml.md`.
FR143: On the item-groups view in day mode, every booked item tile must display the
booker's avatar (or initials fallback when no avatar is available) in a circle, matching
the avatar treatment used on the floor plan. Acceptance: a booked tile shows a circular
avatar with the booker's full display name on hover (desktop) or tap (mobile); when the
booker has no uploaded or synced avatar, initials derived from the display name are shown
inside the circle; available tiles continue to show no avatar.
FR144: On the item-groups view in week mode, each weekday cell that represents a booking
must display the booker's avatar (or initials fallback) in a circle next to or above the
existing booker name abbreviation; the full booker name is shown on hover (desktop) or
tap (mobile). Acceptance: every booked weekday cell renders a circular avatar consistent
with the day-mode and floor-plan presentation; tapping or hovering surfaces the full
display name; free cells continue to render the existing free-state visuals only.
FR145: When the user selects "Book for a colleague", the colleague-selection dropdown
must render inline to the right of the radio group on viewports that have sufficient
horizontal space, instead of wrapping onto a new line and pushing subsequent containers
downward. Acceptance: switching between "Book for myself" and "Book for a colleague" on
a wide viewport does not change the vertical position of the equipment filter or any
booking tile beneath the controls; on narrow viewports where the inline dropdown does
not fit, the dropdown wraps to the next line without any visual jitter mid-interaction.
FR146: Clearing the equipment filter input on the item-groups view (either by clicking
the built-in clear "X" icon or by manually deleting all characters) must reset the
active filter so every tile becomes interactive again. Acceptance: typing a non-matching
keyword blurs/disables all tiles (existing behavior); subsequently clicking the input's
clear icon or backspacing the input to empty restores all tiles to the unfiltered state
without requiring a page reload; the same behavior applies in day mode and week mode.
FR147: The equipment filter input on the weekly desktop table view must filter rows by
the searched equipment. Acceptance: typing a keyword dims and disables every row whose
item lacks at least one equipment entry matching the parsed filter (same matching
semantics as the tile view, per `matchesParsedFilter`); clearing the input restores all
rows to their normal state; the filter input is the same control already present at the
top of the table view, not a new one.
FR148: The booking-cancel popover that opens from the weekly table view when a user
clicks one of their own booked cells must size its container so the "Cancel booking"
button is fully visible without scrolling. Acceptance: opening the popover from any
booked-by-me cell shows the Person, Platz, Datum lines and both the "Schliessen" and
"Buchung stornieren" buttons inside the popover without any of them being clipped at
the bottom edge; the popover remains anchored to the cell that opened it.
FR149: The "Book for myself / Book for a colleague" radio group in the item-groups
booking controls must be replaced by a single "Book for a colleague" checkbox
(unchecked by default). The colleague-selection dropdown must always be rendered in
the layout, in a disabled state when the checkbox is unchecked and enabled when it is
checked. Acceptance: on first render the checkbox is unchecked and the colleague
dropdown is present but disabled; checking the box enables the dropdown without any
layout reflow; unchecking it disables the dropdown and clears the previously selected
colleague; the previous Story 32.3 radio-group layout is removed entirely.
FR150: The booking controls on the item-groups view must be laid out on a single line
on wide viewports, containing in order: the day/week toggle, the date or week
selector, the equipment filter input (with its info icon), the "Book for a colleague"
checkbox, and the colleague-selection dropdown. On narrower viewports the row wraps
naturally to multiple lines. Acceptance: on a wide viewport (desktop) all five
controls are visible on one row without truncation; the booking-controls card is
visibly shorter (less vertical height) than before; on viewports too narrow to fit
the row the controls wrap to subsequent lines without overlapping; the same layout
applies in day mode and week mode.
FR151: The login page must promote Microsoft Entra ID as the primary authentication
option. The Entra ID sign-in button must include the official Microsoft Entra ID
color icon (SVG sourced from
`https://upload.wikimedia.org/wikipedia/commons/8/8c/Microsoft_Entra_ID_color_icon.svg`,
downloaded into the repository and embedded into the binary). The local
username/password form must be hidden on first render and revealed by a clickable
"more login options" link beneath the Entra ID button; revealing the form swaps the
link text to "less login options" which collapses it again. Acceptance: an
unauthenticated user landing on the login page sees only the brand logo, the Entra ID
button (with the official Entra ID icon), and the "more login options" link; clicking
the link expands the local credentials form and swaps the link to "less login
options"; clicking again collapses the form. When Entra ID is not configured on the
server, the local form is shown by default (no toggle link) so users are not locked
out.
FR152: The SitHub brand logo (`private/sithub_logo.svg`, to be downloaded into the
repository and embedded into the binary) must replace the current text-only branding
in two places: a full vertical layout on the login page (icon + "SitHub" wordmark
stacked) and a compact horizontal layout in the application header (icon + "SitHub"
wordmark side-by-side). Acceptance: the login page renders the full logo above the
Entra ID button; the application header renders the compact horizontal variant in
place of any current text-only branding; both variants are served from the embedded
binary asset; the header variant fits within the existing app-bar height without
clipping.

### NonFunctional Requirements

NFR1: For expected usage (<=50 concurrent users), list navigation actions complete within
2 seconds at p95; booking and cancellation complete within 3 seconds at p95.
NFR2: The system can be restarted without data loss; bookings remain intact after restart
and conflicts do not create partial records.
NFR3: All booking data requires authenticated access (Entra ID or local credentials);
unauthenticated requests are denied. Local user passwords are stored as bcrypt hashes;
plaintext passwords are never persisted or logged. Minimum password length is 14 characters.
Data at rest is stored without application-layer encryption; in-transit encryption is managed
outside the application.
NFR4: Single-node deployment is sufficient; no clustering or horizontal scaling is required
for MVP usage levels.
NFR5: Meets WCAG A: all interactive elements have accessible names, keyboard focus is
visible, and form inputs are labeled.

### Additional Requirements

- Go 1.25 with Echo, SQLite (WAL), and JSON:API responses using `application/vnd.api+json`.
- CLI uses cobra; configuration uses viper with TOML and documented keys.
- Migrations handled via golang-migrate.
- Single-binary distribution with embedded frontend assets.
- Real-time availability via WebSockets with polling fallback.
- Booking conflicts handled optimistically with unique constraint on (item_id, booking_date).
- Bookings are full-day only; store a single booking_date per booking.
- Target builds: macOS (arm64) and Linux (amd64) only; Windows out of scope.
- No Docker or Kubernetes workflows.
- Dual-source authentication: Entra ID SSO (optional) and local credentials (always available).
- Central users table storing both Entra ID and local users with `user_source` enum.
- Unified session mechanism: both auth paths produce gorilla/securecookie signed cookies.
- bcrypt password hashing for local users; minimum 14 characters.
- Admin role sync: Entra ID users sync `is_admin` from group membership on every login;
  local admin managed via DB.
- Email uniqueness enforced at DB level across authentication sources.
- `test_auth` mechanism removed; replaced by real local users.
- Demo users SQL file (`tools/database/demo-users.sql`) with 15 users for development and
  testing.
- OpenAPI 3.1 docs in `api-doc/` with per-endpoint files; lint with Redocly.
- Vue 3 + Vuetify + Pinia + Vue Router; Composition API with `<script setup>`.
- Vitest for unit tests, Cypress for E2E with `data-cy` selectors and real API responses.
- Vite dev server proxies `/api` to `http://localhost:9900`.
- Domain-neutral terminology: "items" instead of "desks", "item groups" instead of "rooms"
  throughout codebase, API, and UI.
- Bookings table normalized: user names resolved via JOIN, no denormalized name columns.

### UX Design Requirements

UX-DR1: Remove the left-hand sidebar from the floor plan editor; expand the canvas card to
use the full 12-column width. (FR107)
UX-DR2: Add a subarea v-select dropdown in the toolbar row of the floor plan editor,
replacing the sidebar list for subarea selection. (FR108)
UX-DR3: Add an items v-select dropdown in the toolbar row of the floor plan editor,
replacing the sidebar list for item selection; each option indicates positioned/unpositioned
status; selecting an item enters draw mode or selects its rectangle. (FR109)
UX-DR4: Add a delete action for positioned items accessible from the items dropdown in the
toolbar, replacing the sidebar delete button. (FR109)
UX-DR5: Double the height of the floor plan editor canvas area. (FR110)
UX-DR6: Implement auto-save on pointerup after draw/move/resize interactions, only when
there are actual unsaved changes. (FR111)
UX-DR7: Remove the manual Save button; replace the unsaved changes chip with a
saving/saved indicator reflecting auto-save state. (FR111)
UX-DR8: Remove the Undo button and all undo-related logic from the floor plan editor.
(FR112)
UX-DR9: Reposition the zoom factor percentage label between the minus and plus buttons
instead of next to them. (FR113)
UX-DR10: When drilling into a subarea in the floor plan booking view, enlarge the image to
fill the viewport width without horizontal scrollbars at default zoom. (FR114)
UX-DR11: Hide the "Profile Photo" menu item in the desktop user menu for Entra ID users.
(FR115)
UX-DR12: Hide the "Profile Photo" menu item in the mobile navigation drawer for Entra ID
users. (FR115)
UX-DR13: Make the Entra ID avatar sync asynchronous in the backend CallbackHandler so the
login is not blocked. (FR116)
UX-DR14: Show a loading spinner and disable the "Sign in with Entra ID" button during login
redirect. (FR117)

### FR Coverage Map

FR1: Epic 1 - Dual-source auth
FR2: Epic 1 - Role determination
FR3: Epic 1 - Access control
FR4: Epic 2 - List areas (reworded in Epic 12)
FR5: Epic 2 - List item groups (reworded in Epic 12)
FR6: Epic 2 - List items (reworded in Epic 12)
FR7: Epic 2 - Item equipment (reworded in Epic 12)
FR8: Epic 2 - Booking status (reworded in Epic 12)
FR9: Epic 3 - Single-day booking (reworded in Epic 12)
FR10: Epic 3 - Prevent double-booking (reworded in Epic 12)
FR11: Epic 3 - Item-taken feedback (reworded in Epic 12)
FR12: Epic 4 - My Bookings (reworded in Epic 12)
FR13: Epic 4 - Cancel booking (reworded in Epic 12)
FR14: Epic 4 - Admin cancel (reworded in Epic 12)
FR15: Epic 5 - Item group booking overview (reworded in Epic 12)
FR16: Epic 5 - Today's presence (reworded in Epic 12)
FR17: Epic 6 - Server configuration
FR18: Epic 6 - Space configuration (reworded in Epic 12)
FR19: Epic 6 - Apply on restart
FR20: Epic 7 - Book on behalf
FR21: Epic 7 - Guest bookings
FR22: Epic 7 - Multi-day bookings
FR23: Epic 7 - Booking history
FR24: Epic 7 - Notifications
FR25: Epic 8 - Admin management UI
FR26: Epic 9 - Floor maps
FR27: Epic 9 - Analytics
FR28: Epic 11 - Users table
FR29: Epic 11 - Email uniqueness
FR30: Epic 11 - Local login
FR31: Epic 11 - Password change
FR32: Epic 11 - Admin password reset
FR33: Epic 11 - /me endpoint
FR34: Epic 11 - /users endpoint
FR35: Epic 11 - Demo users
FR36: Epic 13 - Weekly availability preview
FR37: Epic 13 - Booking notes
FR38: Epic 13 - Week booking mode
FR39: Epic 13 - Booker display name
FR40: Epic 13 - Clickable breadcrumbs
FR41: Epic 12 - Schema normalization
FR42: Epic 12 - UI label consistency
FR43: Epic 14 - Simplified action labels
FR44: Epic 14 - Remove redundant "Not available" text
FR45: Epic 14 - Readable font for booker name and notes
FR46: Epic 14 - Icon-based booking result feedback
FR47: Epic 14 - Remove guest booking option
FR48: Epic 14 - Remove multi-day booking checkbox
FR49: Epic 14 - Colleague booking user dropdown
FR50: Epic 15 - Collapsible tiles in week mode
FR51: Epic 15 - Collapsible tiles in day mode
FR52: Epic 15 - Warning icon on folded tiles
FR53: Epic 15 - Past date checkboxes disabled
FR54: Epic 15 - Full booker name on hover
FR55: Epic 16 - Theme selector
FR56: Epic 16 - Show weekends toggle
FR57: Epic 16 - Change Password icon fix
FR58: Epic 17 - Equipment filter
FR59: Epic 18 - Rename [spaces] to [areas] in config
FR60: Epic 18 - Enforce areas config inside data_dir
FR61: Epic 18 - Floor plans directory configuration
FR62: Epic 18 - Validate floor plan image references at startup
FR63: Epic 18 - Authenticated floor plan image serving
FR64: Epic 18 - Area floor plan display
FR65: Epic 18 - Item group floor plan display
FR66: Epic 18 - Connection lost error messaging
FR67: Epic 19 - Cancel dialog bug fix
FR68: Epic 19 - Equipment filter on item-groups view
FR69: Epic 19 - Equipment filter saving
FR70: Epic 19 - Cancel from week view
FR71: Epic 19 - Week selector date range
FR72: Epic 19 - Calendar Monday start
FR73: Epic 19 - Custom icons in areas YAML
FR74: Epic 19 - Favorites
FR75: Epic 20 - Free-busy indicators on favorite tiles
FR76: Epic 20 - Memorize selected week
FR77: Epic 20 - Memorize selected day
FR78: Epic 20 - Consistent snackbar confirmations
FR79: Epic 20 - Interactive floor plan with free/busy
FR80: Epic 20 - First-level floor plan drill-down
FR81: Epic 20 - Floor plan editor (admin)
FR82: Epic 20 - Floor plan positions in SQLite
FR83: Epic 20 - Multi-day floor plan booking with weekday checkboxes
FR84: Epic 20 - Floor plan overlay UX polish
FR85: Epic 21 - Multilanguage UI
FR86: Epic 21 - My Bookings display reorder
FR87: Epic 21 - Equipment filter save icon
FR88: Epic 21 - Floor plan button height and position
FR89: Epic 21 - Booking advance limit
FR90: Epic 21 - Maximum bookings per person with area overrides
FR91: Epic 22 - Translation bug fixes (limit errors, weekdays, n/a)
FR92: Epic 22 - Language selector mobile layout
FR93: Epic 22 - Mobile text truncation fixes
FR94: Epic 22 - Week mode mobile readability
FR95: Epic 22 - Floor plan mobile improvements
FR96: Epic 22 - Favorites heart icon visibility
FR97: Epic 22 - User avatar sync from Entra ID
FR98: Epic 22 - User avatar upload for local users
FR99: Epic 22 - Avatar display integration
FR100: Epic 22 - Reserved areas and items
FR101: Epic 23 - Booking tile heart icon mispositioned in day and week modes
FR102: Epic 23 - Booking limit error hidden at page bottom instead of modal overlay
FR103: Epic 23 - Floor plan wastes space on desktop (does not use full width)
FR104: Epic 24 - Warning confirmation before booking
FR105: Epic 24 - Sequential warning display in week mode
FR106: Epic 24 - Profile and Settings consolidation
FR107: Epic 25 - Remove floor plan editor sidebar, expand canvas to full width
FR108: Epic 25 - Subarea toolbar dropdown replaces sidebar list
FR109: Epic 25 - Items toolbar dropdown replaces sidebar list with delete action
FR110: Epic 25 - Double floor plan editor canvas height
FR111: Epic 25 - Auto-save floor plan editor, remove save button
FR112: Epic 25 - Remove undo function from floor plan editor
FR113: Epic 25 - Reposition zoom factor label between minus and plus
FR114: Epic 25 - Enlarge subarea floor plan images on drill-down
FR115: Epic 25 - Hide Profile Photo menu for Entra ID users
FR116: Epic 25 - Async Entra ID avatar sync
FR117: Epic 25 - Loading spinner on Entra ID login button
FR118: Epic 26 - Subarea selection from dropdown must respect the active tab (Areas/Items)
and not force-switch to Items mode.
FR119: Epic 26 - Items dropdown must be hidden when the Areas tab is active.
FR120: Epic 26 - Selecting an unpositioned subarea on the Areas tab must enter draw mode;
selecting a positioned subarea must select its rectangle.
FR121: Epic 26 - When a subarea is selected for editing, all other subarea rectangles must
be locked (non-interactive) to prevent accidental modification.
FR122: Epic 27 - Avatar sync from Entra ID must handle JPEG, PNG, and other common image
formats; failed syncs must log detailed diagnostics (user ID, HTTP status, content-type,
body size) and fall back to initials avatar.
FR123: Epic 27 - Reserved areas on the interactive floor plan must show correct free/busy
counts and green/red indicators based on actual availability, not reservation status.
FR124: Epic 27 - Users must be able to drill down into reserved areas on the floor plan and
see individual desk availability; free desks the user cannot book must show a "reserved"
overlay.
FR125: Epic 27 - Reserved items in the regular booking list view must show occupancy
(booker names, free/busy status) instead of an opaque veil; free reserved items show a lock
badge and block booking actions.
FR126: Epic 28 - Date selector stays on selected date after booking instead of jumping to
today
FR127: Epic 28 - Floor plan avatar tooltip shows full booker name on hover/tap
FR128: Epic 28 - Floor plan shows booker initials with name tooltip when avatars are disabled
FR137: Epic 30 - Server refuses to start when areas config contains duplicate item
identifiers
FR138: Epic 30 - Floor plan editor image container grows vertically on zoom-in so only
horizontal scrolling is required
FR139: Epic 30 - Floor plan booking view exposes an "Area drill-down" toggle (default off
on large screens, on for small) persisted per device in local storage
FR140: Epic 31 - Live updates for bookings and cancellations across tile, table, and
floor plan views
FR141: Epic 31 - Favorites reworked as a virtual "Favorites" room with consistent heart
affordances across views
FR142: Epic 31 - Areas config location hint in `sithub.example.toml`
FR143: Epic 32 - Booker avatar on item tiles in day mode
FR144: Epic 32 - Booker avatar on item tiles in week mode
FR145: Epic 32 - Stable inline layout for "Book for a colleague" dropdown
FR146: Epic 33 - Equipment filter reset on item-groups view when input is cleared
FR147: Epic 33 - Equipment filter on weekly desktop table view
FR148: Epic 33 - Wider weekly-table booking-cancel popover (no clipping)
FR149: Epic 33 - Replace booking-type radio with single colleague-checkbox and
always-rendered dropdown
FR150: Epic 33 - Single-line compact booking controls layout
FR151: Epic 33 - Entra ID promoted to primary login with official icon and
"more options" toggle for the local credentials form
FR152: Epic 33 - SitHub brand logo embedded for login page (full) and header (compact)

## Epic List

### Epic 1: Dual-Source Authentication & Access Control

Users can authenticate via Entra ID SSO or local credentials, and only authenticated users
can access SitHub.
**FRs covered:** FR1, FR2, FR3

### Epic 2: Space Discovery & Availability

Users can browse areas, rooms, desks, equipment, and availability.
**FRs covered:** FR4, FR5, FR6, FR7, FR8

### Epic 3: Single-Day Booking & Conflict Handling

Users can book a desk for a day and get clear messaging when a desk is taken.
**FRs covered:** FR9, FR10, FR11

### Epic 4: Booking Management & Admin Overrides

Users can view/cancel their bookings; admins can cancel any booking.
**FRs covered:** FR12, FR13, FR14

### Epic 5: Room & Presence Overviews

Users can view room booking summaries and today’s presence.
**FRs covered:** FR15, FR16

### Epic 6: Operator Configuration & Startup

Operators configure server and spaces via config files and changes apply on restart.
**FRs covered:** FR17, FR18, FR19

### Epic 7: Advanced Booking Options (Post-MVP)

Bookings on behalf, guests, recurring, history, notifications.
**FRs covered:** FR20, FR21, FR22, FR23, FR24

### Epic 8: Admin Management UI (Future)

Admins manage rooms/desks in a UI.
**FRs covered:** FR25

### Epic 9: Floor Maps & Analytics (Future)

Graphical floor map booking and analytics.
**FRs covered:** FR26, FR27

### Epic 11: User Management & Local Authentication

User management API, local login, password management, and demo users for development.
**FRs covered:** FR28, FR29, FR30, FR31, FR32, FR33, FR34, FR35

### Epic 12: Domain Rename & Data Normalization

Users interact with consistent, domain-neutral terminology that supports booking any kind
of resource (desks, parking lots, lab benches), not just desks and rooms. The data model is
normalized to eliminate redundant columns.
**FRs covered:** FR41, FR42 (plus codebase migration of reworded FR4-FR19)

### Epic 13: Enhanced Booking Experience

Users get powerful new booking capabilities: weekly availability previews, booking notes,
week-at-a-time booking mode, visibility of who booked what, and clickable navigation
breadcrumbs.
**FRs covered:** FR36, FR37, FR38, FR39, FR40

### Epic 14: UI Cleanup & Booking Simplification

Remove visual clutter, simplify action labels, streamline the booking form, and improve the
display of booked items for a cleaner, faster user experience.
**FRs covered:** FR43, FR44, FR45, FR46, FR47, FR48, FR49

### Epic 15: Collapsible Item Tiles

Introduce fold/unfold tile behavior across day and week booking modes to manage visual
complexity, showing equipment, warnings, and full details on demand.
**FRs covered:** FR50, FR51, FR52, FR53, FR54

### Epic 16: User Preferences & Settings

Let users personalize their booking experience with theme selection, weekend visibility,
and minor menu fixes.
**FRs covered:** FR55, FR56, FR57

### Epic 17: Equipment Filter

Enable users to filter items by equipment keywords to quickly find suitable workspaces
using an advanced search syntax.
**FRs covered:** FR58

### Epic 18: Floor Plan Display & Config Consistency

Users can view floor plan images for areas and item groups while booking, and operators
benefit from consistent configuration terminology and stricter file location validation.
**FRs covered:** FR59, FR60, FR61, FR62, FR63, FR64, FR65, FR66

### Epic 19: User Feedback — Bug Fixes & Feature Requests

Users benefit from a smoother booking experience through bug fixes and new capabilities
including equipment filter enhancements, quick cancellation from week view, customizable
icons, an improved calendar/week selector, and a favorites system.
**FRs covered:** FR67, FR68, FR69, FR70, FR71, FR72, FR73, FR74

### Epic 20: Interactive Floor Plans & UX Consistency

Users can view live free/busy status on floor plan overlays, book items directly from
floor plans, and admins can position items on floor plan images. Navigation state is
preserved across the app and confirmations use a consistent style.
**FRs covered:** FR75, FR76, FR77, FR78, FR79, FR80, FR81, FR82, FR83, FR84

### Epic 21: i18n, UX Improvements & Booking Limits

Users can switch the UI language, benefit from visual refinements (My Bookings layout,
equipment filter icon, floor plan button positioning), and operators can configure booking
limits (advance booking window and maximum bookings per person with per-area overrides).
**FRs covered:** FR85, FR86, FR87, FR88, FR89, FR90

### Epic 22: Bug Fixes, Avatars & Reserved Areas

Mobile UX audit findings are resolved (translation bugs, text truncation, menu overflow,
week mode readability, floor plan usability). User avatars are synced from Entra ID or
uploaded locally and displayed across the app. Areas and items can be reserved for
specific users via YAML configuration.
**FRs covered:** FR91, FR92, FR93, FR94, FR95, FR96, FR97, FR98, FR99, FR100

### Epic 23: UI Bug Fixes

Fix booking tile layout, hidden error messages, and floor plan width on desktop.
**FRs covered:** FR101, FR102, FR103

### Epic 24: Booking Warnings & Profile Consolidation

Users are prompted with a confirmation dialog before booking items that have warnings,
with a "don't show again" option per item. In week mode, warnings for multiple items are
shown sequentially. The Settings menu is removed and consolidated into the Profile menu.
**FRs covered:** FR104, FR105, FR106

### Epic 25: UX/UI Improvements — Floor Plan Editor, Booking & Avatar

The floor plan editor is overhauled for a streamlined editing experience: sidebar replaced
with toolbar dropdowns, canvas enlarged, auto-save replaces manual save, undo removed, and
zoom controls redesigned. Subarea drill-down images are enlarged for usability. Entra ID
avatar sync is made asynchronous with login feedback, and the Profile Photo menu is hidden
for Entra ID users.
**FRs covered:** FR107, FR108, FR109, FR110, FR111, FR112, FR113, FR114, FR115, FR116, FR117

### Epic 26: Floor Plan Editor — Area Drawing Fixes

Fix subarea selection tab switching, items dropdown visibility on Areas tab, draw mode for
subareas, and rectangle locking during subarea editing.
**FRs covered:** FR118, FR119, FR120, FR121

### Epic 27: Avatar Sync Fix & Reserved Item Visibility

Fix Entra ID avatar decoding for non-PNG formats, correct reserved area free/busy display
on floor plans, and show occupancy on reserved items in list view.
**FRs covered:** FR122, FR123, FR124, FR125

### Epic 28: Date Selector Fix & Floor Plan Booker Names

Fix the date picker jumping to today after booking, and add booker identification on floor
plans via name tooltips on avatars and initials display when avatars are disabled.
**FRs covered:** FR126, FR127, FR128

### Epic 30: Operator Validation, Editor Zoom Height & Optional Drill-Down

Make the server fail fast on duplicate items in the areas configuration, fix the floor
plan editor so zooming in expands the canvas vertically (eliminating unnecessary vertical
scrolling), and let large-screen users skip the room drill-down via an "Area drill-down"
toggle persisted per device.
**FRs covered:** FR137, FR138, FR139

### Epic 31: Live Updates, Favorites Rework & Areas Config Hint

Bookings and cancellations made by other users propagate to all open clients in real time
across tile, table, and floor plan views. The favorites feature is reworked into a virtual
"Favorites" room that lists only bookable items, with consistent heart-icon affordances
across views. The example configuration clarifies that the areas config file must live
inside `data_dir`.
**FRs covered:** FR140, FR141, FR142

### Epic 32: Booker Avatars on Item Tiles & Stable Colleague Booking Layout

Show the booker's avatar (with initials fallback and full-name tooltip) on booked item
tiles in both day and week modes so users can recognize who has reserved an item at a
glance, matching the avatar presentation already used on the floor plan. Also stabilize
the "Book for a colleague" layout so that selecting it does not push the equipment filter
and tile grid down a line; the colleague-selection dropdown appears inline next to the
radio group whenever the viewport has room.
**FRs covered:** FR143, FR144, FR145

### Epic 33: Equipment Filter Fixes, Compact Booking Controls & Login Page Rebranding

Fix two equipment-filter bugs (the input on the item-groups view does not reset when
cleared; the input on the weekly desktop table view is non-functional) and widen the
weekly-table booking-cancel popover so the cancel button is no longer clipped.
Replace the booking-type radio group introduced in Story 32.3 with a single
"Book for a colleague" checkbox plus an always-rendered, conditionally-disabled
colleague dropdown, and collapse the booking-controls card to a single line on wide
viewports so it occupies less vertical space. Rebrand the login experience: promote
Entra ID to the primary sign-in option with its official Microsoft icon, hide the
local credentials form behind a "more login options" toggle, and adopt the new
`private/sithub_logo.svg` as the SitHub brand logo on both the login page (full
vertical) and the application header (compact horizontal). The SVG assets must be
downloaded into the repository and embedded into the binary.
**FRs covered:** FR146, FR147, FR148, FR149, FR150, FR151, FR152

<!-- Repeat for each epic in epics_list (N = 1, 2, 3...) -->

## Epic 1 Stories: Dual-Source Authentication & Access Control

Users can authenticate via Entra ID SSO or local credentials, and only authenticated users
can access SitHub.
**FRs covered:** FR1, FR2, FR3

### Story 1.1: Dual-Source Login

**FRs covered:** FR1

As an employee,
I want to sign in via Entra ID or local credentials,
So that I can access SitHub regardless of my company's identity provider.

**Acceptance Criteria:**

**Given** I am not authenticated
**When** I open SitHub
**Then** I see a login page with a username/password form and a "Login via Entra ID" button

**Given** I click "Login via Entra ID"
**When** I complete the Entra ID flow
**Then** I return to SitHub with my name displayed

**Given** I enter valid local credentials in the login form
**When** I submit the form
**Then** I am authenticated and see my name displayed

**Given** I enter invalid local credentials
**When** I submit the form
**Then** I see a descriptive error message

### Story 1.2: Source-Dependent Role Determination

**FRs covered:** FR2

As an admin,
I want my role determined from my authentication source,
So that I see admin-only controls.

**Acceptance Criteria:**

**Given** my Entra ID account is in the admin group
**When** I log in via Entra ID
**Then** the system syncs my admin status from group membership
**And** admin-only controls are visible

**Given** I am a local user with admin role in the database
**When** I log in with local credentials
**Then** admin-only controls are visible

**Given** I am removed from the Entra ID admin group
**When** I log in again
**Then** admin-only controls are no longer visible

### Story 1.3: Access Denied for Unauthenticated Users

**FRs covered:** FR3

As a company operator,
I want unauthenticated users blocked,
So that booking data is protected.

**Acceptance Criteria:**

**Given** I am not authenticated
**When** I attempt to access any protected page
**Then** I am redirected to the login page

**Given** I am not authenticated
**When** I make an API request to a protected endpoint
**Then** the API returns a JSON:API error with 401 status

## Epic 2 Stories: Space Discovery & Availability

Users can browse areas, rooms, desks, equipment, and availability.
**FRs covered:** FR4, FR5, FR6, FR7, FR8

### Story 2.1: List Areas

**FRs covered:** FR4

As an employee,
I want to see the list of areas,
So that I can choose where to book.

**Acceptance Criteria:**

**Given** I am authenticated  
**When** I open the app  
**Then** I see all configured areas  
**And** the list is empty-safe (shows zero areas without error)

### Story 2.2: List Rooms in an Area

**FRs covered:** FR5

As an employee,
I want to see rooms for a selected area,
So that I can choose a room.

**Acceptance Criteria:**

**Given** I am viewing an area  
**When** I select the area  
**Then** I see only rooms belonging to that area  
**And** rooms outside the area are not shown

### Story 2.3: List Desks with Equipment

**FRs covered:** FR6, FR7

As an employee,
I want to see desks and their equipment for a room,
So that I can pick a suitable desk.

**Acceptance Criteria:**

**Given** I am viewing a room  
**When** I open the room  
**Then** I see the list of desks in that room  
**And** each desk shows its equipment list

### Story 2.4: Show Availability Status by Date

**FRs covered:** FR8

As an employee,
I want to see which desks are available for a selected date,
So that I can choose a free desk.

**Acceptance Criteria:**

**Given** I have selected a room and date  
**When** desks are displayed  
**Then** each desk shows available or occupied status for that date  
**And** status updates when the date changes

## Epic 3 Stories: Single-Day Booking & Conflict Handling

Users can book a desk for a day and get clear messaging when a desk is taken.
**FRs covered:** FR9, FR10, FR11

### Story 3.1: Create Single-Day Booking

**FRs covered:** FR9

As an employee,
I want to book a desk for a specific date,
So that I can reserve my workspace.

**Acceptance Criteria:**

**Given** I have selected a desk and date
**When** I confirm the booking
**Then** the booking is created for that date
**And** it appears in "My Bookings"

**Given** I attempt to book a desk for a past date
**When** I submit the booking
**Then** the system rejects the booking
**And** I see a clear error message that past dates cannot be booked

### Story 3.2: Prevent Double-Booking

**FRs covered:** FR10

As an employee,
I want the system to prevent duplicate bookings for the same desk and day,
So that I don’t book a desk that’s already taken.

**Acceptance Criteria:**

**Given** a desk is already booked for a date  
**When** another booking is attempted for the same desk and date  
**Then** the request is rejected  
**And** no duplicate booking is created

### Story 3.3: Desk-Taken Feedback

**FRs covered:** FR11

As an employee,
I want a clear message when the desk becomes unavailable during booking,
So that I can choose another desk.

**Acceptance Criteria:**

**Given** I am booking a desk and it becomes unavailable  
**When** I submit the booking  
**Then** I see a message that the desk is no longer available for that date  
**And** I am prompted to choose another desk

## Epic 4 Stories: Booking Management & Admin Overrides

Users can view/cancel their bookings; admins can cancel any booking.
**FRs covered:** FR12, FR13, FR14

### Story 4.1: View My Bookings

**FRs covered:** FR12

As an employee,
I want to see my upcoming bookings,
So that I can confirm my reservations.

**Acceptance Criteria:**

**Given** I am authenticated
**When** I open "My Bookings"
**Then** I see a list of my future bookings
**And** each entry includes desk, room, area, and date

**Given** I have no upcoming bookings
**When** I open "My Bookings"
**Then** I see an empty state with a helpful message and action to book a desk

### Story 4.2: Cancel My Booking

**FRs covered:** FR13

As an employee,
I want to cancel my booking,
So that I can free the desk if plans change.

**Acceptance Criteria:**

**Given** I have a future booking  
**When** I cancel it  
**Then** the booking is removed from my list  
**And** the desk becomes available for that date

### Story 4.3: Admin Cancel Any Booking

**FRs covered:** FR14

As an admin,
I want to cancel any booking,
So that I can resolve conflicts.

**Acceptance Criteria:**

**Given** I am an admin viewing a room booking overview (Epic 5)
**When** I see another user's booking
**Then** I see an admin cancel action on that booking

**Given** I am an admin
**When** I cancel another user's booking
**Then** the booking is removed from all relevant lists
**And** the affected user sees the cancellation

## Epic 5 Stories: Room & Presence Overviews

Users can view room booking summaries and today’s presence.
**FRs covered:** FR15, FR16

### Story 5.1: Room Booking Overview

**FRs covered:** FR15

As an employee,
I want to see a room-level booking overview for a date,
So that I can understand room utilization.

**Acceptance Criteria:**

**Given** I select a room and date  
**When** I view the overview  
**Then** I see all booked desks and associated users for that date

### Story 5.2: Today’s Presence by Area

**FRs covered:** FR16

As an employee,
I want to see who is in the office today by area,
So that I can coordinate with colleagues.

**Acceptance Criteria:**

**Given** I view today’s presence for an area  
**When** the list is displayed  
**Then** I see all users with bookings in that area for today

## Epic 6 Stories: Operator Configuration & Startup

Operators configure server and spaces via config files and changes apply on restart.
**FRs covered:** FR17, FR18, FR19

### Story 6.1: Load Server Configuration

**FRs covered:** FR17

As an operator,
I want the server to load settings from a config file,
So that I can control listen address, port, and data directory.

**Acceptance Criteria:**

**Given** a valid configuration file  
**When** the server starts  
**Then** the server loads the settings  
**And** invalid settings prevent startup with a clear error

### Story 6.2: Load Space Configuration

**FRs covered:** FR18

As an operator,
I want areas, rooms, desks, and equipment loaded from a config file,
So that space definitions are centrally managed.

**Acceptance Criteria:**

**Given** a valid space configuration file  
**When** the server starts  
**Then** the UI reflects the configured areas, rooms, desks, and equipment

### Story 6.3: Apply Configuration on Restart

**FRs covered:** FR19

As an operator,
I want configuration changes to apply on restart,
So that I can update spaces without manual migration steps.

**Acceptance Criteria:**

**Given** the config file has changed  
**When** the server restarts  
**Then** the new configuration is applied  
**And** no manual data migration steps are required

## Epic 7 Stories: Advanced Booking Options (Post-MVP)

Bookings on behalf, guests, recurring, history, notifications.
**FRs covered:** FR20, FR21, FR22, FR23, FR24

### Story 7.1: Book on Behalf of Another User

**FRs covered:** FR20

As an employee,
I want to book a desk on behalf of another user,
So that we can sit together.

**Acceptance Criteria:**

**Given** I book a desk for another user  
**When** the booking is created  
**Then** it appears in both users’ booking lists  
**And** either user can cancel it

### Story 7.2: Guest Booking

**FRs covered:** FR21

As an employee,
I want to book a desk for a guest,
So that visitors can reserve a spot.

**Acceptance Criteria:**

**Given** I create a guest booking  
**When** the booking is saved  
**Then** the guest name and contact are stored  
**And** the booking is labeled as a guest booking in overviews

### Story 7.3: Multi-Day or Recurring Booking

**FRs covered:** FR22

As an employee,
I want to book multiple days or a recurring schedule,
So that I can plan ahead.

**Acceptance Criteria:**

**Given** I select multiple dates or a recurrence  
**When** I submit the booking  
**Then** the system creates individual daily bookings  
**And** conflicts are reported per day

### Story 7.4: Booking History

**FRs covered:** FR23

As an employee,
I want to view my booking history,
So that I can review past reservations.

**Acceptance Criteria:**

**Given** I open booking history  
**When** I filter by date range  
**Then** I see past bookings within that range

### Story 7.5: Booking Notifications

**FRs covered:** FR24

As an employee,
I want to receive booking notifications,
So that I stay informed about changes.

**Acceptance Criteria:**

**Given** I create or cancel a booking  
**When** the action completes  
**Then** a notification is sent within 5 minutes via the configured channel

## Epic 8 Stories: Admin Management UI (Future)

Admins manage rooms/desks in a UI.
**FRs covered:** FR25

### Story 8.1: Manage Rooms and Desks in Admin UI

**FRs covered:** FR25

As an admin,
I want to add, edit, and remove rooms and desks,
So that I can keep space definitions up to date.

**Acceptance Criteria:**

**Given** I am an admin  
**When** I add, edit, or remove a room or desk  
**Then** changes are saved  
**And** discovery lists reflect the updates

## Epic 9 Stories: Floor Maps & Analytics (Future)

Graphical floor map booking and analytics.
**FRs covered:** FR26, FR27

### Story 9.1: Book via Floor Map

**FRs covered:** FR26

As an employee,
I want to select a desk from a floor map,
So that I can book visually.

**Acceptance Criteria:**

**Given** I open the floor map  
**When** I select a desk and date  
**Then** the desk can be booked for that date

### Story 9.2: Usage Analytics

**FRs covered:** FR27

As an admin,
I want to view usage summaries by area, room, and date range,
So that I can understand utilization.

**Acceptance Criteria:**

**Given** I open the analytics view  
**When** I select area, room, and date range  
**Then** I see usage summaries for the selection

## Epic 10 Stories: UI/UX Redesign

Transform the application from basic Vuetify defaults into a polished, modern desk booking
experience with consistent design language, reusable components, and excellent mobile support.
**FRs covered:** NFR5 (Accessibility), enhances all existing FRs

### Story 10.1: Design System Foundation

**FRs covered:** NFR5

As a user,
I want a visually consistent and branded experience,
So that the application feels professional and trustworthy.

**Acceptance Criteria:**

**Given** the application loads  
**When** I view any page  
**Then** I see consistent colors, typography, and spacing throughout  
**And** the color scheme reflects a professional brand identity  
**And** a logo and favicon are displayed

**Technical Requirements:**

- Custom Vuetify theme in `web/src/plugins/vuetify.ts`
- Color palette: primary, secondary, success, warning, error, surface colors
- Typography: Inter font family with defined scale
- Spacing tokens following 4px base unit
- Logo SVG and favicon in `web/public/`

### Story 10.2: Reusable Component Library

**FRs covered:** NFR5

As a user,
I want a consistent UI experience across all pages,
So that the application feels polished and predictable.

**Acceptance Criteria:**

**Given** I navigate between different views
**When** I interact with common UI patterns (empty states, loading, confirmations)
**Then** they look and behave consistently across the application
**And** all components follow the design system from Story 10.1

**Components to create:**

- `PageHeader.vue` - Page title, breadcrumbs, actions
- `EmptyState.vue` - Illustrated empty states with action
- `LoadingState.vue` - Skeleton loaders matching content layout
- `ConfirmDialog.vue` - Confirmation modal with customizable actions
- `DatePicker.vue` - Vuetify date picker with consistent styling
- `DateRangePicker.vue` - Date range selection for filters
- `StatusChip.vue` - Consistent status indicators (available, booked, etc.)

### Story 10.3: Navigation & Layout Redesign

**FRs covered:** NFR5

As a user,
I want clear navigation and context awareness,
So that I always know where I am and can easily move around.

**Acceptance Criteria:**

**Given** I am on any page  
**When** I look at the navigation  
**Then** I see the current page highlighted in the nav  
**And** I see breadcrumbs showing my location in the hierarchy  
**And** I can access main sections (Areas, My Bookings, History) from any page

**Given** I am on a mobile device  
**When** I open the navigation  
**Then** I see a drawer menu that works well on small screens

**Technical Requirements:**

- Redesigned `App.vue` with improved header
- Breadcrumb component integrated into layout
- Mobile navigation drawer
- User menu with name and logout

### Story 10.4: Space Discovery Views Redesign

**FRs covered:** FR4, FR5, FR6, FR7, FR8, NFR5

As a user,
I want visually appealing space discovery,
So that browsing areas, rooms, and desks is enjoyable and efficient.

**Acceptance Criteria:**

**Given** I am viewing the areas list  
**When** the page loads  
**Then** I see areas displayed as cards with visual hierarchy  
**And** empty state shows an illustration and helpful message

**Given** I am viewing rooms in an area  
**When** the page loads  
**Then** I see room cards with desk availability summary  
**And** breadcrumbs show: Home > [Area Name]

**Given** I am viewing desks in a room  
**When** the page loads  
**Then** I see desks as visual cards with clear status indicators  
**And** equipment and warnings are displayed attractively  
**And** available vs booked desks are visually distinct

### Story 10.5: Booking Flow Redesign

**FRs covered:** FR9, NFR5

As a user,
I want an intuitive and delightful booking experience,
So that reserving a desk feels effortless.

**Acceptance Criteria:**

**Given** I want to book a desk
**When** I click the book button
**Then** I see a clear booking dialog/flow
**And** date selection uses a proper calendar picker

**Given** I complete a booking
**When** the booking succeeds
**Then** I see a success confirmation with booking details
**And** I have clear next actions (view bookings, book another)

### Story 10.6: Booking Management Views Redesign

**FRs covered:** FR12, FR13, NFR5

As a user,
I want my bookings displayed beautifully,
So that managing my reservations is pleasant.

**Acceptance Criteria:**

**Given** I open My Bookings
**When** the page loads
**Then** I see bookings as cards with all relevant info
**And** cancel action has a confirmation dialog

**Given** I have no upcoming bookings
**When** I open My Bookings
**Then** I see a helpful empty state with an action to book a desk

### Story 10.7: Mobile Responsiveness

**FRs covered:** NFR5

As a mobile user,
I want the app to work well on my phone,
So that I can book desks on the go.

**Acceptance Criteria:**

**Given** I access the app on a mobile device  
**When** I view any page  
**Then** the layout adapts to the screen size  
**And** touch targets are appropriately sized (min 44px)  
**And** navigation is accessible via drawer menu  
**And** forms and dialogs are usable on small screens

**Technical Requirements:**

- Responsive breakpoints for all views
- Touch-friendly interactions
- Viewport-appropriate font sizes
- No horizontal scrolling on mobile

## Epic 11 Stories: User Management & Local Authentication

User management API, local login, password management, and demo users for development.
**FRs covered:** FR28, FR29, FR30, FR31, FR32, FR33, FR34, FR35

### Story 11.1: Users Table and Entra ID User Sync

**FRs covered:** FR28

As an operator,
I want all users stored in a central users table,
So that the system has a unified user directory regardless of authentication source.

**Acceptance Criteria:**

**Given** an Entra ID user logs in for the first time
**When** the login completes
**Then** the user is inserted into the users table with source "entraid", name, and email

**Given** an Entra ID user logs in again
**When** the login completes
**Then** the user's name and admin status are updated from Entra ID

**Given** a local user is created via the API
**When** the creation succeeds
**Then** the user exists in the users table with source "internal"

### Story 11.2: Email Uniqueness Across Sources

**FRs covered:** FR29

As an operator,
I want email addresses unique across all authentication sources,
So that identity conflicts are prevented.

**Acceptance Criteria:**

**Given** an Entra ID user exists with email `alex@example.com`
**When** an admin attempts to create a local user with the same email
**Then** the request is rejected with a JSON:API error

**Given** a local user exists with email `dana@example.com`
**When** an Entra ID user with the same email logs in for the first time
**Then** the login fails with a descriptive error

### Story 11.3: Local User Login

**FRs covered:** FR30

As a local user,
I want to log in with my email and password,
So that I can use SitHub without Entra ID.

**Acceptance Criteria:**

**Given** I am a local user with valid credentials
**When** I enter my email and password in the login form and submit
**Then** I am authenticated and see my name displayed

**Given** I enter an incorrect password
**When** I submit the login form
**Then** I see a descriptive error message
**And** no session is created

### Story 11.4: Self-Service Password Change

**FRs covered:** FR31

As a local user,
I want to change my own password,
So that I can maintain my account security.

**Acceptance Criteria:**

**Given** I am authenticated as a local user
**When** I submit a password change via `/me` with a new password of 14+ characters
**Then** the password is updated
**And** the old password no longer works

**Given** I submit a new password shorter than 14 characters
**When** the request is processed
**Then** it is rejected with a validation error

**Given** I am an Entra ID user
**When** I attempt to change my password via `/me`
**Then** the request is rejected (Entra ID passwords are managed externally)

### Story 11.5: Admin Password Reset

**FRs covered:** FR32

As an admin,
I want to reset any local user's password,
So that I can help users who are locked out.

**Acceptance Criteria:**

**Given** I am an admin
**When** I reset a local user's password via `/users/{id}`
**Then** the user can log in with the new password

**Given** I attempt to reset an Entra ID user's password
**When** the request is processed
**Then** it is rejected with a JSON:API error

### Story 11.6: Current User Profile Endpoint

**FRs covered:** FR33

As an authenticated user,
I want to retrieve my profile information,
So that the UI can display my identity and role.

**Acceptance Criteria:**

**Given** I am authenticated
**When** I request `/me`
**Then** I receive my id, email, name, role, and authentication source

**Given** I am not authenticated
**When** I request `/me`
**Then** I receive a 401 JSON:API error

### Story 11.7: User Management API

**FRs covered:** FR34

As an admin,
I want to manage local users via the API,
So that I can create, update, and remove user accounts.

**Acceptance Criteria:**

**Given** I am an admin
**When** I list users via `GET /users`
**Then** I see all users (Entra ID and local) with their source and role

**Given** I am an admin
**When** I create a local user via `POST /users`
**Then** the user is created with source "internal" and a hashed password

**Given** I am an admin
**When** I attempt to create or delete an Entra ID user via `/users`
**Then** the request is rejected with a JSON:API error

**Given** I am a non-admin user
**When** I attempt to create, update, or delete a user via `/users`
**Then** the request is rejected with a 403 JSON:API error

**Given** I am a non-admin user
**When** I list or read users via `/users`
**Then** I receive the data (read access is allowed for all authenticated users)

### Story 11.8: Demo Users SQL File

**FRs covered:** FR35

As a developer,
I want a demo users SQL file,
So that I can quickly set up a development environment with test data.

**Acceptance Criteria:**

**Given** the SQL file at `tools/database/demo-users.sql` exists
**When** it is executed against the database
**Then** 15 users are created: 2 admins and 13 regular users with local credentials
**And** all passwords are bcrypt-hashed

## Epic 12 Stories: Domain Rename & Data Normalization

Users interact with consistent, domain-neutral terminology that supports booking any kind
of resource (desks, parking lots, lab benches), not just desks and rooms. The data model is
normalized to eliminate redundant columns.
**FRs covered:** FR41, FR42 (plus codebase migration of reworded FR4-FR19)

### Story 12.1: Rename Backend Internals and Database Column

**FRs covered:** FR4-FR16 (reworded, partial)

As a developer,
I want the Go packages, structs, and database column to use domain-neutral terminology,
So that the codebase foundation is ready for the public API rename.

**Acceptance Criteria:**

**Given** the Go packages `internal/rooms/` and `internal/desks/` exist
**When** the rename is applied
**Then** they are consolidated or renamed to use "item group" and "item" terminology
**And** all Go struct fields, function names, and variables use the new terminology

**Given** the database has `desk_id` columns in the bookings table
**When** a new migration is applied
**Then** the column is renamed to `item_id`
**And** the unique constraint on (desk_id, booking_date) becomes (item_id, booking_date)
**And** existing bookings are preserved with correct references

**Given** the API routes remain unchanged in this story
**When** the internal rename is applied
**Then** the existing API routes still work (backward compatible)
**And** `go test -race ./...` succeeds
**And** `golangci-lint run ./...` reports no errors

### Story 12.2: Normalize Bookings Table

**FRs covered:** FR41

As a developer,
I want the bookings table to reference users by user_id only,
So that display names are always current and not duplicated across tables.

**Acceptance Criteria:**

**Given** the bookings table contains `user_name`, `booked_by_user_id`, and
`booked_by_user_name` columns
**When** a new database migration is applied
**Then** the redundant columns are removed from the bookings table
**And** existing data is preserved (user_id remains as the foreign key reference)

**Given** a booking query needs to display a user's name
**When** the query is executed
**Then** the display name is resolved via JOIN with the users table
**And** the name reflects the current value in the users table

**Given** the bookings API returns user information
**When** a booking list or detail is requested
**Then** user display names are included in the response via JOIN
**And** the JSON:API response structure remains consistent

**Given** the migration has been applied
**When** `go test -race ./...` is executed
**Then** all existing tests pass with the normalized schema
**And** booking creation and listing continue to work correctly

### Story 12.3: Update API Routes and YAML Configuration

**FRs covered:** FR4-FR16 (reworded, partial), FR18 (reworded)

As an operator,
I want the API routes and YAML configuration to use domain-neutral terminology,
So that the public interface reflects the flexible item model.

**Acceptance Criteria:**

**Given** the API routes use `/rooms/:room_id/desks` and `/areas/:area_id/rooms`
**When** the route rename is applied
**Then** the routes use `/item-groups/:item_group_id/items` and
`/areas/:area_id/item-groups`
**And** JSON:API resource types use the new terminology (e.g., `item-groups`, `items`)

**Given** the YAML configuration uses `rooms` and `desks` keys
**When** the spaces config loader is updated
**Then** it reads `items` keys at both hierarchy levels (item groups and items)
**And** the `sithub_areas.example.yaml` is updated with the new keys
**And** the example includes diverse item types (office desks, parking lots)

**Given** the route rename is applied
**When** `go test -race ./...` is executed
**Then** all tests pass with the new route paths
**And** `golangci-lint run ./...` reports no errors

### Story 12.4: Rename Rooms and Desks to Items in Frontend

**FRs covered:** FR4-FR16 (reworded), FR42

As a user,
I want the UI to use domain-neutral terminology,
So that I see consistent labels regardless of what I am booking.

**Acceptance Criteria:**

**Given** the frontend uses components and routes named `RoomsView`, `DesksView`, etc.
**When** the rename is applied
**Then** components and routes use item-group and item terminology
**And** the Vue Router paths match the new API routes

**Given** the areas list view shows a "VIEW ROOM" button on each area tile
**When** the rename is applied
**Then** the button reads "BOOK"
**And** no UI element references "room" or "desk" in user-facing text

**Given** the item detail view shows a "BOOK THIS DESK" button
**When** the rename is applied
**Then** the button reads "BOOK THIS ITEM"

**Given** a booking is successfully created
**When** the confirmation message is displayed
**Then** the message references the item name from the configuration
(e.g., "Parking Lot 1 booked successfully") rather than the generic term "desk"

**Given** the Pinia stores and API service files reference rooms/desks
**When** the rename is applied
**Then** all store names, API paths, and TypeScript types use the new terminology
**And** `npm run type-check` and `npm run build` succeed without errors
**And** `npx vitest run` passes all unit tests

### Story 12.5: Update Documentation and E2E Tests

**FRs covered:** FR42 (partial)

As a developer,
I want API documentation and E2E tests updated to use the new terminology,
So that the entire codebase is consistent and all tests pass.

**Acceptance Criteria:**

**Given** the OpenAPI documentation references `/rooms` and `/desks` endpoints
**When** the documentation is updated
**Then** all endpoint paths, schemas, and descriptions use the new terminology
**And** `npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml` passes

**Given** the Cypress E2E tests reference rooms/desks in selectors and assertions
**When** the tests are updated
**Then** all E2E tests use the new terminology in routes, selectors, and assertions
**And** `npm run test:e2e -- --browser electron` passes all tests

**Given** the Go code duplication check runs
**When** `npx jscpd --pattern "**/*.go" --ignore "**/*_test.go" --threshold 3` is executed
**Then** the duplication threshold is not exceeded

**Given** the TypeScript code duplication check runs
**When** `npx jscpd --pattern "**/*.ts" --ignore "**/node_modules/**" --threshold 0` is
executed
**Then** the duplication threshold is not exceeded

## Epic 13 Stories: Enhanced Booking Experience

Users get powerful new booking capabilities: weekly availability previews, booking notes,
week-at-a-time booking mode, visibility of who booked what, and clickable navigation
breadcrumbs.
**FRs covered:** FR36, FR37, FR38, FR39, FR40

### Story 13.1: Clickable Breadcrumbs

**FRs covered:** FR40

As a user,
I want breadcrumbs to be clickable links,
So that I can navigate back to any level of the hierarchy quickly.

**Acceptance Criteria:**

**Given** I am viewing items within an item group
**When** I see the breadcrumb showing "Home > Office 1st Floor > Room 101"
**Then** "Home" and "Office 1st Floor" are clickable links
**And** clicking "Home" navigates to the areas list
**And** clicking "Office 1st Floor" navigates to the item groups for that area
**And** "Room 101" (the current page) is not clickable

**Given** I am viewing item groups within an area
**When** I see the breadcrumb showing "Home > Office 1st Floor"
**Then** "Home" is a clickable link that navigates to the areas list
**And** "Office 1st Floor" (the current page) is not clickable

### Story 13.2: Booker Display Name

**FRs covered:** FR39

As a user,
I want to see who has booked an item,
So that I know which colleagues are in the office or using a resource.

**Acceptance Criteria:**

**Given** I am viewing items in an item group for a specific date
**When** an item is booked
**Then** I see the booker's display name alongside the booking status
**And** the display name is resolved from the users table (not stored in the booking)

**Given** an item is available (not booked)
**When** the item is displayed
**Then** no booker name is shown
**And** the item is clearly marked as available

**Given** I am viewing Today's Presence for an area
**When** the presence list is displayed
**Then** each entry shows the user's display name

### Story 13.3: Weekly Availability Preview

**FRs covered:** FR36

As a user,
I want to see a weekly availability preview on item group tiles,
So that I can quickly identify which days have open items without clicking into each group.

**Acceptance Criteria:**

**Given** I am viewing item groups within an area
**When** the page loads
**Then** I see a calendar week selector above the item group tiles
**And** the current week is pre-selected
**And** only the next 8 weeks are available for selection
**And** each week option displays the Monday date and week number in locale-aware format
(e.g., "2026-03-16 - Week 12")

**Given** I have selected a calendar week
**When** the item group tiles are displayed
**Then** each tile shows weekday indicators (MO, TU, WE, TH, FR)
**And** green indicates at least one item is available on that day
**And** red indicates all items in the group are fully booked on that day
**And** each indicator uses a secondary visual cue in addition to color
(e.g., filled circle for available, empty circle for booked) to meet WCAG A

**Given** the backend receives a request for weekly availability
**When** `GET /api/v1/areas/:area_id/item-groups/availability?week=YYYY-Www` is called
**Then** it returns per-day availability counts for each item group in the area
**And** the response includes item group ID, day, total items, and available count
**And** the response completes within the NFR1 performance target (2s at p95)

**Given** I change the selected calendar week
**When** the new week is applied
**Then** the weekday indicators update to reflect availability for the new week

### Story 13.4: Booking Notes

**FRs covered:** FR37

As a user,
I want to add, view, and edit notes on my bookings,
So that I can communicate useful information to colleagues (e.g., "arriving after noon").

**Acceptance Criteria:**

**Given** I have just completed a booking
**When** the success confirmation message is displayed
**Then** I see an "add note" action within the confirmation
**And** clicking it opens a text input where I can type a free-text note
**And** the note is saved to the booking

**Given** I am viewing My Bookings
**When** a booking has a note
**Then** the note is displayed as a single line on the booking card
**And** if the note is longer than the available width, it is truncated
**And** a truncation indicator (icon) signals there is more text
**And** clicking the indicator opens a dialog (desktop) or bottom sheet (mobile)
showing the full note text

**Given** I am viewing My Bookings
**When** I want to edit a note on one of my bookings
**Then** I can modify the note text and save the changes
**And** the updated note is reflected immediately

**Given** I am viewing Today's Presence for an area
**When** a booking has a note
**Then** the note is displayed with the same truncation behavior as My Bookings

**Given** I am viewing items in an item group
**When** a booked item has a note
**Then** the note is displayed alongside the booker's name with truncation behavior

**Given** the backend receives a note update request
**When** the note is saved via the API
**Then** the booking record is updated with the new note text
**And** a `note` field is added to the bookings table via migration
**And** the JSON:API response includes the note in booking attributes

### Story 13.5: Week Booking Mode

**FRs covered:** FR38

As a user,
I want to switch to week booking mode and book multiple days at once,
So that I can reserve my workspace for an entire week efficiently.

**Acceptance Criteria:**

**Given** I am viewing items in an item group
**When** I see the booking mode toggle
**Then** I can switch between "book by day" and "book by week" modes
**And** the selected mode is persisted in browser local storage
**And** on my next visit, the previously selected mode is restored

**Given** I have selected week booking mode
**When** the date selector is displayed
**Then** it becomes a calendar week selector showing the next 8 weeks
**And** each week option displays the Monday date and week number

**Given** I have selected a week in week booking mode
**When** the item tiles are displayed
**Then** each tile shows a per-day breakdown (MO through FR) with checkboxes
**And** days booked by other users show the booker's name in red text
**And** days booked by other users have their checkboxes disabled (cannot be unchecked)
**And** days I have already booked show my name with a checked checkbox
**And** free days show "free" in green text with an unchecked checkbox

**Given** I am viewing week booking mode on a screen narrower than 600px
**When** the item tiles are displayed
**Then** the per-day breakdown uses a compact layout suitable for mobile
**And** touch targets meet the minimum 44px size requirement

**Given** I have checked one or more free days across one or more items
**When** I look below the item tiles
**Then** I see a single "Confirm My Booking" button
**And** the individual "BOOK THIS ITEM" buttons are not shown in week mode

**Given** I click "Confirm My Booking" with multiple days selected
**When** the bookings are submitted
**Then** each day/item combination is submitted as an individual API request
to `POST /api/v1/bookings`
**And** results are collected and displayed per-day
**And** each successful booking appears in My Bookings
**And** if any day fails (e.g., concurrent booking), the error is reported for that day
**And** successful bookings are not rolled back due to a single day's failure

**Given** I switch back to day booking mode
**When** the view updates
**Then** the standard single-day booking interface is restored
**And** the "BOOK THIS ITEM" button reappears on each item tile

## Epic 14 Stories: UI Cleanup & Booking Simplification

Remove visual clutter, simplify action labels, streamline the booking form, and improve the
display of booked items for a cleaner, faster user experience.
**FRs covered:** FR43, FR44, FR45, FR46, FR47, FR48, FR49

### Story 14.1: Simplify Action Labels Across Views

**FRs covered:** FR43

As a user,
I want concise action labels that get me to my destination faster,
So that I spend less time reading and more time booking.

**Acceptance Criteria:**

**Given** I am viewing the areas list (Home)
**When** I see an area tile
**Then** the action button reads "SELECT" instead of "VIEW ITEM GROUPS"

**Given** I am viewing item groups within an area
**When** I see the page
**Then** the page title "Item Groups" and subtitle "Select an item group to view available
items" are removed
**And** the action button on each tile reads "SELECT" instead of "VIEW ITEMS"

**Given** I am viewing items within an item group
**When** I see the page
**Then** the page title "Items" and subtitle "Select an item to book for your chosen date"
are removed

**Given** I am viewing available items in day booking mode
**When** I see an available item tile
**Then** the booking button reads "BOOK" instead of "BOOK THIS ITEM"

### Story 14.2: Streamline Booking Form

**FRs covered:** FR47, FR48, FR49

As a user,
I want a simplified booking form with fewer options,
So that the interface is less overwhelming and common tasks are faster.

**Acceptance Criteria:**

**Given** I am on the booking page
**When** I see the booking type options
**Then** "Book for guest" is not available as an option
**And** only "Book for myself" and "Book for colleague" are shown

**Given** I am on the booking page in day booking mode
**When** I see the booking options
**Then** the "Book multiple days" checkbox is not shown
**And** no additional dates field appears

**Given** I select "Book for colleague"
**When** the colleague fields appear
**Then** I see a dropdown listing existing users by display name (fetched from
`GET /api/v1/users`)
**And** selecting a user from the dropdown sets the booking to use that user's ID
**And** the free-text colleague email and name fields are removed

### Story 14.3: Improve Booked Item Display

**FRs covered:** FR44, FR45, FR46

As a user,
I want booked item details to be clearly readable and booking results to use icons,
So that I can quickly understand who booked what and whether my bookings succeeded.

**Acceptance Criteria:**

**Given** I am viewing items in day booking mode
**When** an item is booked by another user
**Then** the "Not available for \<date\>" message is not shown
**And** the booker name is displayed at body-2 size or larger (not caption)
**And** any booking note is displayed at body-2 size or larger

**Given** I submit bookings (day or week mode)
**When** results are displayed
**Then** each successful booking shows a green checkmark icon with item name and day
**And** each failed booking shows a red warning icon with item name and error detail
**And** raw text labels like "Booked" are replaced by the icons

## Epic 15 Stories: Collapsible Item Tiles

Introduce fold/unfold tile behavior across day and week booking modes to manage visual
complexity, showing equipment, warnings, and full details on demand.
**FRs covered:** FR50, FR51, FR52, FR53, FR54

### Story 15.1: Collapsible Tiles in Week Booking Mode

**FRs covered:** FR50, FR54

As a user,
I want to expand item tiles in week mode to see full details,
So that the default view is compact and I can drill into specifics on demand.

**Acceptance Criteria:**

**Given** I am in week booking mode viewing item tiles
**When** I see a tile
**Then** a chevron-left icon appears in the tile header

**Given** I click the chevron on a folded tile
**When** the tile unfolds
**Then** the chevron rotates to chevron-down
**And** the compact M-F row is replaced by one line per day
**And** each line shows the full day name (Monday, Tuesday, etc.)
**And** each line shows the full booker display name (not truncated)
**And** equipment chips and warning alerts are visible below the daily breakdown

**Given** I click the chevron on an unfolded tile
**When** the tile folds
**Then** the chevron rotates back to chevron-left
**And** the compact M-F row is restored

**Given** I am viewing a folded tile with truncated booker names
**When** I hover over a truncated name
**Then** a tooltip shows the full display name

### Story 15.2: Collapsible Tiles in Day Booking Mode

**FRs covered:** FR51

As a user,
I want booked item tiles in day mode to be collapsed by default,
So that I can focus on available items and expand booked ones only when needed.

**Acceptance Criteria:**

**Given** I am in day booking mode
**When** an item is booked
**Then** the item tile hides equipment and warning details by default
**And** a chevron-left icon appears in the tile header

**Given** I click the chevron on a folded booked item tile
**When** the tile unfolds
**Then** equipment chips and warning alerts become visible
**And** the chevron rotates to chevron-down

**Given** I am in day booking mode
**When** an item is available
**Then** the item tile shows all details (equipment, warnings) without a chevron
**And** the tile is not collapsible

### Story 15.3: Warning Icon on Folded Tiles

**FRs covered:** FR52

As a user,
I want to know about item warnings even when a tile is folded,
So that I can make informed decisions without expanding every tile.

**Acceptance Criteria:**

**Given** a tile is folded and the item has a warning
**When** I see the tile header
**Then** a warning icon is visible

**Given** I click the warning icon on a folded tile
**When** the popup or tooltip appears
**Then** the full warning message is displayed
**And** I do not need to unfold the tile to read the warning

**Given** a tile is folded and the item has no warning
**When** I see the tile header
**Then** no warning icon appears

### Story 15.4: Disable Past Date Checkboxes in Week Mode

**FRs covered:** FR53

As a user,
I want past date checkboxes disabled in week booking mode,
So that I don't waste time selecting dates the backend would reject anyway.

**Acceptance Criteria:**

**Given** I am in week booking mode and the selected week includes past dates
**When** I see the per-day checkboxes
**Then** checkboxes for dates before today are disabled and visually grayed out
**And** I cannot check or uncheck past date checkboxes

**Given** I am in week booking mode and the selected week is entirely in the future
**When** I see the per-day checkboxes
**Then** all free day checkboxes are enabled and interactive

## Epic 16 Stories: User Preferences & Settings

Let users personalize their booking experience with theme selection, weekend visibility,
and minor menu fixes.
**FRs covered:** FR55, FR56, FR57

### Story 16.1: Theme Selector

**FRs covered:** FR55

As a user,
I want to choose between light, dark, and auto themes,
So that the app matches my visual preference or adapts to my system setting.

**Acceptance Criteria:**

**Given** I click my username in the top right corner
**When** I see the user menu
**Then** I see a theme option with choices: auto (default), dark, light

**Given** I select "dark" theme
**When** the selection is applied
**Then** the Vuetify dark theme is activated immediately
**And** my choice is persisted in localStorage
**And** on my next visit, the dark theme is applied without manual selection

**Given** I select "auto" theme
**When** my OS preference is dark mode
**Then** the app uses dark theme
**And** when my OS switches to light mode, the app follows

### Story 16.2: Show Weekends Toggle

**FRs covered:** FR56

As a user,
I want to optionally show weekends on booking pages,
So that I can book Saturday and Sunday if my workplace supports it.

**Acceptance Criteria:**

**Given** I click my username in the top right corner
**When** I see the user menu
**Then** I see a "Show weekends" checkbox (unchecked by default)

**Given** I enable "Show weekends"
**When** I view the booking page in week mode
**Then** Saturday and Sunday columns appear alongside Monday through Friday
**And** my preference is persisted in localStorage

**Given** I enable "Show weekends"
**When** I view the weekly availability indicators on item group tiles
**Then** the indicators include Saturday and Sunday

**Given** I disable "Show weekends"
**When** I view any booking page
**Then** only Monday through Friday are shown

### Story 16.3: Fix Change Password Icon

**FRs covered:** FR57

As a user,
I want the Change Password menu item to show its icon,
So that the menu has a consistent visual appearance.

**Acceptance Criteria:**

**Given** I am a local user viewing the desktop user menu
**When** I see the "Change Password" item
**Then** an icon is displayed next to the text (consistent with other menu items)

**Given** I am a local user viewing the mobile navigation drawer
**When** I see the "Change Password" item
**Then** an icon is displayed next to the text

## Epic 17 Stories: Equipment Filter

Enable users to filter items by equipment keywords to quickly find suitable workspaces
using an advanced search syntax.
**FRs covered:** FR58

### Story 17.1: Equipment Filter with Advanced Search

**FRs covered:** FR58

As a user,
I want to filter items by equipment keywords,
So that I can quickly find a workspace with the tools I need.

**Acceptance Criteria:**

**Given** I am on the booking page (day or week mode)
**When** I see the booking options card
**Then** a text input labeled "Filter equipment" appears below the colleague option
**And** an info icon appears next to the input

**Given** I click the info icon
**When** the explanation popup appears
**Then** it describes the search syntax:
show only items having the filter keyword(s) in any of the equipment items;
multiple keywords are combined with OR;
use plus sign to combine with AND;
use single or double quotation marks for exact matching;
filters are case-insensitive;
example: "27 inch display" + webcam

**Given** I type "webcam" into the filter input
**When** the filter is applied
**Then** items that have "webcam" in any of their equipment are shown normally
**And** items without "webcam" in their equipment are blurred with an "equipment not
available" overlay hint
**And** blurred items are not removed from the list

**Given** I type `"27 inch display" + webcam` into the filter input
**When** the filter is applied
**Then** only items having both "27 inch display" (exact) AND "webcam" in their equipment
are shown normally
**And** all other items are blurred

**Given** I clear the filter input
**When** the filter is removed
**Then** all items are shown normally without blur

**Given** I am in week booking mode
**When** I type a filter
**Then** the same filtering logic applies to the week mode item tiles

## Epic 18 Stories: Floor Plan Display & Config Consistency

Users can view floor plan images for areas and item groups while booking, and operators
benefit from consistent configuration terminology and stricter file location validation.
**FRs covered:** FR59, FR60, FR61, FR62, FR63, FR64, FR65, FR66

### Story 18.1: Rename Config Section from [spaces] to [areas]

**FRs covered:** FR59

As an operator,
I want the configuration to use consistent `[areas]` terminology everywhere,
So that there is no confusion between "spaces" and "areas" in configuration, code,
and documentation.

**Acceptance Criteria:**

**Given** I have a `sithub.toml` with `[areas]` section
**When** the server starts
**Then** it reads configuration from the `[areas]` table

**Given** I use `--areas-config-file` flag or `SITHUB_AREAS_CONFIG_FILE` env var
**When** the server starts
**Then** it applies the override correctly

**Given** the codebase
**When** searching for the term "space" or "spaces"
**Then** no references to the old terminology exist (data models, CLI flags, env vars)

### Story 18.2: Enforce Areas Config Inside data_dir

**FRs covered:** FR60

As an operator,
I want the areas config file path resolved relative to `main.data_dir`,
So that all data files are consistently located in one directory.

**Acceptance Criteria:**

**Given** `areas.config_file` is set to a relative filename like `sithub_areas.yaml`
**When** the server starts
**Then** the file is resolved inside `main.data_dir`

**Given** `areas.config_file` contains an absolute path outside `main.data_dir`
**When** the server starts
**Then** startup fails with a descriptive error message

### Story 18.3: Floor Plans Directory Configuration & Validation

**FRs covered:** FR61, FR62

As an operator,
I want to configure a floor plans directory and have all image references validated
at startup,
So that runtime errors from missing or invalid images are caught early.

**Acceptance Criteria:**

**Given** `areas.floor_plans` is set to a directory name
**When** the server starts
**Then** the directory is resolved inside `main.data_dir` and its existence is validated

**Given** `areas.floor_plans` points to a non-existent directory
**When** the server starts
**Then** the server exits with a descriptive error

**Given** the areas config references floor plan images
**When** the server starts
**Then** each referenced image is checked for existence in the floor plans directory

**Given** a referenced image has an unsupported format (not jpg, png, svg)
**When** the server starts
**Then** the server exits with a descriptive error listing the invalid file

**Given** `areas.floor_plans` is not set
**When** the server starts
**Then** floor plan features are disabled and no validation occurs

### Story 18.4: Authenticated Floor Plan Image Serving

**FRs covered:** FR63

As a user,
I want floor plan images served through an authenticated endpoint,
So that floor plans are protected from unauthorized access.

**Acceptance Criteria:**

**Given** I am authenticated
**When** I request a floor plan image via the API
**Then** the image is returned with the correct content type

**Given** I am not authenticated
**When** I request a floor plan image
**Then** the request is denied with 401

**Given** I request an image that does not exist
**When** the request is processed
**Then** I receive a 404 response

### Story 18.5: Area Floor Plan Display

**FRs covered:** FR64

As a user,
I want to see a "Floor plan" button when viewing an area that has a floor plan,
So that I can see where items are located.

**Acceptance Criteria:**

**Given** I am viewing an area that has a floor plan configured
**When** the page loads
**Then** a "Floor plan" button with an appropriate icon appears next to the calendar
week selector

**Given** I click the "Floor plan" button
**When** the overlay opens
**Then** the floor plan image is displayed with the area name as heading
**And** I can close the overlay

**Given** I am viewing an area without a floor plan
**When** the page loads
**Then** no "Floor plan" button appears

### Story 18.6: Item Group Floor Plan Display

**FRs covered:** FR65

As a user,
I want to see a "Floor plan" button when viewing an item group that has a floor plan,
So that I can see the layout of individual items within the group.

**Acceptance Criteria:**

**Given** I am viewing an item group that has a floor plan configured
**When** the page loads
**Then** a "Floor plan" button with an appropriate icon appears beneath the day/week
selector

**Given** I click the "Floor plan" button
**When** the overlay opens
**Then** the floor plan image is displayed with the item group name as heading
**And** I can close the overlay

**Given** I am viewing an item group without a floor plan
**When** the page loads
**Then** no "Floor plan" button appears

### Story 18.7: Connection Lost Error Messaging

**FRs covered:** FR66

As a user,
I want to see a clear "Connection to server lost" error when the backend is unavailable,
So that I understand the real problem instead of seeing misleading messages like
"no areas available."

**Acceptance Criteria:**

**Given** the backend server is unavailable
**When** the frontend attempts to load data
**Then** a clear error message "Connection to server lost" is displayed

**Given** the backend was available and then goes down
**When** subsequent API calls fail
**Then** the error message is shown instead of empty or misleading content

## Epic 19 Stories: User Feedback — Bug Fixes & Feature Requests

Users benefit from a smoother booking experience through bug fixes and new capabilities
including equipment filter enhancements, quick cancellation from week view, customizable
icons, an improved calendar/week selector, and a favorites system.
**FRs covered:** FR67, FR68, FR69, FR70, FR71, FR72, FR73, FR74

### Story 19.1: Fix Cancel Booking Dialog Not Closing

**FRs covered:** FR67

As a user,
I want the cancel booking confirmation dialog to close after I confirm the cancellation,
So that I am not left with a stale dialog on screen.

**Acceptance Criteria:**

**Given** I am on the My Bookings page and click cancel on a booking
**When** the confirmation dialog appears and I click "Cancel Booking"
**Then** the booking is removed from the list
**And** the confirmation dialog closes automatically

### Story 19.2: Week Selector Date Range Display

**FRs covered:** FR71

As a user,
I want the calendar week selector to show both the first and last day of each week,
So that I can immediately see which date range a calendar week covers.

**Acceptance Criteria:**

**Given** I am on a view with the week selector
**When** I open the week selector dropdown
**Then** each option shows the format "DD.MM.-DD.MM.YYYY - Week N"
(e.g. "23.03.-29.03.2026 - Week 13")

**Given** the show weekends toggle is off
**When** I view the week selector
**Then** the date range still shows Monday through Sunday (full week),
regardless of the weekends setting

### Story 19.3: Calendar Widget Starts on Monday

**FRs covered:** FR72

As a user,
I want the calendar date picker to show Monday as the first day of the week,
So that it matches the European convention I am used to.

**Acceptance Criteria:**

**Given** I am on any view with a date picker
**When** the calendar widget opens
**Then** Monday is displayed as the first (leftmost) column
**And** Sunday is displayed as the last (rightmost) column

### Story 19.4: Cancel Booking from Week View

**FRs covered:** FR70

As a user,
I want to cancel my bookings directly from the week view,
So that I don't have to navigate to My Bookings to undo a booking.

**Acceptance Criteria:**

**Given** I am on the week view and a day/item has my booking (shown with a checkmark)
**When** the page renders
**Then** a small red cancel icon appears next to the checkmark

**Given** I click the red cancel icon
**When** the cancellation is processed
**Then** the booking is cancelled and the checkmark and cancel icon are removed
**And** the item becomes bookable again for that day

**Given** the booking belongs to another user
**When** the page renders
**Then** no cancel icon is shown for that booking

### Story 19.5: Equipment Filter on Item Groups View

**FRs covered:** FR68

As a user,
I want to filter item groups by equipment on the area view,
So that I can quickly find rooms or areas that have the equipment I need.

**Acceptance Criteria:**

**Given** I am on the item-groups view (e.g. `/areas/{areaId}/item-groups`)
**When** I enter an equipment filter keyword
**Then** item groups whose items do not match the filter are blurred and disabled
**And** item groups with at least one matching item are shown normally

**Given** I clear the filter
**When** the filter is removed
**Then** all item groups are shown normally without blur

**Given** I use the advanced filter syntax (AND with `+`, exact match with quotes)
**When** the filter is applied
**Then** the same parsing rules from the existing equipment filter apply

### Story 19.6: Equipment Filter Saving

**FRs covered:** FR69

As a user,
I want to save my equipment filters for reuse,
So that I don't have to retype the same filter keywords every time I book.

**Acceptance Criteria:**

**Given** I have typed a filter into the equipment filter input
**When** I click the save icon next to the input
**Then** the filter is saved to browser local storage
**And** a confirmation is shown

**Given** I have saved filters
**When** I focus the equipment filter input
**Then** a combobox dropdown shows my saved filters alongside free-text input

**Given** I select a saved filter from the combobox
**When** the filter loads
**Then** the save icon becomes a delete icon

**Given** I click the delete icon on a loaded saved filter
**When** the deletion is confirmed
**Then** the filter is removed from local storage
**And** the input is cleared

**Given** I have no saved filters
**When** the page loads
**Then** the input behaves as a regular text field with no dropdown entries

### Story 19.7: Custom Icons in Areas YAML

**FRs covered:** FR73

As an operator,
I want to specify custom MDI icons for areas, item groups, and items in the areas YAML,
So that the UI reflects the actual purpose of each space with meaningful icons.

**Acceptance Criteria:**

**Given** an area, item group, or item in the areas YAML has an `icon` field
(e.g. `icon: mdi-office-building`)
**When** the frontend renders that entity
**Then** the specified MDI icon is displayed instead of the default icon

**Given** an entity does not have an `icon` field
**When** the frontend renders that entity
**Then** the current default icon is used

**Given** the `icon` field contains an invalid or unknown MDI icon name
**When** the server starts
**Then** a warning is logged but the server does not fail to start
**And** the frontend falls back to the default icon

**Given** the areas API returns the `icon` attribute
**When** the frontend receives the response
**Then** the icon value is available for rendering at all three levels
(area, item group, item)

### Story 19.8: Favorites

**FRs covered:** FR74

As a user,
I want to mark item groups and items as favorites,
So that my most-used spaces appear first and are quick to find.

**Acceptance Criteria:**

**Given** I am on the item-groups view (second level)
**When** I see an item group tile
**Then** a heart outline icon is visible on the tile

**Given** I click the heart outline on an item group
**When** the favorite is saved
**Then** a confirmation "{item group name} saved as favorite." is shown
**And** the icon becomes a red-filled heart
**And** the favorite is persisted in browser local storage

**Given** I click a red-filled heart on an item group
**When** the favorite is removed
**Then** a confirmation "{item group name} removed from favorites." is shown
**And** the icon reverts to a heart outline

**Given** I am on the items view (third level)
**When** I see an item tile
**Then** a heart outline icon is visible on the tile
**And** clicking it saves/removes the favorite with confirmation
"{item group name} {item name} saved/removed as favorite."

**Given** I have third-level favorites
**When** I view the item-groups page (second level)
**Then** my third-level favorites appear as bookable tiles on that page

**Given** I am on the item-groups view with favorites
**When** the page renders
**Then** items are ordered: (1) third-level favorites A-Z,
(2) second-level favorites A-Z, (3) remaining item groups in YAML order
with second-level favorites subtracted

## Epic 20 Stories: Interactive Floor Plans & UX Consistency

Users can view live free/busy status on floor plan overlays, book items directly from
floor plans, and admins can position items on floor plan images. Navigation state is
preserved across the app and confirmations use a consistent style.
**FRs covered:** FR75, FR76, FR77, FR78, FR79, FR80, FR81, FR82, FR83, FR84

### Story 20.1: Free-Busy Indicators on Favorite Tiles

**FRs covered:** FR75

As a user,
I want to see weekly availability indicators on my promoted third-level favorite tiles,
So that I can quickly see which days have availability without navigating into the
item group.

**Acceptance Criteria:**

**Given** I have third-level favorites promoted to the item-groups view
**When** the page loads and availability data is fetched
**Then** the favorite tiles show the same MO-TU-WE-TH-FR availability dots as regular
item group tiles

**Given** an item within a favorite's item group is fully booked on a day
**When** the availability dot renders
**Then** the dot shows the booked (red outline) indicator for that day

### Story 20.2: Memorize Selected Week and Day

**FRs covered:** FR76, FR77

As a user,
I want the selected week and day to persist as I navigate between areas and item groups,
So that I don't have to re-select the same date on every page.

**Acceptance Criteria:**

**Given** I select week 16 on the item-groups view
**When** I navigate to an item group and back to the item-groups view
**Then** week 16 is still selected

**Given** I select a specific day on the items view
**When** I navigate to a different item group
**Then** the same day is pre-selected

**Given** the memorized week is in the past
**When** I return to the view
**Then** the week resets to the current week

**Given** I successfully book an item
**When** the booking succeeds
**Then** the memorized day resets to today

### Story 20.3: Consistent Snackbar Confirmations

**FRs covered:** FR78

As a user,
I want all confirmations to use the same bottom snackbar style,
So that the feedback is consistent and predictable across the app.

**Acceptance Criteria:**

**Given** I cancel a booking from My Bookings
**When** the cancellation succeeds
**Then** a bottom snackbar shows "Booking cancelled successfully."
(not an inline alert)

**Given** I perform any action that shows a success confirmation
**When** the confirmation appears
**Then** it uses a bottom snackbar, matching the style used for favorites and filter
confirmations

### Story 20.4: Floor Plan Positions Database Schema and API

**FRs covered:** FR82

As a developer,
I want floor plan item positions stored in SQLite with a CRUD API,
So that the floor plan editor and viewer have a backend to read and write positions.

**Acceptance Criteria:**

**Given** an admin saves item positions for a floor plan
**When** the positions are persisted
**Then** they are stored in a `floor_plan_positions` table with floor plan filename,
item ID, and rectangle coordinates (x, y, width, height)

**Given** a user requests positions for a floor plan
**When** the API responds
**Then** it returns all positions for that floor plan as a JSON:API collection

**Given** an admin updates a position
**When** the PUT request is processed
**Then** the position is updated in the database

**Given** an admin deletes a position
**When** the DELETE request is processed
**Then** the position is removed from the database

### Story 20.5: Floor Plan Editor (Admin)

**FRs covered:** FR81

As an admin,
I want to draw rectangles on floor plan images to mark where items are located,
So that users can see and click items on the interactive floor plan.

**Acceptance Criteria:**

**Given** I am an admin and open the floor plan editor from settings
**When** I select a floor plan
**Then** the floor plan image is displayed with a list of unpositioned items on the left

**Given** I select an item from the list
**When** I draw a rectangle on the floor plan image
**Then** the rectangle is created with a label showing the item name

**Given** I have positioned items on the floor plan
**When** I save
**Then** all positions are persisted via the API

**Given** I want to reposition an item
**When** I drag or resize its rectangle
**Then** the position updates visually and can be saved

**Given** I want to remove a positioned item
**When** I delete it
**Then** the rectangle is removed from the floor plan

### Story 20.6: Interactive Floor Plan Overlay with Free/Busy

**FRs covered:** FR79

As a user,
I want to see free/busy status on the floor plan and book items by clicking them,
So that I can visually find and book available items.

**Acceptance Criteria:**

**Given** I open the floor plan overlay for an item group
**When** the overlay renders
**Then** the floor plan image is shown with positioned items overlaid as rectangles

**Given** a weekday selector appears at the top of the overlay
**When** I select a day
**Then** free items show a green outline and busy items show a red semi-transparent overlay

**Given** the floor plan opens for the current week
**When** today is within the week
**Then** today is pre-selected and past days are disabled

**Given** the floor plan opens for a future week
**When** the overlay renders
**Then** Monday is pre-selected

**Given** I click on a free item
**When** the click is processed
**Then** a booking is created for the selected day and the item status updates to busy

**Given** weekend visibility is off in settings
**When** the weekday selector renders
**Then** Saturday and Sunday are not shown

### Story 20.7: First-Level Floor Plan Drill-Down

**FRs covered:** FR80

As a user,
I want to click on an area in the first-level floor plan to open its detail floor plan,
So that I can drill down from the building overview to individual items.

**Acceptance Criteria:**

**Given** I open the floor plan for an area that has sub-areas with their own floor plans
**When** the overlay renders
**Then** each sub-area is shown with its positioned rectangle and free/busy state

**Given** all items within a sub-area are booked for the selected day
**When** the sub-area renders
**Then** it shows a red semi-transparent overlay

**Given** I click on a sub-area rectangle
**When** the click is processed
**Then** the detail floor plan for that sub-area opens with item-level free/busy state

### Story 20.8: Floor Plan Booking UX Refinements

**FRs covered:** FR83, FR84

As a user,
I want the floor plan booking experience to support multi-day selection, provide precise
feedback, and work reliably on mobile,
So that I can efficiently book items for multiple days and trust the floor plan interaction
on any device.

**Acceptance Criteria:**

**Given** I click on a free item on a detail-level floor plan
**When** the booking dialog opens
**Then** it shows weekday checkboxes (abbreviations only: Mo, Tu, We, Th, Fr) with the
currently selected day pre-checked; past days and already-booked days are disabled

**Given** I select days and click "Book now"
**When** the booking is submitted
**Then** a summary shows "Book [Item] in [Area] for N days starting [date]" and all
selected days are booked; the "Book now" and "Cancel" buttons are always visible

**Given** a booking fails because the item was booked by someone else
**When** the error is displayed
**Then** the message names the specific day: "The selected item is already booked on
[day]."

**Given** the floor plan overlay is open
**When** I click outside the overlay
**Then** the overlay does NOT close; only the close button dismisses it

**Given** I view the floor plan on a small screen
**When** the overlay renders
**Then** it opens fullscreen with "Show labels" and "Close" at the bottom

**Given** I am in a detail floor plan opened via drill-down
**When** I click close/back
**Then** I return to the higher-level floor plan, not the underlying page

**Given** a sub-area on a first-level floor plan has its own detail floor plan
**When** I click anywhere on it
**Then** the detail floor plan opens; direct booking is prevented

## Epic 21 Stories: i18n, UX Improvements & Booking Limits

Users can switch the UI language with auto-detection, benefit from visual refinements
across booking views, and operators can enforce booking limits via configuration.
**FRs covered:** FR85, FR86, FR87, FR88, FR89, FR90

### Story 21.1: i18n Infrastructure and English Baseline

**FRs covered:** FR85

As a developer,
I want vue-i18n configured with all existing UI strings extracted into an English message
file,
So that the app is ready for translation without changing user-visible behavior.

**Acceptance Criteria:**

**Given** the app starts
**When** no language preference is stored
**Then** the UI renders in English, identical to the current behavior

**Given** any component renders text
**When** the text is user-visible (labels, buttons, messages, headings, placeholders)
**Then** the text comes from the i18n message file, not hardcoded in the template

**Given** the English message file exists
**When** a developer inspects it
**Then** all keys are organized by feature area (e.g., `auth.login`, `bookings.cancel`,
`settings.theme`) and use dot-notation nesting

### Story 21.2: Language Selector with Flags and Auto-Detection

**FRs covered:** FR85

As a user,
I want to choose my preferred UI language from the settings page,
So that I can use SitHub in my native language.

**Acceptance Criteria:**

**Given** I open the settings page
**When** I see the language selector
**Then** it shows options: Auto, English, Deutsch, Español, Français, Українська —
each with a colored country flag (UK for English, DE for German, ES for Spanish,
FR for French, UA for Ukrainian)

**Given** I select "Deutsch"
**When** the selection is applied
**Then** the entire UI switches to German immediately without page reload

**Given** I select "Auto"
**When** my browser's preferred language is German
**Then** the UI renders in German

**Given** I select "Auto"
**When** my browser's preferred language is not one of the supported languages
**Then** the UI falls back to English

**Given** I select a language and close the browser
**When** I reopen SitHub
**Then** the previously selected language is restored from local storage

### Story 21.3: German, Spanish, French, and Ukrainian Translations

**FRs covered:** FR85

As a user,
I want all UI text translated into German, Spanish, French, and Ukrainian,
So that I can use SitHub fully in my preferred language.

**Acceptance Criteria:**

**Given** the language is set to German (or Spanish, French, Ukrainian)
**When** I navigate through the app
**Then** all labels, buttons, messages, headings, placeholders, and error messages
appear in the selected language

**Given** translation files exist for all four languages
**When** a developer inspects them
**Then** every key present in the English file has a corresponding entry in each
translation file with no missing keys

**Given** the backend returns error messages (e.g., booking conflicts)
**When** the frontend displays them
**Then** the messages are localized using frontend translation keys, not raw backend
strings

### Story 21.4: My Bookings Display Reorder

**FRs covered:** FR86

As a user,
I want to see the booking date prominently on the first line of each booking card,
So that I can quickly scan my bookings by date.

**Acceptance Criteria:**

**Given** I navigate to My Bookings
**When** the booking cards render
**Then** each card shows the booked date (with calendar icon) on the first line
and the booked item name with area breadcrumb on the second line

**Given** the current layout shows item first and date second
**When** this story is implemented
**Then** the order is swapped: date first, item second

### Story 21.5: Equipment Filter Icon and Floor Plan Button Fixes

**FRs covered:** FR87, FR88

As a user,
I want visual consistency in the booking toolbar,
So that icons communicate their purpose and controls are aligned predictably.

**Acceptance Criteria:**

**Given** I type an equipment filter and see the save button
**When** I look at the save icon
**Then** it shows `mdi-content-save` (floppy disk) instead of the plus icon

**Given** I am on an item group view that has a floor plan
**When** the toolbar renders
**Then** the floor plan button has the same height as the calendar week selector

**Given** I am on an item group view with a detail floor plan
**When** the toolbar renders
**Then** the floor plan button is positioned next to the calendar week selector,
not below the booking mode toggle

**Given** I am on an item group view without a floor plan
**When** the toolbar renders
**Then** no floor plan button is shown (existing behavior preserved)

### Story 21.6: Booking Advance Limit

**FRs covered:** FR89

As an operator,
I want to configure how far in advance users can book,
So that I can prevent booking too far into the future.

**Acceptance Criteria:**

**Given** sithub.toml contains `weeks_in_advanced = 3` under `[bookings]`
**When** the calendar week selector renders
**Then** only the current week plus the next 3 weeks are shown; weeks beyond that
are not available

**Given** sithub.toml does not contain `weeks_in_advanced`
**When** the calendar week selector renders
**Then** the default of 5 weeks ahead applies

**Given** a user attempts to book a date beyond the allowed advance window via API
**When** the request is processed
**Then** it is rejected with a clear error: "Bookings are limited to N weeks in
advance."

**Given** the `[bookings]` section does not exist in sithub.toml
**When** the server starts
**Then** it uses default values without error

### Story 21.7: Maximum Bookings Per Person with Area Overrides

**FRs covered:** FR90

As an operator,
I want to limit how many active bookings a person can hold,
So that shared resources are distributed fairly across users.

**Acceptance Criteria:**

**Given** sithub.toml contains `max_bookings_per_person = 10` under `[bookings]`
**When** a user with 10 active bookings attempts to create another
**Then** the booking is rejected with: "You have reached the maximum of 10 active
bookings."

**Given** the areas YAML sets `max_bookings_per_person: 3` on an item group
**When** a user with 3 active bookings in that item group attempts to book another
item in the same group
**Then** the booking is rejected with: "You have exceeded the maximum of 3 active
bookings for 'Tiefgaragenstellplätze'."

**Given** `max_bookings_per_person` is set at area, item group, and item levels
**When** the system evaluates the limit
**Then** the most specific (deepest) matching limit applies: item overrides item
group, item group overrides area, area overrides global

**Given** `max_bookings_per_person = 0` (or not set) at any level
**When** the system evaluates the limit
**Then** that level imposes no limit; the next higher level's limit applies
(or unlimited if no level sets a limit)

**Given** a booking is rejected due to a limit
**When** the error message is displayed
**Then** it names the exact limit value and the scope where it applies (e.g.,
"You have exceeded the maximum of 2 active bookings for the item
'Tiefgaragenstellplätze, Stellplatz 1'")

## Epic 22 Stories: Bug Fixes, Avatars & Reserved Areas

Mobile UX audit findings are addressed (translations, truncation, menu layout, floor
plan), user avatars are synced and displayed, and areas/items can be reserved for
specific users.
**FRs covered:** FR91, FR92, FR93, FR94, FR95, FR96, FR97, FR98, FR99, FR100

### Story 22.1: Translation and i18n Bug Fixes

**FRs covered:** FR91

As a user,
I want booking limit errors, weekday abbreviations, and availability labels translated
to my selected language,
So that the app feels fully localized.

**Acceptance Criteria:**

**Given** the UI language is German and a booking limit error occurs
**When** the error message is displayed
**Then** the message is in German (e.g., "Sie haben das Maximum von 2 aktiven
Buchungen erreicht"), not English

**Given** the UI language is German
**When** weekday abbreviation dots render on item group tiles
**Then** they show MO, DI, MI, DO, FR, SA, SO (not English MO, TU, WE, TH, FR)

**Given** the UI language is German and a day is free in week mode
**When** the availability label renders
**Then** "n/a" is replaced with a translated label or removed entirely

### Story 22.2: Language Selector and Menu Mobile Layout

**FRs covered:** FR92

As a mobile user,
I want the language and theme buttons to fit the navigation drawer without clipping,
So that I can read and tap each option.

**Acceptance Criteria:**

**Given** I open the hamburger menu on a 390px-wide screen
**When** the language buttons render
**Then** all language names and flags are fully visible (no clipping)

**Given** I open the hamburger menu on a 390px-wide screen
**When** the theme toggle renders
**Then** all three options (Automatisch, Hell, Dunkel) are fully readable

### Story 22.3: Mobile Text Truncation Fixes

**FRs covered:** FR93

As a mobile user,
I want to see full item names and dates without truncation,
So that I can distinguish between similar items.

**Acceptance Criteria:**

**Given** an item name is longer than the card width (e.g., "Tisch 2, Fenster, rechts")
**When** it renders in day mode cards, week mode tile headers, or My Bookings subtitles
**Then** the text wraps to a second line instead of truncating with ellipsis

**Given** I view the booking history page on mobile
**When** the date filter fields render
**Then** the "Von" and "Bis" fields stack vertically with full-width date display

### Story 22.4: Week Mode Mobile Readability

**FRs covered:** FR94

As a mobile user,
I want to see who booked a desk in week mode without overlapping text,
So that I can quickly scan the week grid.

**Acceptance Criteria:**

**Given** a desk is booked in week mode on a mobile screen
**When** the booker's name renders under the day column
**Then** it shows initials (e.g., "AE") or a short abbreviation that fits the column
width without overflow

**Given** I tap on a booked day cell with initials
**When** the interaction completes
**Then** I see the full user name (via tooltip or expanded state)

### Story 22.5: Floor Plan Mobile Improvements

**FRs covered:** FR95

As a mobile user,
I want the floor plan to be readable on my phone and adapt to dark mode,
So that I can use it without squinting or being blinded.

**Acceptance Criteria:**

**Given** I open the floor plan on a 390px-wide screen
**When** it renders
**Then** the zoom level auto-adjusts so the floor plan width fits the viewport

**Given** dark mode is active
**When** the floor plan image renders
**Then** a CSS filter is applied to reduce brightness contrast with the dark UI

**Given** I open the floor plan editor on a narrow viewport
**When** the editor renders
**Then** a banner recommends using a desktop screen for precise positioning

### Story 22.6: Favorites Heart Icon Visibility Fix

**FRs covered:** FR96

As a user,
I want to see the favorite heart icon on all item tiles,
So that I can manage my favorites regardless of other badges shown.

**Acceptance Criteria:**

**Given** an item tile has a warning badge
**When** the tile renders
**Then** the favorite heart icon is still visible and tappable, not hidden behind
or overlapping with the warning badge

### Story 22.7: User Avatars — Backend and Entra ID Sync

**FRs covered:** FR97, FR98

As a user,
I want my profile photo stored and served by SitHub,
So that colleagues can identify me visually across the app.

**Acceptance Criteria:**

**Given** a user logs in via Entra ID
**When** the login completes
**Then** their Microsoft Graph profile photo is downloaded and stored at
`{data_dir}/avatars/{user_id}.png`; if no photo exists, the file is not created

**Given** a local user opens settings
**When** they upload a profile image
**Then** the image is stored at `{data_dir}/avatars/{user_id}.png` with a maximum
file size of 512 KB; the avatars directory is created if missing

**Given** an avatar exists for a user
**When** any authenticated user requests `GET /api/v1/avatars/{user_id}`
**Then** the image is served with appropriate cache headers

**Given** no avatar exists for a user
**When** the avatar endpoint is called
**Then** a 404 is returned and the frontend falls back to initials

### Story 22.8: User Avatars — Frontend Integration

**FRs covered:** FR99

As a user,
I want to see profile photos in the navigation, presence list, and floor plan,
So that I can visually identify colleagues.

**Acceptance Criteria:**

**Given** I am logged in and have an avatar
**When** I see the navigation bar
**Then** my avatar replaces the initials circle in the top-right corner

**Given** I view Today's Presence
**When** the presence list renders
**Then** each user's entry shows their avatar (or initials fallback)

**Given** I open the floor plan with the "Show avatars" checkbox enabled
**When** a desk is booked
**Then** the booker's avatar thumbnail appears on the desk position

### Story 22.9: Reserved Areas and Items — Backend

**FRs covered:** FR100

As an operator,
I want to restrict areas and items to specific users via YAML configuration,
So that shared resources can be reserved for designated teams or individuals.

**Acceptance Criteria:**

**Given** an area has `reserved_for: [anna@sithub.local, tk@system42.io]`
**When** a user not in the list attempts to book any item in that area
**Then** the booking is rejected with a clear error naming the area

**Given** a child item has `reserved_for: [user2@example.com]` but the parent area
does not include `user2@example.com` in its `reserved_for` list
**When** the server starts
**Then** startup fails with a validation error explaining the hierarchical conflict

**Given** `reserved_for` is missing or null on an area/item group/item
**When** a booking is attempted
**Then** no reservation restriction applies at that level

### Story 22.10: Reserved Areas and Items — Frontend

**FRs covered:** FR100

As a user,
I want to see which items I cannot book because they are reserved for others,
So that I do not waste time trying to book restricted items.

**Acceptance Criteria:**

**Given** I am viewing items in an area where some are reserved for other users
**When** the item list renders
**Then** items I cannot book are disabled and visually blurred (similar to the
equipment filter blur pattern)

**Given** a floor plan shows items reserved for other users
**When** the floor plan renders
**Then** reserved items are grayed out or marked with a lock icon

## Epic 23 Stories: UI Bug Fixes

Fix booking tile layout, hidden error messages, and floor plan width on desktop.
**FRs covered:** FR101, FR102, FR103

### Story 23.1: Booking Tile Heart Icon Position

**FRs covered:** FR101

As a user,
I want the favorite heart icon correctly positioned on booking tiles,
so that the tile layout is clean and consistent in both day and week modes.

**Acceptance Criteria:**

**Given** I view items in day booking mode
**When** a tile renders with the heart/favorite icon
**Then** the heart icon is aligned in its designated position on the second line
(after availability, before info/chevron)

**Given** I view items in week booking mode
**When** a tile renders with the heart/favorite icon
**Then** the heart icon is in the same correct position as in day mode

**Given** the tile is rendered
**When** I inspect the layout
**Then** a test exists that detects the wrong layout before the fix and passes
after the fix

### Story 23.2: Booking Limit Error Modal

**FRs covered:** FR102

As a user,
I want booking limit errors shown in a modal overlay,
so that I cannot miss critical error messages when booking by week.

**Acceptance Criteria:**

**Given** I am booking items for multiple days in week mode
**When** the booking exceeds my booking limit
**Then** the error is displayed in a modal dialog overlaying all other content

**Given** the booking limit error modal is displayed
**When** I read the error
**Then** I must actively press a close/dismiss button to continue using the app

**Given** the booking limit error modal is displayed
**When** I dismiss it
**Then** I return to the booking view with my previous selections intact

### Story 23.3: Floor Plan Full-Width Desktop Layout

**FRs covered:** FR103

As a user,
I want the floor plan to use the full available width on desktop,
so that I can see floor plan details without unnecessary whitespace.

**Acceptance Criteria:**

**Given** I am viewing a floor plan on a desktop viewport (>= 960px)
**When** the floor plan renders
**Then** the floor plan container uses the full available width of the content area

**Given** I am viewing a floor plan on a mobile viewport
**When** the floor plan renders
**Then** the existing mobile layout behavior is unchanged

## Epic 24 Stories: Booking Warnings & Profile Consolidation

Users are prompted with a confirmation dialog before booking items that have warnings,
with a "don't show again" option per item. In week mode, warnings for multiple items are
shown sequentially. The Settings menu is removed and consolidated into the Profile menu.
**FRs covered:** FR104, FR105, FR106

### Story 24.1: Warning Confirmation Dialog (Day Mode)

**FRs covered:** FR104

As a user,
I want to see a confirmation dialog with the item's warning before booking,
so that I am aware of restrictions and can decide whether to proceed or choose a
different item.

**Acceptance Criteria:**

**Given** I click BOOK on an item that has a warning in day booking mode
**When** the warning dialog appears
**Then** it displays the item name (truncated with ellipsis if longer than the dialog
width), the warning text, a CONFIRM button, and a CANCEL button

**Given** the warning dialog is displayed
**When** I click CONFIRM
**Then** the booking proceeds as normal

**Given** the warning dialog is displayed
**When** I click CANCEL
**Then** the booking is aborted and I remain on the booking view with no booking created

**Given** the warning dialog is displayed with a "Don't show again" checkbox
**When** I check "Don't show again" and click CONFIRM
**Then** the booking proceeds and the suppression is stored in localStorage keyed by
item ID

**Given** I have previously checked "Don't show again" for an item
**When** I book that same item again
**Then** the warning dialog is skipped and the booking proceeds immediately

**Given** an item has no warning configured
**When** I click BOOK
**Then** no warning dialog is shown and the booking proceeds as before

### Story 24.2: Sequential Warning Dialogs (Week Mode)

**FRs covered:** FR105

As a user,
I want warnings for multiple items shown one after another when booking in week mode,
so that I can review each item's restrictions before confirming the full week booking.

**Acceptance Criteria:**

**Given** I am in week booking mode and have selected days on multiple items that have
warnings
**When** I click "Confirm My Booking"
**Then** the warning dialogs are shown sequentially, one per item with a warning, each
identifying the item by name

**Given** a sequential warning dialog is displayed for item A
**When** I click CONFIRM
**Then** the next item's warning dialog is shown (or booking proceeds if no more warnings
remain)

**Given** a sequential warning dialog is displayed for item B
**When** I click CANCEL
**Then** the entire week booking is aborted and no bookings are created for any item

**Given** I have previously suppressed warnings for some items via "Don't show again"
**When** I book a week that includes those items
**Then** the suppressed items' warning dialogs are skipped; only unsuppressed warnings
are shown

**Given** all items in my week booking have their warnings suppressed
**When** I click "Confirm My Booking"
**Then** the booking proceeds immediately with no warning dialogs

### Story 24.3: Profile and Settings Consolidation

**FRs covered:** FR106

As a user,
I want all settings accessible from a single Profile menu,
so that I don't have to choose between two overlapping menus to find what I need.

**Acceptance Criteria:**

**Given** I am logged in and viewing the app
**When** I look at the navigation
**Then** there is no separate "Settings" menu option; only the Profile menu
(avatar/initials) exists

**Given** I open the Profile menu
**When** I view the menu contents
**Then** all settings are present: theme selector, language selector, show weekends
toggle, and change password option

**Given** I open the Profile menu on mobile
**When** I view the menu contents
**Then** the same settings are available with the current profile layout styling

**Given** I previously accessed a setting via the old Settings menu
**When** I look for it after the consolidation
**Then** the setting is accessible from the Profile menu with no functionality lost

## Epic 25 Stories: UX/UI Improvements — Floor Plan Editor, Booking & Avatar

The floor plan editor is overhauled for a streamlined editing experience: sidebar replaced
with toolbar dropdowns, canvas enlarged, auto-save replaces manual save, undo removed, and
zoom controls redesigned. Subarea drill-down images are enlarged for usability. Entra ID
avatar sync is made asynchronous with login feedback, and the Profile Photo menu is hidden
for Entra ID users.
**FRs covered:** FR107, FR108, FR109, FR110, FR111, FR112, FR113, FR114, FR115, FR116, FR117

### Story 25.1: Editor Layout — Sidebar to Toolbar Dropdowns

**FRs covered:** FR107, FR108, FR109

As an admin,
I want the floor plan editor to use the full page width with controls in the toolbar,
so that I have maximum canvas space for positioning items on the floor plan.

**Acceptance Criteria:**

**Given** I open the floor plan editor as an admin
**When** the editor loads
**Then** there is no left-hand sidebar listing subareas and items; the canvas card uses
the full available width

**Given** the editor is loaded
**When** I look at the toolbar row
**Then** I see a subarea dropdown that lists all subareas for the selected floor plan

**Given** the editor is loaded
**When** I select a subarea from the toolbar dropdown
**Then** the editor switches to that subarea, identical to the old sidebar click behavior

**Given** the editor is loaded
**When** I look at the toolbar row
**Then** I see an items dropdown that lists all items for the current subarea, each
indicating whether it is positioned or unpositioned (e.g., via icon or chip)

**Given** I select an unpositioned item from the items dropdown
**When** the selection is made
**Then** the editor enters draw mode for that item, identical to the old sidebar behavior

**Given** I select a positioned item from the items dropdown
**When** the selection is made
**Then** the editor selects that item's rectangle on the canvas, identical to the old
sidebar behavior

**Given** I have a positioned item selected via the items dropdown
**When** I look for a way to delete it
**Then** I see a delete action accessible from the items dropdown or toolbar that removes
the item's position from the floor plan

### Story 25.2: Canvas Height & Zoom Controls

**FRs covered:** FR110, FR113

As an admin,
I want a taller canvas and compact zoom controls,
so that I can see and edit the floor plan image with less scrolling and a cleaner toolbar.

**Acceptance Criteria:**

**Given** I open the floor plan editor
**When** the editor loads
**Then** the canvas area uses approximately double the vertical space compared to the
previous layout

**Given** I look at the zoom controls in the editor toolbar
**When** I inspect their layout
**Then** the zoom percentage label appears between the minus and plus buttons, not next
to them

**Given** I click the plus or minus zoom buttons
**When** the zoom level changes
**Then** the percentage label between the buttons updates to reflect the current zoom
factor

### Story 25.3: Auto-Save & Remove Undo

**FRs covered:** FR111, FR112

As an admin,
I want the floor plan editor to save automatically and not distract me with undo and
manual save controls,
so that I can focus on positioning items without worrying about losing changes.

**Acceptance Criteria:**

**Given** I draw a new rectangle on the floor plan
**When** I release the mouse button (pointerup)
**Then** the changes are saved automatically without clicking a save button

**Given** I move an existing rectangle on the floor plan
**When** I release the mouse button (pointerup)
**Then** the changes are saved automatically

**Given** I resize an existing rectangle on the floor plan
**When** I release the mouse button (pointerup)
**Then** the changes are saved automatically

**Given** no unsaved changes exist
**When** a pointerup event fires
**Then** no save request is triggered

**Given** the editor is loaded
**When** I look at the toolbar
**Then** there is no manual Save button

**Given** an auto-save is in progress
**When** I look at the toolbar
**Then** I see a brief saving/saved indicator reflecting the auto-save state

**Given** the editor is loaded
**When** I look at the toolbar
**Then** there is no Undo button and the undo keyboard shortcut has no effect

### Story 25.4: Enlarged Subarea Images on Drill-Down

**FRs covered:** FR114

As a user,
I want subarea floor plan images to be displayed enlarged when I drill into them,
so that I can clearly see the layout and available items without zooming.

**Acceptance Criteria:**

**Given** I am viewing the floor plan booking view for an area
**When** I click on a subarea rectangle to drill into it
**Then** the subarea floor plan image renders at an enlarged size that fills the available
viewport width

**Given** I have drilled into a subarea
**When** the subarea floor plan is displayed at default zoom level
**Then** no horizontal scrollbars appear; the image fits within the viewport width

**Given** I manually zoom in beyond the default level
**When** the image exceeds the viewport width
**Then** scrollbars appear as expected to allow navigation

### Story 25.5: Hide Profile Photo for Entra ID Users

**FRs covered:** FR115

As an Entra ID user,
I want the Profile Photo menu option to be hidden,
so that I am not confused by an option that would have no effect since my avatar is
synced from Entra ID.

**Acceptance Criteria:**

**Given** I am logged in via Entra ID on desktop
**When** I open the user menu (Profile menu)
**Then** the "Profile Photo" menu item is not visible

**Given** I am logged in via Entra ID on mobile
**When** I open the navigation drawer
**Then** the "Profile Photo" menu item is not visible

**Given** I am logged in as a local (internal) user on desktop
**When** I open the user menu
**Then** the "Profile Photo" menu item is visible and functional

**Given** I am logged in as a local (internal) user on mobile
**When** I open the navigation drawer
**Then** the "Profile Photo" menu item is visible and functional

### Story 25.6: Async Avatar Sync & Login Spinner

**FRs covered:** FR116, FR117

As an Entra ID user,
I want the login to complete quickly with visual feedback,
so that I am not left waiting on a slow avatar sync with no indication of progress.

**Acceptance Criteria:**

**Given** I click "Sign in with Entra ID" on the login page
**When** the click is registered
**Then** the button shows a loading spinner and is disabled, preventing double-clicks

**Given** the Entra ID OAuth callback is processed by the backend
**When** the avatar sync would normally run
**Then** the avatar sync runs asynchronously in a goroutine; the OAuth callback returns
immediately and redirects the user without waiting for the avatar download

**Given** the async avatar sync completes successfully in the background
**When** I navigate to a page showing my avatar
**Then** my Entra ID profile photo is displayed

**Given** the async avatar sync fails (e.g., no photo in Entra ID, network error)
**When** I navigate to a page showing my avatar
**Then** the fallback initials avatar is displayed; no error is shown to the user

## Epic 26 Stories: Floor Plan Editor — Area Drawing Fixes

The floor plan editor's Areas tab workflow is broken after the sidebar-to-toolbar refactor.
Subarea selection forces a tab switch, drawing areas is impossible, and unrelated subareas
can be accidentally modified. This epic fixes all four interaction bugs.
**FRs covered:** FR118, FR119, FR120, FR121

### Story 26.1: Subarea Selection Respects Active Tab

**FRs covered:** FR118

As an admin,
I want selecting a subarea from the dropdown to stay on the current tab,
so that I can position area rectangles without being forced into Items mode.

**Acceptance Criteria:**

**Given** I am on the Areas tab with an area-level floor plan loaded
**When** I select "Open Space" from the subarea dropdown
**Then** the toggle stays on "Areas" and does not switch to "Items"

**Given** I am on the Items tab
**When** I select a subarea from the dropdown
**Then** the toggle stays on "Items" (existing behavior preserved)

### Story 26.2: Hide Items Dropdown on Areas Tab

**FRs covered:** FR119

As an admin,
I want the items dropdown to be hidden when I am on the Areas tab,
so that I am not confused by irrelevant controls while positioning subareas.

**Acceptance Criteria:**

**Given** I am on the Areas tab
**When** I look at the toolbar
**Then** the "Objekte" (Items) dropdown is not visible

**Given** I switch to the Items tab
**When** I look at the toolbar
**Then** the "Objekte" (Items) dropdown appears

### Story 26.3: Enable Draw Mode for Subareas on Areas Tab

**FRs covered:** FR120

As an admin,
I want to draw a rectangle for an unpositioned subarea when I select it on the Areas tab,
so that I can position subareas on the floor plan.

**Acceptance Criteria:**

**Given** I am on the Areas tab and select an unpositioned subarea from the dropdown
**When** the selection is made
**Then** the editor enters draw mode (crosshair cursor) for that subarea

**Given** I am on the Areas tab and select a positioned subarea from the dropdown
**When** the selection is made
**Then** the editor selects that subarea's rectangle on the canvas

### Story 26.4: Lock Other Rectangles When Subarea Is Selected

**FRs covered:** FR121

As an admin,
I want only the selected subarea to be editable on the canvas,
so that I cannot accidentally move or delete other subareas.

**Acceptance Criteria:**

**Given** I have selected "Open Space" for editing on the Areas tab
**When** I try to click, move, or delete another subarea's rectangle (e.g., "Cube 1")
**Then** the other rectangle does not respond to interaction

**Given** I have a subarea selected
**When** I look at the other subarea rectangles on the canvas
**Then** they appear visually distinct (e.g., dimmed or dashed) to indicate they are locked

## Epic 27 Stories: Avatar Sync Fix & Reserved Item Visibility

Avatar sync from Entra ID fails on certain image formats. Reserved areas and items hide
occupancy information too aggressively — users cannot see who sits where. This epic fixes
avatar decoding, corrects reserved area display on floor plans, and makes reserved items
visible (but non-bookable) in the regular booking view.
**FRs covered:** FR122, FR123, FR124, FR125

### Story 27.1: Fix Avatar Image Sync from Entra ID

**FRs covered:** FR122

As an Entra ID user,
I want my profile photo to sync correctly regardless of image format,
so that my avatar displays properly in SitHub.

**Acceptance Criteria:**

**Given** Microsoft Graph returns a JPEG, PNG, or other common image format
**When** the avatar sync runs
**Then** the image is decoded and re-encoded as PNG successfully

**Given** the avatar download or decoding fails
**When** the error is logged
**Then** the log message includes the user ID, HTTP status, content-type, and body size
for diagnosis

**Given** a user has no profile photo or the sync fails
**When** their avatar is displayed
**Then** the fallback initials avatar is shown (no broken image)

**Given** an admin reads the FAQ
**When** they look for avatar troubleshooting
**Then** the README.md explains common causes and fixes for broken avatar sync

### Story 27.2: Fix Reserved Area Display on Floor Plan

**FRs covered:** FR123, FR124

As a user,
I want reserved areas on the floor plan to show correct availability and allow drill-down,
so that I can see who is in a reserved room even though I cannot book there.

**Acceptance Criteria:**

**Given** I view the floor plan and "People & Finance" is a reserved area with 3 of 4
desks free
**When** I look at the area overlay
**Then** it shows "3/4 free" with a green indicator (not "0/4" red)

**Given** a desk in a reserved area is booked by someone
**When** I look at the floor plan
**Then** their avatar is displayed on that desk

**Given** I click on a reserved area on the floor plan
**When** I drill down into it
**Then** I can see individual desk availability and who has booked

**Given** I see a free desk in a reserved area after drill-down
**When** I look at that desk
**Then** it is blurred/dimmed with a "reserved" message and I cannot book it

### Story 27.3: Show Occupancy on Reserved Items in Regular Booking View

**FRs covered:** FR125

As a user,
I want to see who is sitting where on reserved items in the list view,
so that I know room occupancy even though I cannot book there myself.

**Acceptance Criteria:**

**Given** I view items in a reserved area (e.g., Finance & People Area)
**When** I look at the item list
**Then** I can see free/busy status and booker names — the opaque veil is removed

**Given** I see a reserved item that is free
**When** I look at it
**Then** it shows a badge with a lock icon and "reserved" text

**Given** I try to book a reserved item
**When** I interact with it
**Then** the booking action is blocked (no book button or it is disabled)

## Epic 28 Stories: Date Selector Fix & Floor Plan Booker Names

Fix the date picker jumping to today after booking, and add booker identification on floor
plans via name tooltips on avatars and initials display when avatars are disabled.
**FRs covered:** FR126, FR127, FR128

### Story 28.1: Preserve Selected Date After Booking

**FRs covered:** FR126

As a user,
I want the date selector to stay on my selected date after completing a booking,
so that I can continue browsing availability for the same date without being sent back to
today.

**Acceptance Criteria:**

**Given** I am on the items page for an item group (e.g. Cube) with a future date selected
(e.g. 30 April 2026)
**When** I book an item and the booking confirmation completes
**Then** the date picker still shows 30 April 2026
**And** the displayed items reflect booking status for 30 April 2026 (including my new
booking)

**Given** I am in week booking mode with a future week selected
**When** I confirm bookings for that week
**Then** the week selector stays on the same week
**And** my new bookings are shown as booked in the week view

### Story 28.2: Floor Plan Booker Name Tooltips and Initials

**FRs covered:** FR127, FR128

As a user,
I want to see who has booked a desk on the floor plan by hovering over avatars or seeing
initials,
so that I can identify people without having to drill down into each item.

**Acceptance Criteria:**

**Given** "Show Avatar" is enabled and I view the interactive floor plan
**When** I hover over a booked item's avatar (or tap on mobile)
**Then** a tooltip displays the full display name of the booker (e.g.
"Alexander Seidemann-Klamant")

**Given** "Show Avatar" is disabled and I view the interactive floor plan
**When** I look at booked items
**Then** each booked item shows the booker's initials (e.g. "AS" for
"Alexander Seidemann-Klamant") instead of an avatar image

**Given** "Show Avatar" is disabled and I view the interactive floor plan
**When** I hover over (or tap on mobile) a booked item showing initials
**Then** a tooltip displays the full display name of the booker

**Given** a user's display name has multiple parts (e.g. "Alexander Seidemann-Klamant")
**When** their initials are derived
**Then** they use the first letter of each space-separated name part (e.g. "AS")

## Epic 29 Stories: Desktop Weekly Table View

Add a desktop-only weekly table view on the selected area's item-groups page that restores the
Confluence-style full-week overview while keeping SitHub's stronger booking and permission
model. The view must show all subareas and desks in one long matrix for the currently selected
week, respect user weekend settings, and allow in-place booking and cancellation without
leaving the table.
**FRs covered:** FR129, FR130, FR131, FR132, FR133, FR134, FR135, FR136

### Story 29.1: Weekly Desk Matrix API

**FRs covered:** FR131

As a user,
I want the frontend to load one weekly desk matrix for the selected area,
so that the table can render quickly and consistently without dozens of follow-up requests.

**Acceptance Criteria:**

**Given** the frontend requests weekly table data for an area and week
**When** the backend responds
**Then** the payload contains all subareas and desks for that area in the configured SitHub
order
**And** each desk contains one cell per requested visible day

**Given** the frontend requests 5 visible days because the current user has weekends disabled
**When** the backend responds
**Then** the payload contains Monday through Friday only

**Given** a matrix cell is occupied
**When** the backend responds
**Then** the cell includes the booker's display name and user ID
**And** enough booking metadata is present for permitted cancellation actions

**Given** a matrix cell is free but reserved for other users
**When** the backend responds
**Then** the payload clearly distinguishes between `bookable` free cells and `locked`
reserved cells

### Story 29.2: Table View Entry & View Persistence

**FRs covered:** FR129, FR130

As a desktop user,
I want a Table view option on the area item-groups page and for SitHub to remember my last
desktop view,
so that I can return directly to the weekly matrix without extra setup.

**Acceptance Criteria:**

**Given** I open `/areas/:areaId/item-groups` on a desktop viewport
**When** the page renders
**Then** I see a `Table view` action alongside the existing area actions in the same control
area as `Floor plan`

**Given** I open `/areas/:areaId/item-groups` on a mobile viewport
**When** I look at the `Table view` action
**Then** it is disabled
**And** hovering or long-pressing it explains that the table view is available on desktop only

**Given** I switch from the default card view to `Table view` on desktop
**When** I leave the page and come back later to the same area
**Then** SitHub restores `Table view` as the active desktop view for that area context

### Story 29.3: Weekly Matrix Layout, Sticky Orientation & Room State

**FRs covered:** FR132, FR133

As a user,
I want a dense but readable weekly matrix with collapsible room sections,
so that I can scan the whole floor quickly without losing orientation.

**Acceptance Criteria:**

**Given** I open `Table view`
**When** the matrix renders
**Then** it shows all subareas of the selected area in one long table
**And** the subareas and desks appear in the exact configured SitHub order

**Given** I view the matrix
**When** I scroll vertically
**Then** the weekday header remains sticky
**And** the left desk-name column remains sticky

**Given** I open the matrix for the first time
**When** the room sections render
**Then** all rooms are expanded by default

**Given** I collapse or expand a room using its dedicated chevron
**When** I reopen the table later or switch to another week
**Then** the previous collapsed state is restored from local storage

**Given** a room is collapsed
**When** I look at its header
**Then** I see compact occupied counts for each visible day of the current week

**Given** the selected week contains past days
**When** the matrix renders
**Then** those past-day columns stay visible
**And** their cells are visually muted and non-interactive

### Story 29.4: Cell States, Occupant Identity & Reserved Permissions

**FRs covered:** FR134

As a user,
I want every cell to communicate availability and permissions immediately,
so that I can understand the board without extra clicks.

**Acceptance Criteria:**

**Given** I look at a free cell I am allowed to book
**When** the matrix renders
**Then** the cell looks like a normal bookable free cell with minimal text

**Given** I look at a free cell in a reserved room or desk that I am not allowed to book
**When** the matrix renders
**Then** the cell shows a lock indicator
**And** it is not clickable

**Given** I look at an occupied cell
**When** the matrix renders
**Then** it shows the occupant using avatar plus initials in the compact cell layout
**And** hovering reveals the full person name

**Given** I hover a desk label
**When** that desk has equipment configured
**Then** I see the equipment hints on hover

**Given** I am not an admin and the occupied cell is not my own booking
**When** I interact with it
**Then** it is non-clickable and exposes no extra popup content

### Story 29.5: Direct Booking from Free Cells

**FRs covered:** FR135

As a user,
I want to book directly from a free table cell,
so that I can act from the weekly overview without leaving the matrix.

**Acceptance Criteria:**

**Given** I click a bookable free cell
**When** the booking UI opens
**Then** it appears as a lightweight desktop popover anchored to that cell

**Given** the booking popover is open
**When** I inspect its controls
**Then** `Book for myself` is selected by default
**And** I can switch to `Book for colleague`
**And** the colleague picker appears only when `Book for colleague` is selected

**Given** I previously booked for a colleague from the table
**When** I switch the booking popover to `Book for colleague` again
**Then** the last selected colleague is preselected

**Given** the booking popover is open
**When** I inspect its contents
**Then** the note field is visible immediately
**And** any booking warning is shown inline in the same popover instead of a second dialog

**Given** I enter a note and confirm a booking successfully
**When** the request completes
**Then** the note is stored with the created booking

**Given** I confirm the booking successfully
**When** the request completes
**Then** the popover closes
**And** the cell updates immediately in place without navigation away from the matrix

**Given** another user books the same desk before my confirmation succeeds
**When** the booking request returns a conflict
**Then** the popover stays open
**And** it shows an inline error explaining that the desk is no longer available

### Story 29.6: In-Place Cancellation for Own & Admin-Allowed Bookings

**FRs covered:** FR136

As a user,
I want eligible occupied cells to support lightweight cancellation in place,
so that I can correct bookings without leaving the weekly overview.

**Acceptance Criteria:**

**Given** I click a cell containing my own booking
**When** the cancellation UI opens
**Then** it uses the same lightweight anchored popover pattern as booking
**And** it shows only person, desk, and date before the `Cancel booking` action

**Given** I am an admin and click a cell containing someone else's booking
**When** the cancellation UI opens
**Then** I can cancel that booking from the same anchored popover pattern

**Given** I am not an admin and the occupied cell belongs to another user
**When** I interact with that cell
**Then** no popup opens
**And** the cell remains read-only

**Given** I confirm a cancellation successfully
**When** the request completes
**Then** the popover closes
**And** the cell updates immediately in place to its new free or locked state

## Epic 30 Stories: Operator Validation, Editor Zoom Height & Optional Drill-Down

Detect duplicate items in the areas configuration at startup and refuse to start, make the
floor plan editor canvas grow vertically when zooming in so users only need to scroll
horizontally, and let users on large screens opt out of the room drill-down so they can
book directly from the floor plan.
**FRs covered:** FR137, FR138, FR139

### Story 30.1: Reject Duplicate Items in Areas Config at Startup

**FRs covered:** FR137

As an operator,
I want the server to refuse to start when the areas configuration contains duplicate item
identifiers,
so that misconfigurations are caught immediately instead of producing inconsistent booking
state at runtime.

**Acceptance Criteria:**

**Given** the areas configuration contains the same item identifier (e.g. `desk29`) more
than once across any subareas or rooms
**When** the server starts
**Then** the server logs an error that names the duplicated identifier and the locations
where it appears
**And** the server exits with a non-zero status before opening its listening socket

**Given** the areas configuration contains no duplicate item identifiers
**When** the server starts
**Then** the server starts successfully and serves requests as before

**Given** duplicate detection is implemented
**When** the unit tests run
**Then** a test fixture covers at least one duplicate scenario and asserts the startup
failure path
**And** a separate fixture covers a clean configuration and asserts successful startup

### Story 30.2: Floor Plan Editor Canvas Grows Vertically on Zoom

**FRs covered:** FR138

As an admin editing a floor plan,
I want the image container to grow vertically when I zoom in,
so that I only have to scroll horizontally to inspect the full plan instead of scrolling
in both directions.

**Acceptance Criteria:**

**Given** I open the floor plan editor for an area
**When** the editor first loads
**Then** the height of the box around the floor plan image is derived from the image's
intrinsic height as today

**Given** the floor plan editor is open at the default zoom level
**When** I zoom in (via toolbar controls or keyboard/scroll shortcut)
**Then** the height of the surrounding box grows to match the scaled image height
**And** vertical scrolling inside the editor is no longer required to see the bottom of
the image

**Given** I have zoomed in
**When** the image is wider than the viewport at the current zoom level
**Then** I can scroll horizontally to reach off-screen content
**And** the layout of editor controls outside the image area is unaffected

**Given** I zoom back out to the default level
**When** the editor re-renders
**Then** the surrounding box returns to a height consistent with the displayed image

### Story 30.3: Optional Area Drill-Down Toggle on Floor Plan Booking

**FRs covered:** FR139

As a user on a large screen,
I want to book items directly from the floor plan without drilling into the detailed
room/area view,
so that I can complete bookings faster from the overview I already see.

**Acceptance Criteria:**

**Given** I view a floor plan in the booking view
**When** the page renders
**Then** an "Area drill-down" toggle (checkbox) is visible beneath the room plan

**Given** the viewport is a small screen (mobile breakpoint) and I have not changed the
toggle yet
**When** the floor plan loads
**Then** the "Area drill-down" toggle is enabled by default
**And** clicking on an area or room on the floor plan opens the detailed drill-down view
as today

**Given** the viewport is a large screen (desktop breakpoint) and I have not changed the
toggle yet
**When** the floor plan loads
**Then** the "Area drill-down" toggle is disabled by default
**And** clicking on an item on the floor plan starts the booking flow directly without
loading the detailed room/area view

**Given** I change the "Area drill-down" toggle
**When** the change is applied
**Then** the new value is saved to local storage on the current device
**And** subsequent floor plan booking sessions on the same device use the saved value
regardless of viewport size

**Given** I open the floor plan booking view on a different device
**When** the page renders
**Then** the toggle uses that device's own default and stored value
**And** the choice on this device does not affect any other device

## Epic 31 Stories: Live Updates, Favorites Rework & Areas Config Hint

Push booking and cancellation events to all open clients so users see availability change
in real time without manual reloads, rework favorites into a dedicated virtual "Favorites"
room with consistent heart-icon affordances across tile, table, and floor plan views, and
clarify in the example configuration that the areas config file must live inside
`data_dir`.
**FRs covered:** FR140, FR141, FR142

### Story 31.1: Live Updates for Bookings and Cancellations

**FRs covered:** FR140

As a user browsing room plans, area tiles, or the weekly table view,
I want bookings and cancellations made by other users to appear without manually
reloading,
so that I do not waste time deciding on a desk that has already been taken.

**Acceptance Criteria:**

**Given** I am viewing the area tile, weekly table, or floor plan view for a given date
**When** another user creates a booking that affects an item in my current view
**Then** the corresponding tile, cell, or floor plan marker updates to its new busy state
within a few seconds
**And** I do not need to refresh the page or change any filter for the change to appear

**Given** I am viewing any of the same views
**When** another user cancels a booking that affects an item in my current view
**Then** the corresponding tile, cell, or floor plan marker updates to its new free state
within a few seconds
**And** the change is reflected without page reload

**Given** I have just been viewing live updates and the network connection drops
**When** the connection is restored
**Then** the client reconciles state with the server and reflects any bookings or
cancellations that occurred during the outage
**And** the existing connection-lost messaging from Story 18.7 still surfaces while
disconnected

**Given** I open the app in multiple tabs
**When** an event arrives
**Then** all open tabs reflect the change consistently

### Story 31.2: Favorites Rework as Virtual Room

**FRs covered:** FR141

As a user with favorite desks across multiple areas,
I want my favorites grouped into a dedicated "Favorites" room with clear visual markers
in every view,
so that I can find and book my preferred desks quickly without scanning unrelated areas.

**Acceptance Criteria:**

**Given** I have at least one bookable item marked as a favorite
**When** I open the area/room overview
**Then** a tile labeled "Favorites" appears as the first tile
**And** the tile behaves like any other room tile (drill-down to its items, free/busy
indicators, identical interaction model)

**Given** I have no items marked as favorites
**When** I open the area/room overview
**Then** the "Favorites" tile is not shown

**Given** I am on a screen that previously allowed adding an area or room to favorites
**When** the page renders
**Then** no control to favorite an area or room is available
**And** only bookable items (desks) can be added to favorites

**Given** I am viewing the weekly table view and one or more items are favorites
**When** the table renders
**Then** each favorite item row displays a heart icon
**And** clicking the heart icon removes that item from my favorites
**And** no sorting or filtering by favorites is offered in the table view

**Given** I am viewing a floor plan and one or more items shown on the plan are favorites
**When** an item is in the free state
**Then** a heart icon is rendered with its center positioned exactly at the bottom-right
corner of the item marker (matching the reference image
`epic-31-favorite-heart.png`)
**And** clicking the heart icon removes the item from my favorites

**Given** an item shown on the floor plan is busy
**When** the floor plan renders
**Then** no heart icon is shown for that item, regardless of favorite status

**Given** I remove an item from favorites via any of the heart icons (tile, table, floor
plan)
**When** the change is applied
**Then** the item is removed across all views consistently and the "Favorites" tile
disappears once no favorites remain

### Story 31.3: Areas Config Location Hint in Example TOML

**FRs covered:** FR142

As an operator setting up SitHub,
I want the example configuration to state that the areas config file must live inside
`data_dir`,
so that I do not waste time placing the file in an unsupported location.

**Acceptance Criteria:**

**Given** I read `sithub.example.toml`
**When** I look at the setting that points to the areas configuration file
**Then** its inline comment includes the sentence "Must be inside data_dir."
**And** the wording matches the existing hint used for the floor plans directory setting

**Given** the example file is updated
**When** the TOML linter or formatter runs
**Then** the file remains valid TOML and follows the existing comment style described in
`.claude/rules/toml.md`

## Epic 32 Stories: Booker Avatars on Item Tiles & Stable Colleague Booking Layout

Show booker avatars on booked item tiles in day and week modes so users can recognize who
has reserved an item at a glance, matching the floor plan presentation. Stabilize the
"Book for a colleague" layout so toggling it does not push subsequent containers down a
line on viewports with enough horizontal space.
**FRs covered:** FR143, FR144, FR145

### Story 32.1: Booker Avatar on Day-Mode Item Tiles

**FRs covered:** FR143

As a user browsing the item-groups view in day mode,
I want each booked item tile to show the booker's avatar with their name on hover or tap,
so that I can recognize who has reserved a desk at a glance without opening details.

**Acceptance Criteria:**

**Given** I am viewing an item-groups page in day mode
**When** a tile represents a booked item
**Then** the tile displays a circular avatar of the booker
**And** the avatar treatment matches the one used on the floor plan (size, shape, border)

**Given** I hover over the avatar on a desktop viewport
**When** the tooltip appears
**Then** it shows the booker's full display name

**Given** I tap the avatar on a mobile or touch viewport
**When** the action is recognized
**Then** the booker's full display name is shown without navigating away from the view

**Given** the booker has no synced or uploaded avatar image
**When** the tile renders
**Then** a circular initials avatar is shown using the same derivation rule already used
on the floor plan (per FR128)
**And** the same hover/tap behavior surfaces the full display name

**Given** a tile represents an available (not booked) item
**When** the tile renders
**Then** no booker avatar is shown

### Story 32.2: Booker Avatar on Week-Mode Item Tiles

**FRs covered:** FR144

As a user browsing the item-groups view in week mode,
I want each booked weekday cell on a tile to show the booker's avatar with their name on
hover or tap,
so that I can scan a whole week and immediately see who has booked which day.

**Acceptance Criteria:**

**Given** I switch the item-groups view to week mode
**When** a weekday cell on a tile represents a booking
**Then** the cell displays a circular avatar of that day's booker
**And** the avatar treatment is visually consistent with day mode and the floor plan

**Given** I hover over an avatar in a weekday cell on a desktop viewport
**When** the tooltip appears
**Then** it shows the booker's full display name

**Given** I tap an avatar in a weekday cell on a touch viewport
**When** the action is recognized
**Then** the booker's full display name is shown without navigating away from the view

**Given** the booker has no synced or uploaded avatar image
**When** the weekday cell renders
**Then** an initials-based circular avatar is shown using the same derivation rule used
elsewhere (per FR128)

**Given** a weekday cell is free
**When** the tile renders
**Then** the existing free-state visuals (checkbox, "frei" label) are shown unchanged
**And** no avatar is rendered for that cell

### Story 32.3: Stable Inline Layout for "Book for a Colleague" Dropdown

**FRs covered:** FR145

As a user with permission to book for a colleague,
I want the colleague-selection dropdown to appear next to the radio button without
pushing the rest of the form down,
so that switching the booking target does not cause the UI to feel nervous.

**Acceptance Criteria:**

**Given** I am on the item-groups view with a wide desktop viewport
**When** I select "Book for a colleague"
**Then** the colleague-selection dropdown appears inline, to the right of the radio
group, on the same line
**And** the vertical position of the equipment filter and the tile grid below does not
change compared to the "Book for myself" state

**Given** I toggle back to "Book for myself"
**When** the dropdown disappears
**Then** the row containing the radio group keeps its height and no other element shifts
vertically

**Given** my viewport is narrow enough that the inline dropdown does not fit beside the
radio group
**When** I select "Book for a colleague"
**Then** the dropdown wraps to its own line below the radio group
**And** the transition does not produce additional intermediate layout jumps once the
state has settled

**Given** the layout adjusts between inline and wrapped modes when the viewport is
resized
**When** the breakpoint is crossed
**Then** the form remains usable and no interactive element is hidden behind another

**Given** I open the colleague-selection dropdown (inline or wrapped) and choose a
colleague
**When** the selection is applied
**Then** the existing booking flow continues unchanged
**And** the chosen colleague is used as the booker for the next booking action

## Epic 33 Stories: Equipment Filter Fixes, Compact Booking Controls & Login Page Rebranding

Fix the equipment-filter reset bug and the non-functional table-view filter, widen the
table-view cancel popover, replace the booking-type radio group with a single checkbox
plus always-rendered colleague dropdown, compact the booking controls to a single row,
and rebrand the login page and application header with the new SitHub logo and an
Entra-ID-primary sign-in flow.
**FRs covered:** FR146, FR147, FR148, FR149, FR150, FR151, FR152

### Story 33.1: Equipment Filter Resets on Item-Groups View

**FRs covered:** FR146

As a user filtering items by equipment,
I want clearing the filter input to re-enable every tile,
so that I can recover from a typo or a non-matching search without reloading the page.

**Acceptance Criteria:**

**Given** I am on an item-groups page in day mode and I type a keyword that no item's
equipment matches
**When** I press the input's built-in clear "X" icon
**Then** every tile that was blurred / disabled by the filter returns to its normal
interactive state without a page reload

**Given** I am on the same page with a non-matching filter applied
**When** I backspace the input until it is empty
**Then** every tile returns to its normal interactive state immediately

**Given** I am on the same page in week mode
**When** I clear the filter via either method above
**Then** every week tile returns to its normal interactive state

**Given** the filter is empty
**When** I view the page
**Then** no tile is blurred or disabled by the equipment-filter logic, regardless of
prior filter history

### Story 33.2: Equipment Filter on Weekly Table View

**FRs covered:** FR147

As a user on the weekly desktop table view,
I want the equipment filter at the top of the table to actually filter rows,
so that I can find desks by equipment without leaving the table.

**Acceptance Criteria:**

**Given** I am on the weekly desktop table view and I type a keyword in the equipment
filter input
**When** the filter is applied
**Then** every row whose item does not match the parsed filter is visually dimmed and
its cells become non-interactive
**And** rows whose items match remain at full opacity and remain interactive

**Given** the filter is applied
**When** I clear the input (via the clear icon or by backspacing)
**Then** every row returns to its normal state

**Given** the filter matches every item in the table
**When** the filter is applied
**Then** no row is dimmed

**Given** the filter applies to the table view
**When** I switch back and forth between the table and card views
**Then** the filter input value is preserved within the same session per the existing
saved-filters behavior

### Story 33.3: Widen Weekly-Table Booking-Cancel Popover

**FRs covered:** FR148

As a user cancelling my own booking from the weekly table view,
I want the cancel-confirmation popover to be tall and wide enough to show the cancel
button fully,
so that I can complete the action without scrolling inside the popover.

**Acceptance Criteria:**

**Given** I am on the weekly table view and I click one of my own booked cells
**When** the booking-cancel popover opens
**Then** the popover container is sized so that the Person, Platz, Datum lines and
both the "Schliessen" and "Buchung stornieren" buttons are fully visible without any
internal scrolling
**And** the popover remains anchored to the cell that opened it

**Given** the popover is open
**When** the buttons render
**Then** neither button is clipped at the bottom edge of the container

### Story 33.4: Replace Booking-Type Radio With Single Checkbox and Always-On Dropdown

**FRs covered:** FR149

As a user booking an item,
I want a single "Book for a colleague" checkbox alongside the colleague dropdown,
so that the booking-type intent is expressed with one control and the colleague
dropdown does not appear and disappear as I toggle modes.

**Acceptance Criteria:**

**Given** I open the item-groups view
**When** the booking controls render
**Then** I see a single "Book for a colleague" checkbox (unchecked by default) and
a colleague-selection dropdown (rendered but disabled) — there is no "Book for
myself / Book for a colleague" radio group

**Given** the checkbox is unchecked
**When** I look at the colleague dropdown
**Then** it is visible, occupies the same space it would when enabled, and is
disabled (cannot be opened or typed into)

**Given** I check the "Book for a colleague" checkbox
**When** the change is applied
**Then** the colleague dropdown becomes enabled without any layout reflow
**And** I can pick a colleague exactly as before

**Given** I uncheck the box after having selected a colleague
**When** the change is applied
**Then** the dropdown becomes disabled again
**And** the previously selected colleague is cleared so a stale value is not used
when the box is rechecked

**Given** a booking is being made
**When** the box is unchecked
**Then** the booking is made for the current user (the previous "Book for myself"
behavior)

**Given** the Story 32.3 flex-wrap radio-group layout existed before this story
**When** this story ships
**Then** that layout is removed entirely; this story supersedes Story 32.3 for the
booking-type controls

### Story 33.5: Single-Line Compact Booking Controls Layout

**FRs covered:** FR150

As a user on the item-groups view,
I want all booking controls on a single line on a wide screen,
so that the controls take less vertical space and leave more room for tiles.

**Acceptance Criteria:**

**Given** I am on the item-groups view on a wide desktop viewport
**When** the booking-controls card renders
**Then** the day/week toggle, date or week selector, equipment filter input (with
its info icon), "Book for a colleague" checkbox, and colleague-selection dropdown
all appear on a single horizontal row
**And** none of them is truncated or hidden behind another

**Given** the same wide viewport
**When** I compare the new controls-card height to the previous (multi-row) layout
**Then** the new card is visibly shorter

**Given** I narrow the viewport progressively
**When** the row no longer fits horizontally
**Then** the controls wrap onto additional rows naturally (Vuetify flex-wrap), with
no element overlapping another

**Given** I switch between day mode and week mode
**When** the controls re-render
**Then** the single-line layout is preserved in both modes (the date picker swaps for
the week selector but the rest stays in place)

### Story 33.6: Entra ID Primary Login With Official Icon and "More Options" Toggle

**FRs covered:** FR151

As a corporate user landing on the login page,
I want the Entra ID button to be the obvious primary action with the official Microsoft
icon,
so that I can sign in with one click without being distracted by the local-credentials
form.

**Acceptance Criteria:**

**Given** Entra ID is configured on the server
**When** I open the login page unauthenticated
**Then** I see only the SitHub brand logo, the Entra ID sign-in button (showing the
official Microsoft Entra ID color icon), and a clickable "more login options" link
**And** the local username/password form is hidden by default

**Given** I click "more login options"
**When** the page updates
**Then** the local credentials form expands into view
**And** the link text changes to "less login options"

**Given** the credentials form is expanded
**When** I click "less login options"
**Then** the form collapses again
**And** the link text reverts to "more login options"

**Given** Entra ID is NOT configured on the server
**When** I open the login page
**Then** the local credentials form is shown by default
**And** the Entra ID button and toggle link are not rendered, so users are not locked
out

**Given** the Entra ID icon is rendered on the button
**When** the page loads
**Then** the icon SVG is served from the binary's embedded assets (no external
network request to Wikimedia or any other host); the SVG source was downloaded from
`https://upload.wikimedia.org/wikipedia/commons/8/8c/Microsoft_Entra_ID_color_icon.svg`
and committed to the repository under the existing embedded-assets location

### Story 33.7: SitHub Brand Logo for Login Page and Application Header

**FRs covered:** FR152

As a user of SitHub,
I want a consistent visual brand on the login page and in the app header,
so that the product feels finished and trustworthy.

**Acceptance Criteria:**

**Given** I open the login page
**When** the page renders
**Then** the SitHub logo (icon plus "SitHub" wordmark) is shown in a full vertical
layout above the Entra ID button

**Given** I am authenticated and viewing any page
**When** the application header renders
**Then** the SitHub logo is shown in a compact horizontal layout (icon plus
"SitHub" wordmark side-by-side) in place of any prior text-only branding
**And** the logo fits within the existing app-bar height without clipping

**Given** the logo asset is loaded
**When** any page renders
**Then** the SVG is served from the binary's embedded assets, sourced from
`private/sithub_logo.svg` committed to the repository at the project's existing
embedded-assets location; either the same SVG is reused with CSS layout for both
variants, or a dedicated horizontal variant is generated and embedded

**Given** the logo has been adopted
**When** any prior text-only "SitHub" branding element existed in the login page
or header
**Then** that element is removed so the logo is the single source of branding
