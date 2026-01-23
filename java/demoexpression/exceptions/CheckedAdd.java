package expression.exceptions;

import expression.Action;
import expression.Evaluatable;
import expression.TripleExpression;

public class CheckedAdd extends Action implements TripleExpression {
    public CheckedAdd(Evaluatable left, Evaluatable right) {
        super(left, right, "+");
    }

    private boolean checkOverflow(int l, int r) {
        return (l > 0 && r > 0 && l > Integer.MAX_VALUE - r) ||
                (l < 0 && r < 0 && l < Integer.MIN_VALUE - r);
    }

    @Override
    protected long calc(long l, long r) {
        return l + r;
    }

    @Override
    protected int calc(int l, int r) {
        if (checkOverflow(l, r)) {
            throw new OverflowException("sum");
        }
        return l + r;
    }

    @Override
    protected void checkValues(int l, int r) {
        if (checkOverflow(l, r)) {
            throw new OverflowException("add");
        }
    }
}
