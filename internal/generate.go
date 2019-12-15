package internal

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

// Generate generates the template.
func Generate(input string, lang Language) (string, error) {
	g := &Generator{lang}
	return g.Generate(input)
}

var typeNames = []string{
	"decimal128",
	"double",
	"int32",
	"int64",
}

type Generator struct {
	lang Language
}

func (g *Generator) Generate(input string) (string, error) {
	tmpl, err := template.New("mql.g4.tmpl").Funcs(g.makeFuncMap()).Parse(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var output bytes.Buffer

	if err = tmpl.Execute(&output, nil); err != nil {
		return "", err
	}

	return output.String(), nil
}

// FUNCTIONS
func (g *Generator) makeFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"any_field_any":        g.anyFieldAny,
		"any_start_document":   g.anyStartDocument,
		"definitions":          g.definitions,
		"end_document":         g.endDocument,
		"length":               g.lang.Length,
		"named_field_any":      g.namedFieldAny,
		"named_field_document": g.namedFieldDocument,
		"named_start_document": g.namedStartDocument,
	}

	for _, typeName := range typeNames {
		typeName := typeName
		funcMap[fmt.Sprintf("named_field_%s", typeName)] = func(name string) string {
			return g.namedField(name, typeName)
		}
	}

	return funcMap
}

func (Generator) anyFieldAny() string {
	return "any_field_any"
}

func (g *Generator) anyStartDocument() string {
	return "TYPE_DOCUMENT name=cstring " + g.lang.Length()
}

func (Generator) endDocument() string {
	return "NUL_BYTE"
}

func (g *Generator) namedField(name string, typeName string) string {
	return fmt.Sprintf("TYPE_%s %s %s", strings.ToUpper(typeName), g.convertNameToMarkers(name), typeName)
}

func (g *Generator) namedFieldAny(name string) string {
	result := "("
	for i, typeName := range typeNames {
		if i != 0 {
			result += " | "
		}
		result += g.namedField(name, typeName)
	}
	result += ")"

	return result
}

func (g *Generator) namedFieldDocument(name string) string {
	f := g.namedField(name, "document")
	return f + " " + g.lang.Length()
}

func (g *Generator) namedStartDocument(name string) string {
	return fmt.Sprintf("TYPE_DOCUMENT %s %s", g.convertNameToMarkers(name), g.lang.Length())
}

// HELPERS
func (Generator) convertNameToMarkers(name string) string {
	parts := make([]string, len(name)+1)
	parts[len(name)] = "NUL_BYTE"

	for i := 0; i < len(name); i++ {
		switch name[i] {
		case '$':
			parts[i] = "DOLLAR"
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			parts[i] = strings.ToUpper(string(name[i]))
		default:
			panic(fmt.Sprintf("unknown name character: %s", string(name[i])))
		}
	}

	return strings.Join(parts, " ")
}

func (Generator) definitions() string {
	return `
// fields
any_field_any:
	any_field_decimal128
|   any_field_double
|   any_field_int32
|   any_field_int64
;

any_field_decimal128: TYPE_DECIMAL128 name=cstring value=decimal128;
any_field_double: TYPE_DOUBLE name=cstring value=double;
any_field_int32: TYPE_INT32 name=cstring value=int32;
any_field_int64: TYPE_INT64 name=cstring value=int64;

// values
cstring: 
	non_null_byte* NUL_BYTE
;
decimal128: 
    (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) 
    (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE)
    (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE)
    (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE)
;
double: 
    (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) 
    (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE)
;
int32: 
    (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE)
;
int64: 
    (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) 
    (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE) (non_null_byte | NUL_BYTE)
;

// general
non_null_byte: 
    TYPE_DOUBLE | TYPE_STRING | TYPE_DOCUMENT | TYPE_ARRAY | TYPE_BINARY | TYPE_UNDEFINED | TYPE_FALSE | TYPE_TRUE
|   TYPE_UTCDATETIME | TYPE_NULL | TYPE_REGEX | TYPE_DBPOINTER | TYPE_CODE | TYPE_SYMBOL | TYPE_CODE_WITH_SCOPE
|   TYPE_INT32 | TYPE_TIMESTAMP | TYPE_INT64 | TYPE_DECIMAL128
|   DOLLAR 
|   A | B | C | D | E | F | G | H | I | J | K | L | M | N | O | P | Q | R | S | T | U | V | W | X | Y | Z
|   UNSPECIFIED_NON_NUL_BYTE;

// LEXER
NUL_BYTE: '\u0000';
TYPE_DOUBLE: '\u0001';
TYPE_STRING: '\u0002';
TYPE_DOCUMENT: '\u0003';
TYPE_ARRAY: '\u0004';
TYPE_BINARY: '\u0005';
TYPE_UNDEFINED: '\u0006';
TYPE_FALSE: '\u0007';
TYPE_TRUE: '\u0008';
TYPE_UTCDATETIME: '\u0009';
TYPE_NULL: '\u000A';
TYPE_REGEX: '\u000B';
TYPE_DBPOINTER: '\u000C';
TYPE_CODE: '\u000D';
TYPE_SYMBOL: '\u000E';
TYPE_CODE_WITH_SCOPE: '\u000F';
TYPE_INT32: '\u0010';
TYPE_TIMESTAMP: '\u0011';
TYPE_INT64: '\u0012';
TYPE_DECIMAL128: '\u0013';
fragment NON_NUL_RANGE1: '\u0014'..'\u0023';
DOLLAR: '$'; // \u0024
fragment NON_NUL_RANGE2: '\u0025'..'\u0060';
A: 'a';
B: 'b';
C: 'c';
D: 'd';
E: 'e';
F: 'f';
G: 'g';
H: 'h';
I: 'i';
J: 'j';
K: 'k';
L: 'l';
M: 'm';
N: 'n';
O: 'o';
P: 'p';
Q: 'q';
R: 'r';
S: 's';
T: 't';
U: 'u';
V: 'v';
W: 'w';
X: 'x';
Y: 'y';
Z: 'z';
fragment NON_NUL_RANGE3: '\u007B'..'\u007E';
TYPE_MAXKEY: '\u007F';
TYPE_MINKEY: '\u00FF';

UNSPECIFIED_NON_NUL_BYTE: NON_NUL_RANGE1 | NON_NUL_RANGE2 | NON_NUL_RANGE3;`
}
