package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Login       string   `json:"login"`
	Password    string   `json:"password"`
	Expressions []string `json:"expressions"`
}

type Response struct {
	Success bool `json:"success"`
	User    User `json:"user,omitempty"`
}
type ResponseEntry struct {
	Success     bool     `json:"success"`
	Expressions []string `json:"expressions"`
	UserLogin   string   `json:"userLogin"`
}

func main() {
	http.HandleFunc("/calculate", CalculateHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/entry", EntryHandler)
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("public/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("public/js"))))
	http.ListenAndServe(":8080", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{}

	// Чтение тела запроса для получения данных пользователя
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can not read body", http.StatusInternalServerError)
		return
	}

	// Разбор данных пользователя из тела запроса
	var user User
	el := strings.Split(string(body), ",")
	user.Login = el[0][10 : len(el[0])-1]
	user.Password = el[1][12 : len(el[1])-2]

	// Подключение к базе данных SQLite
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		http.Error(w, "sqlite3 error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Создание таблицы пользователей, если она не существует
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        login TEXT PRIMARY KEY,
        password TEXT,
        expressions TEXT
    )`)
	if err != nil {
		http.Error(w, "sqlite3 error", http.StatusInternalServerError)
		return
	}

	// Вставка данных нового пользователя в базу данных
	_, err = db.Exec("INSERT INTO users (login, password, expressions) VALUES (?, ?, ?)", user.Login, user.Password, "")
	if err != nil {
		http.Error(w, "sqlite3 error", http.StatusInternalServerError)
		return
	}

	// Отправка успешного ответа
	response.Success = true
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func EntryHandler(w http.ResponseWriter, r *http.Request) {
	response := ResponseEntry{}

	// Чтение тела запроса для получения данных пользователя
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can not read body", http.StatusInternalServerError)
		return
	}

	// Разбор данных пользователя из тела запроса
	var user User
	el := strings.Split(string(body), ",")
	user.Login = el[0][10 : len(el[0])-1]
	user.Password = el[1][12 : len(el[1])-2]

	// Подключение к базе данных SQLite
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		http.Error(w, "sqlite3 error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Проверка аутентификации пользователя
	stmt, err := db.Prepare("SELECT expressions FROM users WHERE login = ? AND password = ?")
	if err != nil {
		http.Error(w, "sqlite3 error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Извлечение выражений пользователя из базы данных
	var expressions string
	err = stmt.QueryRow(user.Login, user.Password).Scan(&expressions)
	fmt.Println(expressions)
	if err != nil {
		http.Error(w, "sqlite3 error", http.StatusInternalServerError)
		return
	}

	response.Success = true
	response.Expressions = strings.Split(expressions, ",")
	response.UserLogin = user.Login
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	// Получение выражения от клиента
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can not read body", http.StatusInternalServerError)
		return
	}

	// Сохранение выражения в базе данных
	var data map[string]interface{}
	err = json.Unmarshal([]byte(string(body)), &data)
	if err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return
	}
	expression := data["expression"].(string)
	timePlus := data["timePlus"].(string)
	timeSub := data["timeSub"].(string)
	timeMultiply := data["timeMultiply"].(string)
	timeDivison := data["timeDivision"].(string)
	countPlus, _ := strconv.Atoi(timePlus)
	countSub, _ := strconv.Atoi(timeSub)
	countMultiply, _ := strconv.Atoi(timeMultiply)
	countDivision, _ := strconv.Atoi(timeDivison)
	totalTime := strings.Count(expression, "+")*countPlus + strings.Count(expression, "-")*countSub + strings.Count(expression, "*")*countMultiply + strings.Count(expression, "/")*countDivision
	seconds := totalTime / 1000
	time.Sleep(time.Duration(seconds) * time.Second)
	fmt.Println(expression)
	// Измерение времени выполнения
	startTime := time.Now()
	rpnTokens := convertToRPN(expression)
	result := evaluateRPN(rpnTokens)
	fmt.Println("считаем")
	fmt.Println(result)
	expression += " = " + strconv.FormatFloat(result, 'f', -1, 64)
	// Измерение времени выполнения
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	userLogin, ok := data["userLogin"].(string)
	if !ok {
		fmt.Println("Не удалось извлечь значение userLogin")
		return
	}

	err = saveExpressionToDB(userLogin, expression)
	if err != nil {
		http.Error(w, "Failed to save expression to database", http.StatusInternalServerError)
		return
	}
	response := struct {
		Result float64 `json:"result"`
		Time   int64   `json:"time"`
	}{
		Result: result,
		Time:   duration.Milliseconds(),
	}

	// Кодирование ответа в формат JSON и отправка клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type Stack []string

func (s *Stack) Push(str string) {
	*s = append(*s, str)
}

func (s *Stack) Pop() string {
	if len(*s) == 0 {
		return ""
	}
	str := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return str
}

func isOperator(op string) bool {
	return op == "+" || op == "-" || op == "*" || op == "/"
}

func priority(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func convertToRPN(expression string) []string {
	var output []string
	var stack Stack

	tokens := strings.Fields(expression)

	for _, token := range tokens {
		if !isOperator(token) && token != "(" && token != ")" {
			output = append(output, token)
		} else if token == "(" {
			stack.Push(token)
		} else if token == ")" {
			for stack[len(stack)-1] != "(" {
				output = append(output, stack.Pop())
			}
			stack.Pop()
		} else {
			for len(stack) > 0 && priority(stack[len(stack)-1]) >= priority(token) {
				output = append(output, stack.Pop())
			}
			stack.Push(token)
		}
	}

	for len(stack) > 0 {
		output = append(output, stack.Pop())
	}

	return output
}

func evaluateRPN(tokens []string) float64 {
	stack := Stack{}
	for _, token := range tokens {
		if !isOperator(token) {
			stack.Push(token)
		} else {
			num2, _ := strconv.ParseFloat(stack.Pop(), 64)
			num1, _ := strconv.ParseFloat(stack.Pop(), 64)
			result := applyOperator(token, num1, num2)
			stack.Push(fmt.Sprintf("%f", result))
		}
	}
	result, _ := strconv.ParseFloat(stack.Pop(), 64)
	return result
}

func applyOperator(op string, num1, num2 float64) float64 {
	switch op {
	case "+":
		return num1 + num2
	case "-":
		return num1 - num2
	case "*":
		return num1 * num2
	case "/":
		if num2 != 0 {
			return num1 / num2
		} else {
			panic("Division by zero")
		}
	default:
		panic("Invalid operator")
	}
}

func getExpressionsFromDB(login string) ([]string, error) {
	// Получение выражений пользователя из базы данных
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var expressions string
	err = db.QueryRow("SELECT expressions FROM users WHERE login = ?", login).Scan(&expressions)
	if err != nil {
		return nil, err
	}

	return strings.Split(expressions, ","), nil
}
func saveExpressionToDB(login, expression string) error {
	// Сохранение выражения пользователя в базе данных
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return err
	}
	defer db.Close()
	fmt.Println(expression)
	fmt.Println("1233")
	// Добавление выражения в список выражений пользователя
	_, err = db.Exec("UPDATE users SET expressions = COALESCE(expressions || ',', '') || ? WHERE login = ?", expression, login)
	if err != nil {
		return err
	}

	return nil
}
