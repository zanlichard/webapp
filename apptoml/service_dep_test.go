package apptoml

import (
	"testing"
)

func TestServiceDep(t *testing.T) {
	err := ServiceDepInit("..\\etc\\dependence.xml")
	if err != nil {
		t.Errorf("open service dependence failed for:%+v", err)
		return
	}
	defer ServiceDepFree()
	productServiceUrlPath := "/DEPENDENTSERVERINFO/ProductService/Url"
	productServiceKeyPath := "/DEPENDENTSERVERINFO/ProductService/Key"
	t.Logf("get productServiceUrl:%s", ServiceDepGetField(productServiceUrlPath))
	t.Logf("get productServiceKey:%s", ServiceDepGetField(productServiceKeyPath))

	//os.PathSeparator

}
