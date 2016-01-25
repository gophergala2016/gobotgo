
var gameRoot = "http://localhost:8100/api/v1/game/";
var startGame = gameRoot + "start/";
var size = 19;
var test_data = [];

var sendMove;
var receiveState;
var waitForServer;
var playerColor;
var playRoot;
var gameID;
var lastBoard;

// initialize some sample data and  draw the table containing it
function init() {
    for (var i = 0; i < size; i++){
        test_data[i] = [];
        for (var j = 0; j < size; j++){
            test_data[i][j] = "None";
        }
    }
    drawTable(test_data);
}

function setUpGame(data, status) {
    console.log("Game started. Your Color: " + data["color"] + "The game ID: " + data["ID"]);
    gameID = data["ID"];
    playerColor = data["color"];
    playRoot = gameRoot + "play/" + gameID + "/";
    sendMove = playRoot + "move/";
    waitForServer = playRoot + "wait/";
    receiveState = playRoot + "state/";
    $('.new').hide();
    $.get(receiveState, boardRefresh).fail(connectError);
}

function getState(data, success) {
    okay();
    $.get(receiveState, boardRefresh).fail(connectError);
    waitForMove();
}

function waitForMove(){
    $.get(waitForServer, moveReceived).fail(connectError);
}

function moveReceived(){
    showToast("move recieved", 2000);
    $.get(receiveState, boardRefresh).fail(connectError);
}

function boardRefresh(data, status) {
    var color;
    console.log("data", data);
    console.log(status);
    for ( var i = 0; i < data["board"].length; i++ ) {
        for ( var j = 0; j < data["board"].length; j++ ) {


            if ( data["board"][i][j] == "None") {
                color = "img/null.png"
            }       
            else if ( data["board"][i][j] == "Black" ) {
                color = "img/black.png"
            }  
            else if ( data["board"][i][j] == "White" ) {
                color = "img/white.png"
            }
            $('#GameBoard tr').eq(j).find('td').eq(i).find('img').attr('src', color);
        }
    }
    return;
}

function connectError(err){
    console.log(err);
    showToast("Server Error. Check console.");
}

function ajaxPost(url, inputData) {
    console.log(inputData);
    $.ajax({
    url: url,
    type: 'post',
    data: inputData,
    headers: {},
    dataType: 'json',
    success: function (data) {
        console.info(data);
    }
});
}

// temporary listeners - most events will be based on returns from POST requests
$('#GameBoard').on('click', 'td', function(_evt) {
    console.log("Clicked", this, _evt);
    $.post(sendMove, JSON.stringify([_evt.currentTarget.cellIndex, _evt.currentTarget.parentElement.rowIndex]), getState).fail(connectError);
});

$('#GameBoard').on('mouseenter', 'td', function(_evt) {
    console.log("X: " + _evt.currentTarget.cellIndex + ", Y: " + _evt.currentTarget.parentElement.rowIndex);
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
    $.post(sendMove, "", okay).fail(connectError);
});

function okay() {
    showToast("Send okay", 2000)
}

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

        if ( rowData[j] == "None") {
            color = "<img src=img/null.png>"
        }       
        else if ( rowData[j] == "Black" ) {
            color = "<img src=img/black.png>"
        }  
        else if ( rowData[j] == "White" ) {
            color = "<img src=img/white.png>"
        }

        row.append($("<td>" + color + "</td>"));
    }
}


