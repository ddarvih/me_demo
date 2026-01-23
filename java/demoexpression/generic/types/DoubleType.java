package expression.generic.types;

public class DoubleType implements OperatingBase<Double> {
    @Override
    public Double add(Double left, Double right) {
        return left + right;
    }

    @Override
    public Double divide(Double left, Double right) {
        return left / right;
    }

    @Override
    public Double multiply(Double left, Double right) {
        return left * right;
    }

    @Override
    public Double sqrt(Double content) {
        return Math.sqrt(content);
    }

    @Override
    public Double subtract(Double left, Double right) {
        return left - right;
    }

    @Override
    public Double unMinus(Double operand) {
        return -operand;
    }

    @Override
    public Double parseT(String s) {
        return Double.parseDouble(s);
    }

    @Override
    public Double intToT(int i) {
        return (double) i;
    }

    @Override
    public Double perimeter(Double left, Double right) {
        return (left + right) * 2;
    }

    @Override
    public Double area(Double left, Double right) {
        return left * right / 2;
    }
}
