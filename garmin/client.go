package garmin

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yqt/garmin-intl2cn/util"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type UserInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	ApiServiceHost   = "connect.garmin.com"
	ApiServiceHostCn = "connect.garmin.cn"
	SsoPrefix        = "https://sso.garmin.com"
	SsoPrefixCn      = "https://sso.garmin.cn"
	UserAgent        = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36"
)

type Client struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ApiHost   string `json:"api_host"`
	ApiPrefix string `json:"api_prefix"`
	SsoPrefix string `json:"sso_prefix"`

	client   *util.CookieRequest
	loggedIn bool
}

type Option func(client *Client)

func Credentials(email string, password string) Option {
	return func(c *Client) {
		c.Email = email
		c.Password = password
	}
}

func SetEnv(apiHost string, ssoPrefix string) Option {
	return func(c *Client) {
		c.ApiHost = apiHost
		c.ApiPrefix = "https://" + apiHost
		c.SsoPrefix = ssoPrefix
	}
}

func NewClient(options ...Option) *Client {
	client := &Client{
		client:   util.NewCookieRequest(),
		loggedIn: false,
	}

	client.SetOptions(options...)

	return client
}

func (c *Client) SetOptions(options ...Option) {
	for _, option := range options {
		option(c)
	}
}

func (c *Client) Auth(reLogin bool) error {
	if reLogin {
		c.loggedIn = false
	}
	if c.loggedIn {
		logrus.WithFields(logrus.Fields{
			"loggedIn": c.loggedIn,
		}).Debug()
		return nil
	}
	params := map[string]interface{}{
		"service":                        c.ApiPrefix + "/modern",
		"clientId":                       "GarminConnect",
		"gauthHost":                      c.SsoPrefix + "/sso",
		"generateExtraServiceTicket":     "true",
		"generateTwoExtraServiceTickets": "true",
		"consumeServiceTicket":           "false",
	}

	uri := c.SsoPrefix + "/sso/signin"
	headers := map[string]string{
		"User-Agent": UserAgent,
	}
	c.client.SetHeaders(headers)

	respText, err := c.client.Get(uri, params)
	if err != nil {
		return err
	}
	csrfToken, err := c.extractCSRFToken(respText)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"csrfToken": csrfToken,
	}).Debug()

	formData := map[string]interface{}{
		"username": c.Email,
		"password": c.Password,
		"embed":    "false",
		"_csrf":    csrfToken,
	}
	headers["Origin"] = c.SsoPrefix
	q := url.Values{}
	for key, val := range params {
		q.Set(key, fmt.Sprintf("%v", val))
	}
	headers["Referer"] = uri + "?" + q.Encode()
	c.client.SetHeaders(headers)
	respText, err = c.client.Post(uri, params, formData, nil, false)
	if err != nil {
		return err
	}

	ticketUrl, err := c.extractTicketUrl(respText)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"ticketUrl": ticketUrl,
	}).Debug()

	respText, err = c.client.Get(ticketUrl, nil)
	socialProfileText, err := c.extractSocialProfile(respText)
	if err != nil {
		c.loggedIn = false
		return err
	}
	logrus.WithFields(logrus.Fields{
		"socialProfileText": socialProfileText,
	}).Debug()

	c.loggedIn = true

	return nil
}

func (c *Client) GetActivity(id int64) (Activity, error) {
	uri := c.ApiPrefix + "/proxy/activity-service/activity/" + strconv.FormatInt(id, 10)
	activity := Activity{}
	err := c.client.GetJson(uri, nil, &activity)
	if err != nil {
		return activity, err
	}

	return activity, nil
}

func (c *Client) GetActivityList(start int64, limit int64) ([]ActivityListItem, error) {
	uri := c.ApiPrefix + "/proxy/activitylist-service/activities/search/activities"
	activityList := make([]ActivityListItem, 0)
	params := map[string]interface{}{
		"start": start,
		"limit": limit,
	}
	err := c.client.GetJson(uri, params, &activityList)
	if err != nil {
		return activityList, err
	}

	return activityList, nil
}

func (c *Client) DownloadActivity(id int64) (io.ReadCloser, string, error) {
	uri := c.ApiPrefix + "/modern/proxy/download-service/files/activity/" + strconv.FormatInt(id, 10)

	contentBytes, err := c.client.GetFile(uri, nil)

	zipReader, err := zip.NewReader(bytes.NewReader(contentBytes), int64(len(contentBytes)))
	if err != nil {
		return nil, "", err
	}

	for _, zipFile := range zipReader.File {
		logrus.Debug("Reading file:", zipFile.Name)
		file, err := zipFile.Open()
		if err != nil {
			return nil, "", err
		}
		return file, zipFile.Name, err
	}

	return nil, "", errors.New("not file in zip file")
}

func (c *Client) UploadActivity(fileName string, file io.ReadCloser) error {
	uri := c.ApiPrefix + "/modern/proxy/upload-service/upload/.fit"

	headers := map[string]string{
		"Origin":  c.ApiPrefix,
		"Referer": c.ApiPrefix + "/modern/import-data",
		"Nk":      "NT",
	}
	c.client.UpdateHeaders(headers)

	respText, err := c.client.UploadFile(uri, nil, "file", fileName, file)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"uploadRespText": respText,
	}).Info()

	return nil
}

func (c *Client) extractCSRFToken(respText string) (string, error) {
	fragment := `<input type="hidden" name="_csrf" value="`
	startPos := strings.Index(respText, fragment)
	if startPos == -1 {
		return "", errors.New("CSRF token not found")
	}
	restText := respText[startPos:]
	endPos := strings.Index(restText, `" />`)
	if endPos == -1 {
		return "", errors.New("invalid CSRF token end")
	}
	restText = restText[len(fragment):endPos]
	return restText, nil
}

func (c *Client) extractTicketUrl(respText string) (string, error) {
	t := regexp.MustCompile(`https:\\\/\\\/` + c.ApiHost + `\\\/modern(\\\/)?\?ticket=(([a-zA-Z0-9]|-)*)`)
	//t := regexp.MustCompile(`https(.+?)modern(%2F)?\?ticket=(([a-zA-Z0-9]|-)*)`)
	ticketUrl := t.FindString(respText)

	// NOTE: undo escaping
	ticketUrl = strings.Replace(ticketUrl, "\\/", "/", -1)
	//ticketUrl, err := url.QueryUnescape(ticketUrl)
	//if err != nil {
	//	return "", err
	//}

	if ticketUrl == "" {
		return "", errors.New("wrong credentials")
	}

	return ticketUrl, nil
}

func (c *Client) extractSocialProfile(respText string) (string, error) {
	fragment := `window.VIEWER_SOCIAL_PROFILE = JSON.parse("`
	startPos := strings.Index(respText, fragment)
	if startPos == -1 {
		return "", errors.New("social profile not found")
	}
	restText := respText[startPos:]
	endPos := strings.Index(restText, `");`)
	if endPos == -1 {
		return "", errors.New("invalid social profile end")
	}
	restText = restText[len(fragment):endPos]
	restText = strings.Replace(restText, "\\", "", -1)
	return restText, nil
}
