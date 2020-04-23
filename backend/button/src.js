function display_cmd_options(s){

 var print = function(message) {
     var d = document.createElement("div");
     d.textContent = message;
     output.appendChild(d);
 };


  function f(s){
    console.log(s)
    
    var ws;
    document.getElementById("issue_cmd").onclick = function(evt) {   

    if (ws) {
        return false;
    }
    ws = new WebSocket("{{.}}");
    ws.onopen = function(evt) {
        print("OPEN");
    }
    ws.onclose = function(evt) {
        print("CLOSE");
        ws = null;
    }
    ws.onmessage = function(evt) {
        print("RESPONSE: " + evt.data);
    }
    ws.onerror = function(evt) {
        print("ERROR: " + evt.data);
    }
    return false;
    };
  }
  
  var x = document.getElementById("cmd").value;
  if (x == "add_user"){
    console.log("add user");
    document.getElementById("java_canvas").innerHTML  = ""; // first blank it out
    document.getElementById("java_canvas").innerHTML  = "<label for=\"fname\">First name:</label>";
    document.getElementById("java_canvas").innerHTML += "<input type=\"text\" id=\"fname\" name=\"fname\"><br><br>";

    document.getElementById("java_canvas").innerHTML += "<label for=\"lname\">Last name:</label>";
    document.getElementById("java_canvas").innerHTML += "<input type=\"text\" id=\"lname\" name=\"lname\"><br><br>";
    document.getElementById("java_canvas").innerHTML += "<button id=\"issue_cmd\">Add User</button>";
    f();

  }else{
    document.getElementById("java_canvas").innerHTML = "";
  }
};
