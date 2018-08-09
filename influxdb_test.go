package fmsiot

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"testing"
	"time"
)

func TestExampleNewClient(t *testing.T) {
	/*
		host, err := url.Parse(fmt.Sprintf("http://%s:%d", "localhost", 8086))
		if err != nil {
			log.Fatal(err)
		} */

	conf := client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "",
		Password: "",
	}
	log.Println("test1")

	con, err := client.NewHTTPClient(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	q := fmt.Sprintf("select * from %s order by time desc limit 10", "cpu")
	res, err := queryDB(con, q, "telegraf")
	if err != nil {
		log.Fatal(err)
	}

	for i, row := range res[0].Series[0].Values {
		t, err := time.Parse(time.RFC3339, row[0].(string))
		if err != nil {
			log.Fatal(err)
		}
		val := row[1].(string)
		val2 := row[2].(string)
		log.Printf("[%2d] %s: %s : %s\n", i, t.Format(time.Stamp), val, val2)
	}

	// Count Records
	q = fmt.Sprintf("select count(*) from cpu")
	res, err = queryDB(con, q, "telegraf")
	if err != nil {
		log.Fatal(err)
	}
	count := res[0].Series[0].Values[0][1]
	log.Printf("Found a total of %v records \n", count)

}

func queryDB(clnt client.Client, cmd string, db string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: db,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
