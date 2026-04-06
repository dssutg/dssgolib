package bitfield

import (
	"encoding/hex"
	"reflect"
	"testing"
)

// Tests distinct type aliases.
type (
	aliasByte        byte
	aliasUint16      uint16
	anotherByteAlias byte
)

type uint16Pair struct{ A, B uint16 }

type arrayThenByte struct {
	A [3]byte
	B byte
}

type aliasedByteThenUint16 struct {
	AnotherByte anotherByteAlias
	Uint16      uint16
}

type structSOTB struct {
	A byte `bitfield:"7"`
	B bool `bitfield:"1"`
	C [2]byte
	D byte
}

type union11 struct { // union
	A       uint16Pair
	B       arrayThenByte
	C       aliasedByteThenUint16
	D       aliasUint16
	E       uint16
	F, G, H byte
	I       structSOTB
	J       byte
	K       [4]byte
}

type structBBU struct {
	A byte `bitfield:"1"`
	B byte
	C uint16
}

type structArrBB struct {
	A    arrayThenByte
	B, C byte
}

type structU16AB struct {
	A uint16
	B arrayThenByte
}

type structUBU struct {
	A uint16
	B byte
	C uint16
}

type union11Uint16 struct {
	A union11 `bitfield:"union"`
	B uint16
}

type structBB struct {
	A, B byte
}

type nestedUnion struct { // union
	A union11 `bitfield:"union"`
	B uint32
	C structBBU
	D structArrBB
	E structU16AB
	F structUBU
	G union11Uint16
	H structBB
	I structBB
	J uint16
	K uint32
	L [6]byte
}

type unionOfUnions struct {
	A       bool `bitfield:"1"`
	B       bool `bitfield:"1"`
	C       byte `bitfield:"4"`
	D       bool `bitfield:"1"`
	E       bool `bitfield:"1"`
	F       aliasByte
	G       uint16
	H       nestedUnion `bitfield:"union"`
	I, J, K uint16
}

type structBUint16 struct {
	A byte
	B uint16
}

type bits8Int16Int16 struct {
	A    bool `bitfield:"1"`
	B    bool `bitfield:"1"`
	C    bool `bitfield:"1"`
	D    bool `bitfield:"1"`
	E    bool `bitfield:"1"`
	F    bool `bitfield:"1"`
	G    bool `bitfield:"1"`
	H    bool `bitfield:"1"`
	I, J int16
}

type bitsFloats struct {
	A       bool `bitfield:"1"`
	B       bool `bitfield:"1"`
	C       bool `bitfield:"1"`
	D       bool `bitfield:"1"`
	E       bool `bitfield:"1"`
	F       bool `bitfield:"1"`
	G       bool `bitfield:"1"`
	H       bool `bitfield:"1"`
	I       bool `bitfield:"1"`
	J       bool `bitfield:"1"`
	K       bool `bitfield:"1"`
	L       bool `bitfield:"1"`
	M       bool `bitfield:"1"`
	N       uint16
	O, P, Q float32
	R, S    uint16
	T, U    float32
}

type bits2B struct {
	A    byte `bitfield:"1"`
	B    byte `bitfield:"2"`
	C    bool `bitfield:"1"`
	D    bool `bitfield:"1"`
	E    bool `bitfield:"1"`
	F    bool `bitfield:"1"`
	G    bool `bitfield:"1"`
	H, I byte
}

type flags struct {
	A bool `bitfield:"1"`
	B bool `bitfield:"1"`
	C bool `bitfield:"1"`
	D bool `bitfield:"1"`
	E bool `bitfield:"1"`
	F bool `bitfield:"1"`
	G bool `bitfield:"1"`
	H bool `bitfield:"1"`
}

type bitsABBBU16 struct {
	A bool `bitfield:"1"`
	B bool `bitfield:"1"`
	C byte `bitfield:"2"`
	D bool `bitfield:"1"`
	E bool `bitfield:"1"`
	F bool `bitfield:"1"`
	G arrayThenByte
	H byte
	I byte
	J uint16
}

