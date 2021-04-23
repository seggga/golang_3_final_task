package main

import (
	"fmt"
	"go/scanner"
	"go/token"

	"github.com/seggga/golang_3_final_task/querier"
)

func main() {

	src := []byte(` SELECT name,age FROM file where age>=30 and region=="Europe"`)

	// scanner initialisation
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	fmt.Printf("%s\n", src)

	// run the scanner
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}

	if !querier.CheckQueryPattern(src) {
		fmt.Println("wrong query")
	}

	// фильтруем столбцы для вывода
	fmt.Println(querier.TrimOutput([]string{"name", "age", "region"}, src))

	fmt.Println(querier.GetConditions(src))

	aSlice := []querier.Lexemma{
		{Typ: "operand", Val: "31"},
		{Typ: "operator", Val: ">="},
		{Typ: "operand", Val: "30"},
		{Typ: "operator", Val: "and"},
		{Typ: "operand", Val: "Europe"},
		{Typ: "operator", Val: "=="},
		{Typ: "operand", Val: "europe"},
		{Typ: "operator", Val: "and"},
		{Typ: "operand", Val: "sick"},
		{Typ: "operator", Val: "=="},
		{Typ: "operand", Val: "sick"},
	}

	fmt.Println(querier.Execute(aSlice))
}

/*
		select - name1, - name2, - name3, - nameN
		from file1, file2, ...fileN

		паттерны:
			select - from - where
			select - from

		в начале работы парсера проверяется паттерн запроса
		если нет совпадений - то ошибка


parserSelFrom:
1) выделяем подстроку между select и from, затем разбиваем ее запятыми на имена столбцов
	start

	if p.atStart && isSelect(tok) {
		next lexemma
	}

	if
}
*/
