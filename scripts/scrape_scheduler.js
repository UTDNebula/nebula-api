var i = 0;
// fetch(`https://utdallas.collegescheduler.com/api/terms/2021%20Spring/subjects/CS/courses/4V95`)
var base = `https://utdallas.collegescheduler.com/api/terms/2021%20Spring/subjects`;
var allCourses = [];
fetch(base).then(res => res.json()).then(async (res) => {
    for(var sub of res) {
        await fetch(`${base}/${sub.id}/courses`).then(a=>a.json()).then(async (cs) => {
            for(var course of cs) {
                await fetch(`${base}/${sub.id}/courses/${course.number}`).then(a=>a.json()).then(async (cc) => {
                    cc["id"] = i++;
                    cc["course"] = cc["subjectId"] + " " + cc["number"];
                    allCourses.push(cc);
                });
            }
        })
        console.log(`done with ${sub.id}`);
    }
})
