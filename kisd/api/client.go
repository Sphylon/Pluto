package api

import (
	"encoding/json"
	"fmt"
)

const TimeFormat = "2006-01-02 15:04:05" // Golang의 특별한 시간 형식

type Param map[string]string
type Header map[string]string
type Body map[string]interface{}

func (b Body) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

type Client struct {
	base               string
	appKey             string
	appSecret          string
	accountSetress     string
	accountProductCode string
	token              string
	expired            string
}

func NewClient(url, key, secret, address, productCode string) *Client {
	return &Client{
		base:               url,
		appKey:             key,
		appSecret:          secret,
		accountSetress:     address,
		accountProductCode: productCode,
	}
}

type Oauth2_Approval_Result struct {
	access_token               string `json:"access_token"`
	access_token_token_expired string `json:"access_token_token_expired"`
}

func (c *Client) Oauth2_Approval() error {
	req := NewRequest(*c, "oauth2/tokenP", "POST", false)
	req.SetBody("grant_type", "client_credentials")
	req.SetBody("appkey", c.appKey)
	req.SetBody("appsecret", c.appSecret)

	result := Oauth2_Approval_Result{}
	if err := req.Send(&result); err != nil {
		return fmt.Errorf("error while sending Oauth2_Approval_Result request:%v", err)
	}

	c.token = result.access_token
	c.expired = result.access_token_token_expired

	return nil
}

// func Quote_Inquire_Price(no string, c *Config) (Result, error) {
// 	if req, err := Req(
// 		fmt.Sprintf("%s/uapi/domestic-stock/v1/quotations/inquire-price", c.Url),
// 		"GET",
// 		Header{
// 			"content-type":  "application/json",
// 			"authorization": fmt.Sprintf("bearer %s", c.Access.Token),
// 			"appKey":        c.App.Key,
// 			"appSecret":     c.App.Secret,
// 			"tr_id":         "FHKST01010100",
// 		},
// 		Param{
// 			"FID_COND_MRKT_DIV_CODE": "UN",
// 			"FID_INPUT_ISCD":         no,
// 		},
// 		Body{},
// 	); err != nil {
// 		return nil, fmt.Errorf("error while making request:%v", err)
// 	} else if result, err := Send(req); err != nil {
// 		return nil, fmt.Errorf("error while sending request:%v", err)
// 	} else if result["rt_cd"].(string) != "0" {
// 		return result, fmt.Errorf("request unsuccessful:%v", result["msg1"])
// 	} else {
// 		return result.Parse()
// 	}
// }

// func Trading_Order_Cash(no string, qty, price int, isBuy bool, c *Config) (Result, error) {
// 	var tr_id string
// 	if isBuy {
// 		tr_id = "TTTC0802U"
// 	} else {
// 		tr_id = "TTTC0801U"
// 	}

// 	if req, err := Req(
// 		fmt.Sprintf("%s/uapi/domestic-stock/v1/trading/order-cash", c.Url),
// 		"POST",
// 		Header{
// 			"content-type":  "application/json",
// 			"authorization": fmt.Sprintf("bearer %s", c.Access.Token),
// 			"appKey":        c.App.Key,
// 			"appSecret":     c.App.Secret,
// 			"tr_id":         tr_id,
// 			"custtype":      "P",
// 		},
// 		Param{},
// 		Body{
// 			"CANO":         c.Account.Setress,
// 			"ACNT_PRDT_CD": c.Account.ProductCode,
// 			"PDNO":         no,
// 			"ORD_DVSN":     "00",
// 			"ORD_QTY":      strconv.Itoa(qty),
// 			"ORD_UNPR":     strconv.Itoa(price),
// 		},
// 	); err != nil {
// 		return nil, fmt.Errorf("error while making request:%v", err)
// 	} else if result, err := Send(req); err != nil {
// 		return nil, fmt.Errorf("error while sending request:%v", err)
// 	} else if result["rt_cd"].(string) != "0" {
// 		return result, fmt.Errorf("request unsuccessful:%v", result["msg1"])
// 	} else {
// 		return result.Parse()
// 	}
// }

// func Trading_Inquire_Balance(c *Config) (Result, error) {
// 	var tr_id string
// 	if c.IsMock {
// 		tr_id = "VTTC8434R"
// 	} else {
// 		tr_id = "TTTC8434R"
// 	}

// 	if req, err := Req(
// 		fmt.Sprintf("%s/uapi/domestic-stock/v1/trading/inquire-balance", c.Url),
// 		"GET",
// 		Header{
// 			"authorization": fmt.Sprintf("bearer %s", c.Access.Token),
// 			"appKey":        c.App.Key,
// 			"appSecret":     c.App.Secret,
// 			"tr_id":         tr_id,
// 		},
// 		Param{
// 			"CANO":                  c.Account.Setress,
// 			"ACNT_PRDT_CD":          c.Account.ProductCode,
// 			"AFHR_FLPR_YN":          "N",
// 			"OFL_YN":                "",
// 			"INQR_DVSN":             "01",
// 			"UNPR_DVSN":             "01",
// 			"FUND_STTL_ICLD_YN":     "N",
// 			"FNCG_AMT_AUTO_RDPT_YN": "N",
// 			"PRCS_DVSN":             "01",
// 			"CTX_AREA_FK100":        "",
// 			"CTX_AREA_NK100":        "",
// 		},
// 		Body{},
// 	); err != nil {
// 		return nil, fmt.Errorf("error while making request:%v", err)
// 	} else if result, err := Send(req); err != nil {
// 		return nil, fmt.Errorf("error while sending request:%v", err)
// 	} else if result["rt_cd"].(string) != "0" {
// 		return result, fmt.Errorf("request unsuccessful:%v", result["msg1"])
// 	} else {
// 		return result.Parse()
// 	}
// }
