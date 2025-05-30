'''
SpaceX is testing flight software subroutines (i.e., programs that consist of sequences of instructions) for a custom rocket CPU.
To ensure that the software runs correctly before it's loaded into the rocket, you need to create a CPU simulator.

The CPU has 43 32-bit unsigned integer registers, which are named R00..R42. At the start of the program, all the registers contain 0.
The CPU supports the following instructions:

MOV Rxx,Ryy - copies the value from register Rxx to register Ryy;
MOV d,Rxx - copies the numeric constant d (specified as a decimal) to register Rxx;
ADD Rxx,Ryy - calculates (Rxx + Ryy) MOD 232 and stores the result in Rxx;
DEC Rxx - decrements Rxx by one. Decrementing 0 causes an overflow and results in 2^32-1;
INC Rxx - increments Rxx by one. Incrementing 2^32-1 causes an overflow and results in 0;
INV Rxx - performs a bitwise inversion of register Rxx;
JMP d - unconditionally jumps to instruction number d (1-based). d is guaranteed to be a valid instruction number;
JZ d - jumps to instruction d (1-based) only if R00 contains 0;
NOP - does nothing.
After the last instruction has been executed, the contents of R42 are considered to be the result of the subroutine.

Write a software emulator for this CPU that executes the subroutines and returns the resulting value from R42.

All the commands in the subroutine are guaranteed to be syntactically correct and have valid register numbers, numeric constants, and jump addresses. The maximum program length is 1024 instructions. The maximum total number of instructions that will be executed until the value is returned is 5 · 104. (Keep in mind that the same instruction will be counted as many times as it will be executed.)

Example

For

subroutine = [
  "MOV 5,R00",
  "MOV 10,R01",
  "JZ 7",
  "ADD R02,R01",
  "DEC R00",
  "JMP 3",
  "MOV R02,R42"
]
the output should be
cpuEmulator(subroutine) = "50".

Here is the information about the CPU state after certain steps:

Step	Last executed
command	Non-zero registers	Comment
1	1. MOV 5,R00	R00 = 5	Put 5 into R00
2	2. MOV 10,R01	R00 = 5, R01 = 10	Put 10 into R01
3	3. JZ 7	R00 = 5, R01 = 10	Move to the next instruction
because R00 ≠ 0
4	4. ADD R02,R01	R00 = 5, R01 = 10,
R02 = 10	R02 += R01
5	5. DEC R00	R00 = 4, R01 = 10,
R02 = 10	R00 -= 1
6	6. JMP 3	R00 = 4, R01 = 10,
R02 = 10	Jump to instruction number 3,
i.e. JZ 7
7	3. JZ 7	R00 = 4, R01 = 10,
RO2 = 10	Move to the next instruction
because R00 ≠ 0
Information about 11 steps is skipped
19	3. JZ 7	R00 = 1, R01 = 10,
RO2 = 40	Move to the next instruction
because R00 ≠ 0
20	4. ADD R02,R01	R00 = 1, R01 = 10,
R02 = 50	R02 += R01
21	5. DEC R00	R00 = 0, R01 = 10,
R02 = 50	R00 -= 1
22	6. JMP 3	R00 = 0, R01 = 10,
R02 = 50	Jump to instruction number 3,
i.e. JZ 7
23	3. JZ 7	R00 = 0, R01 = 10,
R02 = 50	Jump to instruction number 7
because R00 = 0
24	7. MOV R02,R42	R00 = 0, R01 = 10,
R02 = 50, R42 = 50	R42 += R02
The subroutine is exited
Input/Output

[execution time limit] 4 seconds (py3)

[input] array.string subroutine

Guaranteed constraints:
1 ≤ subroutine.length ≤ 1024.

[output] string

Return the resulting 32-bit unsigned integer, converted into a string.
'''
import numpy as np

def cpuEmulator(subroutine):
    reg = [0] * 43
    dispatcher = {'MOV': MOV,
                  'ADD': ADD,
                  'DEC': DEC,
                  'INC': INC,
                  'INV': INV,
                  'JMP': JMP,
                  'JZ' : JZ,
                  'NOP': NOP
                  }

    i = 0
    while i < len(subroutine):
        #print(i)
        instruction =  subroutine[i]
        ins = instruction.split(' ')
        #print(ins)
        params = None
        if len(ins) > 1:
            #print(ins[1].split(',')[0])
            params = ins[1].split(',')
        i = dispatcher[ins[0]](params, reg, i)

    return str(reg[42])

def MOV(params, reg, i):
    #print(params)
    x = params[0]
    y = params[1]
    if x[0] == 'R':
        reg[int(y[1:])] = reg[int(x[1:])]
    else:
        reg[int(y[1:])] = int(x)
    #print(reg)
    return i + 1

def ADD(params, reg, i):
    x = params[0]
    y = params[1]
    reg[int(x[1:])] = (reg[int(x[1:])] + reg[int(y[1:])]) % 2**32
    #print(reg)
    return i + 1

def DEC(params, reg = None, i = 0):
    x = params[0]
    if reg[int(x[1:])] == 0:
        reg[int(x[1:])] = 2**32-1
    else:
        reg[int(x[1:])] -= 1
    #print(reg)
    return i + 1

def INC(params, reg, i):
    x = params[0]
    if reg[int(x[1:])] >= 2**32-1:
        reg[int(x[1:])] = 0
    else:
        reg[int(x[1:])] += 1
    #print(reg)
    return i + 1

def INV(params, reg, i):
    x = params[0]
    reg[int(x[1:])] = ~reg[int(x[1:])] & 0xFFFFFFFF
    #print(reg)
    return i + 1

def JMP(params, reg, i):
    x = params[0]
    #print(reg)
    return int(x) - 1

def JZ(params, reg, i):
    x = params[0]
    if reg[0] == 0:
        return int(x) - 1
    #print('skipped')
    #print(reg)
    return i + 1

def NOP(params, reg, i):
    return i + 1



if __name__ == '__main__':
    print('Answer : ' , cpuEmulator([
  "MOV 5,R00",
  "MOV 10,R01",
  "JZ 7",
  "ADD R02,R01",
  "DEC R00",
  "JMP 3",
  "MOV R02,R42"
]))
    #print(~2)
