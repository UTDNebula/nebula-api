// TODO: better course type
import { courseType } from '../../lib/types/types';

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
    <div className="p-8 bg-blue-200 rounded-2xl shadow-lg">
      {course ? (
        <>
          <h1 className="text-2xl text-light">
            {course.title} ({course.id})
          </h1>
          <p className="text-md my-4">{course.description}</p>

          <div className="my-4">
            <p className="text-md">Prerequisites: {course.prerequisites}</p>
            <p className="text-md">Corequisites: {course.corequisites}</p>
          </div>

          <div className="my-4">
            <p className="text-md">In-class Hours: {course.inclass}</p>
            <p className="text-md">Out-class Hours: {course.outclass}</p>
          </div>

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
