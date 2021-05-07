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
 * Parses a prerequisite string into a tree structure.
 */
class PrereqParser extends chevrotain.EmbeddedActionsParser {
  // combine array of expressions into object with "and" key
  // used to generate intuitive prerequisite graph
  generateAnd(children) {
    if (children.length > 1) return { and: children };
    else return children;
  }

  // combine array of expressions into object with "or" key
  generateOr(children) {
    if (children.length > 1) return { or: children };
    else return children;
  }

  // constructs and defines expressions for parsing
  constructor() {
    super(allTokens);
    const $ = this;

    $.RULE('expression', () => {
      let res = $.SUBRULE($.andExpression);
      return res;
    });

    $.RULE('andExpression', () => {
      let value: any = []; // TODO
      // parsing part
      value.push($.SUBRULE($.orExpression));
      $.MANY(() => {
        // consuming 'AdditionOperator' will consume
        // either Plus or Minus as they are subclasses of AdditionOperator
        $.CONSUME(And);
        //  the index "2" in SUBRULE2 is needed to identify the unique
        // position in the grammar during runtime
        value.push($.SUBRULE2($.orExpression));
      });

      return value.length === 1 ? value[0] : this.generateAnd(value);
    });

    $.RULE('orExpression', () => {
      let value: any = [];

      // parsing part
      value.push($.SUBRULE($.atomicBooleanExpression));
      $.MANY(() => {
        $.CONSUME(Or);
        let val = $.SUBRULE2($.atomicBooleanExpression);
        value.push(val);
      });

      return value.length === 1 ? value[0] : this.generateOr(value);
    });

    $.RULE('atomicBooleanExpression', () =>
      $.OR([
        // parenthesisExpression has the highest precedence and thus it
        // appears in the "lowest" leaf in the expression ParseTree.
        { ALT: () => $.SUBRULE($.parenthesisExpression) },
        { ALT: () => $.SUBRULE($.courseExpression) },
        {
          ALT: () => {
            let rand = $.CONSUME(RandomRequest).image;
            return { course: rand, type: 'special' };
          },
        },
      ]),
    );

    $.RULE('courseExpression', () => {
      let course;
      let grade: any = -1;

      course = $.CONSUME(Course);
      $.OPTION(() => {
        grade = $.CONSUME(Grade);
      });

      if (grade == -1) return { course: course.image, grade: '' };
      return { course: course.image, grade: grade.image };
    });

    $.RULE('parenthesisExpression', () => {
      let expValue;
      let grade = '';

      $.CONSUME(LParen);
      expValue = $.SUBRULE($.andExpression);
      $.CONSUME(RParen);

      $.OPTION(() => {
        grade = $.CONSUME(Grade).image;
      });
      let res = { courses: expValue, grade: grade };
      return res;
    });

    // very important to call this after all the rules have been defined.
    // otherwise the parser may not work correctly as it will lack information
    // derived during the self analysis phase.
    this.performSelfAnalysis();
  }
}

// create new parser instance
const parser = new PrereqParser();

/**
 * Parses a prerequisite string into JSON structure
 * @param text prerequisite string
 * @returns parsed structure
 */
export function prettyPrint(text) {
  const lexingResult = CalculatorLexer.tokenize(text);
  parser.input = lexingResult.tokens;
  let res = parser.expression();

  if (parser.errors.length > 0) {
    throw new Error('sad sad panda, Parsing errors detected');
  }
  console.log(JSON.stringify(res, null, 2));
  return res;
}
