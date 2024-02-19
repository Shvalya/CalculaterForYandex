function onClickButtonFirst() {
    var text = document.getElementById('input-number');
    var expression = text.value;
    var response;

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

            listItem.textContent = expression + ' = ' + response.data;
            listRes.appendChild(listItem);
            // console.log(response.data);
            // alert(response.data)
        })
        .catch(function (error) {
            console.error(error);
        });


}
