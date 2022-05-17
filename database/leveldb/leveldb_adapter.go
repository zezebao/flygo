package leveldb

import (
	"errors"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var db *MyLevelDb

func init() {
	//o := &opt.Options{}
	db, _ = OpenFile("leveldb_db1", nil)
}

type MyLevelDb struct {
	DB *leveldb.DB
}

func OpenFile(path string, o *opt.Options) (*MyLevelDb, error) {
	db := &MyLevelDb{}
	var err error
	db.DB, err = leveldb.OpenFile(path, o)
	return db, err
}

func Close() error {
	if nil == db.DB {
		return errors.New("db=nil")
	}
	return db.DB.Close()
}

func GetInt(key string, defaultVal int) int {
	if nil == db.DB {
		return defaultVal
	}
	data, err := db.DB.Get([]byte(key), nil)
	if nil != err {
		return defaultVal
	}
	result, err := strconv.Atoi(string(data))
	if nil != err {
		return defaultVal
	}
	return result
}

func SetInt(key string, val int) error {
	if nil == db.DB {
		return errors.New("db=nil")
	}
	return db.DB.Put([]byte(key), []byte(strconv.Itoa(val)), nil)
}

//改变数字
func SetIntChange(key string, changeVal int) error {
	count := 0
	count = GetInt(key, count)
	count = count + changeVal
	return SetInt(key, count)
}

func GetString(key string, defaultVal string) string {
	if nil == db.DB {
		return defaultVal
	}
	data, err := db.DB.Get([]byte(key), nil)
	if nil != err {
		return defaultVal
	}
	return string(data)
}
func SetString(key string, val string) error {
	if nil == db.DB {
		return errors.New("db=nil")
	}
	return db.DB.Put([]byte(key), []byte(val), nil)
}

func Get(key []byte, ro *opt.ReadOptions) ([]byte, error) {
	return db.DB.Get(key, ro)
}

func Set(key string, value string) error {
	return db.DB.Put([]byte(key), []byte(value), nil)
}

func Delete(key string) error {
	return db.DB.Delete([]byte(key), nil)
}

func Has(key string) (bool) {
	ret, err := db.DB.Has([]byte(key), nil)
	if nil != err {
		return false
	}
	return ret
}
