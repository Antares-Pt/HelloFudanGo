// Package hellofudan client
package hellofudan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

const (
	loginURL         = "https://uis.fudan.edu.cn/authserver/login"
	loginRedirectURL = "http://uis.fudan.edu.cn/authserver/index.do"
	logoutURL        = "https://uis.fudan.edu.cn/authserver/logout?service=/authserver/login"
	dailyFudanURL    = "https://zlapp.fudan.edu.cn/site/ncov/fudanDaily"
	saveURL          = "https://zlapp.fudan.edu.cn/ncov/wap/fudan/save"
	checkURL         = "https://zlapp.fudan.edu.cn/ncov/wap/fudan/get-info"
	userAgent        = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Safari/537.36 Edg/86.0.622.63"
)

// Student struct
type Student struct {
	StudentID string
	Password  string
}

// HelloFudan struct
type HelloFudan struct {
	log    *log.Logger
	stu    Student
	client *http.Client
	info   map[string]interface{}
}

// NewHelloFudan is HelloFudan Constructor
func newHelloFudan(student Student) *HelloFudan {
	jar, _ := cookiejar.New(nil)
	return &HelloFudan{
		log: log.New(os.Stdout, fmt.Sprintf("[StudentsID: %s] ", student.StudentID), log.Ldate|log.Ltime),
		stu: student,
		client: &http.Client{
			Timeout: 5 * time.Second,
			Jar:     jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if strings.Compare(req.URL.String(), loginRedirectURL) == 0 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

func (hf *HelloFudan) initLogin() string {
	req, _ := http.NewRequest("GET", loginURL, nil)
	req.Header.Add("User-Agent", userAgent)

	resp, err := hf.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil && data == nil {
		panic(err)
	}

	return fmt.Sprintf("%s", data)
}

func (hf *HelloFudan) login() {
	hf.log.Println("Start login")
	data := url.Values{
		"username": {hf.stu.StudentID},
		"password": {hf.stu.Password},
		"service":  {dailyFudanURL},
	}
	// add token
	html := hf.initLogin()
	root, _ := htmlquery.Parse(strings.NewReader(html))
	nodes := htmlquery.Find(root, "/html/body/form/input[@type]")
	for _, node := range nodes {
		key := htmlquery.SelectAttr(node, "name")
		value := htmlquery.SelectAttr(node, "value")
		data.Add(key, value)
	}

	body := strings.NewReader(data.Encode())
	req, _ := http.NewRequest("POST", loginURL, body)

	req.Header.Set("Host", "uis.fudan.edu.cn")
	req.Header.Set("Origin", "https://uis.fudan.edu.cn")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", loginURL)
	req.Header.Set("User-Agent", userAgent)

	resp, err := hf.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		hf.log.Println("Login success")
	} else {
		hf.log.Printf("Login failed, Status: %s, Please check your account", resp.Status)
		panic("Login failed")
	}
}

func (hf *HelloFudan) logout() {
	hf.log.Println("Start logout")
	resp, err := hf.client.Get(logoutURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	cookies := resp.Header.Get("Set-Cookie")

	if strings.Contains(cookies, "01-Jan-1970") {
		hf.log.Println("Logout success")
	} else {
		hf.log.Println("Logout failed")
		panic("Logout failed")
	}
}

func (hf *HelloFudan) checkStatus() bool {
	hf.log.Println("Start check status")

	resp, err := hf.client.Get(checkURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if len(body) == 0 {
		panic("Check body is nil")
	}

	var v interface{}
	json.Unmarshal(body, &v)
	info := v.(map[string]interface{})["d"].(map[string]interface{})["info"].(map[string]interface{})

	hf.info = info
	date := info["date"].(string)
	address := info["address"].(string)

	hf.log.Printf("Last check in position: %s???date: %s", address, date)

	cstSH, _ := time.LoadLocation("Asia/Shanghai")
	// cstE8 := time.FixedZone("CST", 8*3600)
	today := time.Now().In(cstSH).Format("20060102")

	hf.log.Printf("Today is %s", today)
	if strings.Compare(date, today) == 0 {
		hf.log.Println("Today already checked in")
		return true
	}
	hf.log.Println("Today haven't checked in")
	return false
}

func (hf *HelloFudan) checkIn() {
	hf.log.Println("Start check in")

	geoInfo := make(map[string]interface{})
	json.Unmarshal([]byte(hf.info["geo_api_info"].(string)), &geoInfo)
	addr := geoInfo["addressComponent"].(map[string]interface{})

	province, has := addr["province"].(string)
	if !has {
		province = ""
	}
	city, has := addr["city"].(string)
	if !has {
		city = ""
	}
	district, has := addr["district"].(string)
	if !has {
		district = ""
	}
	locArr := []string{province, city, district}
	for i, loc := range locArr {
		if len(loc) == 0 {
			locArr = append(locArr[:i], locArr[i+1:]...)
		}
	}
	area := strings.Join(locArr, " ")
	hf.log.Printf("Check in position: %s", area)

	data := url.Values{
		"tw":       {"13"},
		"province": {province},
		"city":     {city},
		"area":     {area},
	}
	// convert map[string]interface{} to url.Values
	for k, v := range hf.info {
		data.Add(k, fmt.Sprint(v))
	}

	body := strings.NewReader(data.Encode())
	req, _ := http.NewRequest("POST", saveURL, body)

	req.Header.Set("Host", "zlapp.fudan.edu.cn")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://zlapp.fudan.edu.cn")
	req.Header.Set("Referer", dailyFudanURL+"?from=history")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("DNT", "1")
	req.Header.Set("TE", "Trailers")

	resp, err := hf.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		ret := make(map[string]interface{})
		json.Unmarshal(data, &ret)
		msg := ret["m"].(string)
		hf.log.Printf("Check in message: %s", msg)
	} else {
		hf.log.Printf("Check in failed: %s", resp.Status)
	}
}
