package util

import "testing"

func TestHiddenPhone(t *testing.T) {
	phone := "15501707783"
	t.Logf("phone: %v, hidden: %v",phone,HiddenPhone(phone))
}
