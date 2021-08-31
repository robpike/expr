package expr

import "testing"

func TestParse(t *testing.T) {
	// The String method adds parens everywhere, so it is an easy
	// check of the parse tree.
	var tests = []struct {
		in, out string
	}{
		// Singletons.
		{"3", "3"},
		{"x", "x"},
		{"x_9", "x_9"},

		// Unary operators.
		{"-3", "(-3)"},
		{"+3", "(+3)"},
		{"!1", "(!1)"},
		{"^1", "(^1)"},

		// Binary arithmetic operators.
		{"x * y", "(x * y)"},
		{"x / y", "(x / y)"},
		{"x % y", "(x % y)"},
		{"x << y", "(x << y)"},
		{"x >> y", "(x >> y)"},
		{"x & y", "(x & y)"},
		{"x &^ y", "(x &^ y)"},
		{"x + y", "(x + y)"},
		{"x - y", "(x - y)"},
		{"x | y", "(x | y)"},
		{"x ^ y", "(x ^ y)"},

		// Binary comparison operators.
		{"x == y", "(x == y)"},
		{"x != y", "(x != y)"},
		{"x < y", "(x < y)"},
		{"x <= y", "(x <= y)"},
		{"x > y", "(x > y)"},
		{"x >= y", "(x >= y)"},

		// Binary logical operators.
		{"x && y", "(x && y)"},
		{"x || y", "(x || y)"},

		// Precedence checks. Left-associative with precedence.

		// && and ||
		{"x || y || z", "((x || y) || z)"},
		{"x && y && z", "((x && y) && z)"},
		{"x && y || z", "((x && y) || z)"},
		{"x || y && z", "(x || (y && z))"},

		// Comparison operators.
		{"x == y == z", "((x == y) == z)"},
		{"x != y != z", "((x != y) != z)"},
		{"x > y > z", "((x > y) > z)"},
		{"x >= y >= z", "((x >= y) >= z)"},
		{"x < y < z", "((x < y) < z)"},
		{"x <= y <= z", "((x <= y) <= z)"},

		// Unaries and binaries at two precedence levels.
		{"+x + y", "((+x) + y)"},
		{"x + +y", "(x + (+y))"},
		{"-x - y", "((-x) - y)"},
		{"x - -y", "(x - (-y))"},
		{"^x + y", "((^x) + y)"},
		{"x + ^y", "(x + (^y))"},
		{"!x + y", "((!x) + y)"},
		{"x + !y", "(x + (!y))"},
		{"+x * y", "((+x) * y)"},
		{"x * +y", "(x * (+y))"},
		{"-x - y", "((-x) - y)"},
		{"x - -y", "(x - (-y))"},
		{"^x * y", "((^x) * y)"},
		{"x * ^y", "(x * (^y))"},
		{"!x * y", "((!x) * y)"},
		{"x * !y", "(x * (!y))"},

		// Grouping of operators at same precedence.
		// Multiplies.
		{"x * y * z", "((x * y) * z)"},
		{"x * y / z", "((x * y) / z)"},
		{"x / y * z", "((x / y) * z)"},
		{"x % y / z", "((x % y) / z)"},
		{"x / y % z", "((x / y) % z)"},
		{"x >> y / z", "((x >> y) / z)"},
		{"x / y >> z", "((x / y) >> z)"},
		{"x << y / z", "((x << y) / z)"},
		{"x / y << z", "((x / y) << z)"},
		{"x & y / z", "((x & y) / z)"},
		{"x / y & z", "((x / y) & z)"},
		{"x &^ y / z", "((x &^ y) / z)"},
		{"x / y &^ z", "((x / y) &^ z)"},

		// Adds.
		{"x + y + z", "((x + y) + z)"},
		{"x + y - z", "((x + y) - z)"},
		{"x - y + z", "((x - y) + z)"},
		{"x + y | z", "((x + y) | z)"},
		{"x | y + z", "((x | y) + z)"},
		{"x + y + z", "((x + y) + z)"},
		{"x + y ^ z", "((x + y) ^ z)"},

		// Multiplies and adds.
		{"x * y + z", "((x * y) + z)"},
		{"x + y * z", "(x + (y * z))"}, // Gotcha!

		// Grouping of comparisons and other operators.
		{"(x < y && z < 3)", "((x < y) && (z < 3))"},
		{"(x < y && z || 1)", "(((x < y) && z) || 1)"},
		{"(u == v && x == y || w == z)", "(((u == v) && (x == y)) || (w == z))"},
		{"(u == v*3 && x == y-2 || w == !z)", "(((u == (v * 3)) && (x == (y - 2))) || (w == (!z)))"},
	}
	for _, test := range tests {
		e, err := Parse(test.in)
		if err != nil {
			t.Errorf("Parsing %s: %v", test.in, err)
			continue
		}
		got := e.String()
		if got != test.out {
			t.Errorf("String for %q: %q, want %q", test.in, got, test.out)
			continue
		}
		// Round-trip. This also tests parentheses.
		e, err = Parse(got)
		if err != nil {
			t.Errorf("Parsing round-trip %s: %v", test.in, err)
			continue
		}
		got2 := e.String()
		if got2 != test.out {
			t.Errorf("String for roundtrip %q: %q, want %q", test.in, got, test.out)
			continue
		}
	}
}

