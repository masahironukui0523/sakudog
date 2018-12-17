package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	datadog "gopkg.in/zorkian/go-datadog-api.v2"
)

type Response struct {
	Data map[string]TempMetric `json:"Data"`
}

type TempMetric struct {
	ReceiveBytesPerSec int64 `json:"receiveBytesPerSec"`
	SendBytesPerSec    int64 `json:"sendBytesPerSec"`
}

// Sakura cloud
const (
	// UserID
	token = "token"
	// Password
	secret = "secret"
	// BaseURL
	url = "url"
	// Metric name
	nameOfRecieveMetric = "sakudog.dx.receive_bytes_per_s"
	nameOfSendMetric    = "sakudog.dx.send_bytes_per_s"
)

// Datadog
const (
	apiKey   = "apiKey"
	appKey   = "appKey"
	screenId = 0
)

func main() {

	lambda.Start(hundler)

}

func hundler() {

	response := GetMetrics()

	PostMetrics(response)

}

// Basic auth
func basicAuth() string {

	var username string = token
	var passwd string = secret

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	req.SetBasicAuth(username, passwd)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)

	s := string(bodyText)

	fmt.Printf("response from sakura: %s", s)

	return s

}

// Get metrics from sakura
func GetMetrics() Response {

	res := basicAuth()

	var response Response

	if err := json.Unmarshal([]byte(res), &response); err != nil {
		log.Fatal(err)
	}

	return response

}

// Send metrics to Datadog
func PostMetrics(response Response) {

	client := datadog.NewClient(apiKey, appKey)
	receiveMetrics := []datadog.Metric{}
	sendMetrics := []datadog.Metric{}

	for key, val := range response.Data {

		receive := datadog.Metric{
			Metric: datadog.String(nameOfRecieveMetric),
			Type:   datadog.String("gauge"),
			Host:   datadog.String("prod-pfm-aws"),
			Points: []datadog.DataPoint{
				// TODO:-convert custom type(val) to float64
				{ConvertStingToFloat64(key), ConvertInt64ToFloat64(val.ReceiveBytesPerSec)},
			},
			Tags: []string{
				"prod-pfm-aws:",
			},
		}

		send := datadog.Metric{
			Metric: datadog.String(nameOfSendMetric),
			Type:   datadog.String("gauge"),
			Host:   datadog.String("prod-pfm-aws"),
			Points: []datadog.DataPoint{
				// TODO:-convert custom type(val) to float64
				{ConvertStingToFloat64(key), ConvertInt64ToFloat64(val.SendBytesPerSec)},
			},
			Tags: []string{
				"prod-pfm-aws:",
			},
		}

		receiveMetrics = append(receiveMetrics, receive)
		sendMetrics = append(sendMetrics, send)

	}

	if err := client.PostMetrics(receiveMetrics); err != nil {
		log.Fatalf("Failed to post metrics to datadog: %v", err)
	}

	fmt.Println("receiveBytesPerSecの送信に成功！")

	if err := client.PostMetrics(sendMetrics); err != nil {
		log.Fatalf("Failed to post metrics to datadog: %v", err)
	}

	fmt.Println("sendBytesPerSecの送信に成功！")

}

// String to float64
func ConvertStingToFloat64(v string) *float64 {

	// Can not parse in this way.
	// layout := "2018-12-12T23:20:00+0900"
	// FIX:- find other way to perse.
	length := len(v)
	s := v[0:length-2] + ":" + v[length-2:length]

	// string to time
	t1, err := time.Parse(time.RFC3339, s)

	if err != nil {
		log.Fatalf("Failed to convert: %v", err)
	}

	// time to float64
	t2 := float64(t1.Unix())

	return &t2

}

// Int64 to float64
func ConvertInt64ToFloat64(v int64) *float64 {

	f := float64(v)

	return &f

}
