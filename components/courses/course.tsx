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

          <div className="flex gap-4">
            <button
              className="btn1"
              onClick={() => {
                editCourse(course);
              }}
            >
              Edit
            </button>
            <button
              className="btn2"
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
