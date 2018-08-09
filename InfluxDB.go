package fmsiot

import (
	"github.com/influxdata/influxdb/client/v2"
	"time"
)

type InfluxDB struct {
	addr     string
	username string
	password string
	db       string
}

func NewInfluxDB(saddr string, user string, pass string, sdb string) InfluxDB {
	return InfluxDB{
		addr:     saddr,
		username: user,
		password: pass,
		db:       sdb,
	}
}

func (c *InfluxDB) Open() (client.Client, error) {
	conf := client.HTTPConfig{
		Addr:     c.addr,
		Username: c.username,
		Password: c.password,
	}
	con, err := client.NewHTTPClient(conf)
	if err != nil {
		return nil, err
	}
	return con, nil

}

func (c *InfluxDB) QueryDB(clnt client.Client, cmd string, db string) (res []client.Result, err error) {
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

func (c *InfluxDB) CloseDB(clnt client.Client) {
	clnt.Close()
}

func (c *InfluxDB) Insert(clnt client.Client, table string, tags map[string]string, fields map[string]interface{}) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  c.db,
		Precision: "s",
	})

	if err != nil {
		return err
	}

	pt, err1 := client.NewPoint(table, tags, fields, time.Now())
	if err1 != nil {
		return err1
	}

	bp.AddPoint(pt)

	err2 := clnt.Write(bp)
	if err2 != nil {
		return err2
	}

	return nil

}
