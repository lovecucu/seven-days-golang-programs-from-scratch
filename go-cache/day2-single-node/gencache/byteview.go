package gencache

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

// 一个值类型，不是指针
type ByteView struct {
	b []byte
	s string
}

// 获取长度
func (v ByteView) Len() int {
	if v.b != nil {
		return len(v.b)
	}
	return len(v.s)
}

// 获取[]byte的一份copy
func (v ByteView) ByteSlice() []byte {
	if v.b != nil {
		return cloneBytes(v.b)
	}
	return []byte(v.s)
}

// 数据String化
func (v ByteView) String() string {
	if v.b != nil {
		return string(v.b)
	}
	return v.s
}

// At 获取索引为i的byte
func (v ByteView) At(i int) byte {
	if v.b != nil {
		return v.b[i]
	}
	return v.s[i]
}

// Slice 获取from, to截取成新的ByteView
func (v ByteView) Slice(from, to int) ByteView {
	if v.b != nil {
		return ByteView{b: v.b[from:to]}
	}
	return ByteView{s: v.s[from:to]}
}

// SliceFrom 从from处截取成ByteView
func (v ByteView) SliceFrom(from int) ByteView {
	if v.b != nil {
		return ByteView{b: v.b[from:]}
	}
	return ByteView{s: v.s[from:]}
}

// Copy 把Byte拷贝进目标[]byte，同时返回拷贝的字符数
func (v ByteView) Copy(dest []byte) int {
	if v.b != nil {
		return copy(dest, v.b)
	}
	return copy(dest, v.s)
}

// Equal 根据值校验两个ByteView是否相等
func (v ByteView) Equal(b2 ByteView) bool {
	if b2.b == nil {
		return v.EqualString(b2.s)
	}
	return v.EqualBytes(b2.b)
}

// EqualString 对比string
func (v ByteView) EqualString(s string) bool {
	if v.b == nil {
		return v.s == s
	}
	l := v.Len()
	if len(s) != l {
		return false
	}
	for i, bi := range v.b {
		if s[i] != bi {
			return false
		}
	}
	return true
}

// EqualBytes 对比byte[]
func (v ByteView) EqualBytes(b []byte) bool {
	if v.b != nil {
		return bytes.Equal(v.b, b)
	}
	l := v.Len()
	if l != len(b) {
		return false
	}
	for i, bi := range b {
		if bi != v.s[i] {
			return false
		}
	}
	return true
}

// Reader 返回一个io.ReadSeeker
func (v ByteView) Reader() io.ReadSeeker {
	if v.b != nil {
		return bytes.NewReader(v.b)
	}
	return strings.NewReader(v.s)
}

// ReadAt 实现io.ReaderAt
func (v ByteView) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errors.New("view: invalid offset")
	}
	if off >= int64(v.Len()) {
		return 0, io.EOF
	}
	n = v.SliceFrom(int(off)).Copy(p)
	if n < len(p) {
		err = io.EOF
	}
	return
}

// WriteTo 实现io.WriterTo
func (v ByteView) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	if v.b != nil {
		m, err = w.Write(v.b)
	} else {
		m, err = io.WriteString(w, v.s)
	}
	if err == nil && m < v.Len() {
		err = io.ErrShortWrite
	}
	n = int64(m)
	return
}
