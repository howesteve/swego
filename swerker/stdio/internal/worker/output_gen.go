package worker

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ErrorMap) DecodeMsg(dc *msgp.Reader) (err error) {
	var zajw uint32
	zajw, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	if (*z) == nil && zajw > 0 {
		(*z) = make(ErrorMap, zajw)
	} else if len((*z)) > 0 {
		for key, _ := range *z {
			delete((*z), key)
		}
	}
	for zajw > 0 {
		zajw--
		var zbai string
		var zcmr string
		zbai, err = dc.ReadString()
		if err != nil {
			return
		}
		zcmr, err = dc.ReadString()
		if err != nil {
			return
		}
		(*z)[zbai] = zcmr
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z ErrorMap) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zwht, zhct := range z {
		err = en.WriteString(zwht)
		if err != nil {
			return
		}
		err = en.WriteString(zhct)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z ErrorMap) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, uint32(len(z)))
	for zwht, zhct := range z {
		o = msgp.AppendString(o, zwht)
		o = msgp.AppendString(o, zhct)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ErrorMap) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zlqf uint32
	zlqf, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	if (*z) == nil && zlqf > 0 {
		(*z) = make(ErrorMap, zlqf)
	} else if len((*z)) > 0 {
		for key, _ := range *z {
			delete((*z), key)
		}
	}
	for zlqf > 0 {
		var zcua string
		var zxhx string
		zlqf--
		zcua, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			return
		}
		zxhx, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			return
		}
		(*z)[zcua] = zxhx
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ErrorMap) Msgsize() (s int) {
	s = msgp.MapHeaderSize
	if z != nil {
		for zdaf, zpks := range z {
			_ = zpks
			s += msgp.StringPrefixSize + len(zdaf) + msgp.StringPrefixSize + len(zpks)
		}
	}
	return
}
