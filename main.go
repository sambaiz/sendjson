package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"encoding/json"

	"net/http"

	"bytes"

	"github.com/sambaiz/sendjson/input"
	"github.com/urfave/cli"
)

func main() {
	var (
		interval, duration float64
		url, method        string
	)

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.Float64Flag{
			Name:        "interval, i",
			Value:       1,
			Usage:       "send interval(sec)",
			Destination: &interval,
		},
		cli.Float64Flag{
			Name:        "duration, d",
			Value:       -1,
			Usage:       "send duration(sec), -1 is infinity",
			Destination: &duration,
		},
		cli.StringFlag{
			Name:        "url, u",
			Value:       "http://127.0.0.1",
			Usage:       "send to url",
			Destination: &url,
		},
		cli.StringFlag{
			Name:        "method, m",
			Value:       "POST",
			Usage:       "send method; POST, PUT etc.",
			Destination: &method,
		},
	}

	app.Action = func(c *cli.Context) error {

		var input string
		if c.NArg() > 0 {
			input = c.Args()[0]
		} else {
			panic("json input required")
		}

		ticker := time.NewTicker(time.Duration(interval*math.Pow(10, 9)) * time.Nanosecond)

		go func() {
			for {
				select {
				case <-ticker.C:
					go func() {
						if b, err := generateJSON([]byte(input)); err != nil {
							panic(err)
						} else if err := send(method, url, b); err != nil {
							panic(err)
						} else {
							fmt.Sprintf("sended: %v", string(b))
						}
					}()
				}
			}
		}()

		if duration != -1 {
			time.Sleep(time.Duration(duration*math.Pow(10, 9)) * time.Nanosecond)
		} else {
			// infinity
			for {
			}
		}

		ticker.Stop()
		return nil
	}

	app.Run(os.Args)
}

func generateJSON(inbyte []byte) ([]byte, error) {

	var input map[string]input.Input
	output := map[string]interface{}{}

	if err := json.Unmarshal(inbyte, &input); err != nil {
		return nil, err
	}

	for k, v := range input {
		ev := v.Eval()
		if ev == nil {
			return nil, fmt.Errorf("eval %v then null!", k)
		}
		output[k] = ev
	}

	outbyte, err := json.Marshal(output)
	fmt.Println(string(outbyte))
	return outbyte, err
}

func send(method, url string, data []byte) error {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
