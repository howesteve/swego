package worker

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Funcs) DecodeMsg(dc *msgp.Reader) (err error) {
	var zbai uint32
	zbai, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zbai) {
		(*z) = (*z)[:zbai]
	} else {
		(*z) = make(Funcs, zbai)
	}
	for zbzg := range *z {
		(*z)[zbzg], err = dc.ReadString()
		if err != nil {
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Funcs) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zcmr := range z {
		err = en.WriteString(z[zcmr])
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Funcs) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zcmr := range z {
		o = msgp.AppendString(o, z[zcmr])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Funcs) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zwht uint32
	zwht, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(zwht) {
		(*z) = (*z)[:zwht]
	} else {
		(*z) = make(Funcs, zwht)
	}
	for zajw := range *z {
		(*z)[zajw], bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Funcs) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zhct := range z {
		s += msgp.StringPrefixSize + len(z[zhct])
	}
	return
}
