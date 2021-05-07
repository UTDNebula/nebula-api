// example representation of degree plans

enum ChooseType {
  hours, // choose x hours from these groups
  count, // choose x groups from these groups
}

/**
 * Course group type
 */
export type courseGroup = {
  name?: string; // optional name
  description?: string; // optional description
  choose: ChooseType; // choose by hours of by count
  pick: number; // how many hours/count to pick
  children: Array<courseGroup | string>; // children can be courseGroups or strings (individual course's title)
};

/**
 * Degree plan type
 */
export type degreePlanType = {
  major: string; // ex. ce, cs, bmen...
  degree: string; // ex. bs, ms, phd...
  creditHours: number; // ex. 126 for ce
  children: Array<courseGroup>; // ex. [Core, Major, Electives] each represented as a courseGroup type
};

// course group constants
let core040: Array<string> = ['PHIL 1111', 'PHIL 2222'];
let CSFreeElectives: Array<string> = ['CS 3333', 'CS 9999'];

let CS: degreePlanType = {
  major: 'CS',
  degree: 'BS',
  creditHours: 124,
  children: [
    {
      name: 'Core Curriculum Requirements',
      choose: ChooseType.count,
      pick: 9,
      children: [
        {
          name: 'Communication',
          choose: ChooseType.hours,
          pick: 6,
          children: ['RHET 1302', 'ECS 3390'],
        },
        {
          name: 'Language, Philosophy and Culture',
          choose: ChooseType.hours,
          pick: 3,
          children: core040,
        },
        // other core requirements
        {
          name: 'Component Area Option',
          choose: ChooseType.count,
          pick: 3,
          children: [
            'MATH 2419',
            'PHYS 2125',
            {
              choose: ChooseType.count,
              pick: 1,
              children: ['MATH 2413, MATH 2417'],
            },
          ],
        },
      ],
    },
    {
      name: 'Major Requirements',
      choose: ChooseType.count,
      pick: 3,
      children: [
        {
          name: 'Major Prep',
          choose: ChooseType.hours,
          pick: 24,
          children: ['ECS 1100', 'CS 2340'], // and more
        },
        {
          name: 'Major Core',
          choose: ChooseType.hours,
          pick: 39,
          children: ['CS 3162', 'CS 3305'], // and more
        },
        {
          name: 'Major Electives',
          choose: ChooseType.hours,
          pick: 9,
          children: ['CS 4314', 'CS 4315'], // and more
        },
      ],
    },
    {
      name: 'Elective Requirements',
      choose: ChooseType.hours,
      pick: 10,
      children: CSFreeElectives,
    },
  ],
};
