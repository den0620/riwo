package apps

import (
	"math/rand"
	"riwo/wm"
	"strconv"
	"syscall/js"
	"time"
)

func init() {
	AppRegistry["Mahjongg"] = mahjonggConstruct
	//rand.Seed(time.Now().UnixNano()) // Deprecated as of Go 1.20
}

var mahjonggLayoutNames = []string{"Classic", "Fortress", "Arena"}
var mahjonggThemeNames = []string{"yellow", "blue", "green", "red", "purple", "aqua", "orange", "pink"}

type mahjonggBrick struct {
	typ     int
	elem    *wm.RiwoObject
	layer   int
	row     int
	col     int
	blocked bool
	removed bool
}

var mahjonggBrickTiles = map[string][]string{
	"Dots":    {"⢀", "⣀", "⣠", "⣤", "⣴", "⣶", "⣾", "⣿", "⑨"},
	"Bamboo":  {"🍩", "🥯", "🥨", "🍕", "🥪", "🌮", "🌭", "🍔", "🍟"},
	"Chars":   {"𝟙", "𝟚", "𝟛", "𝟜", "𝟝", "𝟞", "𝟟", "𝟠", "𝟡"},
	"Winds":   {"←", "↑", "→", "↓"},
	"Dragons": {"🔴", "🟢", "🔵"},
	"Flowers": {"🌺", "🌼", "🌸", "🍀"},
	"Seasons": {"🕑", "🕓", "🕗", "🕙"},
}

func mahjonggGetTileEmoji(tileID int) string {
	categoryOrder := []string{"Dots", "Bamboo", "Chars", "Winds", "Dragons", "Flowers", "Seasons"}

	offset := 0
	for _, category := range categoryOrder {
		tiles := mahjonggBrickTiles[category]
		if tileID < offset+len(tiles) {
			return tiles[tileID-offset]
		}
		offset += len(tiles)
	}
	return "?" // how?
}
func mahjonggGetTotalTileTypes() int {
	total := 0
	for _, tiles := range mahjonggBrickTiles {
		total += len(tiles)
	}
	return total // Should be 42 total tiles
}

var mahjonggLayouts = map[string][][][]int{
	"Classic": {
		{
			{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
			{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
			{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
			{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
			{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
			{87}, // Last index stores amount of bricks
		},
		{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{36},
		},
		{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{16},
		},
		{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{4},
		},
		{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{1},
		},
	},
	"Fortress": {
		{
			{1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1},
			{1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1},
			{1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1},
			{0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0},
			{1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1},
			{1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1},
			{1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1},
			{72},
		},
		{
			{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0},
			{0, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0},
			{0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0},
			{0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 0},
			{0, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0},
			{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0},
			{44},
		},
		{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0},
			{0, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{24},
		},
		{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{4},
		},
	},
	"Arena": {
		{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
			{1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
			{1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
			{1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{80},
		},
		{
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
			{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
			{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
			{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{36},
		},
		{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
			{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
			{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{28},
		},
	},
}

// Check if a brick position is free (can be selected)
func mahjonggIsBrickFree(board [][][]*mahjonggBrick, layer, row, col int) bool {
	brick := board[layer][row][col]
	if brick == nil || brick.removed {
		return false
	}

	// Check if there's a brick on top
	if layer+1 < len(board) {
		if board[layer+1][row][col] != nil && !board[layer+1][row][col].removed {
			return false
		}
	}

	// Check sides
	leftFree := (col == 0) ||
		(board[layer][row][col-1] == nil || board[layer][row][col-1].removed)
	rightFree := (col == len(board[layer][row])-1) ||
		(board[layer][row][col+1] == nil || board[layer][row][col+1].removed)

	// At least one side must be completely free
	return leftFree || rightFree
}

// Get all available (non-removed, typ=-1) brick positions
func mahjonggGetAvailablePositions(board *[][][]*mahjonggBrick) []*mahjonggBrick {
	var available []*mahjonggBrick

	for layer := 0; layer < len(*board); layer++ {
		for row := 0; row < len((*board)[layer]); row++ {
			for col := 0; col < len((*board)[layer][row]); col++ {
				brick := (*board)[layer][row][col]
				if brick != nil && !brick.removed && brick.typ == -1 {
					available = append(available, brick)
				}
			}
		}
	}

	return available
}

// Check if position would be valid for placement (ensuring solvability)
func mahjonggIsValidPlacement(board *[][][]*mahjonggBrick, brick *mahjonggBrick) bool {
	layer, row, col := brick.layer, brick.row, brick.col

	// Check if there's a brick on top that's already assigned
	if layer+1 < len(*board) {
		topBrick := (*board)[layer+1][row][col]
		if topBrick != nil && !topBrick.removed && topBrick.typ != -1 {
			return false // Can't place if brick above is already assigned
		}
	}

	// At least one side must be free or will become free
	leftFree := (col == 0) ||
		((*board)[layer][row][col-1] == nil ||
			(*board)[layer][row][col-1].removed ||
			(*board)[layer][row][col-1].typ == -1)

	rightFree := (col == len((*board)[layer][row])-1) ||
		((*board)[layer][row][col+1] == nil ||
			(*board)[layer][row][col+1].removed ||
			(*board)[layer][row][col+1].typ == -1)

	return leftFree || rightFree
}

// Update which bricks are currently blocked
func mahjonggUpdateBlockedStatus(board *[][][]*mahjonggBrick) {
	for layer := 0; layer < len(*board); layer++ {
		for row := 0; row < len((*board)[layer]); row++ {
			for col := 0; col < len((*board)[layer][row]); col++ {
				brick := (*board)[layer][row][col]
				if brick != nil && !brick.removed {
					brick.blocked = !mahjonggIsBrickFree(*board, layer, row, col)
				}
			}
		}
	}
}

// Guess for now board size will be static 14x8
// no pun intended, this is 9front's size of `games/mahjongg`
func mahjonggBoardCreate(layout string) [][][]*mahjonggBrick {
	/* Example:
	board := { // Board
		{ 	   	// Layer1
			{28},
			{*mahjonggBrick, *mahjonggBrick...}, // blocked = true or false
			...
			{*mahjonggBrick, *mahjonggBrick...}, // removed = true or false
		},
		{       // Layer2
			...
		},
	}
	*/
	layoutData := mahjonggLayouts[layout]
	numLayers := len(layoutData)

	board := make([][][]*mahjonggBrick, numLayers)

	for layer := 0; layer < numLayers; layer++ {
		numRows := len(layoutData[layer]) - 1
		board[layer] = make([][]*mahjonggBrick, numRows)

		for row := 0; row < numRows; row++ {
			numCols := len(layoutData[layer][0])
			board[layer][row] = make([]*mahjonggBrick, numCols)

			for col := 0; col < len(mahjonggLayouts[layout][layer][row]); col++ {
				if layoutData[layer][row][col] == 1 {
					newBrick := mahjonggBrick{
						typ:     -1, // -1 means type not yet assigned
						elem:    nil,
						layer:   layer,
						row:     row,
						col:     col,
						blocked: true,  // Will be calculated later
						removed: false, // It exists initially
					}
					board[layer][row][col] = &newBrick
				}
				// If layout value is 0, the position remains nil, representing empty space.
			}
		}
	}
	return board
}

// Get tile category ranges
func mahjonggGetTileRanges() (flowerStart, flowerEnd, seasonStart, seasonEnd int) {
	offset := 0
	categoryOrder := []string{"Dots", "Bamboo", "Chars", "Winds", "Dragons", "Flowers", "Seasons"}

	for _, category := range categoryOrder {
		tiles := mahjonggBrickTiles[category]
		if category == "Flowers" {
			flowerStart = offset
			flowerEnd = offset + len(tiles)
		} else if category == "Seasons" {
			seasonStart = offset
			seasonEnd = offset + len(tiles)
		}
		offset += len(tiles)
	}
	return
}

// Check if two tiles match (flowers match any flower, seasons match any season)
func mahjonggTilesMatch(typ1, typ2 int) bool {
	flowerStart, flowerEnd, seasonStart, seasonEnd := mahjonggGetTileRanges()

	if typ1 >= flowerStart && typ1 < flowerEnd && typ2 >= flowerStart && typ2 < flowerEnd {
		return true
	}
	if typ1 >= seasonStart && typ1 < seasonEnd && typ2 >= seasonStart && typ2 < seasonEnd {
		return true
	}

	return typ1 == typ2
}

func mahjonggBoardFill(board *[][][]*mahjonggBrick) {
	available := mahjonggGetAvailablePositions(board)
	totalBricks := len(available)

	if totalBricks == 0 || totalBricks%2 != 0 {
		return
	}

	rand.Shuffle(len(available), func(i, j int) {
		available[i], available[j] = available[j], available[i]
	})

	flowerStart, _, seasonStart, _ := mahjonggGetTileRanges()

	tilePool := make([]int, 0, totalBricks)

	// Determine how many flowers and seasons to include
	numFlowerPairs := 4                                       // All 4 flowers
	numSeasonPairs := 4                                       // All 4 seasons
	numSpecialBricks := (numFlowerPairs + numSeasonPairs) * 2 // 16 total
	numRegularBricks := totalBricks - numSpecialBricks

	// Add regular tiles (Dots, Bamboo, Chars, Winds, Dragons) in sets of 4
	regularTileTypes := flowerStart // Everything before flowers
	for i := 0; i < numRegularBricks; i++ {
		tileID := (i / 4) % regularTileTypes
		tilePool = append(tilePool, tileID)
	}

	// Add all 4 flower types as pairs
	for i := 0; i < numFlowerPairs; i++ {
		flowerID := flowerStart + i
		tilePool = append(tilePool, flowerID, flowerID)
	}
	// Add all 4 season types as pairs
	for i := 0; i < numSeasonPairs; i++ {
		seasonID := seasonStart + i
		tilePool = append(tilePool, seasonID, seasonID)
	}

	// Shuffle and assign
	rand.Shuffle(len(tilePool), func(i, j int) {
		tilePool[i], tilePool[j] = tilePool[j], tilePool[i]
	})

	for i := 0; i < len(available) && i < len(tilePool); i++ {
		available[i].typ = tilePool[i]
	}

	mahjonggUpdateBlockedStatus(board)
}

func mahjonggConstruct(window *wm.RiwoWindow) {
	window.Title = "mahjongg"
	themeIdx := 0
	layoutIdx := 0
	theme := wm.GetTheme(mahjonggThemeNames[themeIdx])

	var board [][][]*mahjonggBrick
	var selected *mahjonggBrick
	var startTime time.Time
	var timerStop bool

	container := wm.Create().
		Style("height", "100%").
		Style("display", "flex").
		Style("flexDirection", "column").
		Style("justifyContent", "center").
		Style("alignItems", "center").
		Style("backgroundColor", theme["faded"]).
		Style("object-fit", "fill")

	timerElem := wm.Create().
		Text("Time: 0:00").
		Style("marginBottom", "0.5rem").
		Style("color", theme["vivid"]).
		Style("fontWeight", "bold")

	boardElem := wm.Create().
		Style("position", "relative").
		Style("display", "inline-block").
		Style("width", "36rem").
		Style("height", "26rem")

	container.Append(timerElem, boardElem)
	window.Content.Inner("").Append(container)

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

	// Check win
	checkWin := func() {
		for layer := 0; layer < len(board); layer++ {
			for row := 0; row < len(board[layer]); row++ {
				for col := 0; col < len(board[layer][row]); col++ {
					brick := board[layer][row][col]
					if brick != nil && !brick.removed {
						return
					}
				}
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
			Text("You Won!\nTime: "+strconv.Itoa(min)+":"+secStr).
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
			Mount(boardElem)
	}

	// Check dead end
	checkDeadEnd := func() {
		hasRemaining := false
		for layer := 0; layer < len(board); layer++ {
			for row := 0; row < len(board[layer]); row++ {
				for col := 0; col < len(board[layer][row]); col++ {
					brick := board[layer][row][col]
					if brick != nil && !brick.removed {
						hasRemaining = true
						break
					}
				}
			}
		}
		if !hasRemaining {
			return
		}

		// Check for valid moves
		var free []*mahjonggBrick
		for layer := 0; layer < len(board); layer++ {
			for row := 0; row < len(board[layer]); row++ {
				for col := 0; col < len(board[layer][row]); col++ {
					brick := board[layer][row][col]
					if brick != nil && !brick.removed && !brick.blocked {
						free = append(free, brick)
					}
				}
			}
		}

		hasMove := false
		for i := 0; i < len(free); i++ {
			for j := i + 1; j < len(free); j++ {
				if mahjonggTilesMatch(free[i].typ, free[j].typ) {
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
				Mount(boardElem)
		}
	}

	// Update visual state based on blocked status
	updateBrickVisuals := func() {
		for layer := 0; layer < len(board); layer++ {
			for row := 0; row < len(board[layer]); row++ {
				for col := 0; col < len(board[layer][row]); col++ {
					brick := board[layer][row][col]
					if brick == nil || brick.removed {
						continue
					}

					if brick.blocked {
						brick.elem.Style("filter", "brightness(0.6)").Style("cursor", "not-allowed")
					} else {
						brick.elem.Style("filter", "brightness(1)").Style("cursor", wm.CursorInvertUrl)
					}
				}
			}
		}
	}

	// Click handler
	onBrickClick := func(b *mahjonggBrick) {
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
			if mahjonggTilesMatch(selected.typ, b.typ) {
				selected.removed = true
				selected.elem.Style("display", "none")
				b.removed = true
				b.elem.Style("display", "none")
				selected = nil
				mahjonggUpdateBlockedStatus(&board)
				updateBrickVisuals()
				checkWin()
				checkDeadEnd()
			} else {
				selected.elem.Style("backgroundColor", theme["normal"])
				selected = b
				b.elem.Style("backgroundColor", theme["vivid"])
			}
		}
	}

	// Create visual element for brick
	createBrickElement := func(brick *mahjonggBrick) {
		symbol := mahjonggGetTileEmoji(brick.typ)

		left := float64(brick.col)*2.5 + float64(brick.layer)*0.25
		top := float64(brick.row)*3.125 - float64(brick.layer)*0.25
		zindex := brick.layer*100 + brick.row

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
			Style("border", "0.125em solid "+theme["vivid"]).
			Style("cursor", wm.CursorInvertUrl).
			Style("userSelect", "none").
			Style("boxShadow", "0.125em 0.125em 0.25em rgba(0,0,0,0.3)")

		brick.elem = elem

		elem.Listen("mousedown", func(this js.Value, args []js.Value) interface{} {
			onBrickClick(brick)
			return nil
		})

		boardElem.Append(elem)
	}

	// Start game
	startGame := func() {
		selected = nil
		boardElem.Inner("")
		timerStop = false
		startTime = time.Now()

		// Create and fill board
		board = mahjonggBoardCreate(mahjonggLayoutNames[layoutIdx])
		mahjonggBoardFill(&board)

		// Create visual elements for all bricks
		for layer := 0; layer < len(board); layer++ {
			for row := 0; row < len(board[layer]); row++ {
				for col := 0; col < len(board[layer][row]); col++ {
					brick := board[layer][row][col]
					if brick != nil {
						createBrickElement(brick)
					}
				}
			}
		}

		mahjonggUpdateBlockedStatus(&board)
		updateBrickVisuals()
		updateTimer()
	}

	// Cycle theme
	cycleTheme := func() {
		themeIdx = (themeIdx + 1) % len(mahjonggThemeNames)
		theme = wm.GetTheme(mahjonggThemeNames[themeIdx])

		container.Style("backgroundColor", theme["faded"])
		timerElem.Style("color", theme["vivid"])

		for layer := 0; layer < len(board); layer++ {
			for row := 0; row < len(board[layer]); row++ {
				for col := 0; col < len(board[layer][row]); col++ {
					brick := board[layer][row][col]
					if brick != nil && !brick.removed {
						if brick == selected {
							brick.elem.Style("backgroundColor", theme["vivid"])
						} else {
							brick.elem.Style("backgroundColor", theme["normal"])
						}
						brick.elem.Style("border", "0.125em solid "+theme["vivid"])
					}
				}
			}
		}
	}

	// Cycle layout
	cycleLayout := func() {
		layoutIdx = (layoutIdx + 1) % len(mahjonggLayoutNames)
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
