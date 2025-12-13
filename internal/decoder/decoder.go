package decoder

import (
	"gopkg.in/yaml.v3"
)

type Decoder struct {
	data []byte
}

func New(data []byte) *Decoder {
	return &Decoder{data: data}
}

func (d *Decoder) Decode(v interface{}) error {
	return yaml.Unmarshal(d.data, v)
}

func (d *Decoder) DecodeAll(fn func(interface{}) error) error {
	decoder := yaml.NewDecoder(yaml.NewDecoder(nil))
	
	for {
		var doc interface{}
		err := decoder.Decode(&doc)
		if err != nil {
			break
		}
		
		if err := fn(doc); err != nil {
			return err
		}
	}
	
	return nil
}
