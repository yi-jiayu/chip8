package main

import (
	"math/bits"
)

const VF = 0xF

// 00E0 - CLS
// Clear the display.
func CLS_00E0(ip *Interpreter, instr instruction) {
	for i := 0; i < 32; i++ {
		for j := 0; j < 8; j++ {
			ip.display[i][j] = 0
		}
	}
	ip.pc++
}

// 00EE - RET
// Return from a subroutine.
//
// The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
//
// Implementation note: we're using a 0-based stack pointer, so we decrement the stack pointer first.
func RET_00EE(ip *Interpreter, instr instruction) {
	ip.stackptr--
	ip.pc = ip.stack[ip.stackptr]
}

// 0nnn - SYS addr
// Jump to a machine code routine at nnn.
//
// This instruction is only used on the old computers on which Chip-8 was originally implemented. It is ignored by modern interpreters.
func SYS_0nnn(ip *Interpreter, instr instruction) {}

// 1nnn - JP addr
// Jump to location nnn.
//
// The interpreter sets the program counter to nnn.
func JP_1nnn(ip *Interpreter, instr instruction) {
	ip.pc = instr.addr()
}

// 2nnn - CALL addr
// Call subroutine at nnn.
//
// The interpreter increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn.
//
// Implementation note: we're using a 0-based stack pointer, so we decrement the stack pointer later.
func CALL_2nnn(ip *Interpreter, instr instruction) {
	ip.stack[ip.stackptr] = ip.pc
	ip.stackptr++
	ip.pc = instr.addr()
}

// 3xkk - SE Vx, byte
// Skip next instruction if Vx = kk.
//
// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
func SE_3xkk(ip *Interpreter, instr instruction) {
	if ip.registers[instr.x()] == instr.byte() {
		ip.pc++
	}
	ip.pc++
}

// 4xkk - SNE Vx, byte
// Skip next instruction if Vx != kk.
//
// The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
func SNE_4xkk(ip *Interpreter, instr instruction) {
	if ip.registers[instr.x()] != instr.byte() {
		ip.pc++
	}
	ip.pc++
}

// 5xy0 - SE Vx, Vy
// Skip next instruction if Vx = Vy.
//
// The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
func SE_5xy0(ip *Interpreter, instr instruction) {
	if ip.registers[instr.x()] == ip.registers[instr.y()] {
		ip.pc++
	}
	ip.pc++
}

// 6xkk - LD Vx, byte
// Set Vx = kk.
//
// The interpreter puts the value kk into register Vx.
func LD_6xkk(ip *Interpreter, instr instruction) {
	ip.registers[instr.x()] = instr.byte()
	ip.pc++
}

// 7xkk - ADD Vx, byte
// Set Vx = Vx + kk.
//
// Adds the value kk to the value of register Vx, then stores the result in Vx.
func ADD_7xkk(ip *Interpreter, instr instruction) {
	ip.registers[instr.x()] += instr.byte()
	ip.pc++
}

// 8xy0 - LD Vx, Vy
// Set Vx = Vy.
//
// Stores the value of register Vy in register Vx.
func LD_8xy0(ip *Interpreter, instr instruction) {
	ip.registers[instr.x()] = ip.registers[instr.y()]
	ip.pc++
}

// 8xy1 - OR Vx, Vy
// Set Vx = Vx OR Vy.
//
// Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx.
// A bitwise OR compares the corresponding bits from two values, and if either bit is 1,
// then the same bit in the result is also 1. Otherwise, it is 0.
func OR_8xy1(ip *Interpreter, instr instruction) {
	ip.registers[instr.x()] |= ip.registers[instr.y()]
	ip.pc++
}

// 8xy2 - AND Vx, Vy
// Set Vx = Vx AND Vy.
//
// Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx.
// A bitwise AND compares the corresponding bits from two values, and if both bits are 1,
// then the same bit in the result is also 1. Otherwise, it is 0.
func AND_8xy2(ip *Interpreter, instr instruction) {
	ip.registers[instr.x()] &= ip.registers[instr.y()]
	ip.pc++
}

// 8xy3 - XOR Vx, Vy
// Set Vx = Vx XOR Vy.
//
// Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx.
// An exclusive OR compares the corresponding bits from two values, and if the bits are not both the same,
// then the corresponding bit in the result is set to 1. Otherwise, it is 0.
func XOR_8xy3(ip *Interpreter, instr instruction) {
	ip.registers[instr.x()] ^= ip.registers[instr.y()]
	ip.pc++
}

