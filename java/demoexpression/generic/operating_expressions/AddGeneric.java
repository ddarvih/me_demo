package expression.generic.operating_expressions;

import expression.generic.types.OperatingBase;

public class AddGeneric<T> extends ActionGeneric<T> {
    private final OperatingBase<T> format;

    public AddGeneric(ExpressionGenericBase<T> left, ExpressionGenericBase<T> right, OperatingBase<T> ob) {
        super(left, right, "+");
        this.format = ob;
    }

    @Override
    public T evaluate(T x, T y, T z) {
        return format.add(left.evaluate(x, y, z), right.evaluate(x, y, z));
    }

}
