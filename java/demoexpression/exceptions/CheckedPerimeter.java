package expression.exceptions;

import expression.Action;
import expression.Const;
import expression.Evaluatable;
import expression.TripleExpression;

public class CheckedPerimeter extends Action implements TripleExpression {

    public CheckedPerimeter(Evaluatable left, Evaluatable right) {
        super(left, right, "perimeter");
    }

    @Override
    protected long calc(long l, long r) {
        return (l + r) * 2;
    }

    @Override
    protected int calc(int l, int r) {
        int eva;
        int eval;
        if (l < 0 || r < 0) {
            throw new BadOperandException("triangle can not have negative sides");
        }
        try{
            eva = new CheckedAdd(new Const(l), new Const(r)).calc(l, r);
            eval = new CheckedMultiply(Const.TWO, new Const(eva)).calc(2, eva);
        } catch (OverflowException e) {
            throw new OverflowException("area overflow problems");
        }
        return eval;
    }

}

