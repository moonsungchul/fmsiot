package fmsiot

import (
	"log"
	"testing"
	"time"
)

func TestInfluxDB(t *testing.T) {
	log.Println("test influx db ...")
	saddr := "http://localhost:8086"
	susername := ""
	spassword := ""
	sdb := ""

	dbm := NewInfluxDB(saddr, susername, spassword, sdb)

	con, err := dbm.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer dbm.CloseDB(con)

	q := "select * from cpu order by time desc limit 10"
	res, err2 := dbm.QueryDB(con, q, "telegraf")
	if err2 != nil {
		log.Fatal(err2)
	}

	for i, row := range res[0].Series[0].Values {
		t, err3 := time.Parse(time.RFC3339, row[0].(string))
		if err3 != nil {
			log.Fatal(err)
		}
		val := row[1].(string)
		val2 := row[2].(string)
		log.Printf("[%2d] %s: %s : %s\n", i, t.Format(time.Stamp), val, val2)
	}

}
