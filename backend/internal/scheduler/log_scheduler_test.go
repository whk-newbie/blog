package scheduler

import (
	"testing"

	"github.com/whk-newbie/blog/internal/service"
)

func TestNewLogSchedulerUsesDefaultRetentionDays(t *testing.T) {
	scheduler := NewLogScheduler(nil, 0)

	if service.DefaultLogRetentionDays != 15 {
		t.Fatalf("DefaultLogRetentionDays = %d, want 15", service.DefaultLogRetentionDays)
	}
	if scheduler.retentionDays != service.DefaultLogRetentionDays {
		t.Fatalf("retentionDays = %d, want %d", scheduler.retentionDays, service.DefaultLogRetentionDays)
	}
}
