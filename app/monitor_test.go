package app

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/blackchip-org/retro-cs/mock"
	"github.com/blackchip-org/retro-cs/rcs"
)

type monitorFixture struct {
	cpu *mock.CPU
	out bytes.Buffer
	mon *Monitor
}

func newMonitorFixture() *monitorFixture {
	mach := mock.NewMach()
	cpu := mach.CPU[0].(*mock.CPU)
	f := &monitorFixture{
		mon: NewMonitor(mach),
		cpu: cpu,
	}
	f.mon.out.SetOutput(&f.out)
	return f
}

func testMonitorInput(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}

func testMonitorRun(mon *Monitor) {
	go mon.Run()
	mon.mach.Run()
}

func TestMonitor(t *testing.T) {
	for _, test := range monitorTests {
		t.Run(test.name, func(t *testing.T) {
			f := newMonitorFixture()
			f.mon.in = testMonitorInput(strings.Join(test.in, "\n"))
			testMonitorRun(f.mon)
			have := strings.TrimSpace(f.out.String())
			want := strings.TrimSpace(test.want)
			if have != want {
				t.Errorf("\n have: \n%v \n want: \n%v", have, want)
			}
		})
	}
}

var monitorTests = []struct {
	name string
	in   []string
	want string
}{
	{
		"conversions",
		[]string{"42", "$2a", "%101010", "q"},
		`
42 $2a %101010
42 $2a %101010
42 $2a %101010
		`,
	}, {
		"break",
		[]string{
			"b set $3456",
			"b set $2345",
			"b set $1234",
			"b",
			"b clear $2345",
			"b",
			"b clear-all",
			"b",
			"b set $123456",
			"q",
		},
		`
$1234
$2345
$3456
$1234
$3456
invalid address: $123456
		`,
	}, {
		"dasm",
		[]string{
			"dasm lines 4",
			"poke $10 $09",
			"poke $11 $19 $ab",
			"poke $13 $29 $cd $ab",
			"poke $16 $27 $cd $ab",
			"d $10",
			"q",
		},
		`
$0010:  09        i09
$0011:  19 ab     i19 $ab
$0013:  29 cd ab  i29 $abcd
$0016:  27 cd ab  i27 $abcd
`,
	}, {
		"dasm continue",
		[]string{
			"dasm lines 1",
			"poke $0 $09",
			"poke $1 $19 $ab",
			"poke $3 $29 $cd $ab",
			"poke $6 $27 $cd $ab",
			"d",
			"d",
			"",
			"q",
		},
		`
$0000:  09        i09
$0001:  19 ab     i19 $ab
$0003:  29 cd ab  i29 $abcd
		`,
	}, {
		"dasm range",
		[]string{
			"dasm lines 4",
			"poke $10 $09",
			"poke $11 $19 $ab",
			"poke $13 $29 $cd $ab",
			"poke $16 $27 $cd $ab",
			"d $11 $14",
			"q",
		},
		`
$0011:  19 ab     i19 $ab
$0013:  29 cd ab  i29 $abcd
		`,
	}, {
		"go",
		[]string{"break set $10", "g", "_yield", "cpu reg pc", "q"},
		`
[break]
pc:0010 a:00 b:00 q:false z:false
$0010 +16
		`,
	}, {
		"memory",
		[]string{"mem lines 2", "m", "q"},
		`
$0000  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
$0010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
	`,
	}, {
		"memory page 1",
		[]string{"mem lines 2", "m $100", "q"},
		`
$0100  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
$0110  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
	`,
	}, {
		"memory next page",
		[]string{"mem lines 2", "m $100", "m", "q"},
		`
$0100  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
$0110  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
$0120  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
$0130  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
`,
	}, {
		"memory next page repeat",
		[]string{"mem lines 2", "m $100", "", "q"},
		`
$0100  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
$0110  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
$0120  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
$0130  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
`,
	}, {
		"memory range",
		[]string{"m $140 $15f", "q"},
		`
$0140  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
$0150  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
`,
	}, {
		"memory lines",
		[]string{"mem lines 3", "mem lines", "m", "q"},
		`
3
$0000  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
$0010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
$0020  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................
`,
	}, {
		"memory encoding",
		[]string{
			"mem encoding",
			"poke $0 1 2 3 $41 $42 $43",
			"m 00 $0f",
			"mem encoding az26",
			"m 00 $0f",
			"mem encoding ebcdic",
			"q",
		},
		`
* ascii
  az26
$0000  01 02 03 41 42 43 00 00  00 00 00 00 00 00 00 00  ...ABC..........
$0000  01 02 03 41 42 43 00 00  00 00 00 00 00 00 00 00  ABC.............
no such encoding: ebcdic
		`,
	}, {
		"memory fill",
		[]string{
			"mem lines 1",
			"mem fill 0004 $000b $ff",
			"m 0",
			"mem fill $000b 0004 $aa",
			"m 0",
			"q",
		},
		`
$0000  00 00 00 00 ff ff ff ff  ff ff ff ff 00 00 00 00  ................
$0000  00 00 00 00 ff ff ff ff  ff ff ff ff 00 00 00 00  ................
		`,
	}, {
		"next",
		[]string{"n", "q"},
		`
$0000:  00        i00
		`,
	}, {
		"peek",
		[]string{"poke $1234 $ab", "peek $1234", "q"},
		`
$ab +171
		`,
	}, {
		"registers and flags set",
		[]string{
			"r",
			"cpu reg pc $1234",
			"cpu reg a $56",
			"cpu reg c $ff",
			"cpu flag q on",
			"cpu flag a on",
			"r",
			"q",
		},
		`
[pause]
pc:0000 a:00 b:00 q:false z:false
no such register: c
no such flag: a
[pause]
pc:1234 a:56 b:00 q:true z:false
		`,
	}, {
		"registers and flags list",
		[]string{"cpu reg", "cpu flag", "q"},
		`
a
b
pc
q
z
		`,
	}, {
		"registers and flags set",
		[]string{
			"cpu reg a $56",
			"cpu flag q on",
			"cpu reg a",
			"cpu reg c",
			"cpu flag q",
			"cpu flag a",
			"q"},
		`
$56 +86 %01010110
no such register: c
true
no such flag: a
		`,
	}, {
		"step",
		[]string{"s", "r", "s", "", "r", "q"},
		`
$0001:  00        i00
[pause]
pc:0001 a:00 b:00 q:false z:false
$0002:  00        i00
$0003:  00        i00
[pause]
pc:0003 a:00 b:00 q:false z:false
		`,
	}, {
		"trace",
		[]string{
			"poke 0 $a $b $c",
			"trace",
			"break set 1",
			"break set 2",
			"go",
			"_yield",
			"trace",
			"go",
			"quit",
		},
		`
$0000:  0a        i0a

[break]
pc:0001 a:00 b:00 q:false z:false
		`,
	}, {
		"watch",
		[]string{
			"watch set $10 rw",
			"watch list",
			"poke $10 $22",
			"peek $10",
			"watch clear $10",
			"watch set $10 r",
			"watch list",
			"poke $10 $22",
			"peek $10",
			"watch clear $10",
			"watch set $10 w",
			"watch list",
			"poke $10 $22",
			"peek $10",
			"watch clear-all",
			"watch",
			"q",
		},
		`
$0010 rw
write($0010) => $22
$22 <= read($0010)
$22 +34
$0010 r
$22 <= read($0010)
$22 +34
$0010 w
write($0010) => $22
$22 +34
		`,
	},
}

