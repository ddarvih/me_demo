package expression.generic;

import expression.generic.operating_expressions.ExpressionGenericBase;
import expression.generic.parser.ExpressionParserGeneric;
import expression.generic.types.OperatingBase;
import expression.generic.types.TypeFormation;

public class GenericTabulator implements Tabulator {

    @Override
    public Object[][][] tabulate(String mode, String expression, int x1, int x2,
                                 int y1, int y2, int z1, int z2) throws Exception {
        OperatingBase<?> obT = TypeFormation.chooseType(mode);
        return tabT(obT, expression, x1, x2, y1, y2, z1, z2);
    }

    private <T> Object[][][] tabT(OperatingBase<T> obT, String expression, int x1, int x2,
                                  int y1, int y2, int z1, int z2) {
        ExpressionParserGeneric useToParse = new ExpressionParserGeneric();
        ExpressionGenericBase<T> parsed = useToParse.parse(expression, obT);

        Object[][][] arr = new Object[x2 - x1 + 1][y2 - y1 + 1][z2 - z1 + 1];
        for (int i = 0; i < x2 - x1 + 1; i++) {
            for (int j = 0; j < y2 - y1 + 1; j++) {
                for (int k = 0; k < z2 - z1 + 1; k++) {
                    try {
                        arr[i][j][k] = parsed.evaluate(obT.intToT(x1 + i),
                                obT.intToT(y1 + j), obT.intToT(z1 + k));
                    } catch (Exception e) {
                        arr[i][j][k] = null;
                    }
                }
            }
        }
        return arr;
    }


}
