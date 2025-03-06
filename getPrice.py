import requests
import json
import yaml
from operator import itemgetter

ISREAL = False 
APP_KEY = ""
APP_SECRET = ""
ACCESS_TOKEN = ""
ACCESS_TOKEN_EXPIRED = ""
URL_BASE = "" 

config:dict = {}
configPath = "./config.yaml"

# Auth
def auth():
    global ACCESS_TOKEN, ACCESS_TOKEN_EXPIRED
    if "ACCESS_TOKEN" in config:
        ACCESS_TOKEN, ACCESS_TOKEN_EXPIRED = config["ACCESS_TOKEN"],config["ACCESS_TOKEN_EXPIRED"]
        return 

    headers = {"content-type":"application/json"}
    body = {
        "grant_type":"client_credentials",
        "appkey":APP_KEY, 
        "appsecret":APP_SECRET
        }
    PATH = "oauth2/tokenP"
    URL = f"{URL_BASE}/{PATH}"
    time.sleep(0.1) # 유량제한 예방 (REST: 1초당 20건 제한)
    res:dict = requests.post(URL, headers=headers, data=json.dumps(body)).json()
    
    try: 
        ACCESS_TOKEN, ACCESS_TOKEN_EXPIRED = res["access_token"], res["access_token_token_expired"]
        config["ACCESS_TOKEN"], config["ACCESS_TOKEN_EXPIRED"] = ACCESS_TOKEN, ACCESS_TOKEN_EXPIRED
    except KeyError:
        ACCESS_TOKEN, ACCESS_TOKEN_EXPIRED = config["ACCESS_TOKEN"], config["ACCESS_TOKEN_EXPIRED"]
    except:
        print("unable to get a token ...")
        print(res)
    
    # save changes
    with open(configPath, 'w') as f:
        yaml.dump(config, f)
        

def pre_set():
    global config, APP_KEY, APP_SECRET, URL_BASE
    with open(configPath) as f:
        config = yaml.full_load(f)
        APP_KEY, APP_SECRET = itemgetter("APP_KEY", "APP_SECRET")(config)
        URL_BASE = config["URL_BASE"]["REAL"] if ISREAL else config["URL_BASE"]["MOCK"]
    auth()


# 주식현재가 시세
def get_current_price(stock_no):
    PATH = "uapi/domestic-stock/v1/quotations/inquire-price"
    URL = f"{URL_BASE}/{PATH}"

    # 헤더 설정
    headers = {"Content-Type":"application/json", 
            "authorization": f"Bearer {ACCESS_TOKEN}",
            "appKey":APP_KEY,
            "appSecret":APP_SECRET,
            "tr_id":"FHKST01010100"}

    params = {
        "fid_cond_mrkt_div_code":"J",
        "fid_input_iscd": stock_no
    }

    # 호출
    res = requests.get(URL, headers=headers, params=params)
    if res.status_code == 200 and res.json()["rt_cd"] == "0" :
        return res.json()["output"]
    # 토큰 만료 시
    elif res.status_code == 200 and res["msg_cd"] == "EGW00123" :
        auth()
        get_current_price(stock_no)
    else:
        print("Error Code : " + str(res.status_code) + " | " + res.text)
        return None

def get_current_balance():
    PATH = "/uapi/domestic-stock/v1/trading/inquire-balance"
    URL = f"{URL_BASE}/{PATH}"
    
    headers = {"Content-Type":"application/json", 
            "authorization": f"Bearer {ACCESS_TOKEN}",
            "appKey":APP_KEY,
            "appSecret":APP_SECRET,
            "tr_id":"TTTC8434R" if ISREAL else "VTTC8434R"}

    params = {
        "CANO":config["C_ANO"],
        "ACNT_PRDT_CD": config["ACNT_PRDT_CD"],
        "AFHR_FLPR_YN": "N",
        "OFL_YN": "",
        "INQR_DVSN": "01",
        "UNPR_DVSN": "01",
        "FUND_STTL_ICLD_YN": "N",
        "FNCG_AMT_AUTO_RDPT_YN": "N",
        "PRCS_DVSN" : "00",
        "CTX_AREA_FK100": "",
        "CTX_AREA_NK100" : ""
    }
    
    res = requests.get(URL, headers=headers, params=params).json()

    return res

def order(stock_no):
    PATH = "/uapi/domestic-stock/v1/trading/order-cash"
    URL = f"{URL_BASE}/{PATH}"
    
    headers = {
        "Content-Type":"application/json", 
        "authorization": f"Bearer {ACCESS_TOKEN}",
        "appKey":APP_KEY,
        "appSecret":APP_SECRET,
        "tr_id":"VTTC0802U",
        "custtype": "P"
    }

    params = {
        "CANO":config["C_ANO"],
        "ACNT_PRDT_CD": config["ACNT_PRDT_CD"],
        "PDNO": stock_no,
        "ORD_DVSN": "00",
        "ORD_QTY": "3",
        "ORD_UNPR": "150000"
    } 

    res = requests.post(URL, headers=headers, params=params).json()

    return res


pre_set()
# print(get_current_price("005930"))
# print(get_current_balance())
# for key, value in get_current_balance().items():
#     print(key, value)

for key, value in get_current_price("005930").items():
    print(key, value)

for key, value in order("BBG000BLLVT7").items():
    print(key, value)