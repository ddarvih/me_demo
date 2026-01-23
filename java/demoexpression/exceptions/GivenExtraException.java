package expression.exceptions;

public class GivenExtraException extends ParsingException {
    public GivenExtraException(String message) {
        super("thwart: " + message);
    }
}
