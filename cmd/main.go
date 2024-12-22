package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type CalcResponse struct { // Body ответа от сервера
	Result float64 `json:"result,omitempty"` // тк поле может быть опциональным - либо ошибки нет, есть результат
	Error  string  `json:"error,omitempty"`  // тк поле может быть опциональным - либо ошибка есть, нет результата
}

type CalcRequest struct { // Body запроса на вход
	Expression string `json:"expression"`
}

// Проверить, что символ - это оператор (бинарный)
func isBinaryOperator(symbol string) bool {
	switch symbol {
	case "+":
		return true
	case "-":
		return true
	case "*":
		return true
	case "/":
		return true
	default:
		return false
	}
}

// Получить "вес" оператора
func getWeight(operator string) int {
	switch operator {
	case "(":
		return 0
	case ")":
		return 1
	case "+":
		return 2
	case "-":
		return 2
	case "*":
		return 3
	case "/":
		return 3
	default:
		return -1
	}
}

// Вычислить выражение в инфиксной форме (вида: operand1 operator operand2)
func calcInfix(operator string, operand1, operand2 float64) float64 {
	switch operator {
	case "+":
		return operand1 + operand2
	case "-":
		return operand1 - operand2
	case "*":
		return operand1 * operand2
	case "/":
		return operand1 / operand2
	default:
		return 0
	}
}

// Получить строковое представление числа из строки, начиная с позиции pos
func getNumberFromString(expression string, pos *int) string {
	var output string

	// Идем посимвольно, пока символ это: цифра или "."
	for ; *pos < len(expression); *pos++ {
		symbol := string((expression)[*pos])
		_, err := strconv.Atoi(symbol)

		if err == nil || symbol == "." {
			// Если символ - цифра (ошибок нет), или "." - то добавляем его в результат

			output += symbol
		} else {
			// Иначе - выходим (и уменьшаем pos, чтобы в дальнейшем не пропустить символ)

			*pos--
			break
		}
	}

	return output
}

// Конвертация выражения из инфиксной в постфиксную запись
func infixToPostfix(infixExpression string) (string, error) {
	var postfixExpression string // Выходная строка
	var operationsStack []string // Стек операторов (так-то массив, но будет вести себя как стек)

	if _, err := validateInfix(infixExpression); err != nil {
		// Валидация входной строки
		return postfixExpression, err
	}

	for i := 0; i < len(infixExpression); i++ {
		symbol := string((infixExpression)[i])

		if _, err := strconv.Atoi(symbol); err == nil || symbol == "." {
			// Если символ - цифра (проверка на число без ошибок), или ".",
			// то получаем остальную часть числа, и добавляем его в выходную строку

			postfixExpression += getNumberFromString(infixExpression, &i) + " "
		} else if isBinaryOperator(symbol) {
			// Если символ - оператор, то сравниваем его приоритет
			// с приоритетом оператора на верхушке стека (последний эл-т в operationsStack)

			if len(operationsStack) == 0 ||
				getWeight(operationsStack[len(operationsStack)-1]) < getWeight(symbol) {
				// Если стек пуст, или приоритет symbol > приоритета верхушки стека,
				// то добавляем его в стек

				operationsStack = append(operationsStack, symbol)
			} else {
				// Иначе (приоритет symbol <= приоритета верхушки стека),
				// то вытаскиваем эл-ты из стека, пока не встретим меньший по приоритету или конец стека

				for len(operationsStack) > 0 {
					if getWeight(operationsStack[len(operationsStack)-1]) < getWeight(symbol) {
						// Если приоритет оператора на верхушке стека
						// меньше приоритета symbol - прекращаем вытаскивать элементы
						break
					}

					// После всех проверок, просто вытаскиваем оператор из стека в выходную строку
					postfixExpression += operationsStack[len(operationsStack)-1] + " "
					// И уменьшаем стек, тк оператор мы вытащили
					operationsStack = operationsStack[:len(operationsStack)-1]
				}

				// Ещё раз сравниваем приоритеты,
				// но уже с обновлённым стеком,
				// чтобы не потерять текущий оператор (symbol)
				if len(operationsStack) == 0 ||
					getWeight(operationsStack[len(operationsStack)-1]) < getWeight(symbol) {
					operationsStack = append(operationsStack, symbol)
				}
			}
		} else if symbol == "(" {
			// Если нашли откр. скобку - просто добавляем в стек
			operationsStack = append(operationsStack, symbol)
		} else if symbol == ")" {
			// Если нашли закр. скобку - вытаскиваем всё из стека,
			// пока не найдем откр. скобку или конец стека

			for len(operationsStack) > 0 {
				if operationsStack[len(operationsStack)-1] == "(" {
					// Если нашли откр. скобку - цель достигнута, прекращаем вытаскивать элементы
					break
				}

				// После всех проверок, просто вытаскиваем оператор из стека в выходную строку
				postfixExpression += operationsStack[len(operationsStack)-1] + " "
				// И уменьшаем стек, тк оператор мы вытащили
				operationsStack = operationsStack[:len(operationsStack)-1]
			}

			// После - удаляем из стека саму откр. скобку
			operationsStack = operationsStack[:len(operationsStack)-1]

		}
	}

	// Если входная строка закончилась, а стек ещё не пустой,
	// то просто вытаскиваем всё из стека в выходную строку
	for len(operationsStack) > 0 {
		postfixExpression += operationsStack[len(operationsStack)-1] + " "
		operationsStack = operationsStack[:len(operationsStack)-1]
	}

	// Результат - постфиксное выражение
	return postfixExpression, nil
}

