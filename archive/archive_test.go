package archive

import (
	"fmt"
	"testing"
)

func TestTarFileList(t *testing.T) {
	/*
	if str, err := TarFileList("test/arc.tar.bz2"); err != nil {
		t.Errorf("TarFileList:%v", err)
	} else {
		fmt.Println(str)
	}
	str, err := ArchiveFileList("test/arc.tar.gz")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(str)

	err = UnpackArchive("test/a.zip")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = CreateZip("test/a.zip", []string{"test/archive_test.go", "test/archive.go"})
	if err != nil {
		fmt.Println(err.Error())
	}
	*/

	//err1 := CreateZipByPath("123.zip","E:/_Temp3/20180911142814")
	//if nil != err1 {
	//	fmt.Println("---err1:",err1.Error())
	//}
	err2 := UnpackZip("123.zip", "E:/_Temp3_ouotput")
	if nil != err2 {
		fmt.Println("---err2:",err2.Error())
	}
	fmt.Println("---end")
}
