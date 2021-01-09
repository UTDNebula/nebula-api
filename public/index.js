(function () {
	addModal = true;
	setupEvents();
})();

async function setupEvents() {
	const fields = ["course", "title", "description", "prerequisites", "prerequisiteOrCorerequisites", "corequisites"];
	for (let field of fields) {
		let fg = htmlToElement(`
		<div class="form-group ">
			<label for="add-${field}">Course ${field[0].toUpperCase() + field.slice(1)}</label>
			<input type="text" class="form-control" id="add-${field}">
        </div>
		`);
		document.getElementsByClassName("form-row")[0].appendChild(fg);
	}
	for (let field of fields) {
		let fg = htmlToElement(`
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
			let obj = {};
			for (let field of fields) {
				obj[field] = document.getElementById(`add-${field}`).value;
			}
			let uri = "/courses";
			const method = "POST";
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
			let id = parseInt(document.getElementById("hidden-id").textContent);
			let obj = {id: id};
			for (let field of fields) {
				obj[field] = document.getElementById(`edit-${field}`).value;
			}
			let uri = "/courses/" + id;
			const method = "PUT";
			fetch(uri, {
				method: method,
				headers: { 'Content-type': 'application/json' },
				body: JSON.stringify(obj)
			}).then(res => res.json()).then(async (res) => {
				console.log(res);
				search();
			})
		} else if (event.target.classList.contains("edit-course")) {
			document.querySelector("#edit-button").click();
			let mapping = maps[parseInt(event.target.getAttribute("value"))];
			document.getElementById("hidden-id").textContent = mapping["id"];
			for (let field of fields) {
				document.getElementById(`edit-${field}`).value = mapping[field] ? mapping[field] : "";
			}
			console.log(mapping["prerequisites"]);
		} else if (event.target.classList.contains("prereq-course")) {
			let mapping = maps[parseInt(event.target.getAttribute("value"))];
			if(mapping["prerequisites"] !== "") {
				let res = prettyPrint(mapping["prerequisites"]);
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
			for (let course of res)
				addCard(course);
		} else {
			alert("Course not found!");
		}
	})
}

let container = document.getElementById("cards-container");
let name_input = document.getElementById("search");
let allCards = [];
let maps = {};

function addCard(course) {
	maps[course.id] = course;
	let prerequisites = "";
	let corequisites = "";
	let prerequisiteOrCorequisites = "";
	if(course.prerequisites !== "") prerequisites = `<h6 class="card-subtitle mb-3 text-muted">Prerequisites: ${course.prerequisites}</h6>`;
	if(course.corequisites !== "") corequisites = `<h6 class="card-subtitle mb-3 text-muted">Corequisites: ${course.corequisites}</h6>`;
	if(course.prerequisiteOrCorequisites) prerequisiteOrCorequisites = `<h6 class="card-subtitle mb-3 text-muted">Prerequisite or Corequisite: ${course.prerequisiteOrCorequisites}</h6>`;
	let card = htmlToElement(`
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
	const template = document.createElement('template');
	html = html.trim();
	template.innerHTML = html;
	return template.content.firstChild;
}
