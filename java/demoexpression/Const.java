package expression;

import java.util.Objects;

public class Const extends Evaluatable implements Expression {

    public static final Const TWO = new Const(2);

    private final Number content;

    public Const(int number) {
        content = number;
    }

    public Const(long number) {
        content  = number;
    }

    @Override
    public int evaluate(int put) {
        return content.intValue();
    }

    @Override
    public long evaluateL(long x, long y, long z) {
        return content.longValue();
    }

    @Override
    public int evaluate(int x, int y, int z) {
        return content.intValue();
    }

    @Override
    public String toString() {
        return content.toString();
    }

    @Override
    public boolean equals(Object action) {
        if (action instanceof Const that) {
            return this.content.equals(that.content);
        }
        return false;
    }

    @Override
    public int hashCode() {
        return Objects.hash(content);
    }

}
