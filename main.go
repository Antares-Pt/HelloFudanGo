package main

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	loginURL      = "https://uis.fudan.edu.cn/authserver/login"
	dailyFudanURL = "https://zlapp.fudan.edu.cn/site/ncov/fudanDaily"
	userAgent     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:76.0) Gecko/20100101 Firefox/76.0"
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
	data := url.Values{
		"username": {hf.stu.StudentID},
		"password": {hf.stu.Password},
		"service":  {dailyFudanURL},
	}
	body := strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", loginURL, body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Host", "uis.fudan.edu.cn")
	req.Header.Set("Origin", "https://uis.fudan.edu.cn")
	req.Header.Set("Referer", loginURL)
	req.Header.Set("User-Agent", userAgent)

	resp, err := hf.client.Do(req)
	if err != nil {
		panic(err)
	}

	hf.client.Jar.SetCookies("uis.fudan.edu.cn", resp.Cookies())
}

// Logout HelloFudan
func (hf *HelloFudan) Logout() {

}

// Check HelloFudan status
func (hf *HelloFudan) Check() {

}
