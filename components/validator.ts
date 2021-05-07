const chevrotain = require('chevrotain');

const createToken = chevrotain.createToken;
const tokenMatcher = chevrotain.tokenMatcher;
const Lexer = chevrotain.Lexer;
const EmbeddedActionsParser = chevrotain.EmbeddedActionsParser;

// define the base cases, which are the components that make up expressions
const And = createToken({ name: 'And', pattern: /and/ });
const Or = createToken({ name: 'Or', pattern: /or/ });

const LParen = createToken({ name: 'LParen', pattern: /\(/ });
const RParen = createToken({ name: 'RParen', pattern: /\)/ });

const Course = createToken({ name: 'Course', pattern: /([A-Z]+ [0-9][A-Z0-9][0-9]+)/ });
const RandomRequest = createToken({ name: 'RandomRequest', pattern: /((?! and | or |\(|\)).)+/ });
const Grade = createToken({
  name: 'Grade',
  pattern: /with a (?:minimum )*grade (of )*[ABC]-*\+*( or (?:higher|better))*/,
});

const WhiteSpace = createToken({
  name: 'WhiteSpace',
  pattern: /\s+/,
  group: Lexer.SKIPPED,
});

const Comma = createToken({ name: 'Comma', pattern: /,/, group: Lexer.SKIPPED });

// define all tokens with their order of precedence
const allTokens = [
  WhiteSpace,
  Comma,
  And,
  Or,
  Course,
  Grade,
  RandomRequest,
  LParen,
  RParen, //, IrrelevantWord
];

// create lexer instance
const CalculatorLexer = new Lexer(allTokens);

let taken_courses = [];

class Calculator extends EmbeddedActionsParser {
  constructor() {
    super(allTokens);
    const $ = this;

    $.RULE('expression', () => {
      return $.SUBRULE($.andExpression);
    });

    $.RULE('andExpression', () => {
      let value: boolean = $.SUBRULE($.orExpression);
      $.MANY(() => {
        $.CONSUME(And);
        let res = $.SUBRULE2($.orExpression);
        if (!res) 
            value = false;
      });

      return value;
    });

    $.RULE('orExpression', () => {
      let value: boolean = $.SUBRULE($.atomicBooleanExpression);
      $.MANY(() => {
        $.CONSUME(Or);
        let val = $.SUBRULE2($.atomicBooleanExpression);
        if(val)
            value = true;
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
        for (let taken_course of taken_courses) {
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

const parser = new Calculator();

export function parseInput(text, courses) {
  taken_courses = courses;
  const lexingResult = CalculatorLexer.tokenize(text);
  // "input" is a setter which will reset the parser's state.
  parser.input = lexingResult.tokens;
  let res = parser.expression();
  taken_courses = [];

  if (parser.errors.length > 0) {
    throw new Error('sad sad panda, Parsing errors detected');
  }
  return res;
}

export function verify(text, courses) {
  let res = parseInput(text, courses);
  return res;
}
