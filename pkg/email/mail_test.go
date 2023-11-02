package email

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"testing"
)

func TestEmail(T *testing.T) {
	if err := InitConfig(); err != nil {
		panic(err)
	}
	mail, err := NewMail()
	if err != nil {
		log.Fatal(err)
	}
	err = mail.SendMail(context.Background(), "text@gmail.com", "code")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Send Successful")

}

func InitConfig() error {
	yam, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yam, &config.Config)
	if err != nil {
		return err
	}
	return nil
}
