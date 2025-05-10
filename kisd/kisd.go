package kisd

import (
	"fmt"

	"github.com/Sphylon/Pluto/kisd/api"
)

type KISD struct {
	client    *api.Client
	config    *service.Configuration
	trading   *service.Trading
	quotation *service.Quotation
	savePath  string
}

func NewKISD(isMock bool, savePath string) (*KISD, error) {
	// Load Configuration
	config, err := LoadConfig(savePath, false)
	if err != nil {
		return nil, fmt.Errorf("faile to load config:%v", err)
	}
	kisd := KISD{config: config, savePath: savePath}

	// Authorization
	if kisd.IsAuthNeeded() {
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
