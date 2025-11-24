package code

import "testing"

func TestMake(t *testing.T) {
	testCases := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpGetLocal, []int{255}, []byte{byte(OpGetLocal), 255}},
		{OpClosure, []int{65534, 255}, []byte{byte(OpClosure), 255, 254, 255}},
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

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpGetLocal, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpClosure, 65535, 255),
	}

	expected := `0000 OpAdd
0001 OpGetLocal 1
0003 OpConstant 2
0006 OpConstant 65535
0009 OpClosure 65535 255
`

	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions string wrong. expected=\n%q\ngot=\n%q", expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	testCases := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
		{OpClosure, []int{65535, 255}, 3},
	}
	for _, tC := range testCases {
		t.Run(string(tC.op), func(t *testing.T) {
			instruction := Make(tC.op, tC.operands...)

			def, err := Lookup(byte(tC.op))
			if err != nil {
				t.Fatalf("definition not found: %q\n", err)
			}

			operandsRead, n := ReadOperands(def, instruction[1:])
			if n != tC.bytesRead {
				t.Fatalf("n wrong. want=%d, got=%d", tC.bytesRead, n)
			}

			for i, want := range tC.operands {
				if operandsRead[i] != want {
					t.Errorf("operand wrong. want=%d, got=%d", want, operandsRead[i])
				}
			}
		})
	}
}
