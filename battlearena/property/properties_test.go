package property

import (
	"testing"
)

func TestProperties(t *testing.T) {
	res := MakeIntProperty("toto", 0, FriendlyController, Character)

	var res2 Property = res

	var res3 IntProperty = res

	if res3 != res2 {
		t.Error("res3 != res2")
	}
}