// Проверка на валидность скобок
func checkParentheses(expression string) bool {
	var parenthesesStack []string // стек для скобок (так-то массив, но будет вести себя как стек)

	for i := 0; i < len(expression); i++ {
		symbol := string((expression)[i])

		if symbol == "(" {
			// Нашли откр. скобку - просто добавляем в стек
			parenthesesStack = append(parenthesesStack, symbol)
			continue
		}
		if symbol != ")" {
			// Отсеиваем всё, кроме закр. скобки
			continue
		}

		if len(parenthesesStack) == 0 {
			// Если стек пуст, а цикл идёт,
			// значит для текущей закр. скобки пары нет
			return false
		}

		// На этом моменте остались только закр. скобки,
		// просто вытаскиваем ей пару (откр. скобку) из стека
		parenthesesStack = parenthesesStack[:len(parenthesesStack)-1]
	}

	// Если стек пуст и вся строка пройдена, значит все скобки сбалансированы
	return len(parenthesesStack) == 0
}

// Проверка на валидность бинарных операций
// (что каждым двум операндам соответствует одна операция)
func checkBinaryOperations(expression string) bool {
	operandsCount, operatorsCount := 0, 0

	// Идём по символам в строке
	for i := 0; i < len(expression); i++ {
		symbol := string((expression)[i])
		_, err := strconv.Atoi(symbol)

		if err == nil || symbol == "." {
			// Нашли число (операндов) - увеличиваем счетчик операндов
			getNumberFromString(expression, &i)
			operandsCount++
		} else if isBinaryOperator(symbol) {
			// Нашли знак операции (оператор) - увеличиваем счетчик операций
			operatorsCount++
		}
	}

	// Если такая разница равна единице,
	// то каждым двум операндам соответствует одна операция => бинарность соблюдена
	return operandsCount-operatorsCount == 1
}

// Валидация символов в строке
func checkSymbols(expression string) bool {
	pattern := `^[0-9\.\+\-\*\/()\s]+$` // Пробелы должны быть в начале паттерна
	return regexp.MustCompile(pattern).MatchString(expression)
}

// Валидация инфиксного выражения
func validateInfix(expression string) (bool, error) {
	if len(expression) == 0 {
		// Строка пустая
		return false, errors.New("expression is empty")
	}
	if !checkSymbols(expression) {
		// Есть посторонние символы
		return false, errors.New("expression contains invalid characters")
	}
	if !checkParentheses(expression) {
		// Строка не прошла проверку на валидность скобок
		return false, errors.New("parentheses is not valid")
	}
	if !checkBinaryOperations(expression) {
		// Строка не прошла проверку на бинарность операций
		// (что каждым двум операндам соответствует одна операция)
		return false, errors.New("binary operations is not valid")
	}

	// Если все проверки пройдены - успех
	return true, nil
}

// calc Вычисление выражения (записанного в инфиксной форме)
func calc(expression string) (float64, error) {
	postfixExpression, err := infixToPostfix(expression)

	// Конвертировали в постфиксную форму, проверяем на наличие ошибок
	if err != nil {
		return 0, err
	}

	var operandsStack []float64 // Стек операндов (так-то массив, но будет вести себя как стек)

	// Идем по постфиксной записи посимвольно
	for i := 0; i < len(postfixExpression); i++ {
		symbol := string(postfixExpression[i])

		if _, err := strconv.Atoi(symbol); err == nil || symbol == "." {
			// Если символ - цифра (проверка на число без ошибок), или ".",
			// то получаем остальную часть числа, и добавляем его в стек операндов

			operand, err := strconv.ParseFloat(getNumberFromString(postfixExpression, &i), 64)
			if err != nil {
				return 0, err
			}

			operandsStack = append(operandsStack, operand)
		} else if isBinaryOperator(symbol) {
			// Если символ - оператор, то достаем два верхних операнда из стека,
			// вычисляем результат операции и сохраняем его обратно в стек

			var operand1, operand2 float64

			// Достаем два верхних операнда из стека
			// (в обратном порядке, тк договорились что массив ведет себя как стек)
			if len(operandsStack) > 0 {
				operand2 = operandsStack[len(operandsStack)-1]
				operandsStack = operandsStack[:len(operandsStack)-1]
			}
			if len(operandsStack) > 0 {
				operand1 = operandsStack[len(operandsStack)-1]
				operandsStack = operandsStack[:len(operandsStack)-1]
			}

			// Проверка на деление на ноль,
			// тк может появиться на этапе вычисления
			if symbol == "/" && operand2 == 0.0 {
				return 0, errors.New("division by zero")
			}

			// Добавляем результат операции в стек операндов
			operandsStack = append(operandsStack, calcInfix(symbol, operand1, operand2))
		}
	}

	// Результат - оставшееся число из стека
	return operandsStack[len(operandsStack)-1], nil
}

func startServer() {
	http.HandleFunc("/api/v1/calculate", CalcHandler)

	log.Println("Server started")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Error starting server:", err)
	}
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CalcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Error decoding JSON body:", err)
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(CalcResponse{Error: "Expression is not valid"})
		return
	}
	defer r.Body.Close()

	expression := req.Expression
	result, err := calc(expression)

	if err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(CalcResponse{Error: "Expression is not valid"})
		return
	} else {
		log.Println("Calculated successfully:", result)
	}

	resp := CalcResponse{Result: result}
	jsonResponse, err := json.Marshal(resp)

	if err != nil {
		log.Println("Error marshaling JSON:", err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(CalcResponse{Error: "Expression is not valid"})
		return
	}

	w.WriteHeader(200)
	w.Write(jsonResponse)
}

func main() {
	startServer()
}
