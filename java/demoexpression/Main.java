package expression;

public class Main {
    public static void main(String[] args) {
        int s = new Subtract(
                new Multiply(
                        new Const(2),
                        new Variable("x")
                ),
                new Const(3)
        ).evaluate(5);

        System.out.println(s);

        String ss = new Subtract(
                new Multiply(
                        new Const(2),
                        new Variable("x")
                ),
                new Const(3)
        ).toString();

        System.out.println(ss);

        Evaluatable expr = new Subtract(new Add(new Variable("x"), new Variable("y")), new Const(1L));
        long test2 = expr.evaluateL(2L, 3L, 5L);
        System.out.println(test2);
        String testS = expr.toString();
        System.out.println(testS);

        Evaluatable e = new Subtract(new Add(new Variable("x"), new Variable("y")), new Const(1));
        System.out.println(e.toString());
        System.out.println(e.evaluate(2, 3, 5));

    }
}