// 8xy4 - ADD Vx, Vy
// Set Vx = Vx + Vy, set VF = carry.
//
// The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255,) VF is set to 1,
// otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
func ADD_8xy4(ip *Interpreter, instr instruction) {
	x := instr.x()
	sum, carry := bits.Add(uint(ip.registers[x]), uint(ip.registers[instr.y()]), 0)
	ip.registers[x] = uint8(sum)
	ip.registers[VF] = uint8(carry)
	ip.pc++
}

// 8xy5 - SUB Vx, Vy
// Set Vx = Vx - Vy, set VF = NOT borrow.
//
// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
//
// Implementation note: Shouldn't VF be set to 1 when Vx == Vy as well?
func SUB_8xy5(ip *Interpreter, instr instruction) {
	x := instr.x()
	y := instr.y()
	diff, borrow := bits.Sub(uint(ip.registers[x]), uint(ip.registers[y]), 0)
	ip.registers[x] = uint8(diff)
	ip.registers[VF] = 1 - uint8(borrow)
	ip.pc++
}

// 8xy6 - SHR Vx {, Vy}
// Set Vx = Vx SHR 1.
//
// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
func SHR_8xy6(ip *Interpreter, instr instruction) {
	ip.registers[VF] = ip.registers[instr.x()] & 1
	ip.registers[instr.x()] <<= 2
	ip.pc++
}

// 8xy7 - SUBN Vx, Vy
// Set Vx = Vy - Vx, set VF = NOT borrow.
//
// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
func SUBN_8xy7(ip *Interpreter, instr instruction) {
	x := instr.x()
	y := instr.y()
	diff, borrow := bits.Sub(uint(ip.registers[y]), uint(ip.registers[x]), 0)
	ip.registers[x] = uint8(diff)
	ip.registers[VF] = 1 - uint8(borrow)
	ip.pc++
}

// 8xyE - SHL Vx {, Vy}
// Set Vx = Vx SHL 1.
//
// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
func SHL_8xyE(ip *Interpreter, instr instruction) {
	ip.registers[VF] = ip.registers[instr.x()] >> 7 & 1
	ip.registers[instr.x()] <<= 1
	ip.pc++
}

// 9xy0 - SNE Vx, Vy
// Skip next instruction if Vx != Vy.
//
// The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
func SNE_9xy0(ip *Interpreter, instr instruction) {
	if ip.registers[instr.x()] != ip.registers[instr.y()] {
		ip.pc++
	}
	ip.pc++
}

// Annn - LD I, addr
// Set I = nnn.
//
// The value of register I is set to nnn.
func LD_Annn(ip *Interpreter, instr instruction) {
	ip.i = instr.addr()
	ip.pc++
}

// Bnnn - JP V0, addr
// Jump to location nnn + V0.
//
// The program counter is set to nnn plus the value of V0.
func JP_Bnnn(ip *Interpreter, instr instruction) {
	ip.pc = instr.addr() + uint16(ip.registers[0])
}

// Cxkk - RND Vx, byte
// Set Vx = random byte AND kk.
//
// The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk.
// The results are stored in Vx. See instruction 8xy2 for more information on AND.
func RND_Cxkk(ip *Interpreter, instr instruction) {
	r := ip.rand()
	ip.registers[instr.x()] = r & instr.byte()
	ip.pc++
}

