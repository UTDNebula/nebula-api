import React, { useEffect, useState } from 'react';
import Graph from 'react-graph-vis';
import { prettyPrint } from './parser';
import { verify } from './validator';

// import "./styles.css";
// need to import the vis network css in order to show tooltip
// import "./network.css";

let id = 0;
let rootname = '';
let nodes = [];
let edges = [];

const options = {
  edges: {
    smooth: {
      type: 'cubicBezier',
      forceDirection: 'horizontal',
      roundness: 0.4,
    },
  },
  layout: {
    hierarchical: {
      direction: 'LR',
    },
  },
  physics: true,
};

const events = {
  select: function (event) {
    var { nodes, edges } = event;
  },
};

function setup(obj, level) {
  let myId = id++;
  let label = '';
  if (obj['courses']) {
    label = 'Condition';
    if (obj['grade'] !== '') {
      let nextid = setup(obj['courses'], level + 1);
      edges.push({ from: myId, to: nextid, label: obj['grade'] });
    } else {
      // remove unnecessary group nodes
      let nid = setup(obj['courses'], level);
      return nid;
    }
    if (level == 0) label = rootname;
    nodes.push({ id: myId, label: label, level: level });
    return myId;
  } else if (obj['and'] || obj['or']) {
    label = 'Node';
    let andOr = obj['and'] ? obj['and'] : obj['or'];
    for (let next of andOr) {
      let nextid = setup(next, level + 1);
      edges.push({ from: myId, to: nextid, label: obj['and'] ? 'and' : 'or' });
    }
    if (level == 0) label = rootname;
    nodes.push({ id: myId, label: label, level: level });
    return myId;
  } else if (obj['course']) {
    label = obj['course'];
    if (obj['grade'] && obj['grade'] !== '') label += ' ' + obj['grade'];
    if (level == 0) {
      nodes.push({ id: 0, label: rootname, level: 0 });
      nodes.push({ id: 1, label: label, level: 1 });
      edges.push({ from: 0, to: 1 });
    } else {
      nodes.push({ id: myId, label: label, level: level });
      return myId;
    }
  }
}

/**
 * Prerequisite visualization component using vis.js
 */
export default function Prereq({ name, prereqs }) {
  const [graph, setGraph] = useState(null);

  useEffect(() => {
    if (prereqs) {
      let obj = prettyPrint(prereqs);
      console.log(verify(prereqs, ['CS 2340', 'CS 2305', 'CS 3333', 'CS 2336']));
      nodes = [];
      edges = [];
      rootname = name;
      setup(obj, 0);
      setGraph({
        nodes: nodes,
        edges: edges,
      });
    }
  }, []);

  return (
    <>
      {!!graph ? (
        <Graph
          graph={graph}
          options={options}
          events={events}
          style={{ height: '640px' }}
          getNetwork={(network) => {
            //  if you want access to vis.js network api you can set the state in a parent component using this property
          }}
        />
      ) : (
        <>Graph not initialized</>
      )}
    </>
  );
}
