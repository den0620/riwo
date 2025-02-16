/*
App globals, primarily AppRegistry
*/

package apps

import (
	"riwo/wm"
)

// AppRegistry holds all available app functions.
// Each app should register itself (typically in its init function).
var AppRegistry = make(map[string]func(*wm.Window))
