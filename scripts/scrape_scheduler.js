const base = `https://utdallas.collegescheduler.com/api/terms/2021%20Spring/subjects`;
async function scrape_college_scheduler() {
    let i = 0;
    let allCourses = [];
    fetch(base).then(res => res.json()).then(async (res) => {
        for(let sub of res) {
            await fetch(`${base}/${sub.id}/courses`).then(a=>a.json()).then(async (cs) => {
                for(let course of cs) {
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
}
