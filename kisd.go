package main

import (
	"fmt"
	"time"
)

type KISD struct {
	config   *Config
	stocks   map[string]*Stock
	savePath string
}

func (k *KISD) NeedAuth() bool {
	if k.config.Access.Token == "" {
		return true
	} else {
		location, err := time.LoadLocation("Asia/Seoul")
		if err != nil {
			return false
		}

		expiredTime, err := time.ParseInLocation(TimeFormat, k.config.Access.Expired, location)
		if err != nil {
			return true
		}

		return !expiredTime.After(time.Now().In(location))
	}
}

func (k *KISD) Auth() error {
	data, err := Oauth2_Approval(k.config)
	if err != nil {
		return fmt.Errorf("failed to get access token:%v", err)
	} else if data["access_token"] != nil {
		k.config.Access.Token, k.config.Access.Expired = data["access_token"].(string),
			data["access_token_token_expired"].(string)
	} else {
		return fmt.Errorf("no access token in response:%v", data)
	}

	if err := k.config.Save2File(k.savePath); err != nil {
		return fmt.Errorf("failed to save configuration:%v", err)
	}

	return nil
}

func (k *KISD) Info(no string) error {
	data, err := Quote_Inquire_Price(no, k.config)
	if err != nil {
		return fmt.Errorf("failed to get price:%v", err)
	} else {
		fmt.Println(data)
		return nil
	}
}

func (k *KISD) OrderCash(no string, qty, price int, isBuy bool) error {
	data, err := Trading_Order_Cash(no, qty, price, isBuy, k.config)
	if data == nil {
		return fmt.Errorf("failed to order:%v", err)
	} else if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (k *KISD) AddStock(code string) error {
	stock := NewStock(code)
	if err := stock.Sync(k.config); err != nil {
		return fmt.Errorf("failed to get stock(%s) information:%v", code, err)
	} else {
		k.config.AddCode(code)
		k.stocks[code] = stock
	}

	k.config.Save2File(k.savePath)
	return nil
}

func NewKISD(isMock bool, savePath string) (*KISD, error) {
	// Load Configuration
	config, err := LoadConfig(savePath, false)
	if err != nil {
		return nil, fmt.Errorf("faile to load config:%v", err)
	}
	kisd := KISD{config: config, savePath: savePath}

	// Authorization
	if kisd.NeedAuth() {
		if err = kisd.Auth(); err != nil {
			return nil, fmt.Errorf("failed in authorization:%v", err)
		}
	}

	// Load and Synchronize Information
	kisd.stocks = make(map[string]*Stock)
	for code := range kisd.config.Codes {
		kisd.stocks[code] = NewStock(code)
		kisd.stocks[code].Sync(kisd.config)
	}

	return &kisd, nil
}
