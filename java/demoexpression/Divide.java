package expression;

public class Divide extends Action {
    public Divide (Evaluatable left, Evaluatable right) {
        super(left, right, "/");
    }

    @Override
    protected long calc(long l, long r) {
        return l / r;
    }

    @Override
    protected int calc(int l, int r) {
        if (r != 0) {
            return l / r;
        } else {
            throw new ArithmeticException("division by zero");
        }
    }
}
