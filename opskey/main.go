package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	// "sync"
)

type Cli struct {
	Username     string
	Password     string
	Business     string
	ClientId     string
	SmsCode      string
	ClientSecret string
	GrantType    string
	Timeout      int

	Token string
	Code  string
}

type Comm interface {
	getToken() 	any
	getCode() 	any
	printToken() any
	printCode() any
}

var (
	// token单例
	// onces = &sync.Once{}

	cli = &Cli{
		Username:     "xxxxxxx",
		Password:     "xxxx#",
		Business:     "waxxxtsons",
		ClientId:     "browser",
		SmsCode:      "123456",
		ClientSecret: "xxxxxxxxxxxxxxx",
		GrantType:    "password",
		Timeout:      10,
		Token:        "",
		Code:         "",
	}
)

func init() {
	flag.StringVar(&cli.Username, "u", cli.Username, "set username")
	flag.StringVar(&cli.Password, "p", cli.Password, "set password")
	flag.StringVar(&cli.Business, "b", cli.Business, "set business")
	flag.StringVar(&cli.ClientId, "i", cli.ClientId, "set clientId")
	flag.StringVar(&cli.SmsCode, "c", cli.SmsCode, "set SMS verification code")
	flag.StringVar(&cli.ClientSecret, "s", cli.ClientSecret, "set clientsecret")
	flag.StringVar(&cli.GrantType, "g", cli.GrantType, "set granttype")
	flag.IntVar(&cli.Timeout, "t", cli.Timeout, "set timeout")
	flag.Parse()
}

func MD5(code string) string {
	MD5 := md5.New()
	_, _ = io.WriteString(MD5, code)
	return hex.EncodeToString(MD5.Sum(nil))
}

// GetToken 获取token
func (c *Cli) getToken() any {

	var result = make(map[string]interface{})
	var ReqData = make(url.Values)
	var uB = c.Username + "@" + c.Business

	ReqData["username"] = []string{uB}
	ReqData["password"] = []string{MD5(MD5(c.Password)) + "_" + MD5(c.SmsCode)}
	ReqData["client_id"] = []string{c.ClientId}
	ReqData["client_secret"] = []string{c.ClientSecret}
	ReqData["grant_type"] = []string{c.GrantType}

	// onces.Do(func() {
	res, err := http.PostForm("https://sso.xxxxx.net/oauth/token?", ReqData)
	if err != nil {
		fmt.Printf("Login Error: %v", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	if err != nil {
		panic(err)
	}
	token := result["access_token"].(string)
	if len(token) != 0 {
		c.Token = token
	}
	// })
	return c
}

func (c *Cli) getCode() any {
	var result = make(map[string]interface{})

	codeURL := "https://aiops.xxxxx.net/api/agent_manage/agent/getRegisterCode?access_token=" + c.Token
	res, err := http.Get(codeURL)
	if err != nil {
		fmt.Printf("GetCode Error: %v", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	if err != nil {
		panic(err)
	}
	code := result["registerCode"].(string)

	if len(code) != 0 {
		c.Code = code
	}
	return c
}

func (c *Cli) printCode() any {
	fmt.Println(c.Code)
	return c
}

func (c *Cli) printToken() any {
	fmt.Println(c.Token)
	return c
}

func main() {
	var comm Comm = cli
	comm.
		getToken().(*Cli).getCode().(*Cli).printCode()
}
