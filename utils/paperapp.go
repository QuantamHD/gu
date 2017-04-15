package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"log"

	"github.com/PuerkitoBio/goquery"
)

const baseURL string = "https://paper-app.gonzaga.edu:9192"

/*PaperCutCredentials ...
parameters
	username - the username of the account
	password - the password to papercut
*/
type PaperCutCredentials struct {
	username   string
	password   string
	sessionID  string
	isLoggedIn bool
}

func (p PaperCutCredentials) GetSessionID() string {
	return p.sessionID
}

func CreatePaperCutCredentials(username string, password string) *PaperCutCredentials {
	credentials := PaperCutCredentials{username, password, "", false}
	login(&credentials)
	return &credentials
}

func intitalConnection() (string, *http.Client) {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := netClient.Get(baseURL + "/user")
	if err != nil {
		fmt.Println("Could not contact PaperCutServer")
		os.Exit(1)
	}

	defer resp.Body.Close()
	jsessionID := getCookieByName(resp.Cookies(), "JSESSIONID")

	//io.Copy(os.Stdout, resp.Body)

	return jsessionID, netClient
}

func login(credentials *PaperCutCredentials) {
	jessionid, netClient := intitalConnection()

	loginURL := baseURL + "/app"

	form := url.Values{
		"service":              {"direct/1/Home/$Form$0"},
		"sp":                   {"S0"},
		"Form0":                {"$Hidden$0,$Hidden$1,inputUsername,inputPassword,$PropertySelection$0,$Submit$0"},
		"$Hidden$0":            {"true"},
		"$Hidden$1":            {"X"},
		"inputUsername":        {credentials.username},
		"inputPassword":        {credentials.password},
		"$PropertySelection$0": {"en"},
		"$Submit$0":            {"Log in"}}

	req, _ := http.NewRequest("POST", loginURL, bytes.NewBufferString(form.Encode()))

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("Cookie", "org.apache.tapestry.locale=en;JSESSIONID="+jessionid)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "paper-app.gonzaga.edu:9192")
	req.Header.Add("Origin", "https://paper-app.gonzaga.edu:9192")
	req.Header.Add("Referer", "https://paper-app.gonzaga.edu:9192/user")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36")

	resp, err := netClient.Do(req)

	if err != nil {
		fmt.Println("Could not contact PaperCutServer")
		os.Exit(1)
	}

	defer resp.Body.Close()

	if isLoggedIn(resp) {
		credentials.isLoggedIn = true
		credentials.sessionID = jessionid
	}
}

func isLoggedIn(loginResponse *http.Response) bool {
	loggedIn := false

	doc, err := goquery.NewDocumentFromResponse(loginResponse)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		if title == "PaperCut MF  : Summary" {
			loggedIn = true
		}
	})

	return loggedIn
}

func getCookieByName(cookie []*http.Cookie, name string) string {
	cookieLen := len(cookie)
	result := ""
	for i := 0; i < cookieLen; i++ {
		if cookie[i].Name == name {
			result = cookie[i].Value
		}
	}
	return result
}
