const socket = new WebSocket("ws://" + window.location.host + "/ws");

socket.addEventListener("message", (event) => {
    if (event.data !== "ping") {
        // console.log("Message from server: ", event.data);
        json = JSON.parse(event.data);

        if (json.Type === "monitor") {
            jn = document.getElementsByName(json.Name)
            if (jn && jn.length > 0) {
                mg = jn[0].querySelector('.monitor-grid')
                if (mg.children.length >= 20) {
                    mg.removeChild(mg.children[0])
                }
                mg.appendChild(createMonitorGridItem(json));
            }
        }

        if (json.Type === "traffic-fe" || json.Type === "traffic-be") {
            jn = document.getElementsByName(json.Name + "hits")
            if (jn && jn.length === 1) {
                jn[0].innerText = json.Hits
            }

            if (json.Type === "traffic-be") {
                var cir = document.getElementsByName(json.Name + "cir")
                if (cir) {
                    var rectS = document.getElementsByName(json.FrontendName + "element")[0].getBoundingClientRect()
                    var fromLeft = rectS.right
                    var fromTop = rectS.top + (rectS.height / 2)

                    var rectD = jn[0].getBoundingClientRect()
                    var moveLeft = rectD.left
                    var moveTop = rectD.top + (rectD.height / 2)
                    cir[0].animate([
                        { transform: 'translateX(' + fromLeft + 'px) translateY(' + fromTop + 'px)', opacity: 1 },
                        { transform: 'translateX(' + moveLeft + 'px) translateY(' + moveTop + 'px)', opacity: 0 }
                    ], {
                        duration: 100
                    })
                }
            }
        }
    }
});

function createMonitorGridItem(liveMonitorFrame) {
    let div = document.createElement('div');
    div.classList.add("latency-cell");
    div.classList.add("latency-" + liveMonitorFrame.ColorCode);
    if (liveMonitorFrame.Live) {
        div.title = liveMonitorFrame.Latency
    } else {
        div.title = "unavailable"
    }
    return div;
}
