package expression;

import expression.exceptions.OverflowException;

import java.util.Objects;

public abstract class Action extends Evaluatable implements Expression, LongTripleExpression {
    protected final Evaluatable left;
    protected final Evaluatable right;
    private final String act;

    public Action(Evaluatable left, Evaluatable right, String act) {
        this.left = left;
        this.right = right;
        this.act = act;
    }

    @Override
    public int evaluate(int put) {
        int l = left.evaluate(put);
        int r = right.evaluate(put);
        checkValues(l, r);
        return calc(l, r);
    }
    @Override
    public long evaluateL(long x, long y, long z) {
        return calc(left.evaluateL(x, y, z), right.evaluateL(x, y, z));
    }

    @Override
    public int evaluate(int x, int y, int z) {
        int l = left.evaluate(x, y, z);
        int r = right.evaluate(x, y, z);
        checkValues(l, r);
        return calc(l, r);
    }

    protected abstract long calc(long l, long r);
    protected abstract int calc(int l, int r);

    protected void checkValues(int l, int r) {}

    @Override
    public String toString() {
        return "(" + left.toString() + " " + act + " " + right.toString() + ")";
    }

    @Override
    public boolean equals(Object action) {
         if (action == this) {
             return true;
         }

        if (action instanceof Action that) {
            return this.left.equals(that.left) && this.right.equals(that.right) && this.act.equals(that.act);
        }
        return false;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.left, this.right, this.act);
    }

}
