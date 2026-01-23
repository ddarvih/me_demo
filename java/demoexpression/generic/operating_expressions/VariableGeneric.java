package expression.generic.operating_expressions;

import java.util.Objects;

public class VariableGeneric<T> implements ExpressionGenericBase<T> {
    private final String name;

    public VariableGeneric(String name) {
        this.name = name;
    }

    @Override
    public T evaluate(T x, T y, T z) {
        return switch(name.charAt(name.length()-1)) {
            case 'x' -> x;
            case 'y' -> y;
            case 'z' -> z;
            default -> throw new IllegalArgumentException();
        };
    }

    @Override
    public String toString() {
        return name;
    }

    @Override
    public boolean equals(Object action) {
        if (action instanceof VariableGeneric<?> that) {
            return this.name.equals(that.name);
        }
        return false;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.name);
    }
}
