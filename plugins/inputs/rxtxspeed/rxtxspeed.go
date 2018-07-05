package rxtxspeed

import (
	"io/ioutil"
	"strconv"
	"time"
	"math"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type Rxtxspeed struct {
	Network_interface	string
	Rx_bytes_old		int64
	Tx_bytes_old		int64
	Timestamp_old		int64
}

var RxtxspeedConfig = `
	## Set the network interface
	Network_interface = "enp0s3"
`
func (s *Rxtxspeed) SampleConfig() string {
	return RxtxspeedConfig
}

func (s *Rxtxspeed) Description() string {
	return "Measures rx/tx speed on a chosen network interface"
}

func (s *Rxtxspeed) Gather(acc telegraf.Accumulator) error {	
	rx_bytes_b, _ := ioutil.ReadFile("/sys/class/net/" + s.Network_interface + "/statistics/rx_bytes")
	rx_bytes, _ := strconv.ParseInt(string(rx_bytes_b)[0:len(rx_bytes_b) - 1], 10, 64)
	tx_bytes_b, _ := ioutil.ReadFile("/sys/class/net/" + s.Network_interface + "/statistics/tx_bytes")
	tx_bytes, _ := strconv.ParseInt(string(tx_bytes_b)[0:len(tx_bytes_b) - 1], 10, 64)

	if (s.Timestamp_old == -1){
		s.Timestamp_old = time.Now().Unix() - 1
	}
	if (s.Rx_bytes_old == -1) {
		s.Rx_bytes_old = rx_bytes
		s.Tx_bytes_old = tx_bytes
	}

	timestamp := time.Now().Unix()
	// The final result represents the average speed in MiB on predefined telegraf interval
	var rx_speed float64 = float64(rx_bytes - s.Rx_bytes_old) / 1024.0
	rx_speed /= float64(timestamp - s.Timestamp_old)
	rx_speed = math.Round((rx_speed / 1024.0)*100)/100
	var tx_speed float64 = float64(tx_bytes - s.Tx_bytes_old) / 1024.0
	tx_speed /= float64(timestamp - s.Timestamp_old)
	tx_speed = math.Round((tx_speed / 1024.0)*100)/100

	s.Rx_bytes_old = rx_bytes
	s.Tx_bytes_old = tx_bytes
	s.Timestamp_old = timestamp	

	fields := make(map[string]interface{})
	fields["rx speed"] = rx_speed
	fields["tx speed"] = tx_speed

	tags := make(map[string]string)
	acc.AddFields("rxtxspeed", fields, tags)

	return nil
}

func init() {
	inputs.Add("rxtxspeed", func() telegraf.Input { return &Rxtxspeed{Rx_bytes_old: -1, Tx_bytes_old: -1, Timestamp_old: -1} })
}