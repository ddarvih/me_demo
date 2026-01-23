package expression.generic.operating_expressions;

import java.util.Objects;

public abstract class ActionGeneric<T> implements ExpressionGenericBase<T> {

    protected final ExpressionGenericBase<T> left;
    protected final ExpressionGenericBase<T> right;
    private final String act;

    public ActionGeneric(ExpressionGenericBase<T> left, ExpressionGenericBase<T> right, String act) {
        this.left = left;
        this.right = right;
        this.act = act;
    }

    @Override
    public String toString() {
        return "(" + left.toString() + " " + act + " " + right.toString() + ")";
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.left, this.right, this.act);
    }
}
