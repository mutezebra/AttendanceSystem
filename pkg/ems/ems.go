package ems

import (
	"fmt"
	"github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/config"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/log"
	"github.com/pkg/errors"
)

type Client struct {
	c *dysmsapi20170525.Client
}

func NewEmsClient() *Client {
	c := &client.Config{
		AccessKeyId:     &config.Conf.Aliyun.AccessKeyID,
		AccessKeySecret: &config.Conf.Aliyun.AccessKeySecret,
		Endpoint:        &config.Conf.Aliyun.Endpoint,
	}

	cli, err := dysmsapi20170525.NewClient(c)
	if err != nil {
		log.LogrusObj.Panic(err)
	}
	return &Client{cli}
}

func (c *Client) SendEms(phoneNumber, code string) error {
	name := "点名系统"
	templateCode := "SMS_473010064"
	templateParam := fmt.Sprintf("{\"code\":\"%s\",\"time\":\"3\"}", code)

	req := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  &phoneNumber,
		SignName:      &name,
		TemplateCode:  &templateCode,
		TemplateParam: &templateParam,
	}

	if _, err := c.c.SendSms(req); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed when send %s to %s", code, phoneNumber))
	}

	return nil
}
