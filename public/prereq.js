let nodes = null;
let edges = null;
let id = 0;
let rootname = '';
let network = null;

function destroy() {
    if (network !== null) {
        network.destroy();
        network = null;
    }
}

function setup(obj, level) {
    let myId = id++;
    let label = '';
    if(obj['courses']) {
        label = 'Condition';
        if(obj['grade'] !== '') {
            let nextid = setup(obj['courses'], level + 1);
            edges.push({from: myId, to: nextid, label: obj['grade']});
        } else {
            // remove unnecessary group nodes
            let nid = setup(obj['courses'], level);
            return nid;
        }
        if(level == 0) label = rootname;
        nodes.push({id: myId, label: label, level: level});
        return myId;
    } else if (obj['and'] || obj['or']) {
        label = 'Node';
        let andOr = obj['and'] ? obj['and'] : obj['or'];
        for(let next of andOr)  {
            let nextid = setup(next, level + 1);
            edges.push({from: myId, to: nextid, label: obj['and'] ? 'and' : 'or'});
        }
        if(level == 0) label = rootname;
        nodes.push({id: myId, label: label, level: level});
        return myId;
    } else if (obj['course']) {
        label = obj['course'];
        if(obj['grade'] && obj['grade'] !== '') label += ' ' + obj['grade'];
        if(level == 0) {
            nodes.push({id: 0, label: rootname, level: 0});
            nodes.push({id: 1, label: label, level: 1});
            edges.push({from: 0, to: 1});
        } else {
            nodes.push({id: myId, label: label, level: level});
            return myId;
        }
    }
}

function draw(obj) {
    destroy();
    nodes = [];
    edges = [];
    id = 0;
    setup(obj, 0);

    // create a network
    const container = document.getElementById('mynetwork');
    const data = {
        nodes: nodes,
        edges: edges,
    };

    const options = {
        edges: {
            smooth: {
                type: 'cubicBezier',
                forceDirection:
                    'horizontal',
                roundness: 0.4,
            },
        },
        layout: {
            hierarchical: {
                direction: 'LR'
            },
        },
        physics: true,
    };
    network = new vis.Network(container, data, options);

    // add event listeners
    network.on('select', function (params) {
        document.getElementById('selection').innerHTML =
            'Selection: ' + params.nodes;
    });
}

function drawGraph(obj, course) {
    rootname = course;
    draw(obj);
    setTimeout(() => {
        network.fit();
    }, 1000);
}
