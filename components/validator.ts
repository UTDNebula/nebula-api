const chevrotain = require('chevrotain');
import {
  allTokens,
  And,
  Or,
  Course,
  Grade,
  RandomRequest,
  LParen,
  RParen,
  CalculatorLexer,
} from './tokens';

/**
 * Validates prerequisite strings from the given list of courses
 */
class PrereqValidator extends chevrotain.EmbeddedActionsParser {
  constructor(taken_courses) {
    super(allTokens);
    const $ = this;

    $.taken_courses = taken_courses;

    $.RULE('expression', () => {
      return $.SUBRULE($.andExpression);
    });

    $.RULE('andExpression', () => {
      let value: boolean = $.SUBRULE($.orExpression);
      $.MANY(() => {
        $.CONSUME(And);
        let res = $.SUBRULE2($.orExpression);
        if (!res) value = false;
      });

      return value;
    });

    $.RULE('orExpression', () => {
      let value: boolean = $.SUBRULE($.atomicBooleanExpression);
      $.MANY(() => {
        $.CONSUME(Or);
        let val = $.SUBRULE2($.atomicBooleanExpression);
        if (val) value = true;
      });

      return value;
    });

    $.RULE('atomicBooleanExpression', () =>
      $.OR([
        // parenthesisExpression has the highest precedence and thus it
        // appears in the "lowest" leaf in the expression ParseTree.
        { ALT: () => $.SUBRULE($.parenthesisExpression) },
        { ALT: () => $.SUBRULE($.courseExpression) },
        {
          ALT: () => {
            $.CONSUME(RandomRequest);
            return true;
          },
        },
      ]),
    );

    $.RULE('courseExpression', () => {
      let course = $.CONSUME(Course);
      $.OPTION(() => {
        $.CONSUME(Grade);
      });
      let courseNum = course.image.match(course.tokenType.PATTERN);
      // check if this course is good or not
      // depending on:
      // 1. taken or not
      // 2. grade meets minimum
      if (courseNum != null) {
        // do all checks here
        for (let taken_course of $.taken_courses) {
          if (courseNum[1].includes(taken_course)) {
            console.log(courseNum[1] + ' satisfied');
            return true;
          }
        }
        return false;
      }
      return false;
    });

    $.RULE('parenthesisExpression', () => {
      $.CONSUME(LParen);
      let expValue = $.SUBRULE($.andExpression);
      $.CONSUME(RParen);

      let grade = '';
      $.OPTION(() => {
        grade = $.CONSUME(Grade).image;
      });
      return expValue;
    });

    // very important to call this after all the rules have been defined.
    // otherwise the parser may not work correctly as it will lack information
    // derived during the self analysis phase.
    this.performSelfAnalysis();
  }
}

/**
 * Checks to see if the given prerequisite string is satisfied by a given list of courses
 * @param text prerequisite string
 * @param courses array of courses taken so far
 * @returns whether the prerequisite has been satisfied
 */
export function verify(text, courses) {
  const parser = new PrereqValidator(courses);
  const lexingResult = CalculatorLexer.tokenize(text);
  parser.input = lexingResult.tokens;
  let res = parser.expression();
  if (parser.errors.length > 0) {
    throw new Error('sad sad panda, Parsing errors detected');
  }
  return res;
}
