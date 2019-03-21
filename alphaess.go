package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func (s *AlphaSession) getAuthToken() error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	form := url.Values{"Username": {os.Getenv("ESSUSER")}, "Userpwd": {os.Getenv("ESSPASS")}}
	req, err := http.NewRequest("POST", "http://www.alphaess.com", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	//debugging
	//fmt.Println(resp.Body)
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)

	for _, cookie := range resp.Cookies() {
		switch cookie.Name {
		case "alphacloudsessionid":
			s.sessionID = cookie
			//fmt.Println(cookie.String())
		case "alphacloudgauth":
			s.gAuth = cookie
			//fmt.Println(cookie.String())
		case "CURRENT_LANGUAGE_SESSION_KEY":
			s.Expiry = cookie
			//fmt.Println(cookie.String())
		}
	}

	return err
}

// func (s *AlphaSession) getRTStats(target interface{}) error {
func (s *AlphaSession) getRTStats() (r RTStats, err error)   {

	log.Println("Calling RT stats API")

	client := &http.Client{}
	form := url.Values{"SN": {os.Getenv("ESSSN")}, "amplifyGain": {"1"}}
	req, err := http.NewRequest("POST", "https://www.alphaess.com/Monitoring/VtColdata/GetSecondDataBySn", strings.NewReader(form.Encode()))
	if err != nil {
		return r, err
	}
	req.PostForm = form

	asp := new(http.Cookie)
	asp.Raw = "ASP.NET_SessionId=; path=/; domain=.www.alphaess.com;"

	req.AddCookie(s.gAuth)
	req.AddCookie(s.sessionID)
	req.AddCookie(asp)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	// log.Println(resp.Body)

	defer resp.Body.Close()
	//debugging

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
  err = json.Unmarshal(bodyBytes, &r)
  if err != nil {
    return r, err
  }
  //err = json.NewDecoder(resp.Body).Decode(reta

	return r, nil

}

func (s *AlphaSession) checkSession() (status bool) {
	t := time.Now()
	if s.Expiry == nil {
		return true
	}
	return s.Expiry.Expires.Before(t)
}
