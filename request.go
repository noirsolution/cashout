package cashout

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type RequestCaptcha struct {
	SiteKey string `json:"sitekey"`
	HostKey string `json:"hostkey"`
}

type Bank struct {
	Token     string    `json:"token"`
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Sitekey   string    `json:"sitekey"`
}

func (c *Cashout) getHTTPClient() (*http.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, err
	}
	rand.Seed(time.Now().Unix())
	var client *http.Client
	if len(c.Proxies) != 0 {
		proxyURL, err := url.Parse("http://" + c.Proxies[rand.Intn(len(c.Proxies))])
		if err != nil {
			return nil, errors.New("Bad proxy URL")
		}
		transport := &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			// Disable HTTP/2
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		}

		client = &http.Client{
			Transport: transport,
			Timeout:   5 * time.Second,
			Jar:       jar,
		}
	} else {
		/*
			transport := &http.Transport{
				// Disable HTTP/2
				TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			}
		*/
		client = &http.Client{
			Timeout: 5 * time.Second,
			Jar:     jar,
		}
	}

	return client, nil
}

func requestCaptcha(sitekey, host string) (string, error) {
	var captchaResponse string
	resp, err := http.Get("http://127.0.1.1:4000/fetchCaptcha")

	bodyCaptcha, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var jsonCaptcha []Bank
	json.Unmarshal(bodyCaptcha, &jsonCaptcha)

	var validCaptcha bool
	for _, captcha := range jsonCaptcha {
		if captcha.Host == host {
			captchaResponse = captcha.Token
			validCaptcha = true
		}
	}

	if !validCaptcha {
		bodyKey := &RequestCaptcha{
			SiteKey: sitekey,
			HostKey: host,
		}

		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(bodyKey)
		_, err = http.Post("http://127.0.1.1:4000/setCaptcha", "application/json", buf)
		if err != nil {
			log.Fatal(err)
		}

		for !validCaptcha {
			resp, err = http.Get("http://127.0.1.1:4000/fetchCaptcha")

			bodyCaptcha, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}

			var jsonCaptcha []Bank
			json.Unmarshal(bodyCaptcha, &jsonCaptcha)
			for _, captcha := range jsonCaptcha {
				if captcha.Host == host {
					captchaResponse = captcha.Token
					validCaptcha = true
				}
			}
		}
	}

	return captchaResponse, nil
}
