package main

import (
	"net/http"

	client "github.com/influxdata/influxdb1-client/v2"
)

type App struct {
	Alpha AlphaSession
	DB    client.Client
}

type AlphaSession struct {
	sessionID *http.Cookie
	gAuth     *http.Cookie
	Expiry    *http.Cookie
}

type RTStats struct {
	Status bool       `json:"IsSucceed,omitempty"`
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
