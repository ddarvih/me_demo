package expression.exceptions;

import expression.Action;
import expression.Evaluatable;
import expression.TripleExpression;

public class CheckedDivide extends Action implements TripleExpression {
    public CheckedDivide(Evaluatable left, Evaluatable right) {
        super(left, right, "/");
    }

    @Override
    protected long calc(long l, long r) {
        return l / r;
    }

    @Override
    protected int calc(int l, int r) {
        return l / r;
    }

    @Override
    protected void checkValues(int l, int r) {
        if (r == 0) {
            throw new DivisionByZeroException();
        }
        if (checkOverflow(l, r)) {
            throw new OverflowException("division");
        }
    }

    private boolean checkOverflow(int l, int r) {
        return (l == Integer.MIN_VALUE && r == -1);
    }
}
