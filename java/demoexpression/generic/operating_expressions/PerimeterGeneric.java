package expression.generic.operating_expressions;

import expression.generic.types.OperatingBase;

public class PerimeterGeneric<T> extends ActionGeneric<T> {

    private final OperatingBase<T> format;

    public PerimeterGeneric(ExpressionGenericBase<T> left, ExpressionGenericBase<T> right, OperatingBase<T> ob) {
        super(left, right, "area");
        this.format = ob;
    }

    @Override
    public T evaluate(T x, T y, T z) {
        return format.perimeter(left.evaluate(x, y, z), right.evaluate(x, y, z));
    }
}
