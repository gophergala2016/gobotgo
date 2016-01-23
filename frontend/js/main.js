
var data = [
    ["e", "e", "e"],
    ["e", "e", "e"],
    ["e", "e", "e"]
];

function test() {
    drawTable(data);
}

function drawTable(data) {
    var header = $("<tr>")
    $("#GameBoard").append(header);

    for (var i = 0; i < data.length; i++) {
        header.append("<th>" + i + "</th>");
    }

    header.append("</tr>");

    $("GameBoard").append("<tr>");

    for (var i = 0; i < data.length; i++) {
        drawRow(data[i]);
    }
}

function drawRow(rowData) {
    var row = $("<tr />")
    
    for (var j = 0; j < rowData.length; j++) {
        $("#GameBoard").append(row);
        row.append($("<td>" + rowData[j] + "</td>"));
    }
}
