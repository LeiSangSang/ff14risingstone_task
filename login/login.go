package login

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	qr "github.com/skip2/go-qrcode"
	"image"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"stones/tempsuid"
	"strconv"
	"strings"
	"time"
)

type LoginInfo struct {
	DisplayAccount  string        `json:"displayAccount"`
	CharacterName   string        `json:"character_name"`
	AreaName        string        `json:"area_name"`
	GroupName       string        `json:"group_name"`
	IsActivateUser  int           `json:"isActivateUser"`
	LastLoginTime   string        `json:"lastLoginTime"`
	PunishStatusArr []interface{} `json:"punishStatusArr"`
}

type UserInfo struct {
	Uuid          string `json:"uuid"`
	CharacterId   string `json:"character_id"`
	CharacterName string `json:"character_name"`
	AreaId        int    `json:"area_id"`
	AreaName      string `json:"area_name"`
	GroupId       int    `json:"group_id"`
	GroupName     string `json:"group_name"`
	Profile       string `json:"profile"`
	Experience    string `json:"experience"`
	Status        int    `json:"status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	UpdatedBy     string `json:"updated_by"`
	LastLoginIp   string `json:"last_login_ip"`
	IpLocation    string `json:"ip_location"`
	LastLoginTime string `json:"last_login_time"`
	IsSign        int    `json:"isSign"`
}

type UserData struct {
	client *http.Client
	ticket string
	info   LoginInfo
	user   UserInfo
}

func NewUser() *UserData {
	gCookie, _ := cookiejar.New(nil)
	client := &http.Client{Jar: gCookie}
	return &UserData{
		client: client,
	}
}

func (t *UserData) Login() error {

	err := t.getLoginQRCode()
	if err != nil {
		return err
	}

	err = t.waitScanQRCode()
	if err != nil {
		return err
	}

	err = t.riSingStonesLogin()
	if err != nil {
		return err
	}

	err = t.getLoginInfo()
	if err != nil {
		return err
	}

	err = t.getCharacterBindInfo()
	if err != nil {
		return err
	}

	return nil
}

func (t *UserData) GetClient() *http.Client {
	return t.client
}

func (t *UserData) GetUserInfo() UserInfo {
	return t.user
}

func (t *UserData) getLoginQRCode() error {
	req, err := http.NewRequest("GET", "https://w.cas.sdo.com/authen/getcodekey.jsonp?maxsize=145&appId=6788&areaId=1&r=0.5361605586443803", nil)
	if err != nil {
		return err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	m, _, err := image.Decode(resp.Body)
	if err != nil {
		return err
	}
	bmp, err := gozxing.NewBinaryBitmapFromImage(m)
	if err != nil {
		return err
	}
	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return err
	}
	qrc, err := qr.New(result.GetText(), qr.Low)
	if err != nil {
		return err
	}
	ascii := qrc.ToSmallString(false)
	fmt.Println(ascii)

	fmt.Println("\n如果二维码显示异常,请尝试右键修改控制台字体")
	fmt.Println("github: https://github.com/LeiSangSang/ff14risingstone_task")
	return nil
}

func (t *UserData) waitScanQRCode() error {
	for {
		time.Sleep(time.Second)
		urlPath := `https://w.cas.sdo.com/authen/codeKeyLogin.jsonp?`
		param := `callback=codeKeyLogin_JSONPMethod&appId=6788&areaId=1&code=300&serviceUrl=http%3A%2F%2Fapiff14risingstones.web.sdo.com%2Fapi%2Fhome%2FGHome%2Flogin%3FredirectUrl%3Dhttps%3A%2F%2Fff14risingstones.web.sdo.com%2Fpc%2Findex.html&productId=2&productVersion=3.1.0&authenSource=2&_=`
		currentTime := time.Now().UnixNano()
		milliTime := currentTime / int64(time.Millisecond)
		var numStr string = strconv.FormatInt(milliTime, 10)
		param = param + numStr
		urlP := urlPath + param
		req, err := http.NewRequest("GET", urlP, nil)
		if err != nil {
			return err
		}
		resp, err := t.client.Do(req)
		if err != nil {
			return err
		}
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		res := string(result)

		idx := strings.Index(res, `"ticket":`)
		if idx != -1 {
			t.ticket = res[idx+11 : len(res)-6]
			fmt.Println("登录成功")
			break
		}
		resp.Body.Close()
	}
	return nil
}

func (t *UserData) riSingStonesLogin() error {
	urlLogin := `http://apiff14risingstones.web.sdo.com/api/home/GHome/login?redirectUrl=https://ff14risingstones.web.sdo.com/pc/index.html&ticket=` + t.ticket
	req, err := http.NewRequest("GET", urlLogin, nil)
	if err != nil {
		return err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (t *UserData) getLoginInfo() error {
	urlLogin := `https://apiff14risingstones.web.sdo.com/api/home/GHome/isLogin?tempsuid=`
	tempsUid, err := tempsuid.Get()
	if err != nil {
		return err
	}
	urlLogin = urlLogin + tempsUid
	req, err := http.NewRequest("GET", urlLogin, nil)
	if err != nil {
		return err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type ResultBody struct {
		Code int       `json:"code"`
		Msg  string    `json:"msg"`
		Data LoginInfo `json:"data"`
	}
	result = bytes.TrimSpace(result)
	re := new(ResultBody)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(result, re)
	if err != nil {
		return err
	}
	t.info = re.Data
	fmt.Println(t.info.AreaName, t.info.GroupName, t.info.CharacterName)
	return nil
}

func (t *UserData) getCharacterBindInfo() error {
	urlLogin := `https://apiff14risingstones.web.sdo.com/api/home/groupAndRole/getCharacterBindInfo?platform=1&tempsuid=`
	tempsUid, err := tempsuid.Get()
	if err != nil {
		return err
	}
	urlLogin = urlLogin + tempsUid
	req, err := http.NewRequest("GET", urlLogin, nil)
	if err != nil {
		return err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type ResultBody struct {
		Code int      `json:"code"`
		Msg  string   `json:"msg"`
		Data UserInfo `json:"data"`
	}

	result = bytes.TrimSpace(result)
	re := new(ResultBody)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(result, re)
	if err != nil {
		return err
	}
	t.user = re.Data
	return nil
}
