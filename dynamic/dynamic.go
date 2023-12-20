package dynamic

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"net/url"
	"stones/login"
	"stones/tempsuid"
	"strconv"
	"strings"
)

func Dynamic(user *login.UserData) (int, error) {
	type resultBody struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Id int `json:"id"`
		} `json:"data"`
	}
	re := new(resultBody)

	urlLogin := `https://apiff14risingstones.web.sdo.com/api/home/dynamic/create?tempsuid=`
	tempsUid, err := tempsuid.Get()
	if err != nil {
		return 0, err
	}
	urlLogin = urlLogin + tempsUid
	form := url.Values{}
	form.Add("content", `<p>丝瓜的任务罢了</p>`)
	form.Add("scope", "1")
	form.Add("pic_url", "")

	req, err := http.NewRequest("POST", urlLogin, strings.NewReader(form.Encode()))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	resp, err := user.GetClient().Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	result = bytes.TrimSpace(result)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(result, re)
	if err != nil {
		return 0, err
	}

	fmt.Println(re.Msg)
	if re.Code == 10000 {
		return re.Data.Id, nil
	} else {
		return 0, nil
	}

}

func DelDynamic(user *login.UserData, id int) error {
	type resultBody struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	re := new(resultBody)

	urlLogin := `https://apiff14risingstones.web.sdo.com/api/home/dynamic/deleteDynamic?tempsuid=`
	tempsUid, err := tempsuid.Get()
	if err != nil {
		return err
	}
	urlLogin = urlLogin + tempsUid

	param := `{"dynamic_id":` + strconv.Itoa(id) + `}`

	req, err := http.NewRequest("DELETE", urlLogin, strings.NewReader(param))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	resp, err := user.GetClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	result = bytes.TrimSpace(result)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(result, re)
	if err != nil {
		return err
	}
	fmt.Println(re.Msg)
	return nil
}
