import * as chevrotain from 'chevrotain';

const createToken = chevrotain.createToken;
const Lexer = chevrotain.Lexer;

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
const allTokens = [WhiteSpace, Comma, And, Or, Course, Grade, RandomRequest, LParen, RParen];

// create lexer instance
const CalculatorLexer = new chevrotain.Lexer(allTokens);

/**
 * Exports token and lexer constants for the parser/validator
 */
export { allTokens, And, Or, Course, Grade, RandomRequest, LParen, RParen, CalculatorLexer };
