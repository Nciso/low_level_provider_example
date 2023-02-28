package data_dummy_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/Nciso/low_level_provider_example/internal/acctest"
	"github.com/Nciso/low_level_provider_example/internal/data_dummy"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBlueprintConfig_basic(t *testing.T) {
	var stateInGoStruct data_dummy.DataDummyState
	resourceName := "data.provider_dummy.test"
	experimental := true
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acctest.CreateMuxFactories(experimental),
		Steps: []resource.TestStep{
			{
				Config: testAccDummyData(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDummyConfigExist(resourceName, &stateInGoStruct),
					resource.TestCheckResourceAttrSet(
						resourceName, "regular_block.0.%"),
					resource.TestCheckResourceAttrSet(
						resourceName, "dynamic_attribute"),
					resource.TestCheckResourceAttrSet(
						resourceName, "regular_attribute"),
					resource.TestCheckResourceAttrSet(
						resourceName, "dynamic_block.0.%"),
				),
			},
		},
	})
}

func testAccDummyData() string {
	return `

	data "provider_dummy" "test" {
		dynamic_attribute = "hello"
		regular_attribute = "bye"
		regular_block {
			bar = "bar"
			foo = 3
		}
		
		dynamic_block {
			bar = "bar"
			foo = 4
		}
		
	}
`
}

func testAccCheckDummyConfigExist(resourceName string, dummyConfig *data_dummy.DataDummyState) resource.TestCheckFunc {
	// this is just to explore the state
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		rawConf := rs.Primary.Attributes
		fmt.Println(rawConf)
		// err := json.Unmarshal([]byte(rawConf), dummyConfig)
		// if err != nil {
		// 	return fmt.Errorf("not a valid data config: %s", rawConf)
		// }
		return nil
	}
}

func TestDecodeMSg(t *testing.T) {

	type testCase struct {
		decPack []uint8
	}

	testCases := map[string]testCase{
		"decoded_with_dpt_block": {
			decPack: []uint8("\x85\xb1dynamic_attribute\x92\xc4\b\"string\"\xa5hello\xaddynamic_block\x92\xc46[\"tuple\",[[\"object\",{\"bar\":\"string\",\"foo\":\"number\"}]]]\x91\x82\xa3bar\xa3bar\xa3foo\x04\xa2id\xc0\xb1regular_attribute\xa3bye\xadregular_block\x91\x82\xa3bar\xa3bar\xa3foo\x03"),
		},
		"decoded_without_dpt_block": {
			decPack: []uint8("\x85\xb1dynamic_attribute\x92\xc4\b\"string\"\xa5hello\xaddynamic_block\x91\x82\xa3bar\xa3bar\xa3foo\x04\xa2id\xc0\xb1regular_attribute\xa3bye\xadregular_block\x91\x82\xa3bar\xa3bar\xa3foo\x03"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {

			objTypeDef, ok := data_dummy.GetDataDummyLowLevelSchema().ValueType().(tftypes.Object)
			if !ok {
				t.Fatalf("cant get schema")
			}
			typ := tftypes.Object{
				AttributeTypes: objTypeDef.AttributeTypes,
			}
			// internal function used
			val, err := tftypes.ValueFromMsgPack(test.decPack, typ)
			if err != nil {
				t.Fatalf("fail decoding, testName: %s error -> %s", name, err.Error())
			}
			fmt.Println(val.String())
		})
	}
}

func TestDecodeValue(t *testing.T) {

	objTypeDef, ok := data_dummy.GetDataDummyLowLevelSchema().ValueType().(tftypes.Object)
	if !ok {
		t.Fatalf("cant get schema")
	}
	outputObjectType := map[string]tftypes.Type{
		"bar": tftypes.String,
		"foo": tftypes.Number,
	}

	tfval := tftypes.NewValue(tftypes.Object{
		AttributeTypes: objTypeDef.AttributeTypes,
	}, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, strconv.FormatInt(time.Now().Unix(), 10)),
		"dynamic_block": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: outputObjectType}},

			[]tftypes.Value{
				tftypes.NewValue(tftypes.Object{AttributeTypes: outputObjectType},
					map[string]tftypes.Value{
						"bar": tftypes.NewValue(tftypes.String, "hello world regular"),
						"foo": tftypes.NewValue(tftypes.Number, 101),
					},
				),
			},
		),
		"regular_block": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: outputObjectType}},

			[]tftypes.Value{
				tftypes.NewValue(tftypes.Object{AttributeTypes: outputObjectType},
					map[string]tftypes.Value{
						"bar": tftypes.NewValue(tftypes.String, "hello world regular"),
						"foo": tftypes.NewValue(tftypes.Number, 101),
					},
				),
			},
		),

		"regular_attribute": tftypes.NewValue(tftypes.String, "this is regular string"),
		"dynamic_attribute": tftypes.NewValue(tftypes.String, "this is dynamic string"),
	})
	dynamicValue, err := tfprotov5.NewDynamicValue(objTypeDef, tfval)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(dynamicValue)

}
