package expression.exceptions;

import expression.*;
//import expression.parser.TripleParser;
import expression.parsing.BaseParser;
import expression.parsing.CharSource;
import expression.parsing.StringSource;

public class UseToParse {

    public TripleExpression parse(final String source) {
        return parse(new StringSource(source));
    }

    public static TripleExpression parse(final CharSource source) {
        return new MyExpressionParser(source).startParse();
    }

    private static class MyExpressionParser extends BaseParser {
        public MyExpressionParser(final CharSource source) {
            super(source);
        }

        public Evaluatable startParse() {
            Evaluatable res = parseExpression();
            if (!eof()) {
                throw new ExpressionException("collapsed expressions (more than one without proper connection) -> "
                        + take() + " position: " + getpos());
            }
            return res;
        }

        public Evaluatable parseExpression() {
            return parseExpression(parseTerm());
        }

        public Evaluatable parseExpression(Evaluatable left) {
            boolean wasSpace;
            wasSpace = skipWhitespace(true);
            skipWhitespace();

            if (take('◣') || test('a')) {
                if (test('a')) {
                    expectLine("area");
                    if(test('a')) {
                        take();
                    }
                    skipWhitespace();
                }
                Evaluatable right = parseTerm();
                skipWhitespace();
                return parseExpression(new CheckedArea(left, right));
            } else if (take('▯') || test('p')) {
                if (test('p')) {
                    expectLine("perimeter");
                    skipWhitespace();
                }
                Evaluatable right = parseTerm();
                skipWhitespace();
                return parseExpression(new CheckedPerimeter(left, right));
            } else if (take('+')) {
                Evaluatable right = parseTerm();
                skipWhitespace();
                return parseExpression(new CheckedAdd(left, right));
            } else if (take('-')) {
                Evaluatable right = parseTerm();
                skipWhitespace();
                return parseExpression(new CheckedSubtract(left, right));
            } else {
                return left;
            }
        }

        private Evaluatable parseTerm() {
            return parseTerm(parseOperand());
        }

        private Evaluatable parseTerm(Evaluatable left) {
            skipWhitespace();
            if (take('*')) {
                Evaluatable right = parseOperand();
                skipWhitespace();
                return parseTerm(new CheckedMultiply(left, right));
            } else if (take('/')) {
                Evaluatable right = parseOperand();
                skipWhitespace();
                return parseTerm(new CheckedDivide(left, right));
            } else {
                return left;
            }
        }

        private Evaluatable parseOperand() {
            skipWhitespace();
            if (test('(') || test('{') || test('[')) {
                skipWhitespace();
                char pair_to = take();
                char closing = 0;
                skipWhitespace();
                Evaluatable expr = parseExpression();
                skipWhitespace();
                if (pair_to == '(') {
                    closing = ')';
                } else if (pair_to == '{') {
                    closing = '}';
                } else if (pair_to == '[') {
                    closing = ']';
                } else {
                    throw new GivenLessException("no matching closing, position: " + getpos());
                }
                try {
                    expect(closing);
                } catch (Exception e) {
                    throw new GivenLessException("no matching closing, position: " + getpos());
                }
                return expr;
            } else if (take('√') || test('s')) {
                if (test('s')) {
                    expectLine("sqrt");
                    skipWhitespace();
                }
                return new CheckedSqrt(parseOperand());
            } else if (take('-')) {
                if (between('1', '9')) {
                    return parseConst(true);
                }
                return new CheckedNegate(parseOperand());
            } else if (between('0', '9')) {
                return parseConst(false);
            } else if (between('A', 'Z') || between('a', 'z')) {
                return parseVariable();
            } else {
                int problemPlace = getpos();
                throw new GivenLessException("Expected operand, found: " + (eof() ? "end of input" : ((take()) +
                        " position:" + problemPlace)));
            }
        }

        private Const parseConst(boolean needsMinus) {
            final StringBuilder sb = new StringBuilder();
            if (needsMinus) {
                sb.append('-');
            }
            if (take('0')) {
                sb.append('0');
                if (between('1', '9')) {
                    throw new GivenExtraException("found smth after zero, " + " position: " + getpos());

                }
            } else if (between('1', '9')) {
                takeDigits(sb);
            } else {
                throw new BadElementException("Invalid Const" + " position: " + getpos());
            }
            try {
                return new Const(Integer.parseInt(sb.toString()));
            } catch (NumberFormatException e) {
                throw new BadElementException("not allowed constant, position: " + getpos());
            }

        }

        private Variable parseVariable() {
            char name = take();
            StringBuilder sb = new StringBuilder();
            sb.append(name);
            try {
                switch (name) {
                    case 'a' -> expectLine("area");
                    case 'p' -> expectLine("perimeter");
                    case 's' -> expectLine("sqrt");
                    default -> { ; }
                }
            } catch (Exception e) {
                throw new BadElementException("Incorrect usage of operation (not enough arguments), " + " pos: " + getpos());
            }
            while ((between('A', 'Z') || between('a', 'z'))) {
                name = take();
                sb.append(name);
            }
            if (name == 'x' || name == 'y' || name == 'z') {
                return new Variable(sb.toString());
            } else {
                throw new BadElementException("Neither variable nor operation: " + name + " pos: " + getpos());
            }
        }

        private void takeDigits(final StringBuilder sb) {
            while (between('0', '9')) {
                sb.append(take());
            }
        }

        private void skipWhitespace() {
            while (isNeededToSkip()) {
                take();
            }
        }

        private boolean skipWhitespace(boolean needsInd) {
            char last = 0;
            while (isNeededToSkip()) {
                last = take();
            }
            if (last == ' ') {
                return true;
            }
            return false;
        }

        private void expectLine(String s) {
            for (int i = 0; i < s.length(); i++) {
                char ch = s.charAt(i);
                try {
                    expect(ch);
                } catch (Exception e) {
                    throw new ParsingException(s + "\t" + getpos());
                }
            }
            boolean wasSpace2 = skipWhitespace(true);
            if (!wasSpace2 && !test('-') && between('0', '9') && !test('(')) {
                throw new GivenLessException("no space using " + s + " position: " + getpos());
            }
        }

    }
}
