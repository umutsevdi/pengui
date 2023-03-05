const term = new WebSocket("ws://127.0.0.1:8081/ws/term");
let terminal = {
    in: null,
    out: null,
    button: null,
    started: false
};

document.addEventListener("DOMContentLoaded", function() {
    console.log("Terminal Loaded");
    terminal = {
        in: document.getElementById('term-i'),
        out: document.getElementById('term-o'),
        button: document.getElementById('term-b'),
        started: false
    }
});


term.onopen = (e) => {
    term.send("term#init");
    terminal.in.value = ""
    terminal.out.value = ""
    setInterval(() => socket.send(JSON.stringify({ req: "fetch" })), 250)
};

addEventListener("click", function() {
    if (terminal.started && terminal.in.value != "") {
        term.send(
            JSON.stringify({
                req: "in",
                msg: terminal.in.value
            }));
    }
});

addEventListener("keyup", (event) => {
    if (event.keyCode === 13) {
        event.preventDefault();
        terminal.button.click();
    }
});

/**
* Sending input to server:
{
    "req":"in",
    "msg":"example message"
}
* Sending update request to the server:
{
    "req":"fetch",
}

* receive input response:
{
    "resp":"in",
    "msg":"ok",
}

* receive message:
{
    "resp":"fetch",
    "msg":"example message",
}
*/
term.onmessage = function(event) {
    if (!terminal.started) {
        terminal.started = true
        console.log("terminal handshake completed")
    }
    obj = JSON.parse(event.data)
    if (obj.resp == null) {
        return
    }
    if (obj.resp == "in" && obj.msg == "ok") {
        terminal.out.value += "~$ " + terminal.in.value + "\n"
        return
    }
    if (obj.resp == "fetch" && obj.msg != "") {
        terminal.out.value += event.data.msg + "\n";
        terminal.in.value = "";
        terminal.out.scrollTop = terminal.out.scrollHeight;
    }
}

term.onclose = function(event) {
    closeTerminal()
};

term.onerror = function(error) {
    alert(`[error]`);
    closeTerminal()
};


function closeTerminal() {
    term.send("term#end")
    terminal.started = false;
}


