package expression.generic.types;

import java.math.BigInteger;

public class BigIntegerType implements OperatingBase<BigInteger>{
    @Override
    public BigInteger add(BigInteger left, BigInteger right) {
        return left.add(right);
    }

    @Override
    public BigInteger divide(BigInteger left, BigInteger right) {
        return left.divide(right);
    }

    @Override
    public BigInteger multiply(BigInteger left, BigInteger right) {
        return left.multiply(right);
    }

    @Override
    public BigInteger sqrt(BigInteger content) {
        return content.sqrt();
    }

    @Override
    public BigInteger subtract(BigInteger left, BigInteger right) {
        return left.subtract(right);
    }

    @Override
    public BigInteger unMinus(BigInteger operand) {
        return operand.negate();
    }

    @Override
    public BigInteger parseT(String s) {
        return new BigInteger(s);
    }

    @Override
    public BigInteger intToT(int i) {
        return BigInteger.valueOf(i);
    }

    @Override
    public BigInteger perimeter(BigInteger left, BigInteger right) {
        // :NOTE: use BigInteger.TWO
        return (left.add(right)).multiply(BigInteger.valueOf(2));
    }

    @Override
    public BigInteger area(BigInteger left, BigInteger right) {
        return (left.multiply(right)).divide(BigInteger.valueOf(2));
    }
}
