package expression.exceptions;

public class GivenLessException extends ParsingException {

    public GivenLessException(String message) {
        super("smth is missing: " + message);
    }
}
