console.log("Script initialized")
let started = false
let cpuCreated = false

const socket = new WebSocket("ws://127.0.0.1:8081/ws/ps");

socket.onopen = function(e) {
    socket.send("stream#begin");
    setInterval(() => socket.send("stream#next"), 1000)
};

socket.onmessage = function(event) {
    obj = JSON.parse(event.data)
    if (obj.ready === true && !started) {
        started = true
    } else {
        updateProcessTable(obj)
        updateCPUs(obj)
    }
};

socket.onclose = function(event) {
    closeConnection()
    if (event.wasClean) {
        alert(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
    } else {
        // e.g. server process killed or network down
        // event.code is usually 1006 in this case
        alert('[close] Connection died');
    }
};

socket.onerror = function(error) {
    console.error(error)
    closeConnection()
};


function closeConnection() {
    socket.send("stream#end")
}


function updateProcessTable(data) {
    // Update process table
    var processTableBody = $('#process-table tbody');
    processTableBody.empty();

    for (cpu in data.cpu) {
        cpu = Number.parseFloat(+cpu).toFixed(2)
    }
    $.each(data.proc, function(i, process) {
        processTableBody.append('<tr><td>' + process.pid + '</td><td>' +
            process.name + '</td><td>' + process.cpu.toFixed(2) +
            '%</td><td>' + process.mem.toFixed(2) + '%</td></tr>');
    });
}

function updateCPUs(data) {
    if (!cpuCreated) {
        const cpuUsage = data.cpu;
        const $cpuContainer = $('<div>').addClass('cpu-items');
        for (let i = 0; i < cpuUsage.length; i++) {
            const $cpuBar = $('<progress>').val(cpuUsage[i]).attr('max', 100);
            const $cpuText = $('<p>')
            const $cpubox = $('<div>').prop('id', 'cpu-' + i).addClass('cpu-bar')
                .append($cpuText).append($cpuBar).append();
            $cpuContainer.append($cpubox);
        }
        $('.cpu-items').replaceWith($cpuContainer);
        cpuCreated = true
        return
    }
    const cpuUsage = data.cpu;
    $('.cpu-items').each(
        function(i, _v) { // iterate over all elements with class 'cpu-container'
            const $container = $(this); // wrap the current container in a jQuery object
            // find the progress bar and label inside the current container
            const $progressBar = $container.find('progress');
            const $label = $container.find('p');
            $progressBar.val(cpuUsage[i]); // set the width of the progress bar to the new value
            $label.text(i + ":" + cpuUsage[i].toFixed(2) + '%'); // set the label text to the new value, rounded to 1 decimal place
        });
    // Update memory usage
    $('#mem').val(data.mem * 100)
    $('#mem').find('label').val("Memory: " + (data.mem * 100) + "%");
    // Update disk usage
    $('#disk').val(data.disk)
//    $('#disk').find('label').label("Disk: " + data.disk + "%");
}

window.onbeforeunload = confirmExit;
function confirmExit() {
    closeConnection()
    closeTerminal()
    return true;
}
