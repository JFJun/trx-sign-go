package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"strings"
	"time"
)

func MinimizeSilce(src []byte) []byte {
	dest := make([]byte, len(src))
	copy(dest, src)
	return dest
}

func BytesToHex(byteSlice []byte) string {
	return hex.EncodeToString(byteSlice)
}

func HexToBytes(s string) []byte {
	h, _ := hex.DecodeString(s)
	return h
}

func Uint32ToBytes(src uint32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, src)
	return buf.Bytes()
}

func BytesToUint32(src []byte) uint32 {
	return binary.LittleEndian.Uint32(src)
}

func CurrentTimestamp() uint32 {
	return uint32(time.Now().Unix())
}

type Times []uint32

func (t Times) Len() int           { return len(t) }
func (t Times) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Times) Less(i, j int) bool { return t[i] < t[j] }

func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}

// ZeroMemory erases a slice memory
func ZeroMemory(s interface{}) {
	if v, ok := s.([]interface{}); ok {
		for i := range v {
			v[i] = 0
		}
	}
}

func GetRandomBytes(n int) []byte {
	buf := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		panic("reading rand failed: " + err.Error())
	}
	return buf
}

/*
添加txid标识
*/
func AddTxidFlag(txid string,flag int)string{
	return fmt.Sprintf("%s&%d",txid,flag)
}

/*
去除txid 的flag
*/
func DelTxidFlag(txid_flag string)string{
	ss:=strings.Split(txid_flag,"&")
	if len(ss)==2{
		return ss[0]
	}
	return txid_flag
}