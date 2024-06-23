package main

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestCalculateSum(t *testing.T) {
	cases := []struct {
		name        string
		lines       []string
		expectedSum string
		expectedErr error
	}{
		{
			name:        "Успех",
			lines:       []string{"1.1", "2.2", "3.3"},
			expectedSum: "6.6",
		},
		{
			name:        "Часть строк некорректна",
			lines:       []string{"1.1", "abc", "3.3"},
			expectedSum: "4.4",
			expectedErr: errors.New("Ошибочные строки: abc Ошибка: can't convert abc to decimal. Они не включены в подсчёт суммы."),
		},
		{
			name:        "Все строки некорректны",
			lines:       []string{"abc", "xyz"},
			expectedSum: "0",
			expectedErr: errors.New("Ошибочные строки: abc Ошибка: can't convert abc to decimal xyz Ошибка: can't convert xyz to decimal. Они не включены в подсчёт суммы."),
		},
		{
			name:        "Пустой вход",
			lines:       []string{},
			expectedSum: "0",
		},
		{
			name:        "Большие положительные числа",
			lines:       []string{"999999999999999999999999999999.123456789", "888888888888888888888888888888.123456789"},
			expectedSum: "1888888888888888888888888888887.246913578",
		},
		{
			name:        "Маленькие положительные числа",
			lines:       []string{"0.000000000000000000000000000001", "0.000000000000000000000000000002"},
			expectedSum: "0.000000000000000000000000000003",
		},
		{
			name:        "Большие положительные с маленькими положительными числами",
			lines:       []string{"999999999999999999999999999999.999999999999999999999999999999", "0.000000000000000000000000000001"},
			expectedSum: "1000000000000000000000000000000",
		},
		{
			name:        "Большие положительные с нулем",
			lines:       []string{"999999999999999999999999999999.123456789", "0"},
			expectedSum: "999999999999999999999999999999.123456789",
		},
		{
			name:        "Маленькие положительные с нулем",
			lines:       []string{"0.000000000000000000000000000001", "0"},
			expectedSum: "0.000000000000000000000000000001",
		},
		{
			name:        "Ноль с нулем",
			lines:       []string{"0", "0"},
			expectedSum: "0",
		},
		{
			name:        "Большие отрицательные числа",
			lines:       []string{"-999999999999999999999999999999.123456789", "-888888888888888888888888888888.123456789"},
			expectedSum: "-1888888888888888888888888888887.246913578",
		},
		{
			name:        "Маленькие отрицательные числа",
			lines:       []string{"-0.000000000000000000000000000001", "-0.000000000000000000000000000002"},
			expectedSum: "-0.000000000000000000000000000003",
		},
		{
			name:        "Большие отрицательные с маленькими отрицательными числами",
			lines:       []string{"-999999999999999999999999999999.999999999999999999999999999999", "-0.000000000000000000000000000001"},
			expectedSum: "-1000000000000000000000000000000",
		},
		{
			name:        "Большие отрицательные с нулем",
			lines:       []string{"-999999999999999999999999999999.123456789", "0"},
			expectedSum: "-999999999999999999999999999999.123456789",
		},
		{
			name:        "Маленькие отрицательные с нулем",
			lines:       []string{"-0.000000000000000000000000000001", "0"},
			expectedSum: "-0.000000000000000000000000000001",
		},
		{
			name:        "Ноль с нулем",
			lines:       []string{"0", "0"},
			expectedSum: "0",
		},
		{
			name:        "Большие положительные с маленькими отрицательными",
			lines:       []string{"999999999999999999999999999999.123456789", "-0.000000000000000000000000000001"},
			expectedSum: "999999999999999999999999999999.123456788999999999999999999999",
		},
		{
			name:        "Большие отрицательные с маленькими положительными",
			lines:       []string{"-999999999999999999999999999999.123456789", "0.000000000000000000000000000001"},
			expectedSum: "-999999999999999999999999999999.123456788999999999999999999999",
		},
		{
			name:        "Маленькие отрицательные с маленькими положительными",
			lines:       []string{"-0.000000000000000000000000000001", "0.000000000000000000000000000002"},
			expectedSum: "0.000000000000000000000000000001",
		},
		{
			name:        "Большие отрицательные с большими положительными",
			lines:       []string{"-999999999999999999999999999999.123456789", "999999999999999999999999999999.123456789"},
			expectedSum: "0",
		},
		{
			name:        "Большое отрицательное, большое положительное, маленькое отрицательное, маленькое положительное, нуль",
			lines:       []string{"-999999999999999999999999999999.123456789", "999999999999999999999999999999.123456789", "-0.000000000000000000000000000001", "0.000000000000000000000000000002", "0"},
			expectedSum: "0.000000000000000000000000000001",
		},
	}
	for _, test_case := range cases {
		test_case := test_case
		t.Run(test_case.name, func(t *testing.T) {
			t.Parallel()
			sum, err := calculate_sum(test_case.lines)
			if sum != test_case.expectedSum {
				t.Errorf("Ожидалась сумма %v, получено %v.", test_case.expectedSum, sum)
			}
			if (err != nil && test_case.expectedErr == nil) || (err == nil && test_case.expectedErr != nil) || (err != nil && test_case.expectedErr != nil && err.Error() != test_case.expectedErr.Error()) {
				t.Errorf("Ожидалась ошибка %v, получено %v.", test_case.expectedErr, err)
			}
		})
	}
}

func TestMain(t *testing.T) {
	validNumbersFile, err := createTempFile("1.1\n2.2\n3.3\n")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(validNumbersFile)

	invalidNumbersFile, err := createTempFile("1.1\nabc\n3.3\n")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(invalidNumbersFile)

	cases := []struct {
		name           string
		file           string
		expectedOutput string
	}{
		{
			name:           "Успех",
			file:           validNumbersFile,
			expectedOutput: "Cумма: 6.6\n",
		},
		{
			name:           "Некоторые строки неправильные",
			file:           invalidNumbersFile,
			expectedOutput: "Ошибочные строки: abc Ошибка: can't convert abc to decimal. Они не включены в подсчёт суммы.\nCумма: 4.4\n",
		},
	}

	for _, test_case := range cases {
		test_case := test_case
		t.Run(test_case.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", "main.go", test_case.file)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Ошибка не ожидалась, получна %v", err)
			}
			if string(output) != test_case.expectedOutput {
				t.Errorf("Ожидался вывод %v, получен %v", test_case.expectedOutput, string(output))
			}
		})
	}
}

func createTempFile(content string) (string, error) {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		return "", err
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}
