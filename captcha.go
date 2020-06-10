package cashout

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Task      struct {
	Type       string  `json:"type"`
	WebsiteURL string  `json:"websiteURL"`
	WebsiteKey string  `json:"websiteKey"`
	MinScore   float64 `json:"minScore"`
}
type CapMonster struct {
	ClientKey string `json:"clientKey"`
	Task      Task `json:"task"`
}

type CapMonsterPendingAnswer struct {
	ErrorID *int `json:"errorId"`
	TaskID  *int `json:"taskId"`
}

type CapMonsterGetAnswer struct {
	ClientKey int `json:"clientKey"`
	TaskID  int `json:"taskId"`
}

type CapMonsterAnswer struct {
	ErrorID  int    `json:"errorId"`
	Status   string `json:"status"`
	Solution struct {
		GRecaptchaResponse string `json:"gRecaptchaResponse"`
	} `json:"solution"`
}

type CaptchaClient struct {
	// ApiKey is the API key for the 2captcha.com API.
	// Valid key is required by all the functions of this library
	// See more details on https://2captcha.com/2captcha-api#solving_captchas
	ApiKey string
	// Client is a HTTP client for the api calls to 2captcha
	Client *http.Client
}

// New creates a TwoCaptchaClient instance
func New(apiKey string) *CaptchaClient {
	return &CaptchaClient{
		ApiKey: apiKey,
		Client: http.DefaultClient,
	}
}

// SolvehCaptcha performs a hcaptcha solving request to 2captcha.com
// and returns with the solved captcha if the request was successful.
// Valid ApiKey is required.
// See more details on https://2captcha.com/2captcha-api#solving_hcaptcha
func (c *CaptchaClient) SolveReCaptcha(siteURL, recaptchaKey string) (*CapMonsterAnswer, error) {
	var body = []byte(fmt.Sprintf(`{"clientKey":"%s", "task":{"type":"NoCaptchaTask","websiteURL":"%s","websiteKey":"%s"}}`, c.ApiKey, siteURL, recaptchaKey))
	pendingAnswer, _, err := c.apiRequest(
		"https://api.capmonster.cloud/createTask",
		body,
		0,
		3,
	)

	if err != nil {
		return nil, err
	}

	body = []byte(fmt.Sprintf(`{"clientKey":"%s", "taskId":"%v"}`, c.ApiKey, *pendingAnswer.TaskID))
	_, answer, err := c.apiRequest(
		"https://api.capmonster.cloud/getTaskResult",
		body,
		5,
		20,
	)

	return answer, err
}

func (c *CaptchaClient) apiRequest(URL string, bodyPost []byte, delay time.Duration, retries int) (*CapMonsterPendingAnswer, *CapMonsterAnswer, error) {
	if retries <= 0 {
		return nil, nil, errors.New("Maximum retries exceeded")
	}
	time.Sleep(delay * time.Second)

	req, err := http.NewRequest("POST", URL,  bytes.NewBuffer(bodyPost))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var createTask CapMonsterPendingAnswer
	err = json.Unmarshal(body, &createTask)
	if err != nil {
		return nil, nil, err
	}

	var taskResult CapMonsterAnswer
	err = json.Unmarshal(body, &taskResult)
	if err != nil {
		return nil, nil, err
	}

	resp.Body.Close()
	if strings.Contains(string(body), "processing") {
		return c.apiRequest(URL, bodyPost, delay, retries-1)
	}

	if createTask.TaskID != nil {
		return &createTask, nil, nil
	} else {
		return nil, &taskResult, nil
	}

}