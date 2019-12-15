package bsongen

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Language interface {
	Length() string
}

func NewBSONGenerator(lang Language) *BSONGenerator {
	return &BSONGenerator{lang}
}

var typeNames = []string{
	"decimal128",
	"double",
	"int32",
	"int64",
	"string",
}

type BSONGenerator struct {
	lang Language
}

func (g *BSONGenerator) Generate(input string) (string, error) {
	tmpl, err := template.New("mql").Funcs(g.makeFuncMap()).Parse(input)
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
func (g *BSONGenerator) makeFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"definitions":            g.definitions,
		"end_document":           g.endDocument,
		"field":                  func() string { return "field" },
		"length":                 g.lang.Length,
		"named_field_any":        g.namedFieldAny,
		"named_field":            g.namedField,
		"start_document":         g.startDocument,
		"start_document_no_type": g.startDocumentNoType,
	}

	for _, typeName := range typeNames {
		typeName := typeName
		funcMap[fmt.Sprintf("named_field_%s", typeName)] = func(name string) string {
			return g.namedField(name, typeName)
		}
	}

	return funcMap
}

func (g *BSONGenerator) anyStartDocument() string {
	return "TYPE_DOCUMENT name=cstring " + g.lang.Length()
}

func (BSONGenerator) endDocument() string {
	return "NUL_BYTE"
}

func (g *BSONGenerator) namedField(name string, typeNames ...string) string {
	switch len(typeNames) {
	case 0:
		return g.namedFieldAny(name)
	case 1:
		return g.namedFieldWithType(name, typeNames[0])
	default:
		parts := make([]string, 0, len(typeNames))
		for i := 0; i < len(typeNames); i++ {
			parts = append(parts, g.namedFieldWithType(name, typeNames[i]))
		}
		return strings.Join(parts, " | ")
	}
}

func (g *BSONGenerator) namedFieldWithType(name string, typeName string) string {
	return fmt.Sprintf("TYPE_%s %s %s", strings.ToUpper(typeName), g.convertNameToMarkers(name), typeName)
}

func (g *BSONGenerator) namedFieldAny(name string) string {
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

func (g *BSONGenerator) startDocumentNoType() string {
	return g.lang.Length()
}

func (g *BSONGenerator) startDocument(names ...string) string {
	switch len(names) {
	case 0:
		return "TYPE_DOCUMENT name=cstring " + g.lang.Length()
	default:
		return fmt.Sprintf("TYPE_DOCUMENT %s %s", g.convertNameToMarkers(names[0]), g.lang.Length())
	}
}

// HELPERS
func (BSONGenerator) convertNameToMarkers(name string) string {
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

func (BSONGenerator) definitions() string {
	return `
// fields
field:
	field_decimal128
|   field_double
|   field_int32
|   field_int64
|	field_string
;

field_decimal128: TYPE_DECIMAL128 name=cstring value=decimal128;
field_double: TYPE_DOUBLE name=cstring value=double;
field_int32: TYPE_INT32 name=cstring value=int32;
field_int64: TYPE_INT64 name=cstring value=int64;
field_string: TYPE_STRING name=cstring value=string;

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
string:
	int32 cstring
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
