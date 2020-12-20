var nodes = null;
var edges = null;
var id = 0;
var rootname = "";
var network = null;

function destroy() {
    if (network !== null) {
        network.destroy();
        network = null;
    }
}

function setup(obj, level) {
    var myId = id++;
    var label = "";
    if(obj["courses"]) {
        label = "Condition";
        if(obj["grade"] !== "") {
            var nextid = setup(obj["courses"], level + 1);
            edges.push({from: myId, to: nextid, label: obj["grade"]});
        } else {
            // remove unnecessary group nodes
            var nid = setup(obj["courses"], level);
            return nid;
        }
    } else if (obj["and"] || obj["or"]) {
        label = "Node";
        var andOr = obj["and"] ? obj["and"] : obj["or"];
        for(var next of andOr)  {
            var nextid = setup(next, level + 1);
            edges.push({from: myId, to: nextid, label: obj["and"] ? "and" : "or"});
        }
    } else if (obj["course"]) {
        label = obj["course"];
        if(obj["grade"] && obj["grade"] !== "") label += " " + obj["grade"];
    }
    if(level == 0) label = rootname;
    nodes.push({id: myId, label: label, level: level});
    return myId;
}

function draw(obj) {
    destroy();
    nodes = [];
    edges = [];
    id = 0;
    setup(obj, 0);

    // create a network
    var container = document.getElementById("mynetwork");
    var data = {
        nodes: nodes,
        edges: edges,
    };

    var options = {
        edges: {
            smooth: {
                type: "cubicBezier",
                forceDirection:
                    "horizontal",
                roundness: 0.4,
            },
        },
        layout: {
            hierarchical: {
                direction: "LR"
            },
        },
        physics: true,
    };
    network = new vis.Network(container, data, options);

    // add event listeners
    network.on("select", function (params) {
        document.getElementById("selection").innerHTML =
            "Selection: " + params.nodes;
    });
}

function drawGraph(obj, course) {
    rootname = course;
    draw(obj);
    setTimeout(() => {
        network.fit();
    }, 1000);
}