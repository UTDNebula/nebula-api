// START Firestore init
const admin = require('firebase-admin')
const serviceAccount = require('./firestore_key.json')
admin.initializeApp({
    credential: admin.credential.cert(serviceAccount)
})
const db = admin.firestore();
const increment = admin.firestore.FieldValue.increment(1);
const decrement = admin.firestore.FieldValue.increment(-1);
// END Firestore init

// START server config
const express = require("express")
const app = express()
const bodyParser = require("body-parser")
require('dotenv').config()
// END server config

// DEBUG init database data
const course_info = require("./data/scheduler_prereq.json")

// START server routing
const port = process.env.PORT || 3000

app.use(bodyParser.urlencoded({ extended: false }))
app.use(bodyParser.json())
app.use(express.static(__dirname + "/public"))

app.get("/", (req, res) => {
    res.sendFile(__dirname + "public/index.html")
})

app.get("/courses", async (req, res) => {
    const courses = db.collection("courses")
    const allCourses = await courses.get()
    var result = []
    allCourses.forEach(c => {
        result.push(c.data());
    })
    res.json(result)
})

app.post("/courses", async (req, res) => {
    const course = req.body;
    console.log(course);

    // get counter to set new id
    const courseCol = db.collection("courses");
    const counter = courseCol.doc("_counter");
    const count = await counter.get();
    console.log(count.data());
    var newId = count.data()["count"];
    course["id"] = parseInt(newId);

    await courseCol.doc(course["course"]).set(course);
    await counter.update({ count: increment });

    res.json({ "message": "Course updated and counter updated"});
})

app.get("/courses/id/:id", async (req, res) => {
    const courseId = parseInt(req.params.id);
    const coursesRef = db.collection("courses");
    const idCourses = await coursesRef.where("id", "==", courseId).get();
    if(idCourses.empty) {
        res.json({});
    } else {
        res.json(idCourses.docs[0].data());
    }
})

app.get("/courses/name/:name", async (req, res) => {
    const name = req.params.name.toUpperCase();
    console.log(name);

    const end = name.replace(/.$/, c => String.fromCharCode(c.charCodeAt(0) + 1));
    const courseCol = db.collection("courses");
    const snapshot = await courseCol.where("course", ">=", name)
                                    .where("course",  "<", end)
                                    .get();
    if(snapshot.empty) {
        console.log("not found");
        res.json([]);
    } else {
        console.log("found!");
        var result = [];
        snapshot.forEach(doc => {
            result.push(doc.data());
        })
        res.json(result);
    }
})

app.delete("/courses/:id", async (req, res) => {
    const courseId = parseInt(req.params.id);
    console.log(`Deleting course with id ${courseId}`);
    const result = db.collection("courses").where("id", "==", courseId);
    result.get().then(snapshot => {
        if(snapshot.empty) {
            console.log("not found");
            res.json({"deleted": false});
        } else {
            snapshot.docs[0].ref.delete();
            res.json({"deleted": true});
        }
    })
})

app.put("/courses/:id", async (req, res) => {
    const courseId = parseInt(req.params.id);
    const course = req.body;
    db.collection("courses").where("id", "==", courseId).get()
        .then(snapshot => {
            if(snapshot.empty) {
                res.json({"updated": false})
            } else {
                snapshot.docs[0].ref.update(course);
                res.json({"updated": true})
            }
        })
})

// init database
// app.get("/initDB", async (req, res) => {
//     var batch = db.batch()
//     course_info.forEach((course, index) => {
//         var ref = db.collection("courses").doc(course["course"]);
//         batch.set(ref, course);
//         // Firestore only allows 500 docs every batch write
//         if(index % 400 == 0) {
//             batch.commit()
//             batch = db.batch()
//             console.log(`batch ${index/400} saved.`)
//         }
//     })
//     batch.commit();
//     res.json({"message": "data saved successfully!"})
// })

app.listen(port, () => {
    console.log(`Server started on ${port}`)
})
// END server routing

// read();

// async function read() {
//     const snapshot = await db.collection('users').get();
//     snapshot.forEach((doc) => {
//         console.log(doc.id, '=>', doc.data());
//     });
// }