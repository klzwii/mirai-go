package chatgpt

import (
	"bytes"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync"
	"time"
)

var logger = logrus.WithField("plugin", "Chatgpt")

type config struct {
	MaxUserHistory int               `mapstructure:"max_user_history"`
	Prompt         map[uint64]string `mapstructure:"prompt"`
	AccessKey      string            `mapstructure:"access_key"`
	Model          string            `mapstructure:"model"`
	Temperature    float64           `mapstructure:"temperature"`
	MaxTokens      int               `mapstructure:"max_tokens"`
}

type gpt struct {
	conf        *config
	UserHistory map[string][]*record
	m           *sync.Mutex
}

func (g *gpt) createRequest(message *record, key string, promptID uint64) *request {
	messages := append(g.UserHistory[key], message)
	prompt := g.conf.Prompt[promptID]
	if len(prompt) != 0 {
		messages = append([]*record{{
			Role:    "system",
			Content: prompt,
		}}, messages...)
	}
	return &request{
		Model:       g.conf.Model,
		Messages:    messages,
		Temperature: g.conf.Temperature,
	}
}

const OpenapiUrl string = "https://api.openai.com/v1/chat/completions"

func generateKey(user string, group uint64) string {
	return fmt.Sprintf("%s+%d", user, group)
}

var client = http.Client{Timeout: time.Second * 120}

func (g *gpt) sendRequest(req *request, resp *response) error {
	data, err := jsoniter.Marshal(req)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, OpenapiUrl, bytes.NewReader(data))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.conf.AccessKey))
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(data, resp)
}

var (
	noChoiceError = errors.New("脑袋空空，说不出话")
	internalError = errors.New("大脑断开连接，直接睡觉")
)

func (g *gpt) processResponse(resp *response, req *record, key string) (*record, error) {
	if len(resp.Choices) < 1 {
		return nil, noChoiceError
	}
	g.m.Lock()
	defer g.m.Unlock()
	history := g.UserHistory[key]
	history = append(history, req, resp.Choices[0].Message)
	if len(history) > g.conf.MaxUserHistory*2 {
		history = history[2:]
	}
	g.UserHistory[key] = history
	return resp.Choices[0].Message, nil
}

func (g *gpt) clearHistory() {
	g.m.Lock()
	defer g.m.Unlock()
	for key := range g.UserHistory {
		delete(g.UserHistory, key)
	}
}
