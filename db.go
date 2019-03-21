package main

import (
	"log"
	"os"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

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

func (a *App) createMetrics(s RTStats) {
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
		eventTime,
	)
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp.AddPoint(point)

	err = a.DB.Write(bp)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Battery Charge:%f%% | BattLoad:%fW | Grid Load:%fW | Solar Load:%fW | House Load:%f\n",
		dbStats.Charge, s.Data.Pbat, gridLoad, dbStats.Solar, dbStats.Load)
}
