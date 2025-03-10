package main

import (
	"fmt"
)

type KISD struct {
	config *Config
}

func NewKISD(isMock bool) (*KISD, error) {
	config, err := LoadConfig(".", false)
	if err != nil {
		return nil, fmt.Errorf("faile to load config:%v", err)
	}

	return &KISD{
		config: config,
	}, nil
}

func (k *KISD) Auth() error {
	data, err := Oauth2_Approval(k.config)
	if err != nil {
		return fmt.Errorf("failed to get access token:%v", err)
	} else if data["access_token"] != nil {
		k.config.Access.Token, k.config.Access.Expired = data["access_token"].(string),
			data["access_token_token_expired"].(string)
	} else {
		return fmt.Errorf("no access token in response")
	}

	return nil
}

func (k *KISD) GetPrice(no string) error {
	data, err := Quote_Inquire_Price(no, k.config)
	if err != nil {
		return fmt.Errorf("failed to get price", err)
	} else {
		fmt.Println(data)
		return nil
	}
}

func (k *KISD) Order(no string, isBuy bool) error {
	data, err := Trading_Order_Cash(no, isBuy, k.config)
	if err != nil {
		return fmt.Errorf("failed to order", err)
	} else {
		fmt.Println(data)
		return nil
	}
}
