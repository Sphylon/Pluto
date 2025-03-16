package main

import (
	"fmt"
	"strconv"
)

const (
	Daily = iota
	Weekly
	Monthly
)

type Stock struct {
	Name   string
	Code   string
	Avail  bool
	Prices [3]Price
	Acmls  [3]Acml
}

type Price struct {
	Price   int
	Changes float64
}

type Acml struct {
	Volume  int
	Value   int
	Changes float64
}

func NewStock(code string) *Stock {
	return &Stock{Code: code}
}

func (s *Stock) Sync(config *Config) error {
	rawData, err := Quote_Inquire_Price(s.Code, config)
	if err != nil {
		return fmt.Errorf("failed to get stock(%s) data:%v", s.Code, err)
	}

	data := rawData["output"].(Result)

	// s.Name = data["bstp_kor_isnm"].(string)        // 종목명
	s.Avail = data["temp_stop_yn"].(string) != "Y" // 거래 가능 여부

	if s.Prices[Daily].Price, err = strconv.Atoi(data["stck_prpr"].(string)); err != nil {
		return fmt.Errorf("failed to get current daily stock Price :%v", err)
	}
	if s.Prices[Daily].Changes, err = strconv.ParseFloat(data["prdy_ctrt"].(string), 64); err != nil {
		return fmt.Errorf("failed to get current daily stock Price chagne ratio:%v", err)
	}
	if data["prdy_vrss_sign"].(string) != "2" {
		s.Prices[Daily].Changes *= -1
	}

	if s.Acmls[Daily].Volume, err = strconv.Atoi(data["acml_vol"].(string)); err != nil {
		return fmt.Errorf("failed to get current daily accumulated volume:%v", err)
	}
	if s.Acmls[Daily].Value, err = strconv.Atoi(data["acml_tr_pbmn"].(string)); err != nil {
		return fmt.Errorf("failed to get current daily accumulated value:%v", err)
	}
	if s.Acmls[Daily].Changes, err = strconv.ParseFloat(data["prdy_vrss_vol_rate"].(string), 64); err != nil {
		return fmt.Errorf("failed to get current daily accumlated ratio:%v", err)
	}

	return nil
}
