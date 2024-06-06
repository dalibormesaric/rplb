const socket = new WebSocket("ws://" + window.location.host + "/ws");

socket.addEventListener("message", (event) => {
    if (event.data !== "ping") {
        // console.log("Message from server: ", event.data);
        json = JSON.parse(event.data);
        mg = document.getElementsByName(json.Name)[0].querySelector('.monitor-grid')
        if (mg.children.length >= 20) {
            mg.removeChild(mg.children[0])
        }
        mg.appendChild(createMonitorGridItem(json));
    }
});

function createMonitorGridItem(liveMonitorFrame) {
    let div = document.createElement('div');
    div.classList.add("latency-cell");
    div.classList.add("latency-" + liveMonitorFrame.ColorCode);
    if (liveMonitorFrame.Alive) {
        div.title = liveMonitorFrame.Latency
    } else {
        div.title = "unavailable"
    }
    return div;
}
