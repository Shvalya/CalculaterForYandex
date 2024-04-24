var operationPlus = 200;
var operationMinus = 200;
var operationUmnozhit = 200;
var operationDelit = 200;

function onClickEntry(){
    window.location.replace("/entry.html");
}

function onClickButtonSetting(operation){
    if (operation == "plus"){
        operationPlus = parseInt(document.getElementById('operation-plus').value)
    }
    if (operation == "minus"){
        operationMinus = parseInt(document.getElementById('operation-minus').value)
    }
    if (operation == "umnozhit"){
        operationMinus = parseInt(document.getElementById('operation-umnozhit').value)
    }
    if (operation == "delit"){
        operationMinus = parseInt(document.getElementById('operation-delit').value)
    }
}

var userList = [];
var currentUser = '';

function onClickButtonSetting(operation){
    var value = document.getElementById('operation-' + operation).value;
    axios.post('/settings', {
        operation: operation,
        value: value
    })
        .then(function (response) {
            alert("Параметр успешно изменен");
        })
        .catch(function (error) {
            console.error(error);
            alert("Ошибка изменения параметра");
        });
}
function onClickButtonFirst() {
    var expression = document.getElementById('input-number').value;
    var userLog = localStorage.getItem('userLogin');
    var time = 0;
    var plus = parseInt(document.getElementById('operation-plus').value);
    var minus = parseInt(document.getElementById('operation-minus').value);
    var umnozhit = parseInt(document.getElementById('operation-umnozhit').value);
    var delit = parseInt(document.getElementById('operation-delit').value);
    axios.post('/calculate', {
        expression: expression,
        userLogin: userLog,
        timePlus: document.getElementById('operation-plus').value,
        timeSub: document.getElementById('operation-minus').value,
        timeMultiply: document.getElementById('operation-umnozhit').value,
        timeDivision: document.getElementById('operation-delit').value,
    })
        .then(function (response) {
            var result = response.data.result;
            if (typeof result !== 'undefined') {
                var expressionResult = expression + ' = ' + result;
                var listItem = document.createElement('li');
                listItem.textContent = expressionResult;
                for (var i=0; i<expression.length;i++){
                    if (expression[i] == "+"){
                        time += plus;
                    }
                    if (expression[i] == "-"){
                        time += minus;
                    }
                    if (expression[i] == "*"){
                        time += umnozhit;
                    }
                    if (expression[i] == "/"){
                        time += delit;
                    }
                }
                setTimeout(function() {
                    listItem.textContent = expressionResult +  "  время выполнения: " + time + "ms";
                    document.getElementById('list-result').appendChild(listItem);
                });

            } else {
                console.error("Error calculating expression:", expression);
                alert("Ошибка вычисления выражения");
            }
        })
        .catch(function (error) {
            console.error(error);
            alert("Ошибка вычисления выражения");
        });
}


function registerButton(){
    var login = document.getElementById('register-login').value;
    var password = document.getElementById('register-password').value;

    axios.post('/register', {
        login: login,
        password: password
    })
        .then(function (response) {
            alert("Пользователь успешно зарегистрирован");
            window.location.replace("/entry.html");
        })
        .catch(function (error) {
            console.error(error);
            alert("Ошибка регистрации пользователя");
        });
}
function entryButton() {
    var login = document.getElementById('entry-login').value;
    var password = document.getElementById('entry-password').value;
    axios.post('/entry', {
        login: login,
        password: password
    })
        .then(function(response) {
            if (response.data.success) {
                var user = response.data.userLogin;
                console.log(response.data.expressions);
                var expressions = response.data.expressions.join(", ");
                alert("Вход выполнен успешно.\nЛогин: " + user.login + "\nВыражения: " + expressions);
                localStorage.setItem('expressions', JSON.stringify(response.data.expressions));
                localStorage.setItem('userLogin', response.data.userLogin);
                currentUser = response.data.userLogin;
                console.log(currentUser);
                window.location.replace("/index.html");
            } else {
                alert("Ошибка входа пользователя");
            }
        })
        .catch(function(error) {
            console.error(error);
            alert("Ошибка входа пользователя");
        });
}
function fillListResult() {
    var listResult = document.getElementById('list-result');
    listResult.innerHTML = ''; // Очищаем список перед заполнением

    // Получаем выражения пользователя из localStorage
    var expressions = JSON.parse(localStorage.getItem('expressions'));
    if (expressions) {
        expressions.forEach(function(expression) {
            var mathExpression = expression;
            var listItem = document.createElement('li');
            listItem.textContent = mathExpression;
            listResult.appendChild(listItem);
        });
    }
}
fillListResult();
}

