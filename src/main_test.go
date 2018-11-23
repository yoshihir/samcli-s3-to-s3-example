package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	println("before all...")

	code := m.Run()

	println("after all...")

	os.Exit(code)
}

func TestHandler(t *testing.T) {
	t.Run("handler input test", func(t *testing.T) {
		raw, err := ioutil.ReadFile("../event_file.json")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		var event events.S3Event
		json.Unmarshal(raw, &event)
		err = handler(context.Background(), event)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		println("Test Handler...")
	})

	//
	//
	//t.Run("Unable to get IP", func(t *testing.T) {
	//	DefaultHTTPGetAddress = "http://127.0.0.1:12345"
	//
	//	_, err := handler(events.APIGatewayProxyRequest{})
	//	if err == nil {
	//		t.Fatal("Error failed to trigger with an invalid request")
	//	}
	//})
	//
	//t.Run("Non 200 Response", func(t *testing.T) {
	//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		w.WriteHeader(500)
	//	}))
	//	defer ts.Close()
	//
	//	DefaultHTTPGetAddress = ts.URL
	//
	//	_, err := handler(events.APIGatewayProxyRequest{})
	//	if err != nil && err.Error() != ErrNon200Response.Error() {
	//		t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
	//	}
	//})
	//
	//t.Run("Unable decode IP", func(t *testing.T) {
	//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		w.WriteHeader(500)
	//	}))
	//	defer ts.Close()
	//
	//	DefaultHTTPGetAddress = ts.URL
	//
	//	_, err := handler(events.APIGatewayProxyRequest{})
	//	if err == nil {
	//		t.Fatal("Error failed to trigger with an invalid HTTP response")
	//	}
	//})
	//
	//t.Run("Successful Request", func(t *testing.T) {
	//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		w.WriteHeader(200)
	//		fmt.Fprintf(w, "127.0.0.1")
	//	}))
	//	defer ts.Close()
	//
	//	DefaultHTTPGetAddress = ts.URL
	//
	//	_, err := handler(events.APIGatewayProxyRequest{})
	//	if err != nil {
	//		t.Fatal("Everything should be ok")
	//	}
	//})
}
