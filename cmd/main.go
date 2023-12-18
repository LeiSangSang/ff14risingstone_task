package main

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"stones/login"
	"stones/post"
	"strings"
	"time"
)

func main() {

	User := login.NewUser()

	err := User.Login()
	if err != nil {
		fmt.Println(err)
		return
	}

	if User.GetUserInfo().IsSign == 0 {
		err = sign(User)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("今日已签到,跳过签到")
	}

	fmt.Println("----开始点赞----")
	sum := 0
	page := 1
	for {
		time.Sleep(time.Second)
		list, err := post.GetPostList(User, page)
		if err != nil {
			fmt.Println(err)
		}
		likeNum, err := list.Like()
		if err != nil {
			fmt.Println(err)
		}
		sum += likeNum
		if sum >= 5 {
			fmt.Println("----点赞完成----")
			break
		} else {
			page++
		}
	}

	fmt.Println("----开始发动态----")
	err = create(User)

	if err != nil {
		fmt.Println(err)
	}

	for {
		
	}

}

func sign(user *login.UserData) error {
	signPath := `https://apiff14risingstones.web.sdo.com/api/home/sign/mySignLog?month=2023-12`
	req, err := http.NewRequest("GET", signPath, nil)
	if err != nil {
		return err
	}
	resp, err := user.GetClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("签到成功!")
	return nil
}

func create(user *login.UserData) error {
	type resultBody struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	re := new(resultBody)

	urlLogin := `https://apiff14risingstones.web.sdo.com/api/home/dynamic/create`

	form := url.Values{}
	form.Add("content", `<p>丝瓜的任务罢了</p>`)
	form.Add("scope", "1")
	form.Add("pic_url", "")

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
