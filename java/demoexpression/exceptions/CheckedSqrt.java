package expression.exceptions;

import expression.Evaluatable;
import expression.TripleExpression;

import java.util.Objects;

public class CheckedSqrt extends Evaluatable implements TripleExpression {
    private final Evaluatable content;

    public CheckedSqrt(Evaluatable operand) {
        this.content = operand;
    }

    @Override
    public int evaluate(int x) {
        if (content.evaluate(x) < 0) {
            throw new BadOperandException("sqrt: <0");
        }
        return (int)Math.sqrt(content.evaluate(x));
    }

    @Override
    public long evaluateL(long x, long y, long z) {
        return (int)Math.sqrt(content.evaluateL(x, y, z));
    }

    @Override
    public int evaluate(int x, int y, int z) {
        if (content.evaluate(x, y, z) < 0) {
            throw new BadOperandException("sqrt: <0");
        }
        return (int)Math.sqrt(content.evaluate(x, y, z));
    }

    @Override
    public String toString() {
        return "sqrt(" + content.toString() + ")";
    }

    @Override
    public boolean equals(Object obj) {
        if (obj instanceof CheckedSqrt that) {
            return this.content.equals(that.content);
        } else {
            return false;
        }
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.content, "√");
    }
}