type bitBB832B19U16 struct {
	A bool `bitfield:"1"`
	B byte
	C [8]byte
	D [33]byte
	E [19]byte
	F uint16
}

type mixed11 struct {
	A bitsABBBU16
	B bitBB832B19U16
}

type mixed1 struct {
	A             uint16
	B, C, D, E, F byte
	G             bool `bitfield:"1"`
	H             bool `bitfield:"1"`
	I             bool `bitfield:"1"`
	J             bool `bitfield:"1"`
	K             bool `bitfield:"1"`
	L             bool `bitfield:"1"`
	M             bool `bitfield:"1"`
	N, O          uint16
	P, Q          byte
	R             bool `bitfield:"1"`
	S             bool `bitfield:"1"`
	T             bool `bitfield:"1"`
	U             bool `bitfield:"1"`
	V             bool `bitfield:"1"`
}

type mixed3 struct {
	A bool `bitfield:"1"`
	B bool `bitfield:"1"`
	C byte
	D [5]byte
	E byte
}

type mixed4 struct {
	A                bool `bitfield:"1"`
	B                bool `bitfield:"1"`
	C                bool `bitfield:"1"`
	D                bool `bitfield:"1"`
	E                bool `bitfield:"1"`
	F                bool `bitfield:"1"`
	G                bool `bitfield:"1"`
	H, I, J, K, L, M byte
	N                uint16
	O, P             byte
}

type mixed5 struct {
	A    bool `bitfield:"1"`
	B    byte `bitfield:"3"`
	C    bool `bitfield:"1"`
	D    byte `bitfield:"3"`
	E    byte
	F    structBUint16
	G, H byte
	I    uint16
	J    byte
	K    bool `bitfield:"1"`
}

type mixed6 struct {
	A    bool `bitfield:"1"`
	B    bool `bitfield:"1"`
	C    bool `bitfield:"1"`
	D    bool `bitfield:"1"`
	E    byte `bitfield:"2"`
	F    bool `bitfield:"1"`
	G    bool `bitfield:"1"`
	H, I byte
	J    byte `bitfield:"4"`
	K    byte `bitfield:"4"`
	L, M byte
	N    byte `bitfield:"4"`
	O    byte `bitfield:"4"`
	P    [3]uint16
	Q    byte
	R    byte `bitfield:"4"`
	S    bool `bitfield:"1"`
	T    bool `bitfield:"1"`
	U    bool `bitfield:"1"`
}

type mixed7 struct {
	A bool `bitfield:"1"`
	B byte `bitfield:"2"`
	C bool `bitfield:"1"`
	D bool `bitfield:"1"`
	E bool `bitfield:"1"`
	F bool `bitfield:"1"`
	G bool `bitfield:"1"`
	H uint32
	I arrayThenByte
	J bool `bitfield:"1"`
	K bool `bitfield:"1"`
	L bool `bitfield:"1"`
	M bool `bitfield:"1"`
	N bool `bitfield:"1"`
	O bool `bitfield:"1"`
	P bool `bitfield:"1"`
	Q bool `bitfield:"1"`
	R bool `bitfield:"1"`
	S bool `bitfield:"1"`
	T bool `bitfield:"1"`
	U bool `bitfield:"1"`
	V bool `bitfield:"1"`
}

type mixed8 struct {
	A byte `bitfield:"2"`
}

type mixed9 struct {
	A             bool `bitfield:"1"`
	B             bool `bitfield:"1"`
	C             bool `bitfield:"1"`
	D             bool `bitfield:"1"`
	E             bool `bitfield:"1"`
	F             bool `bitfield:"1"`
	G             bool `bitfield:"1"`
	H             structBUint16
	I             arrayThenByte
	J, K, L, M, N byte
}

