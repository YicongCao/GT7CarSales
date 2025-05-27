package protocol

import (
	"encoding/json"
	"errors"
)

// 主结构体
type SampleData struct {
	UpdateTimestamp string           `json:"updatetimestamp"`
	Used            UsedSection      `json:"used"`
	Legend          LegendSection    `json:"legend"`
	DailyRace       DailyRaceSection `json:"dailyrace"`
}

// Used/Legend 结构体
type UsedSection struct {
	Date string `json:"date"`
	Cars []Car  `json:"cars"`
}
type LegendSection struct {
	Date string `json:"date"`
	Cars []Car  `json:"cars"`
}

// 车辆结构体
type Car struct {
	CarID           string      `json:"carid"`
	Manufacturer    string      `json:"manufacturer"`
	Region          string      `json:"region"`
	Name            string      `json:"name"`
	Credits         int         `json:"credits"`
	State           string      `json:"state"`
	EstimateDays    int         `json:"estimatedays"`
	MaxEstimateDays int         `json:"maxestimatedays"`
	New             bool        `json:"new"`
	RewardCar       *RewardCar  `json:"rewardcar"`
	EngineSwap      *EngineSwap `json:"engineswap"`
	LotteryCar      interface{} `json:"lotterycar"` // string or null
	TrophyCar       interface{} `json:"trophycar"`  // string or null
}

// 奖励车结构体
type RewardCar struct {
	Type        string      `json:"type"`
	Name        string      `json:"name"`
	Requirement interface{} `json:"requirement"`
}

// 引擎互换结构体
type EngineSwap struct {
	CarID        string `json:"carid"`
	Manufacturer string `json:"manufacturer"`
	Region       string `json:"region"`
	Name         string `json:"name"`
	EngineName   string `json:"enginename"`
}

// 日常赛事结构体
type DailyRaceSection struct {
	Date  string      `json:"date"`
	Races []DailyRace `json:"races"`
}
type DailyRace struct {
	CourseID             int           `json:"courseid"`
	CrsBase              string        `json:"crsbase"`
	Track                string        `json:"track"`
	Logo                 string        `json:"logo"`
	Region               string        `json:"region"`
	Laps                 int           `json:"laps"`
	Cars                 int           `json:"cars"`
	StartType            string        `json:"starttype"`
	FuelCons             int           `json:"fuelcons"`
	TyreWear             int           `json:"tyrewear"`
	CarType              string        `json:"cartype"`
	WideBodyBan          bool          `json:"widebodyban"`
	NitrousBan           bool          `json:"nitrousban"`
	Tyres                []interface{} `json:"tyres"`
	RequiredTyres        []interface{} `json:"requiredtyres"`
	BOP                  bool          `json:"bop"`
	CarSettingsSpecified bool          `json:"carsettings_specified"`
	GarageCar            bool          `json:"garagecar"`
	CarUsed              bool          `json:"carused"`
	Damage               bool          `json:"damage"`
	ShortcutPen          bool          `json:"shortcutpen"`
	CarCollisionPen      bool          `json:"carcollisionpen"`
	PitLanePen           bool          `json:"pitlanepen"`
	Time                 int           `json:"time"`
	Offset               int           `json:"offset"`
	Schedule             string        `json:"schedule"`
}

// 反序列化函数
func ParseSampleData(data []byte) (*SampleData, error) {
	var sd SampleData
	if err := json.Unmarshal(data, &sd); err != nil {
		return nil, err
	}
	return &sd, nil
}

// 辅助函数：获取所有二手车
func (sd *SampleData) GetUsedCars() []Car {
	return sd.Used.Cars
}

// 辅助函数：获取所有传奇车
func (sd *SampleData) GetLegendCars() []Car {
	return sd.Legend.Cars
}

// 辅助函数：获取所有车辆（合并二手和传奇）
func (sd *SampleData) GetAllCars() []Car {
	all := make([]Car, 0, len(sd.Used.Cars)+len(sd.Legend.Cars))
	all = append(all, sd.Used.Cars...)
	all = append(all, sd.Legend.Cars...)
	return all
}

// 辅助函数：根据 carid 查找车辆
func (sd *SampleData) FindCarByID(carid string) (*Car, error) {
	for _, c := range sd.GetAllCars() {
		if c.CarID == carid {
			return &c, nil
		}
	}
	return nil, errors.New("car not found")
}
