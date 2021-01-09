const createToken = chevrotain.createToken;
const tokenMatcher = chevrotain.tokenMatcher;
const Lexer = chevrotain.Lexer;
const EmbeddedActionsParser = chevrotain.EmbeddedActionsParser;

const And = createToken({ name: "And", pattern: /and/ });
const Or = createToken({ name: "Or", pattern: /or/ });

const LParen = createToken({ name: "LParen", pattern: /\(/ });
const RParen = createToken({ name: "RParen", pattern: /\)/ });
const Course = createToken({ name: "Course", pattern: /([A-Z]+ [0-9][A-Z0-9][0-9]+)/ });
const RandomRequest = createToken({ name: "RandomRequest", pattern: /((?! and | or |\(|\)).)+/ });
const Grade = createToken({ name: "Grade", pattern: /with a (?:minimum )*grade (of )*[ABC]-*\+*( or (?:higher|better))*/ });

const WhiteSpace = createToken({
    name: "WhiteSpace",
    pattern: /\s+/,
    group: Lexer.SKIPPED
});

const Comma = createToken({ name: "Comma", pattern: /,/, group: Lexer.SKIPPED })

const IrrelevantWord = createToken({
    name: 'IrrelevantWord',
    pattern: /[^\s()]+/,
    group: Lexer.SKIPPED
});

const allTokens = [WhiteSpace, Comma,
    And, Or, Course, Grade, RandomRequest, LParen, RParen, //, IrrelevantWord
];
const CalculatorLexer = new Lexer(allTokens);

function generateAnd(children) {
    if (children.length > 1)
        return { "and": children }
    else return children
}

function generateOr(children) {
    if (children.length > 1)
        return { "or": children }
    else return children
}

var check = false;
var taken_courses = [];

class Calculator extends EmbeddedActionsParser {
    constructor() {
        super(allTokens);
        const $ = this;

        $.RULE("expression", () => {
            var res = $.SUBRULE($.andExpression);
            if (check)
                return res ? "Good" : "Bad";
            else
                return res;
        });

        $.RULE("andExpression", () => {
            let value = [];
            if (check) value = true;
            // parsing part
            if (!check)
                value.push($.SUBRULE($.orExpression));
            else value = $.SUBRULE($.orExpression);
            $.MANY(() => {
                // consuming 'AdditionOperator' will consume
                // either Plus or Minus as they are subclasses of AdditionOperator
                $.CONSUME(And);
                //  the index "2" in SUBRULE2 is needed to identify the unique
                // position in the grammar during runtime
                if (!check)
                    value.push($.SUBRULE2($.orExpression));
                else value &= $.SUBRULE2($.orExpression);
            });

            if (check)
                return value;
            else
                return value.length === 1 ? value[0] : generateAnd(value)
        });

        $.RULE("orExpression", () => {
            let value = [];
            if (check) value = true;

            // parsing part
            if (!check)
                value.push($.SUBRULE($.atomicBooleanExpression));
            else
                value = $.SUBRULE($.atomicBooleanExpression);
            $.MANY(() => {
                $.CONSUME(Or);
                let val = $.SUBRULE2($.atomicBooleanExpression);
                if (!check)
                    value.push(val);
                else value |= val;
            });

            if (check)
                return value;
            else
                return value.length === 1 ? value[0] : generateOr(value)
        });

        $.RULE("atomicBooleanExpression", () => $.OR([
            // parenthesisExpression has the highest precedence and thus it
            // appears in the "lowest" leaf in the expression ParseTree.
            { ALT: () => $.SUBRULE($.parenthesisExpression) },
            { ALT: () => $.SUBRULE($.courseExpression) },
            {
                ALT: () => {
                    var rand = $.CONSUME(RandomRequest).image
                    if(check) return true;
                    else return { "course": rand, "type": "special" }
                }
            }
        ]));

        $.RULE("courseExpression", () => {
            let course;
            let grade = -1;

            course = $.CONSUME(Course);
            $.OPTION(() => {
                grade = $.CONSUME(Grade)
            })
            if (check) {
                var courseNum = course.image.match(course.tokenType.PATTERN);
                // check if this course is good or not
                // depending on:
                // 1. taken or not
                // 2. grade meets minimum
                if (courseNum != null) {
                    // do all checks here
                    for (var taken_course of taken_courses) {
                        if (courseNum[1].includes(taken_course)) {
                            console.log(courseNum[1] + " satisfied");
                            return true;
                        }
                    }
                    return false;
                }
                return false;
            } else {
                if (grade == -1)
                    return { "course": course.image, "grade": "" };
                return { "course": course.image, "grade": grade.image };
            }
        })

        $.RULE("parenthesisExpression", () => {
            let expValue;
            let grade = "";

            $.CONSUME(LParen);
            expValue = $.SUBRULE($.andExpression);
            $.CONSUME(RParen);

            $.OPTION(() => {
                grade = $.CONSUME(Grade).image;
            })
            var res = { "courses": expValue, "grade": grade }
            if(check) return expValue;
            else return res
        });

        // very important to call this after all the rules have been defined.
        // otherwise the parser may not work correctly as it will lack information
        // derived during the self analysis phase.
        this.performSelfAnalysis();
    }
}

const parser = new Calculator()

function parseInput(text) {
    const lexingResult = CalculatorLexer.tokenize(text)
    // "input" is a setter which will reset the parser's state.
    parser.input = lexingResult.tokens
    var res = parser.expression()

    if (parser.errors.length > 0) {
        throw new Error("sad sad panda, Parsing errors detected")
    }
    return res;
}

function prettyPrint(text) {
    console.log(text);
    var res = parseInput(text);
    console.log(JSON.stringify(res, null, 2));
    return res;
}

function verify(text, courses) {
    check = true;
    taken_courses = courses;
    var res = parseInput(text);
    check = false;
    taken_courses = [];
    return res;
}
