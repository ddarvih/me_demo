package expression;

import java.util.Objects;

public class Variable extends Evaluatable implements Expression, LongTripleExpression{
    private final String name;

    public Variable(String name) {
        this.name = name;
    }

    @Override
    public int evaluate(int put) {
        return put;
    }

    @Override
    public long evaluateL(long x, long y, long z) {
        return switch (name.charAt(name.length()-1)) {
            case 'x' -> x;
            case 'y' -> y;
            case 'z' -> z;
            default -> 0;
        };
    }

    @Override
    public int evaluate(int x, int y, int z) {
        return switch(name.charAt(name.length()-1)) {
            case 'x' -> x;
            case 'y' -> y;
            case 'z' -> z;
            default -> 0;
        };
    }

    @Override
    public String toString() {
        return name;
    }

    @Override
    public boolean equals(Object action) {
        if (action instanceof Variable that) {
            return this.name.equals(that.name);
        }
        return false;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.name);
    }
}
