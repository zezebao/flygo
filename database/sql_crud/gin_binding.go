package sql_crud

import (
	"net/http"
	"github.com/gin-gonic/gin/binding"
	"strings"
)

var MyGinQueryBinding = myGinQueryBinding{}

type myGinQueryBinding struct{}

func (myGinQueryBinding) Name() string {
	return "my_query"
}

func (myGinQueryBinding) Bind(req *http.Request, obj interface{}) error {
	tmp := req.URL.Query()
	var tmpValues  map[string][]string = tmp
	values  := make(map[string][]string)
	for k,v := range  tmpValues {
		var key string = strings.ToUpper(k[0:1])+k[1:]
		//gorm.ToColumnName()
		values[key] = v
	}

	if err := mapForm(obj, values); err != nil {
		return err
	}
	return validate(obj)
}

func validate(obj interface{}) error {
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}
