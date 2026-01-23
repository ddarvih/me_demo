#include "textflag.h"

TEXT ·Fibonacci(SB), NOSPLIT, $0-16
    MOVQ n+0(FP), CX

    CMPQ CX, $2
    JB no_count

    XORQ AX, AX
    MOVQ $1, DX
    MOVQ $2, BX

loop:
    MOVQ AX, SI
    ADDQ DX, SI

    MOVQ DX, AX
    MOVQ SI, DX
    INCQ BX

    CMPQ BX, CX
    JBE loop

    MOVQ DX, ret+8(FP)
    RET

no_count:
    MOVQ CX, ret+8(FP)
    RET
