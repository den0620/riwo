package apps

import (
	"riwo/wm"
	"strconv"
	"syscall/js"
	"time"
)

func init() {
	AppRegistry["Mahjongg"] = mahjonggConstruct
}

type brick struct {
	typ     int
	elem    *wm.RiwoObject
	layer   int
	row     int
	col     int
	blocked bool
	removed bool
}

var tiles = map[string][]string{
	"Dots":    {"?", "⢀", "⣀", "⣠", "⣤", "⣴", "⣶", "⣾", "⣿", "⑨"},
	"Bamboo":  {"?", "🍩", "🥯", "🥨", "🍕", "🥪", "🌮", "🌭", "🍔", "🍟"},
	"Chars":   {"?", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
	"Winds":   {"?", "←", "↑", "→", "↓"},
	"Dragons": {"?", "🔴", "🟢", "🔵"},
	"Flowers": {"?", "🌺", "🌼", "🌸", "🍀"},
}

var layouts = map[string][][]int{
	"Classic": {
		{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
		{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
		{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
	},
	"Fortress": {
		{1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1},
		{1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1},
		{1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1},
		{0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1},
		{1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1},
		{1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1},
	},
	"Pyramid": {
		{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
		{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
		{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	},
	"Cross": {
		{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
	},
	"Arena": {
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
		{1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
		{1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	},
}

var layoutNames = []string{"Classic", "Fortress", "Pyramid", "Cross", "Arena"}
var themeNames = []string{"yellow", "blue", "green", "red", "purple", "aqua", "orange", "pink"}

func mahjonggConstruct(window *wm.RiwoWindow) {
	themeIdx := 0
	layoutIdx := 0
	theme := wm.ThemeMap[themeNames[themeIdx]]

	var bricks []*brick
	var selected *brick
	var startTime time.Time
	var timerStop bool

	container := wm.Create().
		Style("height", "100%").
		Style("display", "flex").
		Style("flexDirection", "column").
		Style("justifyContent", "center").
		Style("alignItems", "center").
		Style("backgroundColor", theme["faded"])

	timerElem := wm.Create().
		Text("Time: 0:00").
		Style("marginBottom", "1rem").
		Style("color", theme["vivid"]).
		Style("fontWeight", "bold")

	board := wm.Create().
		Style("position", "relative").
		Style("display", "inline-block").
		Style("width", "36rem").
		Style("height", "24rem") // something like that // TODO make adequate dimensions

	container.Append(timerElem, board)
	window.Content.Inner("").Append(container)

	// Timer update
	var updateTimer func()
	updateTimer = func() {
		if timerStop {
			return
		}
		elapsed := time.Since(startTime)
		min := int(elapsed.Minutes())
		sec := int(elapsed.Seconds()) % 60
		secStr := strconv.Itoa(sec)
		if sec < 10 {
			secStr = "0" + secStr
		}
		timerElem.Text("Time: " + strconv.Itoa(min) + ":" + secStr)

		js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			updateTimer()
			return nil
		}), 1000)
	}

	// Update blocked state
	updateBlocked := func() {
		for _, b := range bricks {
			if b.removed {
				continue
			}
			b.blocked = false

			// Check if covered from above
			for _, other := range bricks {
				if !other.removed && other.layer > b.layer && other.row == b.row && other.col == b.col {
					b.blocked = true
					break
				}
			}

			// Check if both sides blocked
			if !b.blocked {
				left, right := false, false
				for _, other := range bricks {
					if other.removed || other.layer != b.layer || other == b {
						continue
					}
					if other.row == b.row {
						if other.col == b.col-1 {
							left = true
						}
						if other.col == b.col+1 {
							right = true
						}
					}
				}
				b.blocked = left && right
			}

			if b.blocked {
				b.elem.Style("filter", "brightness(0.6)").Style("cursor", "not-allowed")
			} else {
				b.elem.Style("filter", "brightness(1)").Style("cursor", wm.CursorInvertUrl)
			}
		}
	}

	// Check win
	checkWin := func() {
		for _, b := range bricks {
			if !b.removed {
				return
			}
		}
		timerStop = true
		elapsed := time.Since(startTime)
		min := int(elapsed.Minutes())
		sec := int(elapsed.Seconds()) % 60
		secStr := strconv.Itoa(sec)
		if sec < 10 {
			secStr = "0" + secStr
		}

		wm.Create().
			Text("You Won!\nTime: " + strconv.Itoa(min) + ":" + secStr).
			Style("position", "absolute").
			Style("top", "50%").
			Style("left", "50%").
			Style("transform", "translate(-50%, -50%)").
			Style("fontSize", "2rem").
			Style("fontWeight", "bold").
			Style("color", theme["vivid"]).
			Style("backgroundColor", theme["normal"]).
			Style("padding", "1.25rem 2.5rem").
			Style("borderRadius", "0.625rem").
			Style("border", "0.1875rem solid "+theme["vivid"]).
			Style("textAlign", "center").
			Style("whiteSpace", "pre-line").
			Style("boxShadow", "0 0.25rem 0.5rem rgba(0,0,0,0.3)").
			Style("zIndex", "1000").
			Mount(board)
	}

	// Check dead end
	checkDeadEnd := func() {
		hasRemaining := false
		for _, b := range bricks {
			if !b.removed {
				hasRemaining = true
				break
			}
		}
		if !hasRemaining {
			return
		}

		// Check for valid moves
		var free []*brick
		for _, b := range bricks {
			if !b.removed && !b.blocked {
				free = append(free, b)
			}
		}

		hasMove := false
		for i := 0; i < len(free); i++ {
			for j := i + 1; j < len(free); j++ {
				// Check match (flowers match any flower)
				match := false
				if free[i].typ >= 600 && free[i].typ < 700 && free[j].typ >= 600 && free[j].typ < 700 {
					match = true
				} else if free[i].typ == free[j].typ {
					match = true
				}
				if match {
					hasMove = true
					break
				}
			}
			if hasMove {
				break
			}
		}

		if !hasMove {
			timerStop = true
			wm.Create().
				Text("No More Moves!").
				Style("position", "absolute").
				Style("top", "50%").
				Style("left", "50%").
				Style("transform", "translate(-50%, -50%)").
				Style("fontSize", "2rem").
				Style("fontWeight", "bold").
				Style("color", theme["vivid"]).
				Style("backgroundColor", theme["normal"]).
				Style("padding", "1.25rem 2.5rem").
				Style("borderRadius", "0.625rem").
				Style("border", "0.1875rem solid "+theme["vivid"]).
				Style("textAlign", "center").
				Style("whiteSpace", "pre-line").
				Style("boxShadow", "0 0.25rem 0.5rem rgba(0,0,0,0.3)").
				Style("zIndex", "1000").
				Mount(board)
		}
	}

	// Click handler
	onBrickClick := func(b *brick) {
		if b.removed || b.blocked {
			return
		}

		if selected == nil {
			selected = b
			b.elem.Style("backgroundColor", theme["vivid"])
		} else if selected == b {
			selected = nil
			b.elem.Style("backgroundColor", theme["normal"])
		} else {
			// Check match
			match := false
			if selected.typ >= 600 && selected.typ < 700 && b.typ >= 600 && b.typ < 700 {
				match = true
			} else if selected.typ == b.typ {
				match = true
			}

			if match {
				selected.removed = true
				selected.elem.Style("display", "none")
				b.removed = true
				b.elem.Style("display", "none")
				selected = nil
				updateBlocked()
				checkWin()
				checkDeadEnd()
			} else {
				selected.elem.Style("backgroundColor", theme["normal"])
				selected = b
				b.elem.Style("backgroundColor", theme["vivid"])
			}
		}
	}

	// Make brick
	makeBrick := func(typ, layer, row, col int) *brick {
		cat := typ / 100
		idx := typ % 100

		var symbol string
		switch cat {
		case 1:
			symbol = tiles["Dots"][idx]
		case 2:
			symbol = tiles["Bamboo"][idx]
		case 3:
			symbol = tiles["Chars"][idx]
		case 4:
			symbol = tiles["Winds"][idx]
		case 5:
			symbol = tiles["Dragons"][idx]
		case 6:
			symbol = tiles["Flowers"][idx]
		default:
			symbol = "?"
		}

		left := float64(col)*2.5 + float64(layer)*0.25
		top := float64(row)*3.125 + float64(layer)*0.25
		zindex := layer*100 + row

		elem := wm.Create().
			Text(symbol).
			Style("position", "absolute").
			Style("left", strconv.FormatFloat(left, 'f', 3, 64)+"rem").
			Style("top", strconv.FormatFloat(top, 'f', 3, 64)+"rem").
			Style("zIndex", strconv.Itoa(zindex)).
			Style("width", "2rem").
			Style("height", "3rem").
			Style("display", "flex").
			Style("justifyContent", "center").
			Style("alignItems", "center").
			Style("fontSize", "1.5rem").
			Style("backgroundColor", theme["normal"]).
			Style("border", "0.125rem solid "+theme["vivid"]).
			Style("cursor", wm.CursorInvertUrl).
			Style("userSelect", "none").
			Style("boxShadow", "0.125rem 0.125rem 0.25rem rgba(0,0,0,0.3)")

		b := &brick{
			typ:   typ,
			elem:  elem,
			layer: layer,
			row:   row,
			col:   col,
		}

		elem.Listen("mousedown", func(this js.Value, args []js.Value) interface{} {
			onBrickClick(b)
			return nil
		})

		return b
	}

	// Start game
	startGame := func() {
		bricks = []*brick{}
		selected = nil
		board.Inner("")
		timerStop = false
		startTime = time.Now()

		layout := layouts[layoutNames[layoutIdx]]

		// Generate tile pool
		pool := []int{}
		for i := 1; i <= 9; i++ {
			for j := 0; j < 4; j++ {
				pool = append(pool, 100+i) // Dots
			}
		}
		for i := 1; i <= 9; i++ {
			for j := 0; j < 4; j++ {
				pool = append(pool, 200+i) // Bamboo
			}
		}
		for i := 1; i <= 9; i++ {
			for j := 0; j < 4; j++ {
				pool = append(pool, 300+i) // Chars
			}
		}
		for i := 1; i <= 4; i++ {
			for j := 0; j < 4; j++ {
				pool = append(pool, 400+i) // Winds
			}
		}
		for i := 1; i <= 3; i++ {
			for j := 0; j < 4; j++ {
				pool = append(pool, 500+i) // Dragons
			}
		}
		for i := 1; i <= 4; i++ {
			pool = append(pool, 600+i) // Flowers
			pool = append(pool, 600+i)
		}

		// Shuffle
		rng := time.Now().UnixNano()
		for i := len(pool) - 1; i > 0; i-- {
			rng = (rng*1103515245 + 12345) & 0x7fffffff
			j := int(rng) % (i + 1)
			pool[i], pool[j] = pool[j], pool[i]
		}

		// Place tiles
		poolIdx := 0
		for layer := 0; layer < 3; layer++ {
			for row := 0; row < len(layout); row++ {
				for col := 0; col < len(layout[row]); col++ {
					if layout[row][col] == 0 {
						continue
					}
					if layer > 0 && (row < 1 || row > 5 || col < 2 || col > 11) {
						continue
					}
					if poolIdx >= len(pool) {
						poolIdx = 0
					}
					b := makeBrick(pool[poolIdx], layer, row, col)
					bricks = append(bricks, b)
					board.Append(b.elem)
					poolIdx++
				}
			}
		}

		updateBlocked()
		updateTimer()
	}

	// Cycle theme
	cycleTheme := func() {
		themeIdx = (themeIdx + 1) % len(themeNames)
		theme = wm.ThemeMap[themeNames[themeIdx]]

		container.Style("backgroundColor", theme["faded"])
		timerElem.Style("color", theme["vivid"])

		for _, b := range bricks {
			if !b.removed {
				if b == selected {
					b.elem.Style("backgroundColor", theme["vivid"])
				} else {
					b.elem.Style("backgroundColor", theme["normal"])
				}
				b.elem.Style("border", "0.125rem solid "+theme["vivid"])
			}
		}
	}
	// Cycle layout
	cycleLayout := func() {
		layoutIdx = (layoutIdx + 1) % len(layoutNames)
		startGame()
	}

	// Context Menu
	window.MenuEntries = []wm.ContextEntry{
		{Name: "New Game", Callback: startGame},
		{Name: "Cycle Theme", Callback: cycleTheme},
		{Name: "Cycle Layout", Callback: cycleLayout},
	}

	startGame()
}
