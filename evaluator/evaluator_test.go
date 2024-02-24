package evaluator

import (
	"testing"

	"github.com/nga1hte/paite-khawl-laimal/lexer"
	"github.com/nga1hte/paite-khawl-laimal/object"
	"github.com/nga1hte/paite-khawl-laimal/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"69", 69},
		{"-5", -5},
		{"-69", -69},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object is not Integer. got=%T, want=%d", result.Value, expected)
		return false
	}

	return true

}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"tak", true},
		{"zuau", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"tak == tak", true},
		{"zuau == zuau", true},
		{"tak == zuau", false},
		{"tak != zuau", true},
		{"(1 < 2) == zuau", false},
		{"(1 < 2) == tak", true},
		{"(1 > 2) == tak", false},
		{"(1 > 2) == zuau", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!tak", false},
		{"!zuau", true},
		{"!5", false},
		{"!!tak", true},
		{"!!zuau", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"ahihleh (tak) { 10 }", 10},
		{"ahihleh (zuau) { 10 }", nil},
		{"ahihleh (1) { 10 }", 10},
		{"ahihleh (1 < 2) { 10 }", 10},
		{"ahihleh (1 > 2) { 10 }", nil},
		{"ahihleh (1 > 2) { 10 } ahihkeileh { 20 }", 20},
		{"ahihleh (1 < 2) { 10 } ahihkeileh { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"lehkik 10;", 10},
		{"lehkik 10; 9", 10},
		{"lehkik 2 * 5; 9", 10},
		{"9; lehkik 10; 9;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + tak;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + tak;5", "type mismatch: INTEGER + BOOLEAN"},
		{"-tak;", "unknown operator: -BOOLEAN"},
		{"tak + tak;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; tak + zuau; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"ahihleh (10 > 1) { tak + zuau }", "unknown operator: BOOLEAN + BOOLEAN"},
		{`ahihleh (10 > 1) {
			ahihleh (10 > 1) {
				lehkik tak + zuau;
			}
			lehkik 1;
		}`, "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{`"Hello" - "World"`, "unknown operator: STRING - STRING"},
		{`{"name": "Monkey"}[thilhihna(x){ x }];`, "unusable as hash key: FUNCTION"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"huchin a = 5; a;", 5},
		{"huchin a = 5 * 5; a;", 25},
		{"huchin a = 5; huchin b = a; b;", 5},
		{"huchin a = 5; huchin b = a; huchin c = a + b + 5; c", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "thilhihna (x) {x + 2;};"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"huchin identity = thilhihna(x){ x; }; identity(5);", 5},
		{"huchin identity = thilhihna(x){ lehkik x; }; identity(5)", 5},
		{"huchin double = thilhihna(x){ x * 2; }; double(2);", 4},
		{"huchin add = thilhihna(x, y){ x + y; }; add(2,2);", 4},
		{"huchin add = thilhihna(x, y){ x + y; }; add(2 + 2, add(2,2))", 8},
		{"thilhihna(x){x;}(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestEnclosingEnvironments(t *testing.T) {
	input := `
	huchin joyson = 10;
	huchin mungboi = 20;
	huchin lunboi = 30;

	huchin total = thilhihna(first) {
		huchin second = 20;
		first + second + lunboi;
	};

	total(10);
	`

	testIntegerObject(t, testEval(input), 60)

}

func TestClosures(t *testing.T) {
	input := `
	huchin newAdder = thilhihna(x){
		thilhihna(y) {
			x + y;
		};
	};
	huchin addTwo = newAdder(5);
	addTwo(2);
	`
	testIntegerObject(t, testEval(input), 7)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`saudan("")`, 0},
		{`saudan("four")`, 4},
		{`saudan("hello world")`, 11},
		{`saudan(1)`, "argument to `saudan` not supported, got INTEGER"},
		{`saudan("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`saudan([1, 2, 3])`, 3},
		{`saudan([])`, 0},
		{`amasa([1, 2, 3])`, 1},
		{`amasa([])`, nil},
		{`amasa(1)`, "argument to `amasa` must be ARRAY, got INTEGER"},
		{`nanung([1, 2, 3])`, 3},
		{`nanung([])`, nil},
		{`nanung(1)`, "argument to `nanung` must be ARRAY, got INTEGER"},
		{`amasalouteng([1, 2, 3])`, []int{2, 3}},
		{`amasalouteng([])`, nil},
		{`sawnlut([], 1)`, []int{1}},
		{`sawnlut(1, 1)`, "argument to `sawnlut` must be ARRAY, got INTEGER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			testNullObject(t, evaluated)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Elements))
				continue
			}
			for i, expectedElem := range expected {
				testIntegerObject(t, array.Elements[i], int64(expectedElem))
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"huchin i = 0; [1][i]", 1},
		{"[1, 2, 3][1 + 1]", 3},
		{"huchin myArray = [1, 2, 3]; myArray[2]", 3},
		{"huchin myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"huchin myArray = [1, 2, 3]; huchin i = myArray[0]; myArray[i];", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `huchin two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		tak: 5,
		zuau: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}

}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo":5}["foo"]`, 5},
		{`{"foo":5}["bar"]`, nil},
		{`huchin key = "foo"; {"foo":5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5: 5}[5]`, 5},
		{`{tak: 5}[tak]`, 5},
		{`{zuau: 5}[zuau]`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}

}
