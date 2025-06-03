package enum

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Speed struct {
	Speed string
	Up    string
	Down  string
}

var (
	MapStringToSpeed map[string]Speed
	MapTimeToString  map[time.Duration]string
)

func init() {
	MapStringToSpeed = make(map[string]Speed)
	MapTimeToString = make(map[time.Duration]string)
	// 预定义速度 (x,1) x秒 1 封
	for i := 1; i <= 16; i++ {
		//dur := time.Second * time.Duration(i)
		spd := Speed{
			Speed: fmt.Sprintf("%d,1", i),
			Up:    fmt.Sprintf("%d,1", i-1),
			Down:  fmt.Sprintf("%d,1", i+1),
		}
		if i == 1 {
			spd.Up = "1,2"
		}
		if i == 16 {
			spd.Down = "16,1"
		}
		dur, err := SpeedToTime(spd.Speed)
		if err != nil {
			panic(err)
		}
		MapStringToSpeed[spd.Speed] = spd
		MapTimeToString[dur] = spd.Speed
	}

	// 扩展速度
	addSpeed(Speed{
		Speed: "1,2",
		Up:    "1,2",
		Down:  "1,1",
	})
}

func addSpeed(speed Speed) {
	MapStringToSpeed[speed.Speed] = speed
	dur, err := SpeedToTime(speed.Speed)
	if err != nil {
		panic(err)
	}
	MapTimeToString[dur] = speed.Speed
}

func SpeedToTime(speed string) (time.Duration, error) {
	if !check(speed) {
		return time.Second * 16, errors.New("invalid speed format")
	}
	strs := strings.Split(speed, ",")
	seconds, _ := strconv.Atoi(strs[0])
	count, _ := strconv.Atoi(strs[1])
	return time.Second * time.Duration(seconds) / time.Duration(count), nil
}

// 校验 speed 是否合法
func check(speed string) bool {
	strs := strings.Split(speed, ",")
	if len(strs) != 2 {
		return false
	}
	seconds, err := strconv.Atoi(strs[0])
	if err != nil {
		return false
	}
	if seconds == 0 {
		return false
	}
	count, err := strconv.Atoi(strs[1])
	if err != nil {
		return false
	}
	if count == 0 {
		return false
	}
	return true
}

// SlowDown 减速
func SlowDown(speed string) string {
	return MapStringToSpeed[speed].Down
}

func SpeedUp(speed string) string {
	return MapStringToSpeed[speed].Up
}
