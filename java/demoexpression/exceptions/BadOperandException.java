package expression.exceptions;

public class BadOperandException extends ExpressionException {
    public BadOperandException(String reason) {
        super(reason);
    }
}
