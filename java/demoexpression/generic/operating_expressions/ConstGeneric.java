package expression.generic.operating_expressions;

import java.util.Objects;

public class ConstGeneric<T> implements ExpressionGenericBase<T>{
    private final T content;
    public ConstGeneric(T number) {
        content = number;
    }

    @Override
    public T evaluate(T x, T y, T z) {
        return content;
    }

    @Override
    public String toString() {
        return content.toString();
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.content);
    }

    @Override
    public boolean equals(Object action) {
        if (action instanceof ConstGeneric<?> that) {
            return this.content == that.content;
        }
        return false;
    }

}
