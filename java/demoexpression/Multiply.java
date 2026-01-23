package expression;

public class Multiply extends Action {
    public Multiply(Evaluatable left, Evaluatable right) {
        super(left, right, "*");
    }

    @Override
    protected long calc(long l, long r) {
        return l * r;
    }

    @Override
    protected int calc(int l, int r) {
        return l * r;
    }
}
