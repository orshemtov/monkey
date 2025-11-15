package code

import "testing"

func TestMake(t *testing.T) {
	testCases := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
	}
	for _, tC := range testCases {
		t.Run(string(tC.op), func(t *testing.T) {
			instruction := Make(tC.op, tC.operands...)

			if len(instruction) != len(tC.expected) {
				t.Errorf("instruction has wrong length. want=%d, got %d", len(tC.expected), len(instruction))
			}

			for i, b := range tC.expected {
				if instruction[i] != tC.expected[i] {
					t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, b, instruction[i])
				}
			}
		})
	}
}
