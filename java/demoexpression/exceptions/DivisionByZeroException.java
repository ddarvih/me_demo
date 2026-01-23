package expression.exceptions;

public class DivisionByZeroException extends BadOperandException {
    public DivisionByZeroException() {
        super("division by zero");
    }
}
