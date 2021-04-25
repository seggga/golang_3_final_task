package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/seggga/golang_3_final_task/myscanner"
	"github.com/seggga/golang_3_final_task/mytoken"
	"github.com/seggga/golang_3_final_task/querier"
)

func main() {

	src := []byte(` select name,age from "file1.csv" where age>=30 and region=="Europe" and status == "sick" `)
	fmt.Printf("%s\n", src)

	if !querier.CheckQueryPattern(src) {
		fmt.Println("wrong query")
		return
	}

	// scanner initialisation
	var s myscanner.Scanner
	fset := mytoken.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, myscanner.ScanComments)

	var lm querier.LexMachine
	lm.Query = string(src)

	// run the scanner
	for {
		pos, tok, lit := s.Scan()
		if tok == mytoken.EOF {
			break
		}
		querier.AnalyseToken(&lm, lit, tok)
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}

	// check if the query contains at least one file to be read
	if len(lm.From) == 0 {
		fmt.Println("no file has been chosen (section FROM is empty)")
		return
	}

	// check if the query contains at least one column to be written to output
	if len(lm.From) == 0 {
		fmt.Println("no columns has been chosen (section SELECT is empty)")
		return
	}

	// read files
	// открытие файлов
	//		 прочитали заголовок CSV
	//		 проверка, все ли столбцы в блоке select есть в таблице из файла
	//		 проверка, все ли столбцы в блоке where есть в таблице из файла
	//		 считываем строки
	//		 		подставляем данные из таблицы в слайс вычисления
	//				вывод, если строка соответствует условию
	for _, fileName := range lm.From {

		if _, err := os.Stat(fileName); err != nil {
			log.Fatalf("file %s was not found. %v", fileName, err)
		}

		// file opening
		file, err := os.OpenFile(fileName, os.O_RDONLY, 600)
		if err != nil {
			log.Fatal(err)
		}

		// deferred function closing the file created
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		// read the header of the csv-file
		reader := csv.NewReader(file)  // Считываем файл с помощью библиотеки encoding/csv
		fileCols, err := reader.Read() //  Считываем шапку таблицы
		if err != nil {
			log.Fatalf("Cannot read file %s: %v", fileName, err)
		}

		// compare columns sets from the query and the file
		//columns := strings.Split(fileCols, ",")
		err = querier.CheckSelectedColumns(fileCols, lm)
		if err != nil {
			log.Fatalf("Неверный запрос: %v", err)
		}

		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error reading csv-file %s: %v", fileName, err)
			}
			// compose a map holding data of the current row
			rowData := querier.FillTheMap(fileCols, row, lm)
			// create a slice based on the conditions in WHERE-statement
			lexSlice := querier.MakeSlice(rowData, lm)

			if querier.Execute(lexSlice) {
				querier.PrintTheRow(rowData, lm)
			}
		}
	}

	// фильтруем столбцы для вывода
	fmt.Println(querier.TrimOutput([]string{"name", "age", "region"}, src))

	fmt.Println(querier.GetConditions(src))

	// aSlice := []querier.Lexemma{
	// 	{Typ: "operator", Val: "and"},
	// 	{Typ: "operand", Val: "31"},
	// 	{Typ: "operator", Val: ">="},
	// 	{Typ: "operand", Val: "30"},
	// 	{Typ: "operand", Val: "Europe"},
	// 	{Typ: "operator", Val: "=="},
	// 	{Typ: "operand", Val: "europe"},
	// 	{Typ: "operator", Val: "and"},
	// 	{Typ: "operand", Val: "sick"},
	// 	{Typ: "operator", Val: "=="},
	// 	{Typ: "operand", Val: "sick"},
	// }

	// fmt.Println(querier.Execute(aSlice))
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
