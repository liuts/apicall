package call

import "testing"

func TestMake_call(t *testing.T) {
	Init("COM3", 9600)
	Make_call("1337.....")
}
