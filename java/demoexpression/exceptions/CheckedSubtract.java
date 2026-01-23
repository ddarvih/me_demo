package expression.exceptions;

import expression.Action;
import expression.Evaluatable;
import expression.TripleExpression;

public class CheckedSubtract extends Action implements TripleExpression {
    public CheckedSubtract(Evaluatable left, Evaluatable right) {
        super(left, right, "-");
    }

    @Override
    protected long calc(long l, long r) {
        return l - r;
    }

    @Override
    protected int calc(int l, int r) {
        if ((r > 0 && l < Integer.MIN_VALUE + r) || (r < 0 && l > Integer.MAX_VALUE + r)) {
            throw new OverflowException("subtract");
        }
        return l - r;
    }


}
