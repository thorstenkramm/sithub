---
stepsCompleted: [1, 2, 3, 4, 5, 6]
inputDocuments:
  - /Users/thorsten/projects/thorsten/sithub/_bmad-output/planning-artifacts/product-brief-sithub-2026-01-17.md
  - /Users/thorsten/projects/thorsten/sithub/_bmad-output/planning-artifacts/prd.md
  - /Users/thorsten/projects/thorsten/sithub/_bmad-output/planning-artifacts/prd-validation-report.md
  - /Users/thorsten/projects/thorsten/sithub/docs/index.md
  - /Users/thorsten/projects/thorsten/sithub/private/epic-29.md
---

# UX Design Specification sithub

**Author:** Thorsten
**Date:** 2026-04-15

---

<!-- UX design content will be appended sequentially through collaborative workflow steps -->

## Executive Summary

### Project Vision

Design a desktop-only table view for SitHub that gives users a full weekly overview of all sub
areas and desks within the currently selected area, similar to the legacy Confluence matrix,
while preserving SitHub's stronger booking rules and interaction quality.

This table view should be accessed from the selected area context, for example
`/areas/office_2nd_floor/item-groups`, as an optional peer view next to Floor Plan, not as the
default view. It should inherit the currently selected week, respect each user's weekend
preference, and remember the user's last chosen desktop view and room-collapse state.

The core experience goal is to help users see where people are sitting, find a free desk near
friends or colleagues, and book directly from the matrix without leaving it.

### Target Users

The table view serves all user types:

- Individual users booking for themselves
- Users coordinating seating near colleagues
- Users booking on behalf of colleagues
- Admins overseeing occupancy and cancelling bookings when permitted

These users are not primarily looking for drill-down exploration. They want dense, reliable,
week-based situational awareness with immediate action.

### Key Design Challenges

- Preserve high-density weekly overview without falling back into a clumsy spreadsheet experience
- Support direct booking from cells while keeping the interaction lightweight and spatially
  anchored
- Express permissions clearly, especially for reserved desks: forbidden cells must feel sealed,
  allowed cells must feel normal
- Show occupancy clearly in very small cells using avatar plus initials, with full names
  available on hover
- Keep the matrix navigable at scale through sticky headers, sticky desk labels, collapsible room
  sections, and memorized collapse state
- Respect existing SitHub behavior such as selected week, weekend settings, booking rules,
  colleague booking, cancellation permissions, and immediate in-place refresh after changes

### Design Opportunities

- Reclaim the strongest emotional advantage of the old system: "I can see the whole office for
  the whole week in one glance."
- Turn the matrix into a living planning surface instead of a passive report by enabling direct
  booking and cancellation flows in context
- Make occupancy socially useful without adding search complexity by showing names and avatars
  directly in cells
- Use collapsed room summaries to communicate weekly occupied counts at a glance while reducing
  noise in large floor views
- Create a desktop specialist view that complements rather than replaces SitHub's existing
  mobile-friendly flows

## Core User Experience

### Defining Experience

The core experience of the table view is to scan the current week across the full selected area,
spot a suitable desk near colleagues, and book directly from the free cell without leaving the
matrix.

This view is not meant to replace SitHub's other navigation patterns. It is a desktop-only
specialist view for high-density weekly planning. Its value comes from combining overview and
action in the same place.

The experience should remain explicit and trustworthy. No important behavior should happen
automatically behind the user's back. Users should always understand what a cell means before
hover, what they can do with it, and what will happen next.

### Platform Strategy

The table view is for desktop mouse-and-keyboard environments only, with no mobile or tablet
support target. It should be accessed as an optional peer view next to Floor Plan from the
selected area context.

The table should show exactly one selected week at a time, inherited from the surrounding
SitHub context. Weekend columns must respect the user's existing settings: if weekends are
disabled, Saturday and Sunday should be absent entirely.

The matrix should prioritize avoiding horizontal scrolling. Vertical scrolling is acceptable.
Weekday columns may become narrow and dense to preserve the single-board overview.

The UI should remember:
- last selected desktop view for the area context
- room collapse state in local storage across visits and week changes

### Effortless Interactions

The most important effortless behavior is rapid week scanning. Users should be able to
understand the board quickly without hunting, drilling down, or hovering just to decode state.

The following interactions should feel especially seamless:

- switching from Floor Plan to Table View
- visually scanning booked versus free desks across the week
- identifying colleagues through avatar plus initials in booked cells
- spotting free desks near occupied cells belonging to familiar people
- booking directly from a free cell through a lightweight anchored confirmation
- cancelling from an occupied cell through the same lightweight anchored confirmation pattern
- collapsing and reopening room sections without losing context

Cells should be understandable without hover. Hover is only for enrichment, such as full names
and equipment hints.

