package expression;

public class Add extends Action {

    public Add(Evaluatable left, Evaluatable right) {
        super(left, right, "+");
    }

    @Override
    protected long calc(long l, long r) {
        return r + l;
    }

    @Override
    protected int calc(int l, int r) {
        return r + l;
    }
}

