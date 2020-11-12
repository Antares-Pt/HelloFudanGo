package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

const (
	loginURL      = "https://uis.fudan.edu.cn/authserver/login"
	dailyFudanURL = "https://zlapp.fudan.edu.cn/site/ncov/fudanDaily"
	userAgent     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3776.0 Safari/537.36"
)

// Student struct
type Student struct {
	StudentID string
	Password  string
}

// HelloFudan struct
type HelloFudan struct {
	client *http.Client
	stu    Student
}

func main() {
	hf := NewHelloFudan(Student{
		StudentID: "54321",
		Password:  "12345",
	})

	hf.Login()
}

// NewHelloFudan is HelloFudan Constructor
func NewHelloFudan(stu Student) *HelloFudan {
	return &HelloFudan{
		client: &http.Client{},
		stu:    stu,
	}
}

// Login HelloFudan
func (hf *HelloFudan) Login() {
	log.Println("---Start login---")
	data := url.Values{
		"username": {hf.stu.StudentID},
		"password": {hf.stu.Password},
		"service":  {dailyFudanURL},
	}
	// add token
	html := hf.getHTML(loginURL)
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
	req.Header.Set("Referer", loginURL)
	req.Header.Set("User-Agent", userAgent)

	resp, err := hf.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// hf.client.Jar.SetCookies("uis.fudan.edu.cn", resp.Cookies())
}

// Logout HelloFudan
func (hf *HelloFudan) Logout() {

}

// Check HelloFudan status
func (hf *HelloFudan) Check() {

}

func (hf *HelloFudan) getHTML(url string) string {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", userAgent)
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
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
