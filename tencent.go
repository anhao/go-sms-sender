// Copyright 2021 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package go_sms_sender

import (
	"fmt"
	"strconv"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
)

type TencentClient struct {
	core     *sms.Client
	appId    string
	sign     string
	template string
}

func GetTencentClient(accessId, accessKey, sign, region, templateId string, appId []string) (*TencentClient, error) {
	if len(appId) < 1 {
		return nil, fmt.Errorf("missing parameter: appId")
	}

	credential := common.NewCredential(accessId, accessKey)
	config := profile.NewClientProfile()
	config.HttpProfile.ReqMethod = "POST"

	client, err := sms.NewClient(credential, region, config)
	if err != nil {
		return nil, err
	}

	tencentClient := &TencentClient{
		core:     client,
		appId:    appId[0],
		sign:     sign,
		template: templateId,
	}

	return tencentClient, nil
}

func (c *TencentClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	var paramArray []string
	index := 0
	for {
		value := param[strconv.Itoa(index)]
		if len(value) == 0 {
			break
		}
		paramArray = append(paramArray, value)
		index++
	}

	request := sms.NewSendSmsRequest()
	request.SmsSdkAppid = common.StringPtr(c.appId)
	request.Sign = common.StringPtr(c.sign)
	request.TemplateParamSet = common.StringPtrs(paramArray)
	request.TemplateID = common.StringPtr(c.template)
	request.PhoneNumberSet = common.StringPtrs(targetPhoneNumber)

	_, err := c.core.SendSms(request)
	return err
}