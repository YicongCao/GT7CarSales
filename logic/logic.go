package logic

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gt7_car_sales/fetcher"
	"gt7_car_sales/protocol"
)

// FetchAndParseData 从指定 URL 抓取并解析数据
func FetchAndParseData(url string) (*protocol.SampleData, error) {
	data, err := fetcher.FetchJSONFromURL(url)
	if err != nil {
		return nil, err
	}
	return protocol.ParseSampleData(data)
}

// FetchAndParseDataWithHistory 拉取数据并保存到 history 目录，返回今日和昨日数据
func FetchAndParseDataWithHistory(url string) (today *protocol.SampleData, yesterday *protocol.SampleData, err error) {
	data, err := fetcher.FetchJSONFromURL(url)
	if err != nil {
		return nil, nil, err
	}
	today, err = protocol.ParseSampleData(data)
	if err != nil {
		return nil, nil, err
	}

	// 获取今日日期（如 25-05-27）
	dateStr := today.Used.Date
	if len(dateStr) != 8 {
		return today, nil, nil // 日期格式不符，跳过历史
	}

	// 保存今日数据到 history/gt7cars+{日期}.json
	historyDir := "history"
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return today, nil, err
	}
	filename := filepath.Join(historyDir, "gt7cars_"+dateStr+".json")
	_ = os.WriteFile(filename, data, 0644)

	// 计算昨日日期
	todayTime, err := time.Parse("06-01-02", dateStr)
	if err != nil {
		return today, nil, nil // 日期格式不符，跳过历史
	}
	yesterdayTime := todayTime.AddDate(0, 0, -1)
	yesterdayStr := yesterdayTime.Format("06-01-02")
	yesterdayFile := filepath.Join(historyDir, "gt7cars_"+yesterdayStr+".json")

	// 读取昨日数据
	yesterdayData, err := os.ReadFile(yesterdayFile)
	if err != nil {
		return today, nil, nil // 没有昨日数据
	}
	yesterday, err = protocol.ParseSampleData(yesterdayData)
	if err != nil {
		return today, nil, nil // 解析失败
	}
	return today, yesterday, nil
}

// FormatSampleDataTable 将结构化数据输出为字符画表格或文本（不含每日比赛）
func FormatSampleDataTable(sd *protocol.SampleData, style string) string {
	var sb strings.Builder
	if style == "text" {
		sb.WriteString("### Legendary Cars ###\n\n")
		sb.WriteString(drawCarTextSorted(sd.GetLegendCars()))
		sb.WriteString("\n### Used Cars ###\n\n")
		sb.WriteString(drawCarTextSorted(sd.GetUsedCars()))
		return sb.String()
	}
	// 默认表格
	sb.WriteString("### Legendary Cars ###\n\n")
	sb.WriteString(drawCarTableSorted(sd.GetLegendCars()))
	sb.WriteString("\n### Used Cars ###\n\n")
	sb.WriteString(drawCarTableSorted(sd.GetUsedCars()))
	return sb.String()
}

// FormatNewCarsTable 输出今日新上架的二手车、传奇车（昨天没有的），支持表格和文本
func FormatNewCarsTable(today, yesterday *protocol.SampleData, style string) string {
	var sb strings.Builder
	if style == "text" {
		sb.WriteString("### Today New Legendary Cars ###\n\n")
		sb.WriteString(drawCarTextSorted(diffCarList(today.GetLegendCars(), yesterday.GetLegendCars())))
		sb.WriteString("\n### Today New Used Cars ###\n\n")
		sb.WriteString(drawCarTextSorted(diffCarList(today.GetUsedCars(), yesterday.GetUsedCars())))
		return sb.String()
	}
	// 默认表格
	sb.WriteString("### Today New Legendary Cars ###\n\n")
	sb.WriteString(drawCarTableSorted(diffCarList(today.GetLegendCars(), yesterday.GetLegendCars())))
	sb.WriteString("\n### Today New Used Cars ###\n\n")
	sb.WriteString(drawCarTableSorted(diffCarList(today.GetUsedCars(), yesterday.GetUsedCars())))
	return sb.String()
}

// diffCarList 返回 todayList 中昨天没有的车辆（以 carid+credits 唯一标识，防止同车id不同价格重复）
func diffCarList(todayList, yesterdayList []protocol.Car) []protocol.Car {
	yesterdaySet := make(map[string]struct{}, len(yesterdayList))
	for _, c := range yesterdayList {
		key := c.CarID + "_" + fmt.Sprintf("%d", c.Credits)
		yesterdaySet[key] = struct{}{}
	}
	var diff []protocol.Car
	for _, c := range todayList {
		key := c.CarID + "_" + fmt.Sprintf("%d", c.Credits)
		if _, found := yesterdaySet[key]; !found {
			diff = append(diff, c)
		}
	}
	return diff
}

// 排序后输出车辆表格，二手车/传奇车通用
func drawCarTableSorted(cars []protocol.Car) string {
	// 先 new=true，再按价格降序
	sorted := make([]protocol.Car, len(cars))
	copy(sorted, cars)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].New != sorted[j].New {
			return sorted[i].New // true在前
		}
		return sorted[i].Credits > sorted[j].Credits
	})

	// 英文表头
	headers := []string{"Maker", "Name", "Price", "State", "Special"}
	colWidths := []int{10, 25, 10, 8, 20}

	var sb strings.Builder
	sb.WriteString(drawTableLine(colWidths))
	sb.WriteString(drawTableRow(headers, colWidths))
	sb.WriteString(drawTableLine(colWidths))

	for _, car := range sorted {
		name := car.Name
		if car.New {
			name = "*" + name
		}
		special := joinSpecial(car)
		row := []string{
			car.Manufacturer,
			name,
			fmt.Sprintf("%d", car.Credits),
			car.State,
			special,
		}
		sb.WriteString(drawTableRow(row, colWidths))
	}
	sb.WriteString(drawTableLine(colWidths))
	return sb.String()
}

// 文本样式输出车辆列表
func drawCarTextSorted(cars []protocol.Car) string {
	if len(cars) == 0 {
		return "无\n"
	}
	// 先 new=true，再按价格降序
	sorted := make([]protocol.Car, len(cars))
	copy(sorted, cars)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].New != sorted[j].New {
			return sorted[i].New // true在前
		}
		return sorted[i].Credits > sorted[j].Credits
	})
	var sb strings.Builder
	for _, car := range sorted {
		sb.WriteString(fmt.Sprintf("[%s]\n", car.Manufacturer))
		name := car.Name
		if car.New {
			name = "*" + name
		}
		sb.WriteString(fmt.Sprintf("Model: %s, Price: %s\n", name, formatCarPrice(car)))
		if car.State != "normal" {
			sb.WriteString(fmt.Sprintf("Status: %s\n", formatCarState(car.State)))
		}
		remark := joinSpecial(car)
		if strings.TrimSpace(remark) != "" {
			sb.WriteString(fmt.Sprintf("Remarks: %s\n", remark))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// 格式化价格文本
func formatCarPrice(car protocol.Car) string {
	if car.Credits >= 10000 {
		return fmt.Sprintf("%.1f 万", float64(car.Credits)/10000)
	}
	return fmt.Sprintf("%d", car.Credits)
}

// 格式化状态文本
func formatCarState(state string) string {
	switch state {
	case "normal":
		return state
	case "limited":
		return state
	case "soldout":
		return state
	default:
		return state
	}
}

// 合并奖励、奖杯、抽奖为一列
func joinSpecial(car protocol.Car) string {
	var parts []string
	if s := rewardCarStr(car.RewardCar); s != "" {
		parts = append(parts, s)
	}
	if s := nullOrStr(car.TrophyCar); s != "" {
		parts = append(parts, s)
	}
	if s := nullOrStr(car.LotteryCar); s != "" {
		parts = append(parts, s)
	}
	return strings.Join(parts, " / ")
}

// 工具函数
func drawTableLine(widths []int) string {
	var sb strings.Builder
	sb.WriteString("+")
	for _, w := range widths {
		sb.WriteString(strings.Repeat("-", w))
		sb.WriteString("+")
	}
	sb.WriteString("\n")
	return sb.String()
}

func drawTableRow(cols []string, widths []int) string {
	var sb strings.Builder
	sb.WriteString("|")
	for i, col := range cols {
		w := widths[i]
		// 截断过长内容
		colStr := col
		if len([]rune(colStr)) > w-1 {
			colStr = string([]rune(colStr)[:w-2]) + "…"
		}
		sb.WriteString(fmt.Sprintf("%-*s", w, colStr))
		sb.WriteString("|")
	}
	sb.WriteString("\n")
	return sb.String()
}

func boolToStr(b bool) string {
	if b {
		return "是"
	}
	return ""
}

func engineSwapStr(es *protocol.EngineSwap) string {
	if es == nil {
		return ""
	}
	return es.EngineName
}

func rewardCarStr(rc *protocol.RewardCar) string {
	if rc == nil {
		return ""
	}
	return rc.Type + ":" + rc.Name
}

func nullOrStr(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}
