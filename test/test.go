package test

import (
	"encoding/json"
	"fmt"
)

type IT struct {
	Company  string
	Subjects []string
	IsOk     bool
	Price    float64
}

func main() {
	t1 := IT{"郭金龙", []string{"Go", "C++", "Python", "Test"}, true, 666.666}

	//生成一段JSON格式的文本
	//如果编码成功， err 将赋于零值 nil，变量b 将会是一个进行JSON格式化之后的[]byte类型
	//b, err := json.Marshal(t1)
	//输出结果：{"Company":"itcast","Subjects":["Go","C++","Python","Test"],"IsOk":true,"Price":666.666}

	b, err := json.Marshal(t1)

	var v interface{}
	json.Unmarshal(b, &v)

	//println(b)
	msg := v.(map[string]interface{})
	fmt.Println(msg["Company"])
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Println(msg["IsOK"] == true)
}
