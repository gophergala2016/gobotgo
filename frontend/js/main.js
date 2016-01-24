
var data = [];
var size = 19;
var toastValue = "The toast is shown.";

function test() {
    for (var i = 0; i < size; i++){

        data[i] = [];
        for (var j = 0; j < size; j++){
            data[i][j] = Math.floor((Math.random() * 3.0) + 1) - 1;
        }
    }
    drawTable(data);
}

$('.showToast').click(function () {
    $('.error').text(toastValue);
    $('.error').stop().fadeIn(400).delay(3000).fadeOut(400); //fade out after 3 seconds
});

function drawTable(data) {
    var header = $("<tr>")
    $("#GameBoard").append(header);

    for (var i = 0; i < data.length+1; i++) {
        header.append("<th>" + i + "</th>");
    }

    header.append("</tr>");

    $("GameBoard").append("<tr>");

    for (var i = 0; i < data.length; i++) {
        drawRow(data[i], i);
    }
}

function drawRow(rowData, currentRow) {

    var color = ""
    var row = $("<tr />")
    $("#GameBoard").append(row);
    row.append($("<td>" + String.fromCharCode(65+currentRow) + "</td>"));


    for (var j = 0; j < rowData.length; j++) {

        if ( rowData[j] == 0) {
            color = "<img src=img/null.png>"
        }       
        else if ( rowData[j] == 1 ) {
            color = "<img src=img/black.png>"
        }  
        else if ( rowData[j] == 2 ) {
            color = "<img src=img/white.png>"
        }

        row.append($("<td>" + color + "</td>"));
    }
}
