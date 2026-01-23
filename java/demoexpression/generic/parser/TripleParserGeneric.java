package expression.generic.parser;

import expression.generic.operating_expressions.ExpressionGenericBase;

@FunctionalInterface
public interface TripleParserGeneric<T> {
    ExpressionGenericBase<T> parse(String expression);
}
