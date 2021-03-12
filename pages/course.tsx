// TODO: better course type
export type courseType = {
  titleLong: string;
  id: string;
  description: string;
};

type courseProps = {
  course: courseType;
  editCourse: (course: courseType) => any;
  deleteCourse: (id: string) => any;
};

const Course: React.FunctionComponent<courseProps> = (props) => {
  const course = props.course;
  const editCourse = props.editCourse;
  const deleteCourse = props.deleteCourse;

  return (
    <div className="p-4 ring-2 ring-blue-400 rounded-md shadow-md">
      {course ? (
        <>
          <h1 className="text-2xl text-light">
            {course.titleLong} ({course.id})
          </h1>
          <p className="text-md my-4">{course.description}</p>
          <div className="flex">
            <button
              className="mr-2 px-4 font-light rounded-lg bg-blue-300 hover:bg-blue-500"
              onClick={() => {
                editCourse(course);
              }}
            >
              Edit
            </button>
            <button
              className="p-2 font-light rounded-lg bg-red-300 hover:bg-red-500"
              onClick={() => {
                deleteCourse(course.id);
              }}
            >
              Delete
            </button>
          </div>
        </>
      ) : (
        <p>No information available.</p>
      )}
    </div>
  );
};

export default Course;
