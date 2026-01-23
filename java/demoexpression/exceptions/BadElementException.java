package expression.exceptions;

public class BadElementException extends ParsingException {
    public BadElementException(String message) {
        super("don't know what to do with -> " + message);
    }
}
