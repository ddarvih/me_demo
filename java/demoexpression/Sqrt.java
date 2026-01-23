package expression;

import java.util.Objects;

public class Sqrt extends Evaluatable {
    private final Evaluatable content;

    public Sqrt(Evaluatable operand) {
        this.content = operand;
    }

    @Override
    public int evaluate(int x) {
        return (int)Math.sqrt(content.evaluate(x));
    }

    @Override
    public long evaluateL(long x, long y, long z) {
        return (int)Math.sqrt(content.evaluateL(x, y, z));
    }

    @Override
    public int evaluate(int x, int y, int z) {
        return (int)Math.sqrt(content.evaluate(x, y, z));
    }

    @Override
    public String toString() {
        return "√(" + content.toString() + ")";
    }

    @Override
    public boolean equals(Object obj) {
        if (obj instanceof Sqrt that) {
            return this.content.equals(that.content);
        }
        return false;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.content, "√");
    }
}