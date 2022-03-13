import mongoose from 'mongoose'

module schemas {

    export type RequirementType = "course" | "section" | "exam" | "major" | "minor" | "gpa" | "consent" | "collection" | "hours" | "other"

    abstract class Requirement {
        readonly "type": RequirementType;
        constructor(type: RequirementType) { this.type = type }
    }

    export class CourseRequirement extends Requirement {
        "class_reference": mongoose.Types.ObjectId;
        "minimum_grade": string;
        constructor() { super("course") }
    }

    export class SectionRequirement extends Requirement {
        "section_reference": mongoose.Types.ObjectId;
        constructor() { super("section") }
    }

    export class ExamRequirement extends Requirement {
        "exam_reference": mongoose.Types.ObjectId;
        "minimum_score": number;
        constructor() { super("exam") }
    }

    export class MajorRequirement extends Requirement {
        "major": string;
        constructor() { super("major") }
    }

    export class MinorRequirement extends Requirement {
        "minor": string;
        constructor() { super("minor") }
    }

    export class GPARequirement extends Requirement {
        "minimum": number;
        "subset": string;
        constructor() { super("gpa") }
    }

    export class ConsentRequirement extends Requirement {
        "granter": string;
        constructor() { super("consent") }
    }

    export class HoursRequirement extends Requirement {
        "required": number;
        "options": Array<CourseRequirement> = [];
        constructor() { super("hours") }
    }

    export class OtherRequirement extends Requirement {
        "description": string;
        "condition": string;
        constructor() { super("other") }
    }

    export class CollectionRequirement extends Requirement {
        "name": string;
        "required": number;
        "options": Array<Requirement> = [];
        constructor() { super("collection") }
    }

    type MongoStored = {
        _id: mongoose.Types.ObjectId;
    }

    export type AcademicSession = {
        name: string,
        start_date: string,
        end_date: string
    }

    export type Course = MongoStored & {
        course_number: string,
        subject_prefix: string,
        title: string,
        description: string,
        school: string,
        credit_hours: string,
        class_level: string,
        activity_type: string,
        grading: string,
        internal_course_number: string,
        prerequisites: CollectionRequirement,
        corequisites: CollectionRequirement,
        sections: Array<mongoose.Types.ObjectId>,
        lecture_contact_hours: string,
        laboratory_contact_hours: string,
        offering_frequency: string,
        attributes: Object
    }

    export type Section = MongoStored & {
        section_number: string,
        course_reference: mongoose.Types.ObjectId,
        section_corequisites: CollectionRequirement,
        academic_session: AcademicSession,
        professors: Array<mongoose.Types.ObjectId>,
        teaching_assistants: Array<Assistant>,
        internal_class_number: string,
        instruction_mode: string,
        meetings: Array<Meeting>,
        core_flags: Array<string>,
        syllabus_uri: string,
        grade_distribution: Array<number>,
        attributes: Object
    }

    export type DegreeSubtype = "major" | "minor" | "concentration" | "prescribed double major";

    export type Degree = MongoStored & {
        subtype: DegreeSubtype,
        school: string,
        name: string,
        year: string,
        abbreviation: string,
        minimum_credit_hours: number,
        requirements: CollectionRequirement
    }

    export type Location = {
        building: string,
        room: string,
        map_uri: string
    }

    export type ModalityType = "pending" | "traditional" | "hybrid" | "flexible" | "remote" | "online";

    export type Meeting = {
        start_date: string,
        end_date: string,
        meeting_days: Array<string>,
        start_time: string,
        end_time: string,
        modality: ModalityType,
        location: Location
    }

    export type Professor = MongoStored & {
        first_name: string,
        last_name: string,
        title: string,
        email: string,
        phone_number: string,
        office: Location,
        profile_uri: string,
        image_uri: string,
        office_hours: Array<Meeting>,
        section_references: Array<mongoose.Types.ObjectId>
    }

    export type Assistant = {
        first_name: string,
        last_name: string,
        role: string,
        email: string
    }

};

export default schemas;