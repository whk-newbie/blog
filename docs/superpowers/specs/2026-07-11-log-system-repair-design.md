# Log System Repair Design

## Goal

Make persisted system logs visible in the admin UI, retain them for 15 days,
and stop the development server from continuously rebuilding itself.

## Scope

- The admin log page reads the backend response's `items` collection rather
  than the nonexistent `list` collection.
- Database persistence remains restricted to warning and error events. Normal
  successful HTTP requests continue to be written only to the container's
  standard output.
- Both scheduled and manual log cleanup default to a 15-day retention period.
- Air excludes generated Swagger documentation from its watched directories,
  while still generating it during a build.

## Data Flow

`logrus` emits all runtime logs to stdout. Its database hook receives only
`WARN` and `ERROR` records and asynchronously stores those records in
`system_logs`. `GET /api/v1/admin/logs` returns a paginated object containing
`items`; the admin page binds that array directly to its table.

## Verification

- Add a frontend regression test that proves `items` populates the log table
  state and that `list` is not required.
- Add backend tests for the 15-day defaults and the database hook's accepted
  levels.
- Run the relevant frontend and Go test suites.
- Start Air or inspect its watch configuration to verify generated `docs`
  files do not trigger a rebuild loop.

## Out of Scope

- Persisting all successful HTTP requests in `system_logs`.
- Changing the configured stdout destination to a rotating file.
- Altering the unrelated in-progress encryption changes.
