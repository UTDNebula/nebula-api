const { getAll, post, findById, findByName, deleteById, patch, consoleView, logout, auth, authCheck } = require('./crud.js');

// START server config
const express = require('express');
const app = express();
app.engine('html', require('ejs').renderFile);
app.set('view engine', 'html');
app.set('views', __dirname + '/public');

const bodyParser = require('body-parser');
require('dotenv').config();
// END server config

app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());

// START server routing
const port = process.env.PORT || 5000;

app.listen(port, function () {
    console.log(`Server started on ${port}`)
});

app.get('/console', authCheck, consoleView);

app.get('/courses', authCheck, getAll);

app.post('/courses', authCheck, post);

app.get('/courses/id/:id', authCheck, findById);

app.get('/courses/name/:name', authCheck, findByName);

app.delete('/courses/:id', authCheck, deleteById);

app.put('/courses/:id', authCheck, patch);

app.get('/logout', logout);

app.post('/auth', auth)

app.get('/', (req, res) => {
    res.render('login', { ...process.env });
})

app.use(express.static(__dirname + '/public'))
