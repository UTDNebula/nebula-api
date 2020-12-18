const express = require("express");
const server = express();
const body_parser = require("body-parser");
require('dotenv').config()

// used to init database
// const course_info = require('./data/course_info.json')

const db = require("./db");
const dbName = "data";
const collectionName = "courses";

// initializes database (currently MongoDB Atlas) and defines routes

db.initialize(dbName, collectionName, function (dbCollection) { 

    // runs at server start
    // dbCollection.find().toArray(function (err, result) {
    //     if (err) throw err;
    //     // result has all courses' information
    //     // console.log(result);
    // });

    // POST: adds a new course
    server.post("/courses", (request, response) => {
        const course = request.body;
        console.log(course);
        dbCollection.insertOne(course, (error, result) => { 
            if (error) throw error;
            dbCollection.find().toArray((_error, _result) => {
                if (_error) throw _error;
                response.json(_result);
            });
        });
    });

    // used to load database
    // server.get("/init", (request, response) => {
    //     var obj = [];
    //     for(var name in course_info) {
    //         var newObj = course_info[name];
    //         newObj["course"] = name;
    //         obj.push(newObj);
    //     }
    //     dbCollection.insertMany(obj, (error, result) => { // callback of insertOne
    //         if (error) throw error;
    //         // return updated list
    //         dbCollection.find().toArray((_error, _result) => { // callback of find
    //             if (_error) throw _error;
    //             response.json(_result);
    //         });
    //     });
    // });

    // GET: get a course with id
    server.get("/courses/id/:id", (request, response) => {
        const courseId = parseInt(request.params.id);
        dbCollection.findOne({ id: courseId }, (error, result) => {
            if (error) throw error;
            // return course
            response.json(result);
        });
    });

    // GET: get a course with name
    server.get("/courses/name/:name", (request, response) => {
        const name = request.params.name;
        dbCollection.findOne({ course: name }, (error, result) => {
            if (error) throw error;
            response.json(result);
        });
    });

    // GET: get all courses
    server.get("/courses", (request, response) => {
        dbCollection.find().toArray((error, result) => {
            if (error) throw error;
            response.json(result);
        });
    });

    // PUT: edit course with id
    server.put("/courses/:id", (request, response) => {
        const courseId = request.params.id;
        const course = request.body;
        console.log("Original: ", courseId, "; New: ", course);

        dbCollection.updateOne({ id: courseId }, { $set: course }, (error, result) => {
            if (error) throw error;
            // returns updated list TODO: change to limit range (either frontend or backend)
            dbCollection.find().toArray(function (_error, _result) {
                if (_error) throw _error;
                response.json(_result);
            });
        });
    });

    // DELETE: deletes course with id
    server.delete("/courses/:id", (request, response) => {
        const courseId = request.params.id;
        console.log("Deleting course with id: ", courseId);
    
        dbCollection.deleteOne({ id: courseId }, function(error, result) {
            if (error) throw error;
            // send back entire updated list after successful request
            dbCollection.find().toArray(function(_error, _result) {
                if (_error) throw _error;
                response.json(_result);
            });
        });
    });

}, function (err) {
    throw (err);
});

// start listening on 3000 or host port

const port = process.env.PORT || 3000;

server.use(body_parser.urlencoded({ extended: false }))
server.use(body_parser.json());
server.use(express.static(__dirname + '/public'));


server.get("/", (req, res) => {
    res.sendFile(__dirname + '/index.html')
})

server.listen(port, function () {
    console.log('listening on 3000')
})