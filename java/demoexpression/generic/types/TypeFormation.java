package expression.generic.types;

public class TypeFormation {

    public static OperatingBase<?> chooseType(String mode) {
        return switch (mode) {
            case "i" -> new IntegerType();
            case "d" -> new DoubleType();
            case "bi" -> new BigIntegerType();
            default -> throw new IllegalArgumentException("Unknown mode: " + mode);
        };
    }
}
