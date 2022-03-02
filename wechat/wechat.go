package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type (
	Request struct {
		Msgtype string                 `json:"msgtype"`
		Text    map[string]interface{} `json:"text"`
	}

	MarkdownRequest struct {
		Msgtype  string                 `json:"msgtype"`
		Markdown map[string]interface{} `json:"markdown"`
	}

	Build struct {
		Owner   string
		Name    string
		Tag     string
		Event   string
		Number  int
		Commit  string
		Ref     string
		Branch  string
		Author  string
		Message string
		Status  string
		Link    string
		Started int64
		Created int64
	}

	Response struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}

	WeChat struct {
		Build               Build
		Url                 string
		MsgType             string
		MentionedList       string
		MentionedMobileList string
		Content             string
	}
)

func jsonEncode(d interface{}) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	jsonEncoder := json.NewEncoder(buf)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(d)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *WeChat) MarkdownMessage(md string, mentionedList, mentionedMobileList []string) error {
	we := &MarkdownRequest{
		Msgtype: "markdown",
		Markdown: map[string]interface{}{
			"content": md,
		},
	}

	if len(mentionedList) > 0 {
		we.Markdown["mentioned_list"] = mentionedList
	}
	if len(mentionedMobileList) > 0 {
		we.Markdown["mentioned_mobile_list"] = mentionedMobileList
	}

	buf, err := jsonEncode(we)
	if err != nil {
		return err
	}

	return c.call(buf)
}

func (c *WeChat) Template(temp string) ([]byte, error) {
	tmpl, err := template.New("wechat").Parse(temp)
	if err != nil {
		return nil, fmt.Errorf("template parse error %w %s", err, temp)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, c.Build)
	if err != nil {
		return nil, fmt.Errorf("template execute error %w", err)
	}

	return buf.Bytes(), nil
}

func (c *WeChat) call(buf *bytes.Buffer) error {
	resp, err := c.postJson(c.Url, buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if data, err := ioutil.ReadAll(resp.Body); err == nil {
		ret := &Response{}
		err := json.Unmarshal(data, ret)
		if err != nil {
			return err
		}

		if ret.Errcode != 0 {
			return errors.New("ding response error:" + ret.Errmsg + "[" + strconv.Itoa(ret.Errcode) + "]")
		}
	}

	return nil
}

func (c *WeChat) Message(content string, mentionedList, mentionedMobileList []string) error {
	we := &Request{
		Msgtype: "text",
		Text: map[string]interface{}{
			"content": content,
		},
	}

	if len(mentionedList) > 0 {
		we.Text["mentioned_list"] = mentionedList
	}
	if len(mentionedMobileList) > 0 {
		we.Text["mentioned_mobile_list"] = mentionedMobileList
	}

	buf, err := jsonEncode(we)
	if err != nil {
		return err
	}

	return c.call(buf)
}

func (c *WeChat) postJson(url string, body *bytes.Buffer) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	client := &http.Client{}
	client.Timeout = 5 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *WeChat) Send() error {
	var mentionedList, mentionedMobileList []string
	if c.MentionedList != "" {
		mentionedList = strings.Split(c.MentionedList, ",")
	}
	if c.MentionedMobileList != "" {
		mentionedMobileList = strings.Split(c.MentionedMobileList, ",")
	}

	tempBuf, err := c.Template(c.Content)
	if err != nil {
		return err
	}

	if c.MsgType == "text" {
		return c.Message(strings.TrimSpace(string(tempBuf)), mentionedList, mentionedMobileList)
	}

	if c.MsgType == "markdown" {
		return c.MarkdownMessage(strings.TrimSpace(string(tempBuf)), mentionedList, mentionedMobileList)
	}

	return fmt.Errorf("not supported msgtype %s", c.MsgType)
}
