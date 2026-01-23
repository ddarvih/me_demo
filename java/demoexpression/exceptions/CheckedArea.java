package expression.exceptions;

import expression.*;

public class CheckedArea extends Action implements TripleExpression {

    public CheckedArea(Evaluatable left, Evaluatable right) {
        super(left, right, "area");
    }

    @Override
    protected long calc(long l, long r) {
        return l * r / 2;
    }

    @Override
    protected int calc(int l, int r) {
        return calcArea(l, r);
    }

    private int calcArea(int l, int r) {
        int evaluatedArea;
        if (l < 0 || r < 0) {
            throw new BadOperandException("triangle can not have negative sides");
        }
        try {
            boolean dev = false;
            if (l % 2 == 0) {
                dev = true;
                l = l / 2;
            } else if (r % 2 == 0) {
                dev = true;
                r = r / 2;
            }
            evaluatedArea = new CheckedMultiply(new Const(l), new Const(r)).calc(l, r);
            if (!dev) {
                evaluatedArea /= 2;
            }
        } catch (OverflowException e) {
            throw new OverflowException("area overflow problems");
        }
        return evaluatedArea;
    }


}
