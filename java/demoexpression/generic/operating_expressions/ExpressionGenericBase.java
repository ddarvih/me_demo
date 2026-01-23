package expression.generic.operating_expressions;

public interface ExpressionGenericBase<T> {
    T evaluate(T op1, T op2, T op3);
    String toString();
}
