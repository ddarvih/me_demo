package expression;

import java.util.Objects;

public class UnMinus extends Evaluatable {
    private final Evaluatable operand;

    public UnMinus(Evaluatable operand) {
        this.operand = operand;
    }

    @Override
    public int evaluate(int x) {
        return -operand.evaluate(x);
    }

    @Override
    public long evaluateL(long x, long y, long z) {
        return -operand.evaluateL(x, y, z);
    }

    @Override
    public int evaluate(int x, int y, int z) {
        return -operand.evaluate(x, y, z);
    }

    @Override
    public String toString() {
        return "-(" + operand.toString() + ")";
    }

    @Override
    public boolean equals(Object obj) {
        if (obj instanceof UnMinus that) {
            return this.operand.equals(that.operand);
        }
        return false;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.operand, "-");
    }
}
