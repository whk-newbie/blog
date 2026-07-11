package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
)

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
