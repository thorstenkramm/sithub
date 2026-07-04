
## Deferred from: code review of epic-35 (2026-07-04)

- a11y: plain warning icons (weekly matrix + floor plan) surface the warning only on hover and are
  not keyboard-focusable. Pre-existing pattern for the matrix; new for the floor plan. Needs a
  broader accessibility decision (make warning indicators focusable / add an accessible name).
- Possible double tooltip on free floor-plan items: the item wrapper tooltip (name/equipment) and the
  warning-icon tooltip can both appear on hover. Verify visually against private/epic-35.md img_30.
- FR165 coverage gap: the "dismiss → warning text changes → confirmation re-shows" path is unit-tested
  only piecewise (useWarningSuppression hashing). Add an integrated test through the shared dialog.
- Weekly-matrix booking: no busy/loading feedback on the popover Book button during the warning
  confirmation step (submitting is only set inside doBook).
