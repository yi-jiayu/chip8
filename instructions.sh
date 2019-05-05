#!/bin/sh -

ops='CLS_00E0
RET_00EE
SYS_0nnn
JP_1nnn
CALL_2nnn
SE_3xkk
SNE_4xkk
SE_5xy0
LD_6xkk
ADD_7xkk
LD_8xy0
OR_8xy1
AND_8xy2
XOR_8xy3
ADD_8xy4
SUB_8xy5
SHR_8xy6
SUBN_8xy7
SHL_8xyE
SNE_9xy0
LD_Annn
JP_Bnnn
RND_Cxkk
DRW_Dxyn
SKP_Ex9E
SKNP_ExA1
LD_Fx07
LD_Fx0A
LD_Fx15
LD_Fx18
ADD_Fx1E
LD_Fx29
LD_Fx33
LD_Fx55
LD_Fx65'

echo "package main
" > instructions.go

for op in $ops; do cat <<EOF
func $op(ip *Interpreter, instr instruction) {

}

EOF
done >> instructions.go