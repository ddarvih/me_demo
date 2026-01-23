package expression.generic.parser;

import expression.exceptions.*;
import expression.generic.operating_expressions.*;
import expression.generic.types.OperatingBase;
import expression.parsing.BaseParser;
import expression.parsing.CharSource;
import expression.parsing.StringSource;

public class ExpressionParserGeneric {

    public <T> ExpressionGenericBase<T> parse(final String source, OperatingBase<T> format) {
        return parse(new StringSource(source), format);
    }

    public static <T> ExpressionGenericBase<T> parse(final CharSource source, OperatingBase<T> form) {
        return new MyExpressionParserG<T>(source, form).startParse();
    }

    private static class MyExpressionParserG<T> extends BaseParser {
        protected final OperatingBase<T> f;

        public MyExpressionParserG(final CharSource source, OperatingBase<T> f) {
            super(source);
            this.f = f;
        }

        public ExpressionGenericBase<T> startParse() {
            ExpressionGenericBase<T> res = parseExpression();
            if (!eof()) {
                throw new ExpressionException("collapsed expressions (more than one without proper connection) -> "
                        + take() + " position:" + getpos());
            }
            return res;
        }

        public ExpressionGenericBase<T> parseExpression() {
            return parseExpression(parseTerm());
        }

        public ExpressionGenericBase<T> parseExpression(ExpressionGenericBase<T> left) {
            skipWhitespace();
            if (take('◣') || test('a')) {
                if (test('a')) {
                    expectLine("area");
                    if (test('a')) {
                        take();
                    }
                    skipWhitespace();
                }
                ExpressionGenericBase<T> right = parseTerm();
                skipWhitespace();
                return parseExpression(new AreaGeneric<>(left, right, f));
            } else if (take('▯') || test('p')) {
                if (test('p')) {
                    expectLine("perimeter");
                    skipWhitespace();
                }
                ExpressionGenericBase<T> right = parseTerm();
                skipWhitespace();
                return parseExpression(new PerimeterGeneric<>(left, right, f));
            } else if (take('+')) {
                ExpressionGenericBase<T> right = parseTerm();
                skipWhitespace();
                return parseExpression(new AddGeneric<>(left, right, f));
            } else if (take('-')) {
                ExpressionGenericBase<T> right = parseTerm();
                skipWhitespace();
                return parseExpression(new SubtractGeneric<>(left, right, f));
            } else {
                return left;
            }
        }

        private ExpressionGenericBase<T> parseTerm() {
            return parseTerm(parseOperand());
        }

        private ExpressionGenericBase<T> parseTerm(ExpressionGenericBase<T> left) {
            skipWhitespace();
            if (take('*')) {
                ExpressionGenericBase<T> right = parseOperand();
                skipWhitespace();
                return parseTerm(new MultiplyGeneric<>(left, right, f));
            } else if (take('/')) {
                ExpressionGenericBase<T> right = parseOperand();
                skipWhitespace();
                return parseTerm(new DivideGeneric<>(left, right, f));
            } else {
                return left;
            }
        }


        private ExpressionGenericBase<T> parseOperand() {
            skipWhitespace();
            if (test('(') || test('{') || test('[')) {
                skipWhitespace();
                char pair_to = take();
                char closing = 0;
                skipWhitespace();
                ExpressionGenericBase<T> expr = parseExpression();
                skipWhitespace();
                if (pair_to == '(') {
                    closing = ')';
                } else if (pair_to == '{') {
                    closing = '}';
                } else if (pair_to == '[') {
                    closing = ']';
                } else {
                    throw new GivenLessException("no matching closing " + getpos());
                }
                try {
                    expect(closing);
                } catch (Exception e) {
                    throw new GivenLessException("no matching closing " + getpos());
                }
                return expr;
            } else if (take('√') || test('s')) {
                if (test('s')) {
                    expectLine("sqrt");
                    skipWhitespace();
                }
                return new SqrtGeneric<>(parseOperand(), f);
            } else if (take('-')) {
                if (between('1', '9')) {
                    return parseConst(true);
                }
                return new UnMinusGeneric<>(parseOperand(), f);
            } else if (between('0', '9')) {
                return parseConst(false);
            } else if (between('A', 'Z') || between('a', 'z')) {
                return parseVariable();
            } else {
                int problemPlace = getpos();
                throw new GivenLessException("Expected operand, found: " + (eof() ? "end of input" : ((take()) +
                        " position: " + problemPlace)));
            }
        }

        private ConstGeneric<T> parseConst(boolean needsMinus) {
            final StringBuilder sb = new StringBuilder();
            if (needsMinus) {
                sb.append('-');
            }
            if (take('0')) {
                sb.append('0');
                if (between('1', '9')) {
                    throw new GivenExtraException("found smth after zero, " + " position:" + getpos());
                }
            } else if (between('1', '9')) {
                takeDigits(sb);
            } else {
                throw new BadElementException("Invalid Const" + " position:" + getpos());
            }
            try {
                return new ConstGeneric<T>(f.parseT(sb.toString()));
            } catch (NumberFormatException e) {
                throw new BadElementException("not allowed constant, position:" + getpos());
            }
        }

        private VariableGeneric<T> parseVariable() {
            char name = take();
            StringBuilder sb = new StringBuilder();
            sb.append(name);
            try {
                switch (name) {
                    case 'a' -> expectLine("area");
                    case 'p' -> expectLine("perimeter");
                    case 's' -> expectLine("sqrt");
                    default -> {
                        ;
                    }
                }
            } catch (Exception e) {
                throw new BadElementException("Incorrect usage of operation (not enough arguments), " + " pos: " + getpos());
            }
            while ((between('A', 'Z') || between('a', 'z'))) {
                name = take();
                sb.append(name);
            }
            if (name == 'x' || name == 'y' || name == 'z') {
                return new VariableGeneric<T>(sb.toString());
            } else {
                throw new BadElementException("No such variable: " + name);
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

    }
}