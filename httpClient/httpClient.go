package httpClient

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func ParseClient(method, url string, payload *strings.Reader, v interface{}) (*http.Response, error) {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
		MaxIdleConns:        10,
		IdleConnTimeout:     10 * time.Second,
	}
	http_client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Parse-Master-Key", "DEV_MASTER_KEY")
	req.Header.Add("X-Parse-Application-Id", "DEV_APPLICATION_ID")
	req.Header.Add("Content-Type", "application/json")

	resp, err := http_client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		if resp.StatusCode == 400 {
			type temp struct {
				Code  int    `json:"code"`
				Error string `json:"error"`
			}
			var a temp
			_ = json.NewDecoder(resp.Body).Decode(&a)
			fmt.Printf("response error %v", a)
			errR := errors.New(fmt.Sprint(a))
			return nil, errR
		}

		err = fmt.Errorf(fmt.Sprintf("response error from parse client - %v", resp))
		return nil, err
	}
	defer resp.Body.Close()
	errJson := json.NewDecoder(resp.Body).Decode(v)
	if errJson != nil {
		return nil, errJson
	}
	return resp, err
}

func NormalClient(method, url string, payload *strings.Reader, v interface{}) (*http.Response, error) {
	caCert, err := ioutil.ReadFile("rootCA.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
		MaxIdleConns:        10,
		IdleConnTimeout:     10 * time.Second,
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
			
		},
	}
	http_client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http_client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		if resp.StatusCode == 400 {
			type temp struct {
				Code  int    `json:"code"`
				Error string `json:"error"`
			}
			var a temp
			_ = json.NewDecoder(resp.Body).Decode(&a)
			fmt.Printf("response error %v", a)
			errR := errors.New(fmt.Sprint(a))
			return nil, errR
		}

		err = fmt.Errorf(fmt.Sprintf("response error from parse client - %v", resp))
		return nil, err
	}
	defer resp.Body.Close()
	errJson := json.NewDecoder(resp.Body).Decode(v)
	if errJson != nil {
		return nil, errJson
	}
	return resp, err
}