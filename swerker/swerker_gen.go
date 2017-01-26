package swerker

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Call) DecodeMsg(dc *msgp.Reader) (err error) {
	var zbzg uint32
	zbzg, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zbzg != 3 {
		err = msgp.ArrayError{Wanted: 3, Got: zbzg}
		return
	}
	var zbai uint32
	zbai, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap(z.Ctx) >= int(zbai) {
		z.Ctx = (z.Ctx)[:zbai]
	} else {
		z.Ctx = make([]*CtxCall, zbai)
	}
	for zxvk := range z.Ctx {
		if dc.IsNil() {
			err = dc.ReadNil()
			if err != nil {
				return
			}
			z.Ctx[zxvk] = nil
		} else {
			if z.Ctx[zxvk] == nil {
				z.Ctx[zxvk] = new(CtxCall)
			}
			var zcmr uint32
			zcmr, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if zcmr != 2 {
				err = msgp.ArrayError{Wanted: 2, Got: zcmr}
				return
			}
			z.Ctx[zxvk].Func, err = dc.ReadUint8()
			if err != nil {
				return
			}
			err = z.Ctx[zxvk].Args.DecodeMsg(dc)
			if err != nil {
				return
			}
		}
	}
	z.Func, err = dc.ReadUint8()
	if err != nil {
		return
	}
	err = z.Args.DecodeMsg(dc)
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Call) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 3
	err = en.Append(0x93)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Ctx)))
	if err != nil {
		return
	}
	for zxvk := range z.Ctx {
		if z.Ctx[zxvk] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			// array header, size 2
			err = en.Append(0x92)
			if err != nil {
				return err
			}
			err = en.WriteUint8(z.Ctx[zxvk].Func)
			if err != nil {
				return
			}
			err = z.Ctx[zxvk].Args.EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	err = en.WriteUint8(z.Func)
	if err != nil {
		return
	}
	err = z.Args.EncodeMsg(en)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Call) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 3
	o = append(o, 0x93)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Ctx)))
	for zxvk := range z.Ctx {
		if z.Ctx[zxvk] == nil {
			o = msgp.AppendNil(o)
		} else {
			// array header, size 2
			o = append(o, 0x92)
			o = msgp.AppendUint8(o, z.Ctx[zxvk].Func)
			o, err = z.Ctx[zxvk].Args.MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	o = msgp.AppendUint8(o, z.Func)
	o, err = z.Args.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Call) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zajw uint32
	zajw, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zajw != 3 {
		err = msgp.ArrayError{Wanted: 3, Got: zajw}
		return
	}
	var zwht uint32
	zwht, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap(z.Ctx) >= int(zwht) {
		z.Ctx = (z.Ctx)[:zwht]
	} else {
		z.Ctx = make([]*CtxCall, zwht)
	}
	for zxvk := range z.Ctx {
		if msgp.IsNil(bts) {
			bts, err = msgp.ReadNilBytes(bts)
			if err != nil {
				return
			}
			z.Ctx[zxvk] = nil
		} else {
			if z.Ctx[zxvk] == nil {
				z.Ctx[zxvk] = new(CtxCall)
			}
			var zhct uint32
			zhct, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if zhct != 2 {
				err = msgp.ArrayError{Wanted: 2, Got: zhct}
				return
			}
			z.Ctx[zxvk].Func, bts, err = msgp.ReadUint8Bytes(bts)
			if err != nil {
				return
			}
			bts, err = z.Ctx[zxvk].Args.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		}
	}
	z.Func, bts, err = msgp.ReadUint8Bytes(bts)
	if err != nil {
		return
	}
	bts, err = z.Args.UnmarshalMsg(bts)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Call) Msgsize() (s int) {
	s = 1 + msgp.ArrayHeaderSize
	for zxvk := range z.Ctx {
		if z.Ctx[zxvk] == nil {
			s += msgp.NilSize
		} else {
			s += 1 + msgp.Uint8Size + z.Ctx[zxvk].Args.Msgsize()
		}
	}
	s += msgp.Uint8Size + z.Args.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *CtxCall) DecodeMsg(dc *msgp.Reader) (err error) {
	var zcua uint32
	zcua, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zcua != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zcua}
		return
	}
	z.Func, err = dc.ReadUint8()
	if err != nil {
		return
	}
	err = z.Args.DecodeMsg(dc)
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *CtxCall) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 2
	err = en.Append(0x92)
	if err != nil {
		return err
	}
	err = en.WriteUint8(z.Func)
	if err != nil {
		return
	}
	err = z.Args.EncodeMsg(en)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *CtxCall) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 2
	o = append(o, 0x92)
	o = msgp.AppendUint8(o, z.Func)
	o, err = z.Args.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *CtxCall) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zxhx uint32
	zxhx, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zxhx != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zxhx}
		return
	}
	z.Func, bts, err = msgp.ReadUint8Bytes(bts)
	if err != nil {
		return
	}
	bts, err = z.Args.UnmarshalMsg(bts)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *CtxCall) Msgsize() (s int) {
	s = 1 + msgp.Uint8Size + z.Args.Msgsize()
	return
}
