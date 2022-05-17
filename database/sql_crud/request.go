package sql_crud

type Request_Crud struct {
	OpType string `json:"opType"` //操作类型
}

//列表数据
type Request_Crud_List struct {
	Request_Crud
	Request_List
}

//删除数据
type Request_Crud_Del struct {
	Request_Crud

	IdList []int `json:"idList"` //删除渠道Id列表
}

//列表结构请求
type Request_List struct {
	Request_Crud

	State    int `json:"state"`    //请求状态 0表示返回全部  1表示返回设置的部分内容
	Page     int `json:"page"`     //当前页
	PageSize int `json:"pageSize"` //每页展示数量
	//Filters map[string]interface{} `json:"filters"` //过滤选项
	Filters string `json:"filters"` //过滤选项
	Order string `json:"order"` //排序选项
	IdList string `json:"idList"` //id列表
}
