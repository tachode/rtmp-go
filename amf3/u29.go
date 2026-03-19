package amf3

import (
	"fmt"
	"io"
)

// Maximum value that can be encoded as a U29 (2^29 - 1)
const MaxU29 = 0x1FFFFFFF

// readU29 reads a variable-length unsigned 29-bit integer from the reader.
//
// Encoding (from amf3_spec_05_05_08.pdf §1.3.1):
//
//	0x00000000 - 0x0000007F : 0xxxxxxx                                (1 byte,  7 bits)
//	0x00000080 - 0x00003FFF : 1xxxxxxx 0xxxxxxx                      (2 bytes, 14 bits)
//	0x00004000 - 0x001FFFFF : 1xxxxxxx 1xxxxxxx 0xxxxxxx             (3 bytes, 21 bits)
//	0x00200000 - 0x1FFFFFFF : 1xxxxxxx 1xxxxxxx 1xxxxxxx xxxxxxxx    (4 bytes, 29 bits)
func readU29(r io.Reader) (uint32, error) {
	var result uint32
	var buf [1]byte

	for i := 0; i < 3; i++ {
		_, err := io.ReadFull(r, buf[:])
		if err != nil {
			return 0, err
		}
		result = (result << 7) | uint32(buf[0]&0x7F)
		if buf[0]&0x80 == 0 {
			return result, nil
		}
	}

	// Fourth byte uses all 8 bits
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return 0, err
	}
	result = (result << 8) | uint32(buf[0])
	return result, nil
}

// writeU29 writes a variable-length unsigned 29-bit integer to the writer.
func writeU29(w io.Writer, value uint32) error {
	if value > MaxU29 {
		return fmt.Errorf("U29 value %d exceeds maximum %d", value, MaxU29)
	}

	switch {
	case value < 0x80:
		_, err := w.Write([]byte{byte(value)})
		return err
	case value < 0x4000:
		_, err := w.Write([]byte{
			byte(value>>7) | 0x80,
			byte(value & 0x7F),
		})
		return err
	case value < 0x200000:
		_, err := w.Write([]byte{
			byte(value>>14) | 0x80,
			byte(value>>7) | 0x80,
			byte(value & 0x7F),
		})
		return err
	default:
		_, err := w.Write([]byte{
			byte(value>>22) | 0x80,
			byte(value>>15) | 0x80,
			byte(value>>8) | 0x80,
			byte(value),
		})
		return err
	}
}
