(function () {
	addModal = true;
	setupEvents();
	// fetch("/courses").then(res => res.json()).then(res => {
	// 	console.log(res);
	// 	var total = 0;
	// 	var error = 0;
	// 	for(var i = 0; i < res.length; i++) {
	// 		var checkString = res[i]["prerequisites"];
	// 		if(checkString && checkString !== "") {
	// 			try {
	// 				parseInput(checkString);
	// 			} catch(err) {
	// 				console.log(res[i]["course"] + ": " + checkString);
	// 				error++;
	// 			}
	// 			total++;
	// 		}
	// 	}
	// 	console.log(`total: ${total}, errors: ${error}`);
	// })
})();

function setupEvents() {
	var fields = ["course", "title", "description", "prerequisites", "prerequisiteOrCorerequisites", "corequisites"];
	for (var field of fields) {
		var fg = htmlToElement(`
		<div class="form-group ">
			<label for="add-${field}">Course ${field[0].toUpperCase() + field.slice(1)}</label>
			<input type="text" class="form-control" id="add-${field}">
        </div>
		`);
		document.getElementsByClassName("form-row")[0].appendChild(fg);
	}
	for (var field of fields) {
		var fg = htmlToElement(`
		<div class="form-group ">
			<label for="edit-${field}">Course ${field[0].toUpperCase() + field.slice(1)}</label>
			<input type="text" class="form-control" id="edit-${field}">
		</div>
		`);
		document.getElementsByClassName("form-row")[1].appendChild(fg);
	}
	document.body.addEventListener('click', function (event) {
		if (event.target.classList.contains('delete-course')) {
			fetch(`/courses/${event.target.getAttribute("value")}`, {
				method: 'DELETE'
			}).then(res => res.json()).then(res => {
				console.log(res);
				search();
			})
		} else if (event.target.id === 'search-name') {
			event.preventDefault();
			search();
		} else if (event.target.id === 'add-submit') {
			console.log(event.target);
			var obj = {};
			for (var field of fields) {
				obj[field] = document.getElementById(`add-${field}`).value;
			}
			var uri = "/courses";
			var method = "POST";
			if (!addModal) {
				uri = "/courses/" + event.target.getAttribute("value");
				method = "PUT";
			}
			fetch(uri, {
				method: method,
				headers: { 'Content-type': 'application/json' },
				body: JSON.stringify(obj)
			}).then(res => res.json()).then(res => {
				console.log(res);
			})
		} else if (event.target.id === 'edit-submit') {
			console.log(event.target);
			var id = parseInt(document.getElementById("hidden-id").textContent);
			var obj = {id: id};
			for (var field of fields) {
				obj[field] = document.getElementById(`edit-${field}`).value;
			}
			var uri = "/courses/" + id;
			var method = "PUT";
			fetch(uri, {
				method: method,
				headers: { 'Content-type': 'application/json' },
				body: JSON.stringify(obj)
			}).then(res => res.json()).then(res => {
				console.log(res);
			})
		} else if (event.target.classList.contains("edit-course")) {
			document.querySelector("#edit-button").click();
			var mapping = maps[parseInt(event.target.getAttribute("value"))];
			document.getElementById("hidden-id").textContent = mapping["id"];
			for (var field of fields) {
				document.getElementById(`edit-${field}`).value = mapping[field] ? mapping[field] : "";
			}
			console.log(mapping["prerequisites"]);
		} else if (event.target.classList.contains("prereq-course")) {
			var mapping = maps[parseInt(event.target.getAttribute("value"))];
			if(mapping["prerequisites"] !== "") {
				var res = prettyPrint(mapping["prerequisites"]);
				document.querySelector("#graph-button").click();
				drawGraph(res, mapping["course"]);
			} else {
				alert("This course has no prerequisites.");
			}
		}
	});
}

function search() {
	fetch(`/courses/name/${name_input.value}`).then(res => res.json()).then(res => {
		if (res !== null) {
			container.innerHTML = "";
			for (var course of res)
				addCard(course);
		} else {
			alert("Course not found!");
		}
	})
}

var container = document.getElementById("cards-container");
var name_input = document.getElementById("search");
var allCards = [];
var maps = {};

function addCard(course) {
	maps[course.id] = course;
	var prerequisites = "";
	var corequisites = "";
	var prerequisiteOrCorequisites = "";
	if(course.prerequisites !== "") prerequisites = `<h6 class="card-subtitle mb-3 text-muted">Prerequisites: ${course.prerequisites}</h6>`;
	if(course.corequisites !== "") corequisites = `<h6 class="card-subtitle mb-3 text-muted">Corequisites: ${course.corequisites}</h6>`;
	if(course.prerequisiteOrCorequisites) prerequisiteOrCorequisites = `<h6 class="card-subtitle mb-3 text-muted">Prerequisite or Corequisite: ${course.prerequisiteOrCorequisites}</h6>`;
	var card = htmlToElement(`
	<div class="card">
		<div class="card-body">
			<h5 class="card-title">${course.course}</h5>
			<h6 class="card-subtitle mb-2 text-muted">${course.title}</h6>
			<p class="card-text">${course.description}</p>
			${prerequisites}
			${corequisites}
			${prerequisiteOrCorequisites}
			<button value="${course.id}" class="edit-course btn btn-primary">Edit</button>
			<button value="${course.id}" class="prereq-course btn btn-primary">Show Prereq</button>
			<button value="${course.id}" class="delete-course btn btn-primary">Delete</button>
		</div>
	</div>`);
	container.appendChild(card);
}

function htmlToElement(html) {
	var template = document.createElement('template');
	html = html.trim(); // Never return a text node of whitespace as the result
	template.innerHTML = html;
	return template.content.firstChild;
}