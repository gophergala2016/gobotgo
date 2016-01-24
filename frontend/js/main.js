
var gameRoot = "http://localhost:8100/api/v1/game/";
var startGame = gameRoot + "start/";
var size = 19;
var test_data = [];

var sendMove;
var receiveState;
var waitForServer;
var playerColor;
var gameID;

// initialize some sample data and  draw the table containing it
function init() {
    for (var i = 0; i < size; i++){
        test_data[i] = [];
        for (var j = 0; j < size; j++){
            test_data[i][j] = 0;
        }
    }
    drawTable(test_data);
}

function setUpGame(data, status) {
    console.log("Game started. Your Color: " + data[1] + "The game ID: " + data[0]);
    gameID = data[0];
    playerColor = data[1];
    sendMove = gameRoot + "/" + gameID + "/move/";
    waitForServer = gameRoot + "/" + gameID + "/wait/";
    receiveState = gameRoot + "/" + gameID + "/state/";
    $('.new').hide();
}

function boardRefresh(data, status) {
    console.log(data);
    console.log(status);
    drawTable(data);
    return;
}

function connectError(err){
    console.log(err);
    showToast("Server Error. Check console.");
}

// temporary listeners - most events will be based on returns from POST requests
$('#GameBoard').on('click', 'td', function(_evt) {
    console.log("Clicked", this, _evt);
});

$('#GameBoard').on('mouseenter', 'td', function(_evt) {
    console.log("Over", this, _evt);
    $('feedbackBox').text("X: " + _evt.currentTarget.cellIndex + ", Y: " + _evt.currentTarget.parentElement.rowIndex);
});

// New Game
$('.new').click(function () {
    $.get(startGame, setUpGame).fail(connectError);
});

$('.refresh').click(function () {
    $.get(receiveState, boardRefresh).fail(connectError);
});

$('.pass').click(function () {
    $.post();
});

// Activate the temporary notification 'toast' for _time ms with _message
function showToast(_message, _time) {
    $('.error').text(_message);
    $('.error').stop().fadeIn(400).delay(_time).fadeOut(400); //fade out after 3 seconds
}

function mouseEnter(_fn) {
    return function(_evt) {
        var relTarget = _evt.relatedTarget;
        if (this === relTarget || isAChildOf(this, relTarget)) {
            return;
        }
        _fn.call(this, _evt);
    }
}

function isAChildOf(_parent, _child) {
    if (_parent === _child) { return false; }

    while (_child && _child !== _parent) { 
        _child = _child.parentNode;
    }
    return _child === _parent;
}

// Render the board
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

// Generate one full row given the data and the row to be generated
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
