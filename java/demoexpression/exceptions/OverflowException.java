package expression.exceptions;

public class OverflowException extends ExpressionException {
    public OverflowException(String name) {
        super(name + "overflow");
    }
}
