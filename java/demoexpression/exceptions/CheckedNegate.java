package expression.exceptions;

import expression.Evaluatable;
import expression.TripleExpression;

import java.util.Objects;

public class CheckedNegate extends Evaluatable implements TripleExpression {
    private final Evaluatable operand;

    public CheckedNegate(Evaluatable operand) {
        this.operand = operand;
    }

    @Override
    public int evaluate(int x) {
        if (operand.evaluate(x) == Integer.MIN_VALUE) {
            throw new OverflowException("negate");
        }
        return -operand.evaluate(x);
    }

    @Override
    public long evaluateL(long x, long y, long z) {
        return -operand.evaluateL(x, y, z);
    }

    @Override
    public int evaluate(int x, int y, int z) {
        if (operand.evaluate(x, y, z) == Integer.MIN_VALUE) {
            throw new OverflowException("negate");
        }
        return -operand.evaluate(x, y, z);
    }

    @Override
    public String toString() {
        return "-(" + operand.toString() + ")";
    }

    @Override
    public boolean equals(Object obj) {
        if (obj instanceof CheckedNegate that) {
            return this.operand.equals(that.operand);
        } else {
            return false;
        }
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.operand, "-");
    }
}
