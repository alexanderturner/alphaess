package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

// DBURL, DBUSER, DBPASS, DBNAME, ESSUSER, ESSPASS, ESSSN

type AlphaSession struct {
	sessionID *http.Cookie
	gAuth     *http.Cookie
	Expiry    *http.Cookie
}

type RTStats struct {
	Status bool       `json:"IsSucceed"`
	Data   RTStatData `json:"data"`
}

type RTStatData struct {
	Id       string  `json:"Id"`
	Soc      float32 `json:"Soc"`
	Pbat     float32 `json:"Pbat"`
	Sva      float32 `json:"Sva"`
	Ppv1     float32 `json:"Ppv1"`
	Ppv2     float32 `json:"Ppv2"`
	PmeterDc float32 `json:"PmeterDc"`
	PmeterL1 float32 `json:"PmeterL1"`
	PmeterL2 float32 `json:"PmeterL2"`
	PmeterL3 float32 `json:"PmeterL3"`
}

type BatteryStats struct {
	Solar         float32
	BattCharge    float32
	BattDischarge float32
	GridIn        float32
	GridOut       float32
	Load          float32
	Charge        float32
}

func main() {
	done := make(chan bool)
	c := influxDBClient()
	stats := new(RTStats)
	s := new(AlphaSession)

	go s.apiPoller(c, stats)
	<-done

}

func influxDBClient() client.Client {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     os.Getenv("DBURL"),
		Username: os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASS"),
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	return c
}

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

	resp, err := client.Do(req)
	defer resp.Body.Close()

	// //debugging
	// //fmt.Println(resp.Body)
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)

	for _, cookie := range resp.Cookies() {
		switch cookie.Name {
		case "alphacloudsessionid":
			s.sessionID = cookie
			fmt.Println(cookie.String())
		case "alphacloudgauth":
			s.gAuth = cookie
			fmt.Println(cookie.String())
		case "CURRENT_LANGUAGE_SESSION_KEY":
			s.Expiry = cookie
			fmt.Println(cookie.String())
		}
	}

	return err
}

func (s *AlphaSession) getRTStats(target interface{}) error {

	log.Println("Calling RT stats API")

	client := &http.Client{}
	form := url.Values{"SN": {os.Getenv("ESSSN")}, "amplifyGain": {"1"}}
	req, err := http.NewRequest("POST", "https://www.alphaess.com/Monitoring/VtColdata/GetSecondDataBySn", strings.NewReader(form.Encode()))
	if err != nil {
		return err
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
	defer resp.Body.Close()
	//debugging

	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)
	// fmt.Println()

	return json.NewDecoder(resp.Body).Decode(target)

}

func (s *AlphaSession) checkSession() (status bool) {
	t := time.Now()
	if s.Expiry == nil {
		return true
	}
	return s.Expiry.Expires.Before(t)
}

func (s *AlphaSession) apiPoller(client client.Client, stats *RTStats) error {
	for {
		time.Sleep(2 * time.Second)
		if s.checkSession() {
			log.Println("Session not valid, creating new session")
			s.getAuthToken()
		}
		s.getRTStats(stats)
		createMetrics(client, stats)
		// fmt.Printf("Battery power draw: %f watts\n", stats.Data.BatteryDraw)

	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func createMetrics(c client.Client, s *RTStats) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  os.Getenv("DBNAME"),
		Precision: "s",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	dbStats := new(BatteryStats)
	dbStats.Charge = s.Data.Soc
	if s.Data.Pbat >= 0 {
		dbStats.BattDischarge = s.Data.Pbat
	} else {
		dbStats.BattCharge = s.Data.Pbat * -1
	}
	gridLoad := s.Data.PmeterL1 + s.Data.PmeterL2 + s.Data.PmeterL3
	if gridLoad >= 0 {
		dbStats.GridIn = gridLoad
	} else {
		dbStats.GridOut = gridLoad * -1
	}
	dbStats.Solar = s.Data.Ppv1 + s.Data.Ppv2 + s.Data.PmeterDc
	dbStats.Load = s.Data.Pbat + dbStats.Solar + gridLoad

	// eventTime := time.Now().Add(time.Second * -20)
	eventTime := time.Now()

	tags := map[string]string{}

	fields := map[string]interface{}{
		"Charge":        dbStats.Charge,
		"BattCharge":    dbStats.BattCharge,
		"BattDischarge": dbStats.BattDischarge,
		"GridIn":        dbStats.GridIn,
		"GridOut":       dbStats.GridOut,
		"Solar":         dbStats.Solar,
		"Load":          dbStats.Load,
	}

	point, err := client.NewPoint(
		"batt_stat",
		tags,
		fields,
		// eventTime.Add(time.Second*10),
		eventTime,
	)
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp.AddPoint(point)

	err = c.Write(bp)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Battery Charge:%f%% | BattLoad:%fW | Grid Load:%fW | Solar Load:%fW | House Load:%f\n",
		dbStats.Charge, s.Data.Pbat, gridLoad, dbStats.Solar, dbStats.Load)
}
