package xml

import (
	"os"

	xmlpath "gopkg.in/xmlpath.v1"
)

var file *os.File = nil
var root *xmlpath.Node = nil

func XmlInit(fileName string) error {
	var err error
	file, err = os.OpenFile(fileName, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	root, err = xmlpath.Parse(file)
	if err != nil {
		return err
	}
	return nil
}

func XmlFree() error {
	if file != nil {
		return file.Close()
	}
	return nil
}

func XmlGetField(accessPath string) string {
	if path := xmlpath.MustCompile(accessPath); path != nil {
		it := path.Iter(root)
		for it.Next() {
			return it.Node().String()
		}
	}
	return ""
}
