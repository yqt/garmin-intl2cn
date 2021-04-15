package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strconv"
	"time"
)

type Request interface {
	Get(string, map[string]interface{}) (string, error)
	GetJson(string, map[string]interface{}, interface{}) error
	Post(string, map[string]interface{}, map[string]interface{}, []byte, bool) (string, error)
	PostJson(string, map[string]interface{}, map[string]interface{}, []byte, bool, interface{}) error
	GetFile(string, map[string]interface{}) ([]byte, error)
	UploadFile(string, map[string]interface{}, string, string, io.ReadCloser) (string, error)
	SetHeaders(map[string]string)
	UpdateHeaders(map[string]string)
}

type CookieRequest struct {
	headers map[string]string
	client  *http.Client
}

func NewCookieRequest() *CookieRequest {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: time.Second * 30,
	}
	return &CookieRequest{
		headers: make(map[string]string),
		client:  client,
	}
}

func (c *CookieRequest) Get(url string, params map[string]interface{}) (string, error) {
	return c.requestText(url, http.MethodGet, params, nil, nil, false)
}

func (c *CookieRequest) GetJson(url string, params map[string]interface{}, dataOut interface{}) error {
	respText, err := c.requestText(url, http.MethodGet, params, nil, nil, true)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(respText), dataOut)
}

func (c *CookieRequest) Post(url string, params map[string]interface{}, data map[string]interface{}, rawBody []byte, sendJson bool) (string, error) {
	return c.requestText(url, http.MethodPost, params, data, rawBody, sendJson)
}

func (c *CookieRequest) PostJson(url string, params map[string]interface{}, data map[string]interface{}, rawBody []byte, sendJson bool, dataOut interface{}) error {
	respText, err := c.requestText(url, http.MethodPost, params, data, rawBody, sendJson)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(respText), dataOut)
}

func (c *CookieRequest) GetFile(url string, params map[string]interface{}) ([]byte, error) {
	resp, err := c.request(url, http.MethodGet, params, nil, nil, false)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return body, nil
}

func (c *CookieRequest) UploadFile(url string, params map[string]interface{}, fileParamName string, fileName string, file io.ReadCloser) (string, error) {
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fileParamName, fileName)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	for key, val := range params {
		_ = writer.WriteField(key, fmt.Sprintf("%v", val))
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}

	c.headers["Content-Type"] = writer.FormDataContentType()
	return c.requestText(url, http.MethodPost, nil, nil, body.Bytes(), false)
}

func (c *CookieRequest) SetHeaders(headers map[string]string) {
	c.headers = headers
}

func (c *CookieRequest) UpdateHeaders(headers map[string]string) {
	for k, v := range headers {
		c.headers[k] = v
	}
}

func (c *CookieRequest) requestText(url string, method string, params map[string]interface{}, data interface{}, rawBody []byte, sendJson bool) (string, error) {
	resp, err := c.request(url, method, params, data, rawBody, sendJson)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		logrus.Errorf("invalid status code[%d]", resp.StatusCode)
		return "", errors.New("invalid status code: " + strconv.Itoa(resp.StatusCode))
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	respText := string(bodyBytes)
	logrus.WithFields(logrus.Fields{
		"respText": respText,
	}).Debug()
	return respText, nil
}

func (c *CookieRequest) request(url string, method string, params map[string]interface{}, data interface{}, rawBody []byte, sendJson bool) (*http.Response, error) {
	var buffer *bytes.Buffer

	if data != nil {
		var body []byte
		if sendJson {
			var err error
			body, err = json.Marshal(data)
			if err != nil {
				logrus.Error(err)
				return nil, errors.New("convert http post body to json failed")
			}
		} else {
			dataParam := neturl.Values{}
			for k, v := range data.(map[string]interface{}) {
				dataParam.Set(fmt.Sprintf("%v", k), fmt.Sprintf("%v", v))
			}
			body = []byte(dataParam.Encode())
		}
		buffer = bytes.NewBuffer(body)
	}

	if rawBody != nil {
		buffer = bytes.NewBuffer(rawBody)
	}

	var req *http.Request
	var err error
	if buffer == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, buffer)
	}
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("new request is fail: %v \n")
	}
	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Set(key, fmt.Sprintf("%v", val))
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if c.headers != nil {
		for key, val := range c.headers {
			req.Header.Set(key, val)
		}
	}

	if sendJson {
		req.Header.Set("Content-type", "application/json")
	} else if req.Header.Get("Content-type") == "" && method == http.MethodPost {
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	}

	logrus.Debugf("Go %s URL : %s \n", method, req.URL.String())
	return c.client.Do(req)
}
