package expression.exceptions;

import expression.TripleExpression;


public class ExpressionParser implements TripleParser {
    @Override
    public TripleExpression parse(String expression) throws Exception {
        return new UseToParse().parse(expression);
    }
}
