package log

import "github.com/gookit/color"

func Highlight(s interface{}) string {
	return color.Gray.Sprintf("%s", s)
}

func Good(s interface{}) string {
	return color.LightGreen.Sprintf("%s", s)
}

func Warning(s interface{}) string {
	return color.LightYellow.Sprintf("%s", s)
}

func Danger(s interface{}) string {
	return color.LightRed.Sprintf("%s", s)
}
