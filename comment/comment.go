package comment

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"stones/login"
	"stones/tempsuid"
	"strconv"
	"strings"
)

func Comment(user *login.UserData) error {
	type resultBody struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	re := new(resultBody)

	urlLogin := `https://apiff14risingstones.web.sdo.com/api/home/posts/comment?tempsuid=`
	tempsUid, err := tempsuid.Get()
	if err != nil {
		return err
	}
	urlLogin = urlLogin + tempsUid

	form := url.Values{}
	emoji := rand.Intn(29)
	form.Add("content", `<p><span class="at-emo">[emo`+strconv.Itoa(emoji)+`]</span>&nbsp;</p>`)
	form.Add("posts_id", "9365")
	form.Add("parent_id", "0")
	form.Add("root_parent", "0")
	form.Add("comment_pic", "")

	req, err := http.NewRequest("POST", urlLogin, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
