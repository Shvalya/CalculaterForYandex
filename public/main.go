package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/calculate", CalculateHandler)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("public/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("public/js"))))
	http.ListenAndServe(":8888", nil)
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can not read body", http.StatusInternalServerError)
		return
	}

	expression := string(body)
	result := evalExpression(expression)

	w.Write([]byte(strconv.FormatFloat(result, 'f', -1, 64)))
	return
}

func evalExpression(expression string) float64 {
	expression = expression[15:]
	expression = expression[:len(expression)-2]

	tokens := strings.Split(expression, " ")

	operands := []float64{}

	operators := []string{}

	for _, token := range tokens {
		if num, err := strconv.Atoi(token); err == nil {
			operands = append(operands, float64(num))
		} else {
			operators = append(operators, token)
		}
	}

	for len(operators) > 0 {
		op1 := operands[len(operands)-1]
		operands = operands[:len(operands)-1]
		op2 := operands[len(operands)-1]
		operands = operands[:len(operands)-1]

		operator := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		var result float64
		switch operator {
		case "*":
			result = op1 * op2
		case "/":
			result = op1 / op2
		case "+":
			result = op1 + op2
		case "-":
			result = op1 - op2
		}

		operands = append(operands, result)
	}
	return operands[0]
}
