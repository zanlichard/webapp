package toolkit

import "testing"

func TestUniq(t *testing.T) {
	basicInfo, id, err := GetUniqId("check_app_version")
	if err != nil {
		t.Logf("generate uniq id failed for:%+v", err)
	}
	t.Logf("origin: %s, id: %s", basicInfo, id)
}
