# Log System Repair Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Display persisted system logs in the admin UI, retain them for 15 days, and stop Air from rebuilding on generated Swagger files.

**Architecture:** Keep the existing separation between stdout application logging and the database Hook. Introduce one shared service-level retention default so the scheduler, router, and manual cleanup endpoint cannot drift. Keep the API's existing `items` pagination contract and make the Vue view consume it directly.

**Tech Stack:** Go 1.25, Gin, logrus, robfig/cron, Vue 3, Vite, Air.

---

### Task 1: Make the admin page consume the log API contract

**Files:**
- Modify: `frontend/src/views/admin/LogManage.vue:227-229`
- Test: `frontend` production build

- [ ] **Step 1: Confirm the current UI/API mismatch**

The API response data uses `items`:

```go
type LogListResponse struct {
    Items []*LogResponse `json:"items"`
}
```

The current Vue page incorrectly uses `response.list`:

```js
logs.value = response.list || []
```

- [ ] **Step 2: Read the returned `items` collection**

```js
const response = await api.log.getLogs(params)
logs.value = response.items || []
pagination.total = response.total || 0
```

- [ ] **Step 3: Build the frontend**

Run: `npm run build`

Expected: Vite completes successfully and its production verification script exits with status 0.

### Task 2: Centralize and apply the 15-day retention default

**Files:**
- Modify: `backend/internal/service/log_service.go:13-25`
- Modify: `backend/internal/scheduler/log_scheduler.go:22-33`
- Modify: `backend/internal/router/router.go:286-296`
- Modify: `backend/internal/handler/log_handler.go:130-150`
- Modify: `frontend/src/views/admin/LogManage.vue:188-190, 300-303`
- Create: `backend/internal/scheduler/log_scheduler_test.go`
- Create: `backend/internal/handler/log_handler_test.go`

- [ ] **Step 1: Write failing scheduler and handler tests**

```go
func TestNewLogSchedulerUsesDefaultRetentionDays(t *testing.T) {
    scheduler := NewLogScheduler(nil, 0)
    if service.DefaultLogRetentionDays != 15 {
        t.Fatalf("DefaultLogRetentionDays = %d, want 15", service.DefaultLogRetentionDays)
    }
    if scheduler.retentionDays != service.DefaultLogRetentionDays {
        t.Fatalf("retentionDays = %d, want %d", scheduler.retentionDays, service.DefaultLogRetentionDays)
    }
}
```

```go
func TestCleanupLogsUsesDefaultRetentionDaysForMissingBody(t *testing.T) {
    logService := &recordingLogService{}
    handler := NewLogHandler(logService)
    recorder := httptest.NewRecorder()
    context, _ := gin.CreateTestContext(recorder)
    context.Request = httptest.NewRequest(http.MethodPost, "/admin/logs/cleanup", nil)

    handler.CleanupLogs(context)

    if logService.retentionDays != service.DefaultLogRetentionDays {
        t.Fatalf("retentionDays = %d, want %d", logService.retentionDays, service.DefaultLogRetentionDays)
    }
}
```

```go
type recordingLogService struct{ retentionDays int }

func (*recordingLogService) Log(models.LogLevel, string, string, map[string]interface{}, *uint, string) error {
    return nil
}
func (*recordingLogService) GetLogs(*service.LogQueryRequest) (*service.LogListResponse, error) {
    return nil, nil
}
func (*recordingLogService) GetLogByID(uint) (*service.LogResponse, error) {
    return nil, nil
}
func (s *recordingLogService) CleanupOldLogs(days int) (*service.CleanupResponse, error) {
    s.retentionDays = days
    return &service.CleanupResponse{}, nil
}
```

- [ ] **Step 2: Run the new tests and confirm the old 90-day default fails**

Run: `go test ./internal/scheduler ./internal/handler`

Expected: compilation fails because `service.DefaultLogRetentionDays` does
not exist yet.

- [ ] **Step 3: Define one shared default and replace every hard-coded 90**

In `internal/service/log_service.go`:

```go
const DefaultLogRetentionDays = 15
```

Use `service.DefaultLogRetentionDays` in the scheduler fallback, router
manager construction, cleanup endpoint JSON-bind fallback, Swagger example,
and the Vue cleanup form's initial/reset value.

- [ ] **Step 4: Re-run backend tests**

Run: `go test ./internal/scheduler ./internal/handler`

Expected: both packages pass and the tests prove missing retention input uses
15 days.

### Task 3: Preserve the database-hook level boundary

**Files:**
- Create: `backend/internal/pkg/logger/database_hook_test.go`
- Modify: `backend/internal/pkg/logger/database_hook.go:18-28` only if needed for testability

- [ ] **Step 1: Write a regression test for the default Hook levels**

```go
func TestNewDatabaseHookOnlyPersistsWarningsAndErrors(t *testing.T) {
    hook := NewDatabaseHook(nil)
    want := []logrus.Level{logrus.WarnLevel, logrus.ErrorLevel}
    got := hook.Levels()
    if len(got) != len(want) {
        t.Fatalf("Levels() length = %d, want %d", len(got), len(want))
    }
    for i := range want {
        if got[i] != want[i] {
            t.Fatalf("Levels()[%d] = %s, want %s", i, got[i], want[i])
        }
    }
}
```

This uses direct length-and-position comparison and adds no dependency.

- [ ] **Step 2: Run the logger test**

Run: `go test ./internal/pkg/logger`

Expected: PASS, proving successful `INFO` requests remain stdout-only.

### Task 4: Stop Air from watching generated Swagger output

**Files:**
- Modify: `backend/.air.toml:12`
- Test: Air configuration inspection and development-container logs

- [ ] **Step 1: Exclude the generated `docs` directory from Air watches**

```toml
exclude_dir = ["assets", "docs", "tmp", "vendor", "testdata", "uploads", "backups", "logs"]
```

`swag init` remains in the build command. Its generated files no longer
trigger the next build.

- [ ] **Step 2: Format and run all focused Go tests**

Run: `gofmt -w internal/service/log_service.go internal/scheduler/log_scheduler.go internal/scheduler/log_scheduler_test.go internal/handler/log_handler.go internal/handler/log_handler_test.go internal/pkg/logger/database_hook_test.go`

Run: `go test ./internal/scheduler ./internal/handler ./internal/pkg/logger`

Expected: all packages pass.

- [ ] **Step 3: Verify the dev server is no longer in a Swagger rebuild loop**

Run: `docker logs --tail 120 blog-backend-dev`

Expected: a source change produces one `building...`/`running...` sequence,
not repeated sequences caused by `docs/docs.go has changed`.

### Task 5: Verify the completed behavior

**Files:**
- Verify: `frontend/src/views/admin/LogManage.vue`
- Verify: `backend/internal/service/log_service.go`
- Verify: `backend/internal/scheduler/log_scheduler.go`
- Verify: `backend/internal/handler/log_handler.go`
- Verify: `backend/.air.toml`

- [ ] **Step 1: Run production frontend verification**

Run: `npm run build`

Expected: exit status 0.

- [ ] **Step 2: Run all backend tests**

Run: `go test ./...`

Expected: exit status 0.

- [ ] **Step 3: Inspect the final diff**

Run: `git diff --check` and `git diff --stat`

Expected: only log UI, retention, Air configuration, focused tests, and this
plan are changed; the unrelated encryption work remains untouched.
