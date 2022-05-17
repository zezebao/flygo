package sql_crud

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"fmt"
	"reflect"
	"strings"

	"com.ffl/common"
	"com.ffl/lkpush"

	"com.ffl/common/models"
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
)

//是否初始化
var _inited bool = false
var _crudMap map[string]*CRUD = make(map[string]*CRUD)

var _engine *gin.Engine
var _db *gorm.DB

//操作码定义
const (
	OP_NEW    = 1
	OP_UPDATE = 2
	OP_DELETE = 3
	OP_LIST   = 4
)

//回调
var _callback func(string)

func SetCallback(callback func(string)) {
	_callback = callback
}

func Init(ginEngine *gin.Engine, db *gorm.DB) {
	if _inited {
		return
	}

	_engine = ginEngine
	_db = db

	_inited = true
	ginEngine.GET("/com/list", onCrudSelectListHandler)
	ginEngine.POST("/com/new", onCrudNewHandler)
	ginEngine.POST("/com/edit", onCrudEditHandler)
	ginEngine.POST("/com/del", onCrudDeleteHandler)
}

//注册crud支持
//sql_crud.RegisterCRUD("keyname").
func RegisterCRUD(key string) *CRUD {
	fmt.Println("---curd.RegisterCRUD:" + key)
	tmp := &CRUD{}
	_crudMap[key] = tmp
	return tmp
}

//检测opType
func _checkOpType(opType string) (*CRUD, bool, error) {
	if len(opType) <= 0 {
		return nil, false, errors.New("类型异常")
	}
	if v, ok := _crudMap[opType]; !ok {
		return nil, false, errors.New("未注册类型")
	} else {
		return v, true, nil
	}
}

//新增
func onCrudNewHandler(c *gin.Context) {
	response := &models.Response_Common{}
	response.Msg = "新增成功"

	jsonStr := common.UtilGin.Gin_ParseDataFromContext(c)
	requestData := &Request_Crud{}
	json.Unmarshal([]byte(jsonStr), requestData)

	if v, ok, err := _checkOpType(requestData.OpType); !ok {
		response.Code = 1
		response.Msg = err.Error()
	} else {
		tmp := v.GetNewModelDb()
		json.Unmarshal([]byte(jsonStr), tmp)

		if v, ok := tmp.(lkpush.IInfo_Common); ok {
			v.CreateTime()
		}
		if v, ok := tmp.(ICrud_Insert); ok {
			_, err := v.BeforeInsert()
			if err != nil {
				response.Code = 1
				response.Msg = err.Error()
			}
		}

		if 0 == response.Code {
			common.LogConsole.Debug("=======Insert:%#v", tmp)
			//_, err := lkpush.MOrmer.Insert(tmp)
			err := _db.Create(tmp).Error
			if nil != err {
				common.LogConsole.Debug("---insert err:%s", err.Error())

				response.Code = 1
				response.Msg = err.Error()
			} else {
				if nil != v.callbackNew {
					v.callbackNew(tmp)
				}
			}
		}
	}
	c.String(http.StatusOK, common.UtilCommon.Parse_ToJson(&response))

	if nil != _callback && 0 == response.Code {
		go _callback(requestData.OpType)
	}
}

//更新
func onCrudEditHandler(c *gin.Context) {
	response := &models.Response_Common{}
	response.Msg = "更新成功"

	jsonStr := common.UtilGin.Gin_ParseDataFromContext(c)
	requestData := &Request_Crud{}
	json.Unmarshal([]byte(jsonStr), requestData)

	model := &gorm.Model{}
	json.Unmarshal([]byte(jsonStr), model)

	if v, ok, err := _checkOpType(requestData.OpType); !ok {
		response.Code = 1
		response.Msg = err.Error()
	} else {
		tmp := v.GetNewModelDb()

		updateMap := make(map[string]interface{}, 0)
		json.Unmarshal([]byte(jsonStr), &updateMap)
		err00 := json.Unmarshal([]byte(jsonStr), tmp)
		if nil != err00 {
			fmt.Println("===err00:", err00.Error())
		}
		//if v, ok := tmp.(lkpush.IInfo_Common); ok {
		//	v.UpdateTime()
		//}
		//if v, ok := tmp.(ICrud_Update); ok {
		//	v.BeforeUpdate()
		//}

		fmt.Println("====update json:", jsonStr)
		fmt.Printf("====update v:%#v\n", tmp)
		fmt.Printf("====update map v:%#v, %d\n", updateMap, len(updateMap))

		//_db.Model(&tmp).Update(updateMap)
		_db.Table(v.Table).Where("id=?", model.ID).Save(tmp)
		//lkpush.MOrmer.Update(tmp)

		if nil != v.callbackUpdate {
			v.callbackUpdate(tmp)
		}
	}
	c.String(http.StatusOK, common.UtilCommon.Parse_ToJson(&response))

	if nil != _callback && 0 == response.Code {
		go _callback(requestData.OpType)
	}
}

