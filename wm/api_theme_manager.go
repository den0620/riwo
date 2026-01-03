package wm

var themeMap = map[string]map[string]string{
	"monochrome": {
		"faded":  "#ffffff",
		"normal": "#777777",
		"vivid":  "#000000",
	},
	"red": {
		"faded":  "#ffeaea",
		"normal": "#df9595",
		"vivid":  "#bb5d5d",
	},
	"green": {
		"faded":  "#eaffea",
		"normal": "#88cc88",
		"vivid":  "#448844",
	},
	"blue": {
		"faded":  "#c0eaff",
		"normal": "#00aaff",
		"vivid":  "#0088cc",
	},
	"yellow": {
		"faded":  "#ffffea",
		"normal": "#eeee9e",
		"vivid":  "#99994c",
	},
	"aqua": {
		"faded":  "#eaffff",
		"normal": "#9eeeee",
		"vivid":  "#8888cc",
	},
	"gray": {
		"faded":  "#eeeeee",
		"normal": "#cccccc",
		"vivid":  "#888888",
	},
}

func GetTheme(key string) map[string]string {
	return themeMap[key]
}

// GetBackgroundColorStr -> map[key][FADED] "#COLOR"
func GetBackgroundColorStr(key string) string {
	return themeMap[key]["faded"]
}

// GetFontColorStr -> map[key][NORMAL] "#color"
func GetFontColorStr(key string) string {
	return themeMap[key]["normal"]
}

// GetBorderColorStr -> map[key][VIVID] "#color"
func GetBorderColorStr(key string) string {
	return themeMap[key]["vivid"]
}

func ApplyThemeToWindow(window *RiwoWindow, key string) {
	// extract the colors
	bg := themeMap[key]["faded"]
	mg := themeMap[key]["vivid"]
	// fg := wm.themeMap["green"]["normal"]
	window.Content.
		Style("background", bg). // <-- apply the borders/foreground font style
		Style("borderColor", mg)

}

func ApplyThemeToObject(object *RiwoObject, key string) {
	bg := GetBackgroundColorStr(key)
	fg := GetFontColorStr(key)
	object.
		Style("background", bg).
		Style("foreground", fg)
}

func GetThemesMap() *map[string]map[string]string {
	return &themeMap
}
