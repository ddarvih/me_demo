package expression.generic.types;

import expression.exceptions.*;

public class IntegerType implements OperatingBase<Integer> {
    @Override
    public Integer add(Integer left, Integer right) {
        if (checkOverflowSum(left, right)) {
            throw new OverflowException("sum");
        }
        return left + right;
    }

    private boolean checkOverflowSum(int l, int r) {
        return (l > 0 && r > 0 && l > Integer.MAX_VALUE - r) ||
                (l < 0 && r < 0 && l < Integer.MIN_VALUE - r);
    }

    @Override
    public Integer divide(Integer left, Integer right) {
        if (right != 0) {
            if (checkOverflowDiv(left, right)) {
                throw new OverflowException("division");
            }
            return left / right;
        } else {
            throw new DivisionByZeroException();
        }
    }

    private boolean checkOverflowDiv(int l, int r) {
        return (l == Integer.MIN_VALUE && r == -1);
    }

    @Override
    public Integer multiply(Integer left, Integer right) {
        if (checkOverflowMul(left, right)) {
            throw new OverflowException("multiply");
        }
        return left * right;
    }

    private boolean checkOverflowMul(int l, int r) {
        return (l > 0 && r > 0 && l > Integer.MAX_VALUE / r) ||
                (l < 0 && r < 0 && l < Integer.MAX_VALUE / r) ||
                ((l > 0 && r < 0 && r < Integer.MIN_VALUE / l) || (l < 0 && r > 0 && l < Integer.MIN_VALUE / r));
    }

    @Override
    public Integer sqrt(Integer content) {
        if (content < 0) {
            throw new BadOperandException("sqrt: <0");
        }
        return (int) Math.sqrt(content);
    }

    @Override
    public Integer subtract(Integer left, Integer right) {
        if ((right > 0 && left < Integer.MIN_VALUE + right) || (right < 0 && left > Integer.MAX_VALUE + right)) {
            throw new OverflowException("subtract");
        }
        return left - right;
    }

    @Override
    public Integer unMinus(Integer operand) {
        if (operand == Integer.MIN_VALUE) {
            throw new OverflowException("negate");
        }
        return -operand;
    }

    @Override
    public Integer parseT(String s) {
        return Integer.parseInt(s);
    }

    @Override
    public Integer intToT(int i) {
        return i;
    }

    @Override
    public Integer perimeter(Integer left, Integer right) {
        if (left < 0 || right < 0) {
            throw new BadOperandException("triangle can not have negative sides");
        }
        if (checkOverflowSum(left, right) || checkOverflowMul(2, left + right)) {
            throw new OverflowException("area overflow problems");
        }
        return (left + right) * 2;
    }

    @Override
    public Integer area(Integer left, Integer right) {
        if (left < 0 || right < 0) {
            throw new BadOperandException("triangle can not have negative sides");
        }
        int evaluatedArea;
        int l = left;
        int r = right;
        boolean dR = false;
        boolean dev = false;
        if (left % 2 == 0) {
            dev = true;
            left = left / 2;
        } else if (right % 2 == 0) {
            dev = true;
            right = right / 2;
        } else {
            if (left > right) {
                left = (left - 1) / 2;
            } else {
                dR = true;
                right = (right - 1) / 2;
            }
        }
        if (checkOverflowMul(left, right)) {
            throw new OverflowException("area overflow problems");
        }
        evaluatedArea = left * right;
        if (!dev) {
            evaluatedArea += dR ? l / 2 : r / 2;
        }
        return evaluatedArea;
    }
}