After a successful booking or cancellation, the user must remain in the matrix and see the cell
update immediately. The popover should close, and the table itself should become the
confirmation.

### Critical Success Moments

The primary success moment is when a user sees a free cell and books directly from it. That is
the moment where SitHub clearly becomes better than Confluence.

Other critical moments include:

- opening the table and immediately understanding the week without horizontal scrolling
- recognizing where people are sitting from compact booked cells
- understanding permissions at a glance, especially for locked reserved desks
- trusting that room collapse, selected view, and week context persist predictably
- completing booking or cancellation in-place without being pulled into another page or flow

The interaction that would most damage this experience is excessive horizontal scrolling. If
users lose their place or need to pan sideways to understand the week, the table fails its core
purpose.

### Experience Principles

- **Overview first:** The table must win on whole-week scanability before it wins on feature richness.
- **Action in place:** Free cells should turn directly into booking, and eligible occupied cells should turn directly into cancellation, without leaving the matrix.
- **Clarity before hover:** Every cell's basic meaning must be obvious without hover or extra clicks.
- **Desktop density with discipline:** The design may compress columns and content, but never at the cost of orientation.
- **Stable context:** Week selection, view selection, and room-collapse memory should make the board feel persistent and dependable.
- **One interaction language:** Booking and cancellation should use the same lightweight anchored confirmation style so the matrix feels coherent.

## Desired Emotional Response

### Primary Emotional Goals

The primary emotional goal of the table view is confidence. Users should feel that they
understand the office state, trust what the matrix is telling them, and can act without
hesitation.

A strong supporting emotional goal is power. When users open the table, they should feel that
they have a complete, useful command of the week across the whole selected area.

During direct booking, the strongest emotional quality should be speed. The interaction should
feel immediate and decisive, not ceremonial.

### Emotional Journey Mapping

**On entry to the table view**  
Users should feel powerful. The board should communicate, at a glance, that they can see the
full situation and do not need to jump through multiple views to understand the week.

**During scanning**  
Users should feel confident. The state of each room and desk should feel legible, stable, and
easy to interpret without hover or extra clicks.

**During booking**  
Users should feel fast. Selecting a free cell and confirming the booking should feel like a
short, direct action with no loss of context.

**During cancellation**  
Users should still feel in control. The anchored confirmation should feel lightweight and
reversible, not dramatic.

**When encountering restricted reserved cells**  
Users should feel gently informed, not punished. The interaction should communicate boundaries
in a soft and polite way through passive visual cues such as lock indicators and non-clickable
behavior.

**On return visits**  
Users should feel continuity and trust. Remembered view choice and room-collapse state should
make the system feel dependable and familiar.

### Micro-Emotions

The most important micro-emotions in this experience are:

- confidence instead of confusion
- trust instead of skepticism
- speed instead of drag
- control instead of friction
- clarity instead of ambiguity
- polite constraint instead of rejection

These micro-emotions matter because the table is dense. In a dense interface, even small
moments of hesitation or misreading quickly become emotional fatigue.

### Design Implications

To create confidence:
- states must be understandable without hover
- sticky headers and sticky left columns must preserve orientation
- the selected week and visible room structure must feel stable and predictable

To create power:
- the matrix must reveal the entire weekly situation across all sub areas in one board
- collapsed room summaries must still communicate meaningful occupancy signals
- the view must preserve whole-context awareness instead of fragmenting users into drill-down
  screens

To create speed:
- free cells must support direct booking from the table
- booking confirmation must be lightweight and anchored near the point of action
- success must update in place immediately without navigation away from the board

To create soft and polite constraint:
- forbidden reserved cells should use passive signals such as lock icons and non-clickable
  behavior
- the design should avoid harsh error moments where possible
- restrictions should feel clear, calm, and unsurprising

### Emotional Design Principles

- **Confidence through legibility:** users should trust the board because it is easy to read.
- **Power through overview:** the table should feel like a full-week command surface, not a narrow list.
- **Speed through locality:** actions should happen exactly where attention already is.
- **Trust through stability:** remembered state and in-place updates should make the interface feel dependable.
- **Constraint without hostility:** permission limits should feel soft and polite, never punitive.

## UX Pattern Analysis & Inspiration

### Inspiring Products Analysis

The primary inspiration source for this table view is the existing Confluence booking page
currently used by SitHub users.

What the Confluence page does well:

- presents the entire selected week in one dense board
- groups desks by sub area and room in a way users already understand
- allows users to visually compare rooms and desks without drill-down
- supports rapid social scanning by showing who sits where across the week
- gives users a strong sense of control through overview

Why it remains compelling:

- it is comprehensive rather than fragmented
- it compresses many decisions into one visible surface
- users trust it because the structure is stable and familiar

Where it fails:

