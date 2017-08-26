package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"net/url"
)

func httpPostForm() {
	//resp, err := http.PostForm("http://www.01happy.com/demo/accept.php", url.Values{"key": {"Value"}, "id": {"123"}})
	resp, err := http.PostForm("http://192.168.3.243:8080/post_test", url.Values{"key": {"test"}, "id": {"123"}})
	if err != nil{
		fmt.Printf("err:%s", err)
	}
	
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Printf("err:%s", err)
	}
	fmt.Println(string(body))
}

func main(){
	httpPostForm()
}