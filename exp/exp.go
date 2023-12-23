package exp

import (
	"fmt"
	"stones/comment"
	"stones/dynamic"
	"stones/login"
	"strconv"
	"time"
)

func Exp(user *login.UserData) error {

	fmt.Println("\n----开始获取经验----")

	fmt.Println("\n----发动态五条并删除----")

	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		id, err := dynamic.Dynamic(user)
		if err != nil {
			return err
		}
		if id != 0 {
			fmt.Println("发动态成功,id为", strconv.Itoa(id))
			time.Sleep(5 * time.Second)
			fmt.Println("删除id为", strconv.Itoa(id), "的动态")
			err = dynamic.DelDynamic(user, id)
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("\n----发评论五条----")
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		err := comment.Comment(user)
		if err != nil {
			return err
		}
	}

	fmt.Println("----获取经验结束,按任意键退出----")
	return nil
}