- it does not allow direct booking from a free cell
- the overview stops short of action, forcing users into a separate mental or manual step

This means the new SitHub table view should not try to reinvent the weekly matrix pattern. It
should preserve the overview behavior users already value and attach a better action model to
it.

### Transferable UX Patterns

The following patterns from the Confluence matrix should be transferred into SitHub:

**Structure Patterns**
- grouped sections for sub areas and rooms
- desks as stable row entries within each room
- weekdays as the primary columns
- one selected week visible at a time

**Scanning Patterns**
- dense whole-week visibility across all rooms
- direct visual comparison between desks and days
- occupancy visible directly inside cells
- minimal reliance on secondary navigation

**Recognition Patterns**
- users should be able to recognize where people sit by reading the board itself
- the table should preserve the familiar planning-wall feeling users already know
- collapsed sections should still summarize weekly room occupancy in a compact way

**Interaction Patterns to Adapt**
- keep the matrix structure familiar
- add direct booking from free cells
- add direct cancellation for permitted users from occupied cells
- keep actions anchored to the clicked cell so the user never leaves the board

### Anti-Patterns to Avoid

The primary anti-pattern to avoid is preserving the old Confluence limitation where the board is
only informative but not actionable.

Specifically, avoid:

- forcing users to leave the matrix to complete a booking
- turning a free cell into a dead end instead of a direct action point
- breaking the see-it-and-act-on-it flow that should define SitHub's advantage
- introducing interaction layers that make the matrix feel less direct than the old overview

No additional anti-patterns were identified from the legacy table beyond this action gap. The
old overview model itself is considered strong and worth preserving.

### Design Inspiration Strategy

**What to Adopt**
- the Confluence table's weekly matrix structure
- room-based grouping and desk-row organization
- the strong single-glance overview across an entire week
- the familiar visual planning-board mental model

**What to Adapt**
- replace passive free cells with directly bookable free cells
- replace static occupancy text with clearer booked-cell identity signals such as avatar plus
  initials and hover name
- add reserved-state permission signaling through soft lock-based cues
- add collapsible room sections with remembered state and compact occupancy summaries

**What to Avoid**
- any redesign that sacrifices overview for visual novelty
- any interaction pattern that introduces horizontal-scroll dependence
- any booking flow that pulls the user away from the matrix
- any permission treatment that feels harsh, noisy, or punitive

The design strategy is therefore conservative in structure and progressive in action: keep the
familiar weekly board, preserve the overview users already trust, and let SitHub win by making
the board directly useful.

## Design System Foundation

### 1.1 Design System Choice

The table view should remain visually inside SitHub's existing design system and styling
language. It must be implemented strictly with Vue and Vuetify. No alternative UI framework,
table library, or styling system should be introduced.

For the table itself, the recommended foundation is a purpose-built custom matrix component
built from Vue and Vuetify primitives rather than a generic data-table pattern.

### Rationale for Selection

This feature has a highly specific interaction model:

- grouped room sections
- desks as stable rows
- weekdays as dynamic columns based on user settings
- collapsible room blocks with memorized state
- sticky weekday header and sticky left label column
- free, occupied, and reserved permission-aware cell states
- anchored booking and cancellation popovers
- immediate in-place updates after action

A generic table or data-table abstraction would create friction because the behavior is not
really tabular CRUD. It is a dense planning matrix with custom interaction rules.

Using a purpose-built component keeps the UX honest and the code more maintainable:
- the structure matches the product model directly
- logic is easier to isolate and reason about
- Vuetify still provides visual consistency, tokens, spacing, overlays, icons, and form controls
- the implementation stays aligned with the existing SitHub stack

### Implementation Approach

Implement the feature as a focused matrix-view module composed of smaller Vue components built
with Vuetify primitives.

Recommended component strategy:

- a page-level table-view container integrated into the existing area context
- a weekly matrix header for weekday columns
- room section components with collapse controls and compact weekly summaries
- desk row components for stable row rendering
- booking cell components for state-specific rendering
- lightweight anchored booking and cancellation popovers
- composables for week-column generation, collapse persistence, remembered view state, and
  permission/state mapping

This approach keeps custom complexity localized while preserving consistency with the broader
SitHub interface.

### Customization Strategy

Customization should stay inside the existing SitHub visual language rather than introducing a
new specialist aesthetic.

The matrix may be denser than the rest of the product, but it should still use:
- Vuetify components and overlay patterns
- existing spacing, typography, iconography, and feedback conventions where possible
- SitHub-consistent color semantics for free, occupied, muted past days, and restricted states

The customization focus should be structural and behavioral, not stylistic:
- tune density for desktop scanability
- preserve strong orientation through sticky elements
- express permission and occupancy clearly within compact cells
- keep the interaction model specialized while the visual system remains familiar