//删除
func onCrudDeleteHandler(c *gin.Context) {
	response := &models.Response_Common{}
	response.Msg = "删除成功"

	jsonStr := common.UtilGin.Gin_ParseDataFromContext(c)
	requestData := &Request_Crud_Del{}
	json.Unmarshal([]byte(jsonStr), requestData)

	if v, ok, err := _checkOpType(requestData.OpType); !ok {
		response.Code = 1
		response.Msg = err.Error()
	} else {
		var sql string
		sql = `DELETE FROM ` + v.Table + ` WHERE id IN (`
		idLen := len(requestData.IdList)
		for i := 0; i < idLen; i++ {
			sql += strconv.Itoa(requestData.IdList[i])
			if i < idLen-1 {
				sql += `,`
			}
		}
		sql += `)`

		common.LogConsole.Debug("----sql:%s", sql)
		//lkpush.MOrmer.Raw(sql).Exec()
		_db.Exec(sql)
	}
	c.String(http.StatusOK, common.UtilCommon.Parse_ToJson(&response))

	if nil != _callback && 0 == response.Code {
		go _callback(requestData.OpType)
	}
}

//忽略列表
var _ignoreTypes = []string{"time.Time"}

func _isInIgnoreList(typeStr string) bool {
	for _, v := range _ignoreTypes {
		if strings.Contains(typeStr, v) {
			return true
		}
	}
	return false
}

//获取字段类型字符串
func getFieldTypeStr(v interface{}, key string, searchSub bool) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ""
	}

	var e reflect.Value
	if (searchSub) {
		e = reflect.ValueOf(v).Elem()
	}

	fieldNum := t.NumField()
	//fmt.Println("----fieldNum:", fieldNum)
	for i := 0; i < fieldNum; i++ {
		//fmt.Println("----", i, t.Field(i).Name, t.Field(i).Type.String())
		if t.Field(i).Type.Kind() == reflect.Struct {
			//fmt.Println("--is struct:", t.Field(i).Name)

			if !_isInIgnoreList(t.Field(i).Type.String()) {
				if (searchSub) {
					subStruct := e.FieldByName(t.Field(i).Name).Interface()
					subTypeStr := getFieldTypeStr(subStruct, key, false)

					//fmt.Println("--sub type str:", subTypeStr)
					if len(subTypeStr) > 0 {
						return subTypeStr
					}
				}
			}
		}

		if strings.EqualFold(key, t.Field(i).Name) {
			return t.Field(i).Type.String()
		}
	}
	return ""
}

//处理列表查询筛选条件
func handleListFilters(tx *gorm.DB, v *CRUD, filterMap map[string]string, request *Request_List) (*gorm.DB) {
	tmp := v.GetNewModelDb()
	for k, v := range filterMap {
		if len(v) > 0 {
			keyType := getFieldTypeStr(tmp, k, true)
			condition := ""

			if len(keyType) > 0 {
				k := gorm.ToColumnName(k)
				switch keyType {
				case "string":
					condition = fmt.Sprintf("%s LIKE '%%%s%%' ", k, v)
				case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64", "float32", "float64", "bool":
					condition = fmt.Sprintf(" %s=%s ", k, v)
				}
			}
			common.LogConsole.Debug("--[0]-query list add filter key:%s,value:%s, keyType=%s", k, v, keyType)
			if len(condition) > 0 {
				tx = tx.Where(condition)
			}
		}
	}
	if len(request.IdList) > 0 {
		tx = tx.Where(fmt.Sprintf("id IN (%s)", request.IdList))
	}
	return tx
}

//查询列表
func onCrudSelectListHandler(c *gin.Context) {
	response := models.Response_List{}

	requestList := &Request_List{}
	c.ShouldBindWith(requestList, MyGinQueryBinding)
	common.LogConsole.Debug("----request list:%#v", requestList)

	if v, ok, err := _checkOpType(requestList.OpType); !ok {
		response.Code = 1
		response.Msg = err.Error()
	} else {
		//var querySeter orm.QuerySeter
		//querySeter = lkpush.MOrmer.QueryTable(v.Table)

		//根据filter查询
		filterMap := make(map[string]string)
		json.Unmarshal([]byte(requestList.Filters), &filterMap)

		tx := _db.Table(v.Table)
		tx = handleListFilters(tx, v, filterMap, requestList)

		//查询总数
		total := 0
		tx.Count(&total)

		tx = _db.Table(v.Table)
		tx = handleListFilters(tx, v, filterMap, requestList)

		//排序
		//规则：有外部传参，则忽略默认值
		if len(requestList.Order) > 0 {
			tx = tx.Order(requestList.Order)
		} else if len(v.orderBy) > 0 {
			for _, v := range v.orderBy {
				tx = tx.Order(v)
			}
		}
		tx = tx.Limit(requestList.PageSize).Offset(requestList.PageSize * (requestList.Page - 1))

		//返回结构注入
		resultList := v.GetNewModelResponseItemSlice(requestList)
		fmt.Printf("====result list:%#v, %d", resultList, requestList.PageSize)
		tx.Find(resultList)

		response.Total = int(total)
		response.List = resultList
	}
	c.String(http.StatusOK, common.UtilCommon.Parse_ToJson(&response))
}