// Dxyn - DRW Vx, Vy, nibble
// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
//
// The interpreter reads n bytes from memory, starting at the address stored in
// I. These bytes are then displayed as sprites on screen at coordinates (Vx,
// Vy). Sprites are XORed onto the existing screen. If this causes any pixels
// to be erased, VF is set to 1, otherwise it is set to 0. If the sprite is
// positioned so part of it is outside the coordinates of the display, it wraps
// around to the opposite side of the screen. See instruction 8xy3 for more
// information on XOR, and section 2.4, Display, for more information on the
// Chip-8 screen and sprites.
func DRW_Dxyn(ip *Interpreter, instr instruction) {
	n := instr.nibble()
	x := ip.registers[instr.x()]
	y := ip.registers[instr.y()]
	sprite := ip.memory[ip.i : ip.i+uint16(n)]
	var collision uint8
	for i := uint8(0); i < n; i++ {
		rowIdx := (y + i) & 0x1F
		colGrpIdx := x >> 3
		rem := x & 0x7
		if rem > 0 {
			collision |= ip.display[rowIdx][colGrpIdx] & (sprite[i] >> rem)
			ip.display[rowIdx][colGrpIdx] ^= sprite[i] >> rem
			collision |= ip.display[rowIdx][(colGrpIdx+1)&0x7] & (sprite[i] << (8 - rem))
			ip.display[rowIdx][(colGrpIdx+1)&0x7] ^= sprite[i] << (8 - rem)
		} else {
			collision |= ip.display[rowIdx][colGrpIdx] & sprite[i]
			ip.display[rowIdx][colGrpIdx] ^= sprite[i]
		}
	}
	if collision > 0 {
		ip.registers[VF] = 1
	} else {
		ip.registers[VF] = 0
	}
	ip.pc++
}

// Ex9E - SKP Vx
// Skip next instruction if key with the value of Vx is pressed.
//
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 2.
func SKP_Ex9E(ip *Interpreter, instr instruction) {
	key := ip.registers[instr.x()]
	var mask uint16 = 1 << key
	keypad := <-ip.keypadch
	if keypad&mask > 0 {
		ip.pc++
	}
	ip.pc++
}

// ExA1 - SKNP Vx
// Skip next instruction if key with the value of Vx is not pressed.
//
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2.
func SKNP_ExA1(ip *Interpreter, instr instruction) {
	key := ip.registers[instr.x()]
	var mask uint16 = 1 << key
	keypad := <-ip.keypadch
	if keypad&mask == 0 {
		ip.pc++
	}
	ip.pc++
}

// Fx07 - LD Vx, DT
// Set Vx = delay timer value.
//
// The value of DT is placed into Vx.
func LD_Fx07(ip *Interpreter, instr instruction) {
	ip.registers[instr.x()] = <-ip.dtget
	ip.pc++
}

// Fx0A - LD Vx, K
// Wait for a key press, store the value of the key in Vx.
//
// All execution stops until a key is pressed, then the value of that key is stored in Vx.
func LD_Fx0A(ip *Interpreter, instr instruction) {
	keypad := <-ip.keypadch
	zeros := bits.LeadingZeros16(keypad)
	if zeros < 16 {
		ip.registers[instr.x()] = uint8(0xF - zeros)
		ip.pc++
	}
}

// Fx15 - LD DT, Vx
// Set delay timer = Vx.
//
// DT is set equal to the value of Vx.
func LD_Fx15(ip *Interpreter, instr instruction) {
	ip.dtset <- ip.registers[instr.x()]
	ip.pc++
}

// Fx18 - LD ST, Vx
// Set sound timer = Vx.
//
// ST is set equal to the value of Vx.
func LD_Fx18(ip *Interpreter, instr instruction) {
	ip.stset <- ip.registers[instr.x()]
	ip.pc++
}

// Fx1E - ADD I, Vx
// Set I = I + Vx.
//
// The values of I and Vx are added, and the results are stored in I.
func ADD_Fx1E(ip *Interpreter, instr instruction) {
	ip.i += uint16(ip.registers[instr.x()])
	ip.pc++
}

// Fx29 - LD F, Vx
// Set I = location of sprite for digit Vx.
//
// The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx.
// See section 2.4, Display, for more information on the Chip-8 hexadecimal font.
func LD_Fx29(ip *Interpreter, instr instruction) {
	// TODO: implement me
}

// Fx33 - LD B, Vx
// Store BCD representation of Vx in memory locations I, I+1, and I+2.
//
// The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location in I,
// the tens digit at location I+1, and the ones digit at location I+2.
func LD_Fx33(ip *Interpreter, instr instruction) {
	// TODO: implement me
}

// Fx55 - LD [I], Vx
// Store registers V0 through Vx in memory starting at location I.
//
// The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
func LD_Fx55(ip *Interpreter, instr instruction) {
	// TODO: implement me
}

// Fx65 - LD Vx, [I]
// Read registers V0 through Vx from memory starting at location I.
//
// The interpreter reads values from memory starting at location I into registers V0 through Vx.
func LD_Fx65(ip *Interpreter, instr instruction) {
	// TODO: implement me
}
