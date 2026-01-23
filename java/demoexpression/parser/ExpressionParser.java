package expression.parser;

import expression.*;
import expression.parsing.BaseParser;
import expression.parsing.CharSource;
import expression.parsing.StringSource;

public class ExpressionParser implements TripleParser {

    public TripleExpression parse(final String source) {
        return parse(new StringSource(source));
    }

    public static TripleExpression parse(final CharSource source) {
        return new ExpressionParser.MyExpressionParser(source).parseExpression();
    }

    private static class MyExpressionParser extends BaseParser {
        public MyExpressionParser(final CharSource source) {
            super(source);
        }

        public Evaluatable parseExpression() {
            return parseExpression(parseTerm());
        }

        public Evaluatable parseExpression(Evaluatable left) {
            skipWhitespace();
            if (take('+')) {
                Evaluatable right = parseTerm();
                skipWhitespace();
                return parseExpression(new Add(left, right));
            } else if (take('-')) {
                Evaluatable right = parseTerm();
                skipWhitespace();
                return parseExpression(new Subtract(left, right));
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
                return parseTerm(new Multiply(left, right));
            } else if (take('/')) {
                Evaluatable right = parseOperand();
                skipWhitespace();
                return parseTerm(new Divide(left, right));
            } else {
                return left;
            }
        }

        private Evaluatable parseOperand() {
            skipWhitespace();
            if (take('(')) {
                skipWhitespace();
                Evaluatable expr = parseExpression();
                skipWhitespace();
                expect(')');
                return expr;
            } else if (take('√')) {
                return new Sqrt(parseOperand());
            } else if (take('-')) {
                if (between('1', '9')) {
                    return parseConst(true);
                }
                return new UnMinus(parseOperand());
            } else if (between('0', '9')) {
                return parseConst(false);
            } else if (between('A', 'Z') || between('a', 'z')) {
                return parseVariable();
            } else {
                throw error("Expected operand, found: " + (eof() ? "end of input" : take()));
            }
        }

        private Const parseConst(boolean needsMinus) {
            final StringBuilder sb = new StringBuilder();
            if (needsMinus) {
                sb.append('-');
            }
            if (take('0')) {
                sb.append('0');
            } else if (between('1', '9')) {
                takeDigits(sb);
            } else {
                throw error("Invalid Const");
            }
            return new Const(Integer.parseInt(sb.toString()));
        }

        private Variable parseVariable() {
            char name = take();
            StringBuilder sb = new StringBuilder();
            sb.append(name);
            while ((between('A', 'Z') || between('a', 'z'))) {
                name = take();
                sb.append(name);
            }
            if (name == 'x' || name == 'y' || name == 'z') {
                return new Variable(sb.toString());
            } else {
                throw error("No such variable: " + name);
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
    }
}
