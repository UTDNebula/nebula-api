(function () {
	setupEvents();
	/*fetch("/courses").then(res => res.json()).then(res => {
		console.log(res);
		allCards = res;
		for(var i = 0; i < 10; i++) {
			addCard(allCards[i]);
		}
	})*/
})();

function setupEvents() {
	document.body.addEventListener('click', function (event) {
		if (event.target.classList.contains('delete-course')) {
			fetch(`/courses/${event.target.getAttribute("value")}`, {
				method: 'DELETE'
			}).then(res => res.json()).then(res => {
				console.log(res);
			})
		} else if (event.target.id === 'search-name') {
			fetch(`/courses/name/${name_input.value}`).then(res => res.json()).then(res => {
				if (res !== null) {
					container.innerHTML = "";
					addCard(res);
				} else {
					alert("Course not found!");
				}
			})
		} else if (event.target.id === 'add-submit') {
			console.log(event.target);
			var number = document.getElementById("add-number").value;
			var name = document.getElementById("add-name").value;
			var description = document.getElementById("add-description").value;
			fetch(`/courses`, {
				method: 'POST',
				headers: { 'Content-type': 'application/json' },
				body: JSON.stringify({
					"id": 5000, "course": number,
					"name": name, "description": description
				})
			}).then(res => res.json()).then(res => {
				console.log(res);
			})
		}
	});
}

var container = document.getElementById("cards-container");
var name_input = document.getElementById("search");
var allCards = [];

function addCard(course) {
	var card = htmlToElement(`
	<div class="card">
		<div class="card-body">
			<h5 class="card-title">${course.course}</h5>
			<h6 class="card-subtitle mb-2 text-muted">${course.name}</h6>
			<p class="card-text">${course.description}</p>
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