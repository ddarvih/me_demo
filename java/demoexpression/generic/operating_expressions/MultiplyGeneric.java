package expression.generic.operating_expressions;

import expression.generic.types.OperatingBase;

public class MultiplyGeneric <T> extends ActionGeneric<T>{
    private final OperatingBase<T> format;

    public MultiplyGeneric(ExpressionGenericBase<T> left, ExpressionGenericBase<T> right, OperatingBase<T> ob) {
        super(left, right, "*");
        this.format = ob;
    }

    @Override
    public T evaluate(T x, T y, T z) {
        return format.multiply(left.evaluate(x, y, z), right.evaluate(x, y, z));
    }
}
