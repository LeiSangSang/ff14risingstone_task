package post

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"net/url"
	"stones/login"
	"strconv"
	"strings"
	"time"
)

type P struct {
	PostsId string `json:"posts_id"`
	Title   string `json:"title"`
	IsLike  int    `json:"is_like"`
}

type List struct {
	client *http.Client
	Data   []P `json:"data"`
}

func GetPostList(User *login.UserData, page int) (*List, error) {
	list := new(List)
	list.client = User.GetClient()
	urlLogin := `https://apiff14risingstones.web.sdo.com/api/home/posts/postsList?type=1&is_top=0&is_refine=0&part_id=&hotType=postsHotNow&order=&page=` + strconv.Itoa(page) + `&limit=15`
	req, err := http.NewRequest("GET", urlLogin, nil)
	if err != nil {
		return nil, err
	}
	resp, err := User.GetClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type resultBody struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Rows     []P    `json:"rows"`
			PageTime string `json:"pageTime"`
		} `json:"data"`
	}

	result = bytes.TrimSpace(result)
	re := new(resultBody)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(result, re)
	if err != nil {
		return nil, err
	}
	list.Data = re.Data.Rows

	return list, nil

}

func (t *List) Like() (int, error) {
	likeNum := 0
	for k, _ := range t.Data {

		if t.Data[k].Like(t.client) {
			likeNum++
		}

		if likeNum >= 5 {
			break
		}
	}

	return likeNum, nil
}

func (t *P) Like(client *http.Client) bool {

	if t.IsLike == 1 {
		fmt.Println(t.Title, "--已点赞跳过")
		return false
	}

	for {
		code, msg, err := t.l(client)
		if err != nil {
			fmt.Println(err)
			return false
		}
		switch code {
		case 10000:
			fmt.Println(t.Title, `----`, msg)
			return true
		case 10301:
			fmt.Println(t.Title, `----`, msg, `----等待3s后重试`)
			break
		default:
			fmt.Println(t.Title, `----`, msg)
			return false
		}
		time.Sleep(3 * time.Second)
	}
}

func (t *P) l(client *http.Client) (int, string, error) {

	type resultBody struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data int    `json:"data"`
	}
	re := new(resultBody)

	urlLogin := `https://apiff14risingstones.web.sdo.com/api/home/posts/like`

	form := url.Values{}
	form.Add("id", t.PostsId)
	form.Add("type", "1")

	req, err := http.NewRequest("POST", urlLogin, strings.NewReader(form.Encode()))
	if err != nil {
		return 0, ``, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return 0, ``, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, ``, err
	}
	result = bytes.TrimSpace(result)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(result, re)
	if err != nil {
		return 0, ``, err
	}
	return re.Code, re.Msg, nil
}
