package playwright

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string {
	return &v
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool {
	return &v
}

// Int is a helper routine that allocates a new int32 value
// to store v and returns a pointer to it.
func Int(v int) *int {
	return &v
}

// Float is a helper routine that allocates a new float64 value
// to store v and returns a pointer to it.
func Float(v float64) *float64 {
	return &v
}

// Null will be used in certain scenarios where a strict nil pointer
// check is not possible
func Null() interface{} {
	return "PW_NULL"
}

type Rect struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	X      int `json:"x"`
	Y      int `json:"y"`
}
