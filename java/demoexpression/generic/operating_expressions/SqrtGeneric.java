package expression.generic.operating_expressions;

import expression.generic.types.OperatingBase;

import java.util.Objects;

public class SqrtGeneric<T> implements ExpressionGenericBase<T> {
    private final ExpressionGenericBase<T> content;
    private final OperatingBase<T> format;

    public SqrtGeneric(ExpressionGenericBase<T> operand, OperatingBase<T> format) {
        this.content = operand;
        this.format = format;
    }

    @Override
    public T evaluate(T x, T y, T z) {
        return format.sqrt(content.evaluate(x, y, z));
    }

    @Override
    public String toString() {
        return "√(" + content.toString() + ")";
    }

    @Override
    public boolean equals(Object obj) {
        if (obj instanceof SqrtGeneric<?> that) {
            return this.content.equals(that.content);
        }
        return false;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.content, "√");
    }
}
