package ems

import (
	. "github.com/bytedance/mockey"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/config"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEms(t *testing.T) {
	akID := "your ak id"
	aks := "your aks"
	edp := "dysmsapi.aliyuncs.com"
	conf := &config.Config{Aliyun: &config.Aliyun{}}
	MockValue(&config.Conf).To(conf)

	MockValue(&config.Conf.Aliyun.AccessKeyID).To(akID)
	MockValue(&config.Conf.Aliyun.AccessKeySecret).To(aks)
	MockValue(&config.Conf.Aliyun.Endpoint).To(edp)

	Convey("Test New Ems client", t, func() {
		client := NewEmsClient()
		So(client, ShouldNotBeNil)

		phoneNumber := "you phone number"
		code := "123456"

		Convey("Test SendEms", func() {
			err := client.SendEms(phoneNumber, code)
			So(err, ShouldBeNil)
		})

		Convey("Given a wrong phone number", func() {
			err := client.SendEms("", code)
			So(err, ShouldNotBeNil)
		})
	})
}
