package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gt7_car_sales/logic"
	"gt7_car_sales/wxwork"
)

type Config struct {
	WxBotAPIKey string `json:"wxBotAPIKey"`
	EnableWxBot bool   `json:"enableWxBot"`
}

const dataURL = "https://ddm999.github.io/gt7info/data.json"

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 拉取今日和昨日数据
	today, yesterday, err := logic.FetchAndParseDataWithHistory(dataURL)
	if err != nil {
		log.Fatalf("拉取数据失败: %v", err)
	}

	// 输出全部数据表格
	fullTable := logic.FormatSampleDataTable(today)
	fmt.Println(fullTable)
	if cfg.EnableWxBot {
		if err := wxwork.SendBotMarkdown(cfg.WxBotAPIKey, "GT7车辆信息\n```\n"+fullTable+"\n```"); err != nil {
			log.Printf("推送企业微信失败: %v", err)
		}
	}

	// 输出新车表格
	if yesterday != nil {
		newCarsTable := logic.FormatNewCarsTable(today, yesterday)
		fmt.Println(newCarsTable)
		if cfg.EnableWxBot {
			if err := wxwork.SendBotMarkdown(cfg.WxBotAPIKey, "GT7今日新上架\n```\n"+newCarsTable+"\n```"); err != nil {
				log.Printf("推送企业微信失败: %v", err)
			}
		}
	}
}
