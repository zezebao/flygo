package sql_crud

import (
	"github.com/astaxie/beego/orm"
)

//crud钩子

//插入
type ICrud_Insert interface {
	BeforeInsert()(bool,error)//插入前操作
}
//更新
type ICrud_Update interface {
	BeforeUpdate()(bool,error) //更新前操作
}
//删除
type ICrud_Delete interface {

}
//列表查询
type ICrud_List interface {
	QueryList()(bool,error)
}
//列表查询过滤
type ICrud_List_Filter interface {
	QueryListFilter(jsonStr string,qs orm.QuerySeter)(orm.QuerySeter)
}


//TODO 2018-12-31新增 ICrud接口
type ICrud interface {

}