package def

import (
	"testing"
)

func TestHP(t *testing.T) {
	skp := HP()
	if skp == nil {
		t.Error("HP() should not return nil")
	}
}
