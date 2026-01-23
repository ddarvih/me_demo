package expression.generic.operating_expressions;

import expression.generic.types.OperatingBase;

import java.util.Objects;

public class UnMinusGeneric<T> implements ExpressionGenericBase<T> {
    private final ExpressionGenericBase<T> operand;
    private final OperatingBase<T> format;

    public UnMinusGeneric(ExpressionGenericBase<T> operand, OperatingBase<T> format) {
        this.operand = operand;
        this.format = format;
    }

    @Override
    public T evaluate(T x, T y, T z) {
        return format.unMinus(operand.evaluate(x, y, z));
    }

    @Override
    public String toString() {
        return "-(" + operand.toString() + ")";
    }


    @Override
    public boolean equals(Object obj) {
        if (obj instanceof UnMinusGeneric<?> that) {
            return this.operand.equals(that.operand);
        }
        return false;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.operand, "-");
    }
}
