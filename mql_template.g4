grammar mql;

start:
    int32 (TYPE_DOCUMENT field_name int32 stage NUL_BYTE)+ NUL_BYTE
;

// stages
stage:
    stage_limit
|   stage_match
|   stage_skip
;

stage_limit:
    {{named_field_decimal128 "$limit"}}
|   {{named_field_double "$limit"}}
|   {{named_field_int32 "$limit"}}
|   {{named_field_int64 "$limit"}}
;

stage_match:
    TYPE_DOCUMENT DOLLAR M A T C H NUL_BYTE int32 match_expr* NUL_BYTE
;

stage_skip:
    {{named_field_decimal128 "$skip"}}
|   {{named_field_double "$skip"}}
|   {{named_field_int32 "$skip"}}
|   {{named_field_int64 "$skip"}}
;

// match expressions
match_expr:
    match_eq_no_op
|   match_multi_op
;

match_eq_no_op: field_element;
match_multi_op:
    TYPE_DOCUMENT field_name int32 (
        {{named_field_any "$eq"}}
    |   {{named_field_any "$gt"}}
    |   {{named_field_any "$gte"}}
    |   {{named_field_any "$lt"}}
    |   {{named_field_any "$lte"}}
    |   {{named_field_any "$ne"}}
    |   {{named_field_any "$not"}}
    )*
    NUL_BYTE
;

// fields
field_element:
    field_element_decimal128
|   field_element_double
|   field_element_int32
|   field_element_int64
;

field_element_decimal128: TYPE_DECIMAL128 field_name decimal128;
field_element_double: TYPE_DOUBLE field_name double;
field_element_int32: TYPE_INT32 field_name int32;
field_element_int64: TYPE_INT64 field_name int64;

field_name: non_null_byte* NUL_BYTE;

// values
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

UNSPECIFIED_NON_NUL_BYTE: NON_NUL_RANGE1 | NON_NUL_RANGE2 | NON_NUL_RANGE3;