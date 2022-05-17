package sql_crud

import (
	"reflect"
)

func newStruct(elem reflect.Type) interface{} {
	return reflect.New(elem).Interface()
}

//--CURD操作公共逻辑-------------------------------------
type CRUD struct {
	Table             string      //表名
	modelDb           interface{} //数据库中的结构体
	modelResponseItem interface{} //请求列表结构体
	orderBy           []string    //排序依据
	callbackNew       func(v interface{})
	callbackUpdate       func(v interface{})
	callbackDelete    func(id int)
	callbackList      func(v interface{})
}

//--getter-----------------------------------------------------------------

//获取数据库模型
func (c *CRUD) GetNewModelDb() interface{} {
	val := reflect.ValueOf(c.modelDb)
	typ := reflect.Indirect(val).Type()
	return newStruct(typ)
	//return newStruct(reflect.TypeOf(c.modelDb).Elem())
}

//获取返回item数组
func (c *CRUD) GetNewModelResponseItemSlice(list *Request_List) interface{} {
	if 1 == list.State {
		if nil != c.modelResponseItem {
			val := reflect.ValueOf(c.modelResponseItem)
			typ := reflect.Indirect(val).Type()

			return newStruct(reflect.SliceOf(typ))
		}
	}
	val := reflect.ValueOf(c.modelDb)
	typ := reflect.Indirect(val).Type()
	return newStruct(reflect.SliceOf(typ))
}

//--setter------------------------------------------------------------------
func (c *CRUD) SetCallbackNew(val func(v interface{})) *CRUD {
	c.callbackNew = val
	return c
}

func(c *CRUD) SetCallbackUpdate(val func(v interface{})) *CRUD {
	c.callbackUpdate = val
	return c
}

func (c *CRUD) SetCallbackList(val func(v interface{})) *CRUD {
	c.callbackList = val
	return c
}

func (c *CRUD) SetCallbackDelete(val func(id int)) *CRUD {
	c.callbackDelete = val
	return c
}

//设置数据库对应Model
func (c *CRUD) SetModelDb(model interface{}) *CRUD {
	c.modelDb = model
	c.Table = getTableName(reflect.ValueOf(model))
	return c
}

//设置返回数据对应模型
func (c *CRUD) SetModelResponseItem(model interface{}) *CRUD {
	c.modelResponseItem = model
	return c
}

//设置orderby
func (c *CRUD) SetOrderBy(orderBy ...string) *CRUD {
	c.orderBy = orderBy
	return c
}
