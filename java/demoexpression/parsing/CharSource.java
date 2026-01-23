package expression.parsing;

/**
 * @author Georgiy Korneev (kgeorgiy@kgeorgiy.info)
 */
public interface CharSource {
    boolean hasNext();
    char next();
    int getpos();
    IllegalArgumentException error(String message);
}
