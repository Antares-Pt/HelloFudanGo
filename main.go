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

// DailyFudan struct
type DailyFudan struct {
	client *http.Client
	stu    Student
}

func main() {

}

// NewDailyFudan is DailyFudan Constructor
func NewDailyFudan(stu Student) *DailyFudan {
	return &DailyFudan{
		client: &http.Client{},
		stu:    stu,
	}
}

// Login DailyFudan
func (df *DailyFudan) Login() {
	data := url.Values{
		"username": {df.stu.StudentID},
		"password": {df.stu.Password},
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

	resp, err := df.client.Do(req)
	if err != nil {
		panic(err)
	}

	cookie := resp.Cookies()
}

// Logout DailyFudan
func (df *DailyFudan) Logout() {

}

// Check DailyFudan status
func (df *DailyFudan) Check() {

}