func TestParseError(t *testing.T) {
	var tests = []struct {
		expr string
		err  string
	}{
		{"x x", `syntax error at "x"`},
		{"(x + ", `unexpected eof`},
		{"(x + 1", `unclosed paren at eof`},
		{"(x + 1))", `syntax error at ")"`},
		{"x + >4", `bad expression at ">4"`},
		{"x @ 4", `syntax error at "@ 4"`},
	}
	for _, test := range tests {
		_, err := Parse(test.expr)
		if err == nil {
			t.Errorf("Parsing %s: no error", test.expr)
			continue
		}
		got := err.Error()
		if got != test.err {
			t.Errorf("Wrong error for %s: got %q, want %q", test.expr, got, test.err)
			continue
		}
	}
}

func TestEval(t *testing.T) {
	var tests = []struct {
		x, y   int
		expr   string
		result int
	}{
		// Singletons.
		{0, 0, "3", 3},
		{0, 0, "3+4", 7},
		{10, 0, "x", 10},
		{10, 20, "y", 20},

		// Unary operators.
		{0, 0, "+3", 3},
		{0, 0, "-3", -3},
		{0, 0, "!0", 1}, // No booleans here, just 0 and 1.
		{0, 0, "!3", 0},
		{0, 0, "^0", -1},
		{0, 0, "^1", -2},
		{7, 0, "+x", 7},
		{7, 0, "-x", -7},

		// Binary operators.

		// Arithmetic.
		// To be sure it's working, aim for a different answer for each when feasible.
		{7, 3, "x * y", 21},
		{7, 3, "x / y", 2},
		{7, 0, "x / y", 0}, // Special case.
		{7, 3, "x % y", 1},
		{7, 0, "x % y", 0}, // Special case.
		{7, 3, "x << y", 56},
		{7, 1, "x >> y", 3},
		{7, 4, "x & y", 4},
		{17, 5, "x &^ y", 16},
		{7, 3, "x + y", 10},
		{7, 2, "x - y", 5},

		// Bits.
		{4, 2, "x | y", 6},
		{7, 7, "x ^ y", 0},
		{14, 7, "x ^ y", 9},
		{14, 7, "x &^ y", 8},
		{144, 3, "x >> y", 18},
		{14, 1, "x << y", 28},

		// Comparison.
		{14, 7, "x == y", 0},
		{14, 7, "x != y", 1},
		{14, 7, "x > y", 1},
		{7, 14, "x > y", 0},
		{14, 7, "x < y", 0},
		{7, 14, "x < y", 1},
		{14, 7, "x >= y", 1},
		{7, 14, "x >= y", 0},
		{14, 14, "x >= y", 1},
		{14, 7, "x <= y", 0},
		{14, 14, "x <= y", 1},
		{7, 14, "x <= y", 1},

		// Mix it up.
		{4, 4, "x > y || y == x", 1},
		{4, 4, "x > y && y == x", 0},
		{5, 4, "x > y && y == 4", 1},
		{5, 4, "x*x + y*y", 41},

		// No errors.
		{4, 0, "x / 0", 0},
		{4, 0, "x % 0", 0},
		{4, 0, "x << -1", 0},
		{4, 0, "x >> -1", 0},
		{4, 0, "v >> -1", 0},

		// Check that ReturnZero works mid-expression.
		{0, 0, "zz", 0},
		{0, 0, "4 + zz", 4},
		{1, 0, "3 + x / y", 3},
		{1, 0, "3 + x % y", 3},
		{1, -3, "3*x + x << y", 3},
		{1, -3, "3*y + x >> y", -9},
	}
	vars := make(map[string]int)
	for _, test := range tests {
		e, err := Parse(test.expr)
		if err != nil {
			t.Errorf("Parsing %s: %v", test.expr, err)
			continue
		}
		vars["x"] = test.x
		vars["y"] = test.y
		got, err := e.Eval(vars, ReturnZero)
		if err != nil {
			t.Errorf("Evaluating %s: %v", test.expr, err)
			continue
		}
		if got != test.result {
			t.Errorf("Evaluating %s: got %d, expected %d", test.expr, got, test.result)
		}
	}
}

func TestEvalError(t *testing.T) {
	var tests = []struct {
		expr string
		err  string
	}{
		{"y", `undefined variable y`},
		{"x / 0", `division by zero`},
		{"x % 0", `modulo by zero`},
		{"x << -1", `negative left shift amount`},
		{"x >> -1", `negative right shift amount`},
	}
	vars := map[string]int{"x": 1}
	for _, test := range tests {
		e, err := Parse(test.expr)
		if err != nil {
			t.Errorf("Parsing %s: %v", test.expr, err)
			continue
		}
		_, err = e.Eval(vars, ReturnError)
		if err == nil {
			t.Errorf("Evaluating %s: no error", test.expr)
			continue
		}
		got := err.Error()
		if got != test.err {
			t.Errorf("Wrong error for %s: got %q, want %q", test.expr, got, test.err)
			continue
		}
	}
}
