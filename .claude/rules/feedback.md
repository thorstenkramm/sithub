# User Feedback Conventions

All user-facing feedback (success, error, warning) for transient operations
must follow these rules to ensure a consistent experience.

## Snackbar for Transient Feedback

Use Vuetify's `v-snackbar` for all operation feedback: bookings, cancellations, saves,
deletions, filter actions, and any other user-initiated action that succeeds or fails.

### Rules

- Position: `location="bottom"`
- Success: `color="success"`, timeout `3000`ms
- Error: `color="error"`, timeout `6000`ms (longer — users need time to read errors),
  add `closable` for manual dismiss
- Warning: `color="warning"`, timeout `4000`ms
- Always include a `data-cy` attribute for testability

### Do

```vue
<v-snackbar v-model="showFeedback" :color="feedbackColor" location="bottom" :timeout="timeout">
  {{ feedbackMessage }}
</v-snackbar>
```

### Do Not

```vue
<!-- WRONG: inline v-alert for transient operation feedback -->
<v-alert v-if="errorMessage" type="error" class="mb-4" closable>
  {{ errorMessage }}
</v-alert>
```

## When to Use v-alert

Use `v-alert` only for **persistent page-level states** that are not
triggered by a user action, for example:

- Connection lost (shown until reconnected)
- Page load failures (area not found, unauthorized)
- Empty states with guidance

These are not transient — they stay visible until the underlying condition changes.

## Summary

| Scenario | Component | Color | Timeout |
| --- | --- | --- | --- |
| Operation succeeded | `v-snackbar` | `success` | 3000ms |
| Operation failed | `v-snackbar` | `error` | 6000ms |
| Validation warning | `v-snackbar` | `warning` | 4000ms |
| Page-level persistent state | `v-alert` | varies | none |
