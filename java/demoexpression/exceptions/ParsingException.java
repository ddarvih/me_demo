package expression.exceptions;

public class ParsingException extends IllegalArgumentException {
    public ParsingException (String message) {
        super(message);
    }
}
