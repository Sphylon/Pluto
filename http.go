package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func Send(req *http.Request) (Result, error) {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while requesting auth:%v", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading %v:%v", res, err)
	}

	var result Result
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal:%v", err)
	}

	return result, nil
}

func Req(base string, method string, header Header, param Param, body Body) (*http.Request, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("error while parsing %s:%v", base, err)
	}

	q := u.Query()
	for key, content := range param {
		q.Add(key, content)
	}
	u.RawQuery = q.Encode()

	b, err := body.Marshal()
	if err != nil {
		return nil, fmt.Errorf("error while marshallig body:%v", err)
	}

	req, _ := http.NewRequest("POST", u.String(), bytes.NewBuffer(b))
	for key, content := range header {
		req.Header.Add(key, content)
	}

	return req, nil
}

func Oauth2_Approval(c *Config) (Result, error) {
	if req, err := Req(
		fmt.Sprintf("%s/oauth2/tokenP", c.Url),
		"POST",
		Header{
			"content-type": "application/json",
		},
		Param{},
		Body{
			"grant_type": "client_credentials",
			"appkey":     c.App.Key,
			"appsecret":  c.App.Secret,
		},
	); err != nil {
		return nil, fmt.Errorf("error while making request:%v", err)
	} else {
		return Send(req)
	}
}

func Quote_Inquire_Price(no string, c *Config) (Result, error) {
	if req, err := Req(
		fmt.Sprintf("%s/uapi/domestic-stock/v1/quotations/inquire-price", c.Url),
		"GET",
		Header{
			"content-type":  "application/json",
			"authorization": fmt.Sprintf("bearer %s", c.Access.Token),
			"appKey":        c.App.Key,
			"appSecret":     c.App.Secret,
			"tr_id":         "FHKST01010100",
		},
		Param{
			"FID_COND_MRKT_DIV_CODE": "UN",
			"FID_INPUT_ISCD":         no,
		},
		Body{},
	); err != nil {
		return nil, fmt.Errorf("error while making request:%v", err)
	} else {
		return Send(req)
	}
}

func Trading_Order_Cash(no string, qty, price int, isBuy bool, c *Config) (Result, error) {
	var tr_id string
	if isBuy {
		tr_id = "TTTC0802U"
	} else {
		tr_id = "TTTC0801U"
	}

	if req, err := Req(
		fmt.Sprintf("%s/uapi/domestic-stock/v1/trading/order-cash", c.Url),
		"POST",
		Header{
			"content-type":  "application/json",
			"authorization": fmt.Sprintf("bearer %s", c.Access.Token),
			"appKey":        c.App.Key,
			"appSecret":     c.App.Secret,
			"tr_id":         tr_id,
			"custtype":      "P",
		},
		Param{},
		Body{
			"CANO":         c.Account.Address,
			"ACNT_PRDT_CD": c.Account.ProductCode,
			"PDNO":         no,
			"ORD_DVSN":     "00",
			"ORD_QTY":      strconv.Itoa(qty),
			"ORD_UNPR":     strconv.Itoa(price),
		},
	); err != nil {
		return nil, fmt.Errorf("error while making request:%v", err)
	} else if result, err := Send(req); err != nil {
		return nil, fmt.Errorf("error while sending request:%v", err)
	} else if result["rt_cd"].(string) != "0" {
		return result, fmt.Errorf("request unsuccessful:%v", result["msg1"])
	} else {
		return result, nil
	}
}
