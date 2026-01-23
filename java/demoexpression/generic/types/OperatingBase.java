package expression.generic.types;

public interface OperatingBase<T> {
    T add(T left, T right);
    T divide(T left, T right);
    T multiply(T left, T right);
    T sqrt(T content);
    T subtract(T left, T right);
    T unMinus(T operand);
    T parseT(String s);
    T intToT(int i);

    T perimeter(T left, T right);
    T area(T left, T right);
}
