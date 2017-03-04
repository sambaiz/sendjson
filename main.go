package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/sambaiz/sendjson/input"
	"github.com/urfave/cli"
)

func main() {
	var (
		interval, duration time.Duration
		url, method        string
		check              bool
	)

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.DurationFlag{
			Name:        "interval, i",
			Value:       1 * time.Second,
			Usage:       "send interval",
			Destination: &interval,
		},
		cli.DurationFlag{
			Name:        "duration, d",
			Value:       60 * time.Second,
			Usage:       "send duration",
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
		cli.BoolFlag{
			Name:        "check",
			Usage:       "for checking output without sending",
			Destination: &check,
		},
	}

	app.Action = func(c *cli.Context) error {

		var input string
		if c.NArg() > 0 {
			input = c.Args()[0]
		} else {
			panic("json input required")
		}

		ticker := time.NewTicker(interval)

		go func() {
			for {
				select {
				case <-ticker.C:
					go func() {
						if b, err := generateJSON([]byte(input)); err != nil {
							panic(err)
						} else if check {
							fmt.Printf("not sended: %v\n", string(b))
						} else if err := send(method, url, b); err != nil {
							panic(err)
						} else {
							fmt.Printf("%v\n", string(b))
						}
					}()
				}
			}
		}()

		time.Sleep(duration)
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

	if res.StatusCode/100 != 2 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		} else {
			fmt.Printf("[send error] status code: %d, response: %s\n", res.StatusCode, string(body))
		}
	}

	return nil
}
