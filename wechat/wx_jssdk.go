package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"io/ioutil"

	"com.ffl/common"

	"github.com/gin-gonic/gin"
	"github.com/zezebao/flygo/logger"
)

//token与jsticket返回数据示例
//{"access_token":"9__14xOdEeztzcwgawDf3vWFvbkZgWf76DL83Vz0hI_y2J75DQxJV9WWl0r9Ie4uST9nGExapH41jTv4Jcjpav1kUbgmczngzJOHiy8PSt9IQ6GhAq7MsxsdrZlgoEzxaW21bIIAvnCduc4cgfARQfCEAGXW","expires_in":7200}
//{"errcode":0,"errmsg":"ok","ticket":"sM4AOVdWfPE4DxkXGEs8VEd6QHyc9Pkd9zArjGYTkVN-NptRJZ209rEkrOS3rPs19XvNNnI3DL62cKMP7Wf91A","expires_in":7200}

const NEXT_DELTA = 6000 //下一次访问间隔

type WxJsSdkSign struct {
	Noncestr  string `json:"noncestr"`
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
}

//微信SDK处理
type WxJsSdk struct {
	AppId         string `json:"appId"`
	AppSecret     string `json:"appSecret"`
	Token         string `json:"token"`
	JsTiket       string `json:"jsTiket"`
	LastTokenTime int64  `json:"lastTokenTime"`
	LastTiketTime int64  `json:"lastTiketTime"`
}

//初始化
func (sdk *WxJsSdk) Init() {
	sdk.LastTiketTime = -1
	sdk.LastTokenTime = -1
	sdk.GetToken()
	sdk.GetTiket()

	if nil != Module.ModuleConfig {
		Module.ModuleConfig.GinEngine.GET(API_GET_JS_SIGN, sdk.OnGetJsSignHandler)
	} else {
		fmt.Printf("##WARN-NO ModeleConfig::wx\n")
	}
}

//获取token
func (sdk *WxJsSdk) GetToken() string {

	if len(sdk.Token) > 0 {
		now := time.Now().Unix()
		delta := now - sdk.LastTokenTime

		if delta < NEXT_DELTA {
			return sdk.Token
		}
	}

	type Tmp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	tokenUrl := fmt.Sprintf(`https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s`, sdk.AppId, sdk.AppSecret)
	result := common.UtilCommon.HttpGet(tokenUrl)

	logger.Main.Debugf("---token url:%s", tokenUrl)
	logger.Main.Debugf("---token result:%s", result)

	if len(result) > 0 {
		tmp := &Tmp{}
		json.Unmarshal([]byte(result), tmp)

		sdk.Token = tmp.AccessToken
		sdk.LastTokenTime = time.Now().Unix()

		fmt.Println("====WX JSSDK GetToken():%s", sdk.Token)
	}

	return sdk.Token
}

//获取js票据
func (sdk *WxJsSdk) GetTiket() string {
	if len(sdk.JsTiket) > 0 {
		now := time.Now().Unix()
		delta := now - sdk.LastTiketTime

		if delta < NEXT_DELTA {
			return sdk.JsTiket
		}
	}

	type Tmp struct {
		Ticket    string `json:"ticket"`
		ExpiresIn int    `json:"expires_in"`
	}
	jsTiketUrl := fmt.Sprintf(`https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi`, sdk.GetToken())
	result := common.UtilCommon.HttpGet(jsTiketUrl)

	common.LogConsole.Debug("----js ticket url:%s", jsTiketUrl)
	common.LogConsole.Debug("----js ticket result:%s", result)

	if len(result) > 0 {
		tmp := &Tmp{}
		json.Unmarshal([]byte(result), tmp)

		sdk.JsTiket = tmp.Ticket
		sdk.LastTiketTime = time.Now().Unix()

		fmt.Println("====WX JSSDK GetTiket():%s", sdk.JsTiket)
	}
	return sdk.JsTiket
}

//计算JSSDK签名
func (sdk *WxJsSdk) GetJsSign(url string) *WxJsSdkSign {
	var tmp = &WxJsSdkSign{}
	tmp.Noncestr = common.UtilCommon.GetRandomString(16)
	tmp.Timestamp = time.Now().Unix()

	tmp1 := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", sdk.GetTiket(), tmp.Noncestr, tmp.Timestamp, url)
	tmp.Sign = common.Crypt.Sha1(tmp1)

	return tmp
}

//获取js签名
func (sdk *WxJsSdk) OnGetJsSignHandler(c *gin.Context) {
	response := common.Response_Common{}

	url := c.Query("url")
	signData := sdk.GetJsSign(url)

	response.Data = signData

	c.String(http.StatusOK, common.UtilCommon.Parse_ToJson(&response))
}

func (sdk *WxJsSdk) HttpGetOpenId(code string) string {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", sdk.AppId, sdk.AppSecret, code)
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		common.LogErr.Debug("----opid err:%s", err.Error())
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		common.LogErr.Debug("----opid body err:%s", err.Error())
		return ""
	}
	return string(body)
}

//新建JSSDK
func NewWxJsSdk(appid string, appSecret string) *WxJsSdk {
	var wxJsSdk *WxJsSdk = &WxJsSdk{AppId: appid, AppSecret: appSecret}
	wxJsSdk.Init()
	return wxJsSdk
}
