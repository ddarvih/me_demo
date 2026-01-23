#include "textflag.h"

TEXT ·Sum(SB), NOSPLIT, $0-32
    MOVQ x+0(FP), AX
    MOVQ x+8(FP), DX
    XORQ R10, R10

    CMPQ DX, $0
    JE done

loop:
    MOVLQSX (AX), R9
    ADDQ R9, R10

    ADDQ $4, AX
    DECQ DX
    JNZ  loop

done:
    MOVQ R10, ret+24(FP)
    RET
