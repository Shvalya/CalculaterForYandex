package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/calculate", CalculateHandler)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("public/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("public/js"))))
	http.ListenAndServe(":8888", nil)
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
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
	rpn := infixToRPN(expression)
	result := evaluateRPN(rpn)

	w.Write([]byte(strconv.FormatFloat(result, 'f', -1, 64)))
	return
}

func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}

func infixToRPN(expression string) []string {
	var result []string
	var stack []rune

	for _, char := range expression {
		switch {
		case char >= '0' && char <= '9':
			result = append(result, strconv.Itoa(int(char-'0')))
		case char == '(':
			stack = append(stack, char)
		case char == ')':
			for stack[len(stack)-1] != '(' {
				result = append(result, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]
		case char == '+' || char == '-' || char == '*' || char == '/':
			for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(char) {
				result = append(result, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, char)
		}
	}

	for len(stack) > 0 {
		result = append(result, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return result
}

func evaluateRPN(tokens []string) float64 {
	var stack []float64

	for _, token := range tokens {
		if num, err := strconv.Atoi(token); err == nil {
			stack = append(stack, float64(num))
		} else {
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				stack = append(stack, a/b)
			}
		}
	}

	return stack[0]
}
