#include "textflag.h"

TEXT ·LowerBound(SB), NOSPLIT, $0-40
    MOVQ slice+0(FP), SI
    MOVQ slice+8(FP), CX
    MOVQ value+24(FP), DX

    XORQ AX, AX
    MOVQ CX, BX

loop:
    CMPQ AX, BX
    JGE  done

    MOVQ BX, DI
    SUBQ AX, DI
    SHRQ $1, DI
    ADDQ AX, DI

    MOVQ (SI)(DI*8), R8
    CMPQ R8, DX
    JGE else
    LEAQ 1(DI), AX
    JMP loop
else:
    MOVQ DI, BX
    JMP loop

done:
    MOVQ AX, ret+32(FP)
    RET
