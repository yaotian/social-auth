// Copyright 2014 beego authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
//
// Maintain by https://github.com/slene

package apps

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"

	"github.com/yaotian/social-auth"
)

type QQ struct {
	BaseProvider
}

func (p *QQ) GetType() social.SocialType {
	return social.SocialQQ
}

func (p *QQ) GetName() string {
	return "QQ"
}

func (p *QQ) GetPath() string {
	return "qq"
}

func (p *QQ) GetIndentify(tok *social.Token) (string, error) {
	vals := make(map[string]interface{})

	uri := "https://graph.qq.com/oauth2.0/me?access_token=" + url.QueryEscape(tok.AccessToken)
	req := httplib.Get(uri)
	req.SetTransport(social.DefaultTransport)

	body, err := req.String()
	if err != nil {
		return "", err
	}
	beego.Info("get indentify body:", body)
	body = strings.Replace(body, "callback( ", "", -1)
	body = strings.Replace(body, " );", "", -1)
	beego.Info("get indentify body2:", body)

	//
	//	vals, err := url.ParseQuery(body)
	//	if err != nil {
	//		return "", err
	//	}

	//	resp, err := req.Response()
	//	if err != nil {
	//		return "", err
	//	}
	//	defer resp.Body.Close()

	decoder := json.NewDecoder(strings.NewReader(body))
	decoder.UseNumber()

	if err := decoder.Decode(&vals); err != nil {
		return "", err
	}

	beego.Info(vals)

	if opid, ok := vals["openid"]; ok {
		return fmt.Sprint(opid), nil
	}else{
		return "", fmt.Errorf("can't get qq identify")
	}
}

var _ social.Provider = new(QQ)

func NewQQ(clientId, secret string) *QQ {
	p := new(QQ)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "get_user_info"
	p.AuthURL = "https://graph.qq.com/oauth2.0/authorize"
	p.TokenURL = "https://graph.qq.com/oauth2.0/token"
	p.RedirectURL = social.DefaultAppUrl + "login/qq/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}
