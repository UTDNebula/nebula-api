export type announcementType = {
  id: string;
  title: string;
  link: string;
  description: string;
};

export type courseType = {
  id: string;
  course: string;
  description: string;
  title: string;
  prerequisites: string;
  corequisites: string;
  hours: string;
  inclass: string;
  outclass: string;
};

// This is a section that contains a selection of courses. 
// Pick courses from "courses" up until the hours for those courses equal "hoursToPick"
export type degreeSectionType = {
  name: string;
  hoursToPick: number; // pick this much out of the total number of hours for the courses
  hoursAvailable: number; // this is just calculation of total hours in courses
  courses: Array<string>; // array of course IDs
};

export type degreeTrackType = {
  sections: Array<degreeSectionType>;
}

export type degreePlanType = {
  major: string; // ex. ce, cs, bmen...
  degree: string; // ex. bs, ms, phd...
  creditHours: number; // ex. 126 for ce
  tracks: Array<degreeTrackType>; // something like CE will only have 1 track in this array
  majorCourses: Array<courseType>; // a list of all major courses that a student in this major/degree can take
  // TODO: concentration: string; ??
};
