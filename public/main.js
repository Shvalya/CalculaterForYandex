var operationPlus = 200;
var operationMinus = 200;
var operationUmnozhit = 200;
var operationDelit = 200;

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

function onClickButtonFirst() {
    var text = document.getElementById('input-number');
    var expression = text.value;
    var time = 0;

    var plus = operationPlus;
    var minus = operationMinus;
    var umnozhit = operationUmnozhit;
    var delit = operationDelit;

    for (var i=0; i<expression.length;i++){
        if (expression[i]=="+"){
            time += plus; 
        }
        if (expression[i]=="-"){
            time += minus; 
        }
        if (expression[i]=="*"){
            time += umnozhit; 
        }
        if (expression[i]=="/"){
            time += delit; 
        }
    }


    axios.post('/calculate', {
        expression: expression
    })
        .then(function (response) {
            var listRes = document.getElementById('list-result');
            var listItem = document.createElement('li');

            if (response.data == null){
                listItem.textContent = expression + "= ? выражение неккоректно";
                listRes.appendChild(listItem);
                return;
            }

            setTimeout(function() {
                listItem.textContent = expression + ' = ' + response.data + "  время выполнения: " + time + "ms";
                listRes.appendChild(listItem);
                text.value = "";
            });
            text.value = ""
        })
        .catch(function (error) {
            console.error(error);
        });


}

