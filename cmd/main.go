package main

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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

	err = sign(User)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\n----开始点赞----")
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

	time.Sleep(time.Second)
	fmt.Println("\n----开始发评论----")
	err = comment(User)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("----评论完成----")

	fmt.Println("\n----将在5秒后开始发动态----")
	go func() {
		time.Sleep(5 * time.Second)
		err = dynamic(User)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	fmt.Print("\n如果不需要发动态,请按任意键退出")
	b := make([]byte, 1)
	os.Stdin.Read(b)

}

func sign(user *login.UserData) error {
	signPath := `https://apiff14risingstones.web.sdo.com/api/home/sign/signIn`
	req, err := http.NewRequest("POST", signPath, nil)
	if err != nil {
		return err
	}
	resp, err := user.GetClient().Do(req)
	if err != nil {
		return err
	}
	type resultBody struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			SqMsg          string `json:"sqMsg"`
			ContinuousDays int    `json:"continuousDays"`
			TotalDays      string `json:"totalDays"`
			SqExp          int    `json:"sqExp"`
			ShopExp        int    `json:"shopExp"`
		} `json:"data"`
	}
	re := new(resultBody)
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
	if re.Code == 10000 {
		fmt.Println(re.Data.SqMsg)
	}
	return nil
}

func dynamic(user *login.UserData) error {
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

	fmt.Println("\n", re.Msg)
	return nil
}

func comment(user *login.UserData) error {
	type resultBody struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	re := new(resultBody)

	urlLogin := `https://apiff14risingstones.web.sdo.com/api/home/posts/comment`

	form := url.Values{}
	form.Add("content", `<p><span class="at-emo">[emo2]</span>&nbsp;</p>`)
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