func TestBreakpointOffset(t *testing.T) {
	f := newMonitorFixture()
	f.cpu.OffsetPC = 1
	cmds := `
b set $10
g
_yield
cpu reg pc
q
	`
	f.mon.in = testMonitorInput(cmds)
	testMonitorRun(f.mon)
	have := strings.TrimSpace(f.out.String())
	want := strings.TrimSpace(`
[break]
pc:000f a:00 b:00 q:false z:false
$000f +15
`)
	if have != want {
		t.Errorf("\n have: \n%v \n want: \n%v", have, want)
	}
}

func TestDump(t *testing.T) {
	var dumpTests = []struct {
		name     string
		start    int
		data     func() []int
		showFrom int
		showTo   int
		want     string
	}{
		{
			"one line", 0x10,
			func() []int { return []int{} },
			0x10, 0x20, "" +
				"$0010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................",
		}, {
			"two lines", 0x10,
			func() []int { return []int{} },
			0x10, 0x30, "" +
				"$0010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................\n" +
				"$0020  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................",
		}, {
			"jagged top", 0x10,
			func() []int { return []int{} },
			0x14, 0x30, "" +
				"$0010              00 00 00 00  00 00 00 00 00 00 00 00      ............\n" +
				"$0020  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................",
		}, {
			"jagged bottom", 0x10,
			func() []int { return []int{} },
			0x10, 0x2b, "" +
				"$0010  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  ................\n" +
				"$0020  00 00 00 00 00 00 00 00  00 00 00 00              ............",
		},
		{
			"single value", 0x10,
			func() []int { return []int{0, 0x41} },
			0x11, 0x11, "" +
				"$0010     41                                              A",
		},
		{
			"$40-$5f", 0x10,
			func() []int {
				data := make([]int, 0)
				for i := 0x40; i < 0x60; i++ {
					data = append(data, i)
				}
				return data
			},
			0x10, 0x30, "" +
				"$0010  40 41 42 43 44 45 46 47  48 49 4a 4b 4c 4d 4e 4f  @ABCDEFGHIJKLMNO\n" +
				"$0020  50 51 52 53 54 55 56 57  58 59 5a 5b 5c 5d 5e 5f  PQRSTUVWXYZ[\\]^_",
		},
	}

	mock.ResetMemory()
	m := mock.TestMemory
	for _, test := range dumpTests {
		t.Run(test.name, func(t *testing.T) {
			for i, value := range test.data() {
				m.Write(test.start+i, uint8(value))
			}
			have := dump(m, test.showFrom, test.showTo, rcs.ASCIIDecoder)
			have = strings.TrimSpace(have)
			if have != test.want {
				t.Errorf("\n have: \n%v \n want: \n%v \n", have, test.want)
			}
		})
	}
}
