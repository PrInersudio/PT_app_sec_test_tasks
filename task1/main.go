package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shopspring/decimal"
)

// Чтение строк из файла и сохранение их в срез
func parse_lines(file *os.File) []string {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

/*
Принимает на вход срез строк, представляющих числа с плавающей точкой.
Переводит их в числа с фиксированной точкой.
Считает сумму.
Переводит сумму в строку и возвращает.
Если часть строк не удалось преобразовать в числа,
возвращает ошибку со списком этих строк.
Иначе ошибка равна nil.
В подсчёт суммы ошибочные строки не включаются.
*/
func calculate_sum(lines []string) (string, error) {
	var failed_lines []string
	sum := decimal.Zero
	for _, line := range lines {
		num, err := decimal.NewFromString(line)
		if err != nil {
			failed_lines = append(failed_lines,
				fmt.Sprintf("%v Ошибка: %v", line, err))
			continue
		}
		sum = sum.Add(num)
	}
	if len(failed_lines) != 0 {
		return sum.String(),
			errors.New("Ошибочные строки: " +
				strings.Join(failed_lines, " ") +
				". Они не включены в подсчёт суммы.")
	}
	return sum.String(), nil
}

func main() {
	logfile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Ошибка создания файла логгирования:", err)
	}
	log.SetOutput(logfile)
	defer logfile.Close()
	log.Println("Программа запущена.")
	if len(os.Args) < 2 {
		log.Printf("Неправильное количество аргуметов командной строки. "+
			"Должно быть 2. Дано %v.\n", len(os.Args))
		fmt.Printf("Использовать: %v <файл с числами>\n", os.Args[0])
		os.Exit(1)
	}
	log.Printf("Открытие файла %v.\n", os.Args[1])
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Не удалось открыть файл %v. Ошибка: %v\n", os.Args[1], err)
		log.Fatalf("Не удалось открыть файл %v. Ошибка: %v\n", os.Args[1], err)
	}
	log.Printf("Открыт файл %v.\n", os.Args[1])
	log.Println("Чтение строк с числами.")
	lines := parse_lines(file)
	file.Close()
	log.Println("Строки считаны, файл закрыт.")
	log.Println("Подсчёт суммы.")
	sum_string, err := calculate_sum(lines)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
	}
	log.Printf("Cумма: %v. Будет выведена в консоль.\n", sum_string)
	fmt.Printf("Cумма: %v\n", sum_string)
	log.Println("Программа выполнена.")
}
