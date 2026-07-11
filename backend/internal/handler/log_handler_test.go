package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/whk-newbie/blog/internal/models"
	"github.com/whk-newbie/blog/internal/service"
)

type recordingLogService struct {
	retentionDays int
}

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

func TestCleanupLogsUsesDefaultRetentionDaysForMissingBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
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
