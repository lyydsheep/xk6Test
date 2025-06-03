package enum

import (
	"fmt"
	"testing"
)

func TestSpeedToTime(t *testing.T) {
	duration, _ := SpeedToTime("1,2")
	fmt.Println(duration)
	duration, _ = SpeedToTime("16,1")
	fmt.Println(duration)
}

func TestSpeed(t *testing.T) {
	newSpeed := SpeedUp("1,2")
	fmt.Println(SpeedToTime(newSpeed))
	res := "16,1"
	for range 20 {
		res = SpeedUp(res)
		fmt.Println(res)
		fmt.Println(SpeedToTime(res))
	}
}
