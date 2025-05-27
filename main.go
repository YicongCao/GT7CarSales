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
	Style       string `json:"style"`
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

	// 输出全部数据表格或文本
	usedTable, legendTable := logic.FormatSampleDataTable(today, cfg.Style)
	fmt.Println(legendTable)
	fmt.Println(usedTable)
	if cfg.EnableWxBot {
		if err := wxwork.SendBotMarkdown(cfg.WxBotAPIKey, "GT7 在售二手车\n\n"+usedTable+"\n"); err != nil {
			log.Printf("推送企业微信失败: %v", err)
		}
		if err := wxwork.SendBotMarkdown(cfg.WxBotAPIKey, "GT7 在售传奇车\n\n"+legendTable+"\n"); err != nil {
			log.Printf("推送企业微信失败: %v", err)
		}
	}

	// 输出新车表格或文本
	if yesterday != nil {
		usedNew, legendNew := logic.FormatNewCarsTable(today, yesterday, cfg.Style)
		fmt.Println(legendNew)
		fmt.Println(usedNew)
		if cfg.EnableWxBot {
			if err := wxwork.SendBotMarkdown(cfg.WxBotAPIKey, "GT7 今日新上架二手车\n\n"+usedNew+"\n"); err != nil {
				log.Printf("推送企业微信失败: %v", err)
			}
			if err := wxwork.SendBotMarkdown(cfg.WxBotAPIKey, "GT7 今日新上架传奇车\n\n"+legendNew+"\n"); err != nil {
				log.Printf("推送企业微信失败: %v", err)
			}
		}
	}
}
