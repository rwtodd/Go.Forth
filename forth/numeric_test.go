package forth

import (
	"testing"
)

func TestAdd(t *testing.T) {
	tstRunForth(t, `2 3 +  2 3.1 +`, 5, 5.1)
	tstRunForth(t, `" hi" "  there" +`, "hi there")
}

func TestMul(t *testing.T) {
	tstRunForth(t, "2 3 *   2 .25 *", 6, 0.5)
	tstRunForth(t, `" hi" 3 *  3 " yo" *`,
		"hihihi", "yoyoyo")

}
