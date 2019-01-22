package z80

import (
	"fmt"

	"github.com/blackchip-org/retro-cs/rcs"
)

// CPU is the Zilog Z80 processor.
type CPU struct {
	pc uint16 // Program counter
	A  uint8  // Accumulator
	F  uint8  // Flags
	B  uint8
	C  uint8
	D  uint8
	E  uint8
	H  uint8
	L  uint8

	A1 uint8 // Shadow registers
	F1 uint8
	B1 uint8
	C1 uint8
	D1 uint8
	E1 uint8
	H1 uint8
	L1 uint8

	I   uint8 // Interrupt vector base
	R   uint8 // DRAM refresh counter
	IXH uint8
	IXL uint8
	IYH uint8
	IYL uint8
	SP  uint16 // Stack pointer

	IFF1 bool // Interrupt flip flops
	IFF2 bool
	IM   uint8 // Interrupt mode
	Halt bool  // Halted by instruction

	mem   *rcs.Memory
	ops   map[uint8]func(*CPU)
	delta uint8
	// address used to load on the last (IX+d) or (IY+d) instruction
	iaddr int
}

const (
	// FlagC is the carry flag
	FlagC = uint8(1 << 0)

	// FlagN is set after subtraction
	FlagN = uint8(1 << 1)

	// FlagV is the overflow flag (also parity)
	FlagV = uint8(1 << 2)

	// FlagP is the parity flag (also overflow)
	FlagP = uint8(1 << 2)

	// Flag3 is undefined
	Flag3 = uint8(1 << 3)

	// FlagH is the half-carry flag
	FlagH = uint8(1 << 4)

	// Flag5 is undefined
	Flag5 = uint8(1 << 5)

	// FlagZ is the zero flag
	FlagZ = uint8(1 << 6)

	// FlagS is the sign flag
	FlagS = uint8(1 << 7)
)

func New(mem *rcs.Memory) *CPU {
	return &CPU{mem: mem, ops: opcodes}
}

// FIXME: return value is a testing crutch
func (c *CPU) Next() bool {
	// here := c.PC()
	opcode := c.fetch()
	execute, ok := c.ops[opcode]
	c.refreshR()
	if !ok {
		return false
		//log.Printf("%04x: illegal instruction: 0x%02x", here, opcode)
		//return
	}
	execute(c)
	return ok
}

func (c *CPU) PC() int {
	return int(c.pc)
}

func (c *CPU) SetPC(pc int) {
	c.pc = uint16(pc)
}

func (c *CPU) fetch() uint8 {
	c.pc++
	return c.mem.Read(int(c.pc - 1))
}

func (c *CPU) fetch2() int {
	return int(c.fetch()) + (int(c.fetch()) << 8)
}

func (c *CPU) refreshR() {
	// Lower 7 bits of the refresh register are incremented on an instruction
	// fetch
	bit7 := c.R & 0x80
	c.R = (c.R+1)&0x7f | bit7
}

func (c *CPU) String() string {
	b := func(v uint8, ch string) string {
		if v != 0 {
			return ch
		}
		return "."
	}

	iff1 := ""
	if c.IFF1 {
		iff1 = "iff1"
	}
	iff2 := ""
	if c.IFF2 {
		iff2 = "iff2"
	}

	return fmt.Sprintf(""+
		" pc   af   bc   de   hl   ix   iy   sp   i  r\n"+
		"%04x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %04x  %02x %02x %v\n"+
		"im %v %02x%02x %02x%02x %02x%02x %02x%02x      %v %v %v %v %v %v %v %v  %v\n",
		// line 1
		c.pc,
		c.A, c.F,
		c.B, c.C,
		c.D, c.E,
		c.H, c.L,
		c.IXH, c.IXL,
		c.IYH, c.IYL,
		c.SP,
		c.I,
		c.R,
		iff1,
		// line 2
		c.IM,
		c.A1, c.F1,
		c.B1, c.C1,
		c.D1, c.E1,
		c.H1, c.L1,
		// flags
		b(c.F&FlagS, "S"),
		b(c.F&FlagZ, "Z"),
		b(c.F&Flag5, "5"),
		b(c.F&FlagH, "H"),
		b(c.F&Flag3, "3"),
		b(c.F&FlagV, "V"),
		b(c.F&FlagN, "N"),
		b(c.F&FlagC, "C"),
		iff2,
	)
}