type mixed10 struct {
	A       bool `bitfield:"1"`
	B       bool `bitfield:"1"`
	C       bool `bitfield:"1"`
	D       bool `bitfield:"1"`
	E       bool `bitfield:"1"`
	F       byte `bitfield:"3"`
	G       structBUint16
	H, I, J byte
	K       bool `bitfield:"1"`
	L       bool `bitfield:"1"`
	M       bool `bitfield:"1"`
	N       bool `bitfield:"1"`
	O       bool `bitfield:"1"`
	P       bool `bitfield:"1"`
	Q       bool `bitfield:"1"`
	R       bool `bitfield:"1"`
	S       bool `bitfield:"1"`
}

type mixed2 struct {
	A mixed1
	B mixed3
	C mixed4
	D mixed5
	E mixed6
	F mixed7
	G mixed8
	H mixed9
	I mixed10
	J [15]byte
}

type variousBitFields struct {
	A bool `bitfield:"1"`
	B bool `bitfield:"1"`
	C bool `bitfield:"1"`
	D bool `bitfield:"1"`
	E byte `bitfield:"2"`
	F bool `bitfield:"1"`
	G bool `bitfield:"1"`
}

type largeStructA struct {
	A uint16Pair
	B byte `bitfield:"5"`
	C byte `bitfield:"3"`
	D arrayThenByte
	E byte
	F variousBitFields
	G arrayThenByte
	H byte
	I structBUint16
	J byte
	K mixed2 `bitfield:"union"`
	L byte
	M uint16
	N uint64 `bitfield:"-"`
	O uint64 `bitfield:"-"`
	P uint16 `bitfield:"-"`
	Q string `bitfield:"-"`
	R string `bitfield:"-"`
	S string `bitfield:"-"`
}

type bThenU16 struct {
	A byte
	B uint16
}

type byteStruct struct {
	A byte
}

type bStructThen32B struct {
	A byteStruct
	B [32]byte
}

type bytes32 struct {
	A [32]byte
}

type largeStructB struct {
	A    bool `bitfield:"1"`
	B, C byte
	D    bStructThen32B `bitfield:"union"`
	E    bThenU16
	F    byte
	G    bytes32 `bitfield:"union"`
	H    [32]byte
	I    [55]byte
	J    uint16
}

type largeStructsAB struct {
	A largeStructA
	B largeStructB
}

type avu16 struct {
	A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V uint16
}

type fourUint32 struct {
	A, B, C, D uint32
}

type u32ThenBLU16 struct {
	A                               uint32
	B, C, D, E, F, G, H, I, J, K, L uint16
}

type u32Pair struct {
	A, B uint32
}

type avu16WAAByte struct {
	A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V uint16
	W, X, Y, Z, AA                                                   byte
}

type u16Then8U16 struct {
	A, B uint16
	C    [8]uint16
}

type b16U16 struct {
	A [16]byte
	B uint16
}

type simpleUnion struct {
	A uint8
	B uint16
}

type innerUnion struct {
	A uint8
	B simpleUnion `bitfield:"union"`
	C uint8
}

var tests = []struct {
	name string
	data any
	hex  string
}{
	{
		name: "mixed",
		data: mixed11{
			A: bitsABBBU16{
				A: false,
				B: true,
				C: 2,
				D: true,
				E: false,
				F: true,
				G: arrayThenByte{
					A: [3]byte{0xab, 0xcd, 0xef},
					B: 0x23,
				},
				H: 0,
				I: 0xDE,
				J: 0xCAFE,
			},
			B: bitBB832B19U16{
				A: true,
				B: 0xBC,
				C: [8]byte{0x56, 0x6f, 0x68, 0x77, 0x32, 0x7a, 0x65, 0x69},
				D: [33]byte{
					0x30, 0x52, 0x59, 0x69, 0x4a, 0x75, 0x71, 0x41, 0x39,
					0x67, 0x56, 0x55, 0x51, 0x48, 0x36, 0x6c, 0x54, 0x4e,
					0x63, 0x66, 0x56, 0x71, 0x64, 0x71, 0x69, 0x48, 0x72,
					0x32, 0x78, 0x38, 0x6b, 0x70, 0x74,
				},
				E: [19]byte{
					0x55, 0x4f, 0x4d, 0x68, 0x31, 0x76, 0x62, 0x6f, 0x41,
					0x6e, 0x53, 0x72, 0x70, 0x54, 0x77, 0x6d, 0x30, 0x4b,
					0x4e,
				},
				F: 0xBAAF,
			},
		},
		hex: "5aabcdef2300defeca01bc566f6877327a6569305259694a757141396756555148366c544e6366567164716948723278386b7074554f4d683176626f416e53727054776d304b4eafba",
	},
	{
		name: "mixed",
		data: struct {
			A [2]bits2B
			B bool `bitfield:"1"`
			C byte `bitfield:"7"`
			D bool `bitfield:"1"`
			E bool `bitfield:"1"`
			F bool `bitfield:"1"`
			G bool `bitfield:"1"`
			H bool `bitfield:"1"`
			I bool `bitfield:"1"`
			J bool `bitfield:"1"`
			K bool `bitfield:"1"`
		}{
			A: [2]bits2B{
				{
					A: 1,
					B: 0b11,
					C: true,
					D: false,
					E: false,
					F: true,
					G: true,
					H: 0xc3,
					I: 0x4e,
				},
				{
					A: 1,
					B: 0b10,
					C: false,
					D: true,
					E: false,
					F: true,
					G: true,
					H: 0x4d,
					I: 0x59,
				},
			},
			B: true,
			C: 0x1f,
			D: false,
			E: false,
			F: true,
			G: false,
			H: true,
			I: false,
			J: true,
			K: false,
		},
		hex: "cfc34ed54d593f54",
	},
	{
		name: "only bit flags",
		data: flags{
			A: true,
			B: false,
			C: false,
			D: true,
			E: false,
			F: true,
			G: true,
			H: false,
		},
		hex: "69",
	},
	{
		name: "mixed",
		data: bits8Int16Int16{
			A: false,
			B: false,
			C: true,
			D: false,
			E: false,
			F: false,
			G: false,
			H: true,
			I: 0x68f2 & 0x7FFF,
			J: 0xa661 & 0x7FFF,
		},
		hex: "84f2686126",
	},
	{
		name: "ArrayThenByte",
		data: arrayThenByte{
			A: [3]byte{0x12, 0x34, 0x5a},
			B: byte(0x04),
		},
		hex: "12345a04",
	},
	{
		name: "BitsFloats",
		data: bitsFloats{
			A: false,
			B: true,
			C: true,
			D: false,
			E: false,
			F: true,
			G: true,
			H: true,
			I: false,
			J: true,
			K: false,
			L: true,
			M: true,
			N: 0xf774,
			O: 0.17032479267946243,
			P: 17032.479267946243,
			Q: 3.235944344633299,
			R: 0x5118,
			S: 0xfff0,
			T: -10.528768872959308,
			U: 3.14,
		},
		hex: "e61a74f79f692e3ef5108546b6194f401851f0ffd67528c1c3f54840",
	},
	{
		name: "mixed",
		data: struct {
			A    bool `bitfield:"1"`
			B    bool `bitfield:"1"`
			C    bool `bitfield:"1"`
			D    bool `bitfield:"1"`
			E    bool `bitfield:"1"`
			F    bool `bitfield:"1"`
			G, H uint16
			I    byte
			J    uint16
		}{
			A: false,
			B: true,
			C: false,
			D: true,
			E: true,
			F: true,
			G: 0xdef6,
			H: 0xdcdc,
			I: 0x91,
			J: 0x09d9,
		},
		hex: "3af6dedcdc91d909",
	},
	{
		name: "mixed",
		data: struct {
			A int8
			B b16U16
		}{
			A: 0x19,
			B: b16U16{
				A: [16]byte{0xd1, 0x97, 0x9e, 0x6d, 0xf3, 0xac, 0x38, 0x31, 0x99, 0x6f, 0x30, 0xb1, 0x63, 0xbf, 0xa5, 0x27},
				B: 0x982d,
			},
		},
		hex: "19d1979e6df3ac3831996f30b163bfa5272d98",
	},
	{
		name: "mixed",
		data: struct {
			A structBB
			B structBB
			C uint16
			D uint32
			E byte
			F [32]byte
		}{
			A: structBB{
				A: 0x1a,
				B: 0x9d,
			},
			B: structBB{
				A: 0xe4,
				B: 0xcb,
			},
			C: 0xd5c9,
			D: 0xbeb1a16c,
			E: 0x51,
			F: [32]byte{
				0x7d, 0x93, 0x23, 0x0d, 0xff, 0xc1, 0x4b, 0xc8,
				0x0b, 0xc7, 0x0f, 0x5f, 0x6f, 0x9c, 0x65, 0xe5,
				0x78, 0x4e, 0x3d, 0x3f, 0xeb, 0xb7, 0xa8, 0xf4,
				0xb6, 0x53, 0x1d, 0xb0, 0x3c, 0x85, 0x0b, 0x9d,
			},
		},
		hex: "1a9de4cbc9d56ca1b1be517d93230dffc14bc80bc70f5f6f9c65e5784e3d3febb7a8f4b6531db03c850b9d",
	},
	{
		name: "mixed",
		data: struct {
			A avu16
			B fourUint32
			C u32ThenBLU16
			D u32Pair
			E avu16WAAByte
			F u16Then8U16
		}{
			A: avu16{
				A: 0x68c4,
				B: 0x4f35,
				C: 0x6f98,
				D: 0xe447,
				E: 0x91a9,
				F: 0xc1e7,
				G: 0x67d9,
				H: 0x1f85,
				I: 0x6cc1,
				J: 0x7e60,
				K: 0x877f,
				L: 0x4e4b,
				M: 0xda2b,
				N: 0x9ee6,
				O: 0x8a03,
				P: 0x6b35,
				Q: 0x39dd,
				R: 0xe2d5,
				S: 0xa06b,
				T: 0x513d,
				U: 0xeadb,
				V: 0x6a15,
			},
			B: fourUint32{
				A: 0x367eed0f,
				B: 0x54cc7069,
				C: 0xa2e2dcd1,
				D: 0x9309586d,
			},
			C: u32ThenBLU16{
				A: 0xfef7f27a,
				B: 0xb40d,
				C: 0x5910,
				D: 0x0f51,
				E: 0x4da8,
				F: 0xe90e,
				G: 0x07be,
				H: 0x5b5f,
				I: 0x549d,
				J: 0x0b4b,
				K: 0xe4bd,
				L: 0xdc10,
			},
			D: u32Pair{
				A: 0xe1ae9448,
				B: 0x4e4ff3a2,
			},
			E: avu16WAAByte{
				A:  0xfc86,
				B:  0xece1,
				C:  0x0424,
				D:  0x699a,
				E:  0x60ae,
				F:  0x0250,
				G:  0xf475,
				H:  0x680a,
				I:  0xcb93,
				J:  0x7409,
				K:  0xb9f6,
				L:  0x06f6,
				M:  0x307c,
				N:  0x1a5d,
				O:  0xa1e3,
				P:  0x97a1,
				Q:  0x45b5,
				R:  0xdd61,
				S:  0xbbb8,
				T:  0xa373,
				U:  0x747d,
				V:  0xe2a0,
				W:  0xb7,
				X:  0x5f,
				Y:  0x98,
				Z:  0x00,
				AA: 0x4e,
			},
			F: u16Then8U16{
				A: 0xe6d4,
				B: 0xea73,
				C: [8]uint16{
					0x0752,
					0x8640,
					0x6bca,
					0x712f,
					0x017e,
					0xa029,
					0x12f9,
					0x25dc,
				},
			},
		},
		hex: "c468354f986f47e4a991e7c1d967851fc16c607e7f874b4e2bdae69e038a356bdd39d5e26ba03d51dbea156a0fed7e366970cc54d1dce2a26d5809937af2f7fe0db41059510fa84d0ee9be075f5b9d544b0bbde410dc4894aee1a2f34f4e86fce1ec24049a69ae60500275f40a6893cb0974f6b9f6067c305d1ae3a1a197b54561ddb8bb73a37d74a0e2b75f98004ed4e673ea52074086ca6b2f717e0129a0f912dc25",
	},
	{
		name: "boolean fields",
		data: struct {
			A bool `bitfield:"1"`
			B bool `bitfield:"1"`
		}{
			A: true,
			B: false,
		},
		hex: "01",
	},
	{
		name: "boolean fields",
		data: struct {
			A int8 `bitfield:"5"`
			B int8 `bitfield:"3"`
		}{
			A: 0b10101,
			B: 0b011,
		},
		hex: "75",
	},
	{
		name: "slice of boolean",
		data: []struct {
			A int8 `bitfield:"5"`
			B int8 `bitfield:"3"`
		}{
			{A: 0b10101, B: 0b011},
			{A: 0b01110, B: 0b101},
		},
		hex: "75ae",
	},
	{
		name: "a regular value (byte)",
		data: byte(0xa4),
		hex:  "a4",
	},
	{
		name: "a regular value (uint16)",
		data: uint16(0xcafe),
		hex:  "feca",
	},
	{
		name: "a regular value (uint32)",
		data: uint32(0xceedcafe),
		hex:  "fecaedce",
	},
	{
		name: "a regular value (uint64)",
		data: uint64(0xabcdef123456),
		hex:  "563412efcdab0000",
	},
	{
		name: "struct with pointers",
		data: struct {
			A *struct{ Value byte }
		}{
			A: &struct{ Value byte }{Value: 0x54},
		},
		hex: "54",
	},
	{
		name: "signed byte",
		data: int8(-1),
		hex:  "ff",
	},
	{
		name: "union",
		data: largeStructsAB{
			A: largeStructA{
				A: uint16Pair{
					A: 0x4e4b,
					B: 0xda2b,
				},
				B: 0b10111,
				C: 0b011,
				D: arrayThenByte{
					A: [3]byte{0x12, 0x34, 0x5a},
					B: 0x04,
				},
				E: 0xdc,
				F: variousBitFields{
					A: false,
					B: true,
					C: false,
					D: false,
					E: 2,
					F: true,
					G: true,
				},
				G: arrayThenByte{
					A: [3]byte{0x21, 0x43, 0xa5},
					B: 0x40,
				},
				H: 0xfc,
				I: structBUint16{
					A: 0x39,
					B: 0x6b35,
				},
				J: 0x54,
				K: mixed2{
					B: mixed3{
						A: true,
						B: false,
						C: 0x60,
						D: [5]byte{
							0x6f, 0x30, 0xb1, 0x63,
							0xbf,
						},
						E: 0xc4,
					},
				},
				L: 0x0a,
				M: 0xe1ae,
			},
			B: largeStructB{
				A: false,
				B: 0xff,
				C: 0x25,
				D: bStructThen32B{
					A: byteStruct{
						A: 0x0b,
					},
				},
				E: bThenU16{
					A: 0x8f,
					B: 0xfc92,
				},
				F: 0x96,
				G: bytes32{
					A: [32]byte{
						0xa7, 0x6e, 0x67, 0x1c, 0x3c, 0x16, 0x7e, 0x1c,
						0xdd, 0xb5, 0xd8, 0x5a, 0x52, 0xe1, 0xd9, 0xf3,
						0x73, 0x60, 0x4d, 0x42, 0xb1, 0xcb, 0x56, 0x1f,
						0x6f, 0x76, 0x8c, 0xc5, 0x38, 0x7b, 0xe1, 0x25,
					},
				},
				H: [32]byte{
					0xce, 0x97, 0xba, 0xb5, 0x32, 0x16, 0x5a, 0x5d,
					0x22, 0x96, 0xc7, 0x12, 0x91, 0x55, 0xb3, 0xb5,
					0x0e, 0x45, 0xb0, 0xfa, 0xd6, 0x3b, 0x89, 0x49,
					0xde, 0x08, 0xdd, 0xe4, 0xdb, 0xe5, 0x50, 0x5d,
				},
				I: [55]byte{
					0x0d, 0x2c, 0xc1, 0x5a, 0xf5, 0x68, 0x99, 0xaa,
					0x39, 0xe2, 0x9c, 0xd5, 0x0d, 0xb2, 0xc6, 0xa0,
					0x9e, 0xfb, 0xde, 0x4b, 0xb0, 0x1f, 0x25, 0x00,
					0xa8, 0x57, 0xf5, 0x95, 0xc9, 0xb9, 0x9c, 0x6b,
					0xf2, 0x6d, 0xb6, 0xf2, 0xb6, 0x3f, 0x82, 0xf3,
					0xe7, 0xf0, 0x3e, 0x29, 0x8b, 0x5f, 0x35, 0x03,
					0x1f, 0xb1, 0xfe, 0x00, 0x5f, 0xee, 0xf1,
				},
				J: 0x56a6,
			},
		},
		hex: "4b4e2bda7712345a04dce22143a540fc39356b5401606f30b163bfc4000000000000000aaee100ff250b000000000000000000000000000000000000000000000000000000000000008f92fc96a76e671c3c167e1cddb5d85a52e1d9f373604d42b1cb561f6f768cc5387be125ce97bab532165a5d2296c7129155b3b50e45b0fad63b8949de08dde4dbe5505d0d2cc15af56899aa39e29cd50db2c6a09efbde4bb01f2500a857f595c9b99c6bf26db6f2b63f82f3e7f03e298b5f35031fb1fe005feef1a656",
	},
	{
		name: "union",
		data: largeStructsAB{},
		hex:  "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	},
	{
		name: "nested union",
		data: struct {
			A byte
			B innerUnion `bitfield:"union"`
			C byte
		}{
			A: 0xAA,
			B: innerUnion{
				A: 0,
				B: simpleUnion{A: 0, B: 0x3344},
				C: 0,
			},
			C: 0xBB,
		},
		hex: "aa4433bb",
	},
	{
		name: "large nested union",
		data: unionOfUnions{},
		hex:  "00000000000000000000000000000000",
	},
}

func TestSerializeStruct(t *testing.T) {
	t.Parallel()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			serialized, err := Marshal(tc.data)
			if err != nil {
				t.Fatalf("SerializeStruct error: %v", err)
			}

			gotHex := hex.EncodeToString(serialized)
			if gotHex != tc.hex {
				t.Fatalf("For type %T, want hex %v, got %v", tc.data, tc.hex, gotHex)
			}
		})
	}
}

func TestSizeof(t *testing.T) {
	t.Parallel()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			serialized, err := Marshal(tc.data)
			if err != nil {
				t.Fatalf("TestSizeof error: %v", err)
			}

			wantSize := len(serialized)

			gotSize, err := Sizeof(tc.data)
			if err != nil {
				t.Fatalf("TestSizeof error: %v", err)
			}

			if gotSize != wantSize {
				t.Fatalf("For type %T, want size %v, got %v", tc.data, wantSize, gotSize)
			}
		})
	}
}

func TestDeserializeStruct(t *testing.T) {
	t.Parallel()

	for _, testCase := range tests {
		tc := testCase // Copy value for modification
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotHex, err := MarshalToHexString(tc.data)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}

			if gotHex != tc.hex {
				t.Fatalf("For type %T, want hex %v, got %v", tc.data, tc.hex, gotHex)
			}

			// Create a zero value of the same type as tc.B and unmarshal into it.
			out := reflect.New(reflect.TypeOf(tc.data)).Interface()
			_, err = UnmarshalHexString(tc.hex, out, 0)
			if err != nil {
				t.Fatalf("UnmarshalHexString error: %v", err)
			}

			gotHexAfterMarshal, err := MarshalToHexString(out)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			if gotHexAfterMarshal != tc.hex {
				t.Fatalf("For type %T, want hex %v, got %v", tc.data, tc.hex, gotHexAfterMarshal)
			}
		})
	}
}
