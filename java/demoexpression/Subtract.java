package expression;

public class Subtract extends Action {
    public Subtract(Evaluatable left, Evaluatable right) {
        super(left, right, "-");
    }

    @Override
    protected long calc(long l, long r) {
        return l - r;
    }

    @Override
    protected int calc(int l, int r) {
        if ((r >= 0 && l < Integer.MIN_VALUE + r) || (r < 0 && l > Integer.MAX_VALUE + r)) {
            return l - r;
        } else {
            throw new ArithmeticException("subtract overflow");
        }
    }
}
