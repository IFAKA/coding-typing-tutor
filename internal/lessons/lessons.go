// Package lessons provides lesson definitions and progress tracking for
// structured touch-typing progression.
package lessons

// Lesson defines a single lesson level and which keys are allowed.
// A nil AllowedKeys means all keys are allowed (free-code mode).
type Lesson struct {
	Number      int
	Name        string
	AllowedKeys []rune // base keys (lowercase); nil = unrestricted
}

// All is the ordered list of lessons from easiest to hardest.
var All = []Lesson{
	{
		1, "home row",
		[]rune("asdfghjkl;"),
	},
	{
		2, "+ top row",
		[]rune("asdfghjkl;qwertyuiop"),
	},
	{
		3, "+ bottom row",
		[]rune("asdfghjkl;qwertyuiopzxcvbnm"),
	},
	{
		4, "+ numbers",
		[]rune("asdfghjkl;qwertyuiopzxcvbnm1234567890"),
	},
	{
		5, "+ common symbols",
		[]rune("asdfghjkl;qwertyuiopzxcvbnm1234567890-=[]"),
	},
	{
		6, "+ operators",
		[]rune("asdfghjkl;qwertyuiopzxcvbnm1234567890-=[].,/"),
	},
	{
		7, "+ all symbols",
		nil, // allow everything
	},
	{
		8, "free code",
		nil,
	},
}
