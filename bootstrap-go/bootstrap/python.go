package bootstrap

import (
	"fmt"
	"io"
	"reflect"

	"github.com/whyrusleeping/cbor/go"
)

const CBOR_TAG = 0xef

type PyObject interface {
	GetPyObjectName() string
}

func ToCBOR(obj PyObject, w io.Writer, enc *cbor.Encoder) error {
	_, err := w.Write([]byte{0xd8, CBOR_TAG})
	if err != nil {
		return err
	}
	m := map[string]interface{}{}
	m["py/object"] = obj.GetPyObjectName()
	v := reflect.ValueOf(obj).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("pyobj")
		if tag != "" {
			m[tag] = v.Field(i).Interface()
		}
	}
	return enc.Encode(m)
}

type PyObjectDecoder struct{}

func (self *PyObjectDecoder) GetTag() uint64 {
	return CBOR_TAG
}

func (self *PyObjectDecoder) DecodeTarget() interface{} {
	return make(map[string]interface{})
}

func (self *PyObjectDecoder) PostDecode(v interface{}) (interface{}, error) {
	m := v.(map[string]interface{})
	var res interface{}
	if m["py/object"] == "golem.network.p2p.node.Node" {
		res = &Node{}
	} else {
		return nil, fmt.Errorf("Unsupported py/object %s", m["py/object"])
	}

	elem := reflect.ValueOf(res).Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)
		val := elem.Field(i)
		tag := field.Tag.Get("pyobj")
		if tag != "" {
			if vv, ok := m[tag]; ok && vv != nil {
				val.Set(reflect.ValueOf(vv))
			}
		}
	}
	return res, nil
}

type Node struct {
	NodeName     string        `pyobj:"node_name"`
	Key          string        `pyobj:"key"`
	PrvPort      uint64        `pyobj:"prv_port"`
	PubPort      uint64        `pyobj:"pub_port"`
	P2pPrvPort   uint64        `pyobj:"p2p_prv_port"`
	P2pPubPort   uint64        `pyobj:"p2p_pub_port"`
	PrvAddr      string        `pyobj:"prv_addr"`
	PubAddr      string        `pyobj:"pub_addr"`
	PrvAddresses []interface{} `pyobj:"prv_addresses"`
	NatType      string        `pyobj:"nat_type"`
}

func (self *Node) GetPyObjectName() string {
	return "golem.network.p2p.node.Node"
}

func (self *Node) ToCBOR(w io.Writer, enc *cbor.Encoder) error {
	return ToCBOR(self, w, enc)
}
