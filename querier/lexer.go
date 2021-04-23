package querier

import (
	"regexp"
	"strings"
)

type Lexemma struct {
	Typ string
	Val string
}

// функция проверяет, содержатся ли указанные в select столбцы в переданном файле
// input parameters are:
//		slice of conumns names, obtainet from the *.csv file
//		user's query
// output: a slice of column names that where found in the query
func TrimOutput(allColumns []string, b []byte) []string {

	theQuery := string(b)
	theQuery = strings.ToLower(theQuery)
	theQuery = strings.TrimSpace(theQuery)

	theQuery = strings.TrimLeft(theQuery, "select")
	theQuery = strings.Split(theQuery, "from")[0]

	theQuery = strings.TrimSpace(theQuery)
	outColumns := strings.Split(theQuery, ",")

	counter := len(outColumns)
	for _, colInQuery := range outColumns {
		for _, colInTable := range allColumns {
			if colInQuery == colInTable {
				counter--
				break
			}
		}
	}

	if counter > 0 {
		return nil // неверный набор столбцов в запросе. Имеются столбцы...allColumns
	}

	return outColumns
}

// CheckQueryPattern - checks the query pattern
// if there is no matcing pattern, the query is incorrect
func CheckQueryPattern(b []byte) bool {

	theQuery := string(b)
	theQuery = strings.ToLower(theQuery)
	theQuery = strings.TrimSpace(theQuery)

	// проверка на первый key_word
	if !strings.HasPrefix(theQuery, "select") {
		return false
	}

	for _, patt := range QueryPatterns {
		matched, _ := regexp.Match(patt, []byte(theQuery))
		if matched {
			return true // также надо запомнить, какой паттерн подошел
		}
	}
	return false
}

func GetConditions(b []byte) []string {

	// obtain the substring that contains conditions only
	theQuery := string(b)
	theQuery = strings.ToLower(theQuery)
	theQuery = strings.TrimSpace(theQuery)

	theQuery = strings.Split(theQuery, "where")[1]

	return nil
}

func Execute(sl []Lexemma) bool {
	for i := 0; i < len(sl); i += 4 {
		// финальное вычисление ??
		if i+3 >= len(sl) {
			res := calculator(sl[i : i+3])
			if res.Val == "true" {
				return true
			}
			return false
		}
		// нефинальное вычисление
		sl = append(sl, calculator(sl[i:i+3]))
		if sl[i+3].Typ == "operator" {
			sl = append(sl, sl[i+3])
		} else {
			i -= 1
		}
	}
	return false
}

func calculator(ops []Lexemma) Lexemma {

	for i, op := range ops {
		if op.Typ == "operator" {
			return calculate(i, ops)
		}
	}

	return Lexemma{}
}

func calculate(i int, ops []Lexemma) Lexemma {

	var operand1, operand2 Lexemma
	switch i {
	case 0:
		operand1 = ops[1]
		operand2 = ops[2]
	case 1:
		operand1 = ops[0]
		operand2 = ops[2]
	case 2:
		operand1 = ops[0]
		operand2 = ops[1]
	}

	var result bool
	switch ops[i].Val {
	case ">":
		result = operand1.Val > operand2.Val
	case ">=":
		result = operand1.Val >= operand2.Val
	case "==":
		result = operand1.Val == operand2.Val
	case "<":
		result = operand1.Val < operand2.Val
	case "<=":
		result = operand1.Val <= operand2.Val
	case "and":
		result = (operand1.Val == "true") && (operand2.Val == "true")
	case "or":
		result = (operand1.Val == "true") || (operand2.Val == "true")
	}

	if result {
		return Lexemma{"bool", "true"}
	}
	return Lexemma{"bool", "false"}
}

/*
 slice += execute(slice[i], slice[i+1], slice[i+2])  // внутри функции идет проверка ситуации, есть ли четвертый элемент в группе.
 //если четвертого в группе нет, значит это финальное вычисление
 if slice[i+3] isOperator {
	 slice += slice[i+3]
 } else { // закончились операторы второго уровня
	i--
 }
	return false
}

/*   age >= 30 AND region=="Europe" AND status == "sick"

1) сформировать слайс операндов и действий с ними и выполнять по 3


age, >=, 30, AND, | region, ==, europe, AND, | status, ==, sick, | true, AND, false, AND, | true | false AND  >= false

===================================

/*   age >= 30 AND region=="Europe"

age, >=, 30, AND, | region, ==, europe | true, AND, false, | false

====================================

for i := 0; i < len(operators); i += 4

 slice += execute(slice[i], slice[i+1], slice[i+2])  // внутри функции идет проверка ситуации, есть ли четвертый элемент в группе.
 //если четвертого в группе нет, значит это финальное вычисление
 if slice[i+3] isOperator {
	 slice += slice[i+3]
 } else { // закончились операторы второго уровня
	i--
 }

if slice[len(slice)] -> print the string

/*


1) age >= 30 region == "europe" status == "sick"
2) result1 AND result2 AND result3
в перенос AND на позицию 7 - еще перенос на позицию +4 = 11
второй перенос AND на позицию

	lexerOut := []struct{
		lexType columnName / condition / value / operator /
		lexText string
		lexPriority int //
		resultPriority int // для condition

				columnName - 0
				condition - 1
				value - 0
				operator - 2

col - age - 0
cond - >= - 1
val - 30 - 0


for trippleToken = range trippleTokens {


}

1)
func exec (token1, condition, token3 ) token {
	if condition == ">=" {
		return token1 >= token3
	}

	if condition == "==" {
		return token1 == token3
	}

	if condition == "AND" {
		return token AND token
	}
}


call[0]
call[1]
call[2]


equation = GOE(col, val)



	}



	return nil
}
*/
