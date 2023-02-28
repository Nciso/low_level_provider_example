package data_dummy

import (
	"math/big"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type Object struct {
	Foo *big.Float
	Bar string
}

type DataDummyState struct {
	Id               string
	DynamicAttribute interface{}
	RegularAttribute string
	RegularBlock     []Object
	DynamicBlock     []Object
}

func (in *DataDummyState) FromTerraform5Value(val tftypes.Value) error {
	v := map[string]tftypes.Value{}
	err := val.As(&v)
	if err != nil {
		return err
	}

	err = v["id"].As(&in.Id)
	if err != nil {
		return err
	}

	err = v["regular_attribute"].As(&in.RegularAttribute)
	if err != nil {
		return err
	}
	return nil
}
