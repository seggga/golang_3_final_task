package querier

type Token int

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	COLUMN
	//	COMMENT

	literal_beg
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	INT    // 12345
	STRING // "abc"
	literal_end

	operator_beg
	// Operators and delimiters
	EQL // ==
	LSS // <
	GTR // >

	NEQ // !=
	LEQ // <=
	GEQ // >=

	COMMA  // ,
	PERIOD // .

	operator_end

	keyword_beg
	// Keywords
	SELECT
	FROM
	WHERE

	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	COLUMN: "COLUMN",
	INT:    "INT",
	STRING: "STRING",

	EQL: "==",
	LSS: "<",
	GTR: ">",

	NEQ: "!=",
	LEQ: "<=",
	GEQ: ">=",

	COMMA:  ",",
	PERIOD: ".",

	SELECT: "select",
	FROM:   "from",
	WHERE:  "where",
}

// A set of constants for precedence-based expression parsing.
// Non-operators have lowest precedence, followed by operators
// starting with precedence 1 up to unary operators. The highest
// precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.
//
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
//
func (op Token) Precedence() int {
	switch op {
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3
	}
	return LowestPrec
}

type lexemma struct {
	lexType Token
}
