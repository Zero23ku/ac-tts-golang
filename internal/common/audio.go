package common

import "math/rand"

var Pitch = 1.2
var IsPitchRandom = false

func GetRandomPitch() float64 {

	steps := int((2.0-0.4)/0.1) + 1

	randomSteps := rand.Intn(steps)

	return 0.4 + float64(randomSteps)*0.1

}
