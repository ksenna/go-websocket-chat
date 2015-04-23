try {
  var socket = new WebSocket("ws://localhost:8000/socket");
  console.log("Websocket Status: " + socket.readyState);
}

catch(exception) {
  console.log(exception);
}

socket.onopen = function(m) { 
  console.log("Connection opened:" + this.readyState);
}

socket.onmessage = function(m) { 
  $('#mainChatArea').append('<p>' + m.data + '</p>');
}

socket.onerror = function(m) {
  console.log("Error occured sending:" + m.data);
}

socket.onclose = function(m) { 
  console.log("Disconnected Status: " + this.readyState);
}


$('#clientText').val("");
$('#send').on('click', function(event) {
  socket.send($('#clientText').val());
  $('#clientText').val("");
});
