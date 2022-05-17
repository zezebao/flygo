package baidu_ai

import (
	"fmt"
	"os"
	"com.ffl/common"
	"github.com/zezebao/flygo/net"
	"encoding/json"
	"time"
	"github.com/zezebao/flygo/database/leveldb"
	"github.com/gin-gonic/gin"
	"net/http"
	"bytes"
	"io/ioutil"
)

const (
	KEY_TOKEN        = "baidu_ai_token"
	KEY_TOKEN_EXPIRE = "baidu_ai_expire"
)

type configBaiduAI struct {
	AppID     string
	APIKey    string
	SecretKey string
}

type BaiduAssessTokenResponse struct {
	RefreshToken  string `json:"refresh_token"`
	ExpiresIn     int    `json:"expires_in"` // Access Token的有效期(秒为单位，有效期30天)
	SessionKey    string `json:"session_key"`
	AccessToken   string `json:"access_token"`
	Scope         string `json:"scope"`
	SessionSecret string `json:"session_secret"`
}

var ConfigBaiduAI *configBaiduAI

func InitBaiduAI(ginEngine *gin.Engine) {
	fmt.Println("--dir:", os.Args[0])

	config := new(configBaiduAI)
	err := common.ConfToml.ReadConf("conf/conf_baidu_ai.toml", config)
	if err != nil {
		common.LogErr.Debug("conf配置文件读取错误：%v", err)
		ConfigBaiduAI = new(configBaiduAI)
		return
	} else {
		ConfigBaiduAI = config
	}
	common.LogConsole.Debug("ConfigBaiduAI %#v", ConfigBaiduAI)

	ginEngine.GET("/api/ai/baidu/token", func(c *gin.Context) {
		res := &common.Response_Common{}
		res.Data = GetAssessToken()
		c.JSON(http.StatusOK, res)
	})

	ginEngine.POST("/baidu/face/merge", func(c *gin.Context) {
		res := &common.Response_Common{}

		target := c.PostForm("target")
		template := c.PostForm("template")

		sendData := make(map[string]interface{})
		sendData["version"] = "2.0"

		image_target := make(map[string]interface{})
		image_target["image_type"] = "BASE64"
		image_target["image"] = target
		image_template := make(map[string]interface{})
		image_template["image_type"] = "URL"
		image_template["image"] = template

		sendData["image_target"] = image_target
		sendData["image_template"] = image_template

		bytesData, err := json.Marshal(sendData)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		reader := bytes.NewReader(bytesData)
		url := "https://aip.baidubce.com/rest/2.0/face/v1/merge?access_token=" + GetAssessToken()
		request, err := http.NewRequest("POST", url, reader)
		if err != nil {
			fmt.Println(err.Error())
			res.Code = 1
			res.Msg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
		request.Header.Set("Content-Type", "application/json;charset=UTF-8")
		client := http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			fmt.Println(err.Error())
			res.Code = 1
			res.Msg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}

		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			res.Code = 1
			res.Msg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}

		res.Data = string(respBytes)
		c.JSON(http.StatusOK, res)
	})
}

func GetAssessToken() string {
	now := int(time.Now().Unix())
	expire := leveldb.GetInt(KEY_TOKEN_EXPIRE, -1)
	if expire > 0 && now < expire {
		return leveldb.GetString(KEY_TOKEN, "")
	} else {
		url := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", ConfigBaiduAI.APIKey, ConfigBaiduAI.SecretKey)
		result := net.Get(url)

		rsp := &BaiduAssessTokenResponse{}
		json.Unmarshal([]byte(result), rsp)

		if len(rsp.AccessToken) > 0 {
			leveldb.SetInt(KEY_TOKEN_EXPIRE, now+rsp.ExpiresIn)
			leveldb.SetString(KEY_TOKEN, rsp.AccessToken)
		} else {
			fmt.Println("-百度AI token获取异常:", result)
		}
		return rsp.AccessToken
	}
}