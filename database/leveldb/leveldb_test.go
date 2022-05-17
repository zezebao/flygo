package leveldb

import (
	"testing"
	"fmt"
)

func TestDB(t *testing.T) {
	//db, _ := OpenFile("leveldb_db1", nil)
	//db.PutString("aaa", "aaa")
	//db.PutInt("bbb", 123)

	fmt.Println("-----", GetString("aaa", ""), GetInt("bbb", 0))
}
