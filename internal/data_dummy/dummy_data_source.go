package data_dummy

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	//autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	//"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
	//"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func GetDataDummyLowLevelSchema() *tfprotov5.Schema {

	// return map[string]tftypes.Type{
	// 	"id":                tftypes.String,
	// 	"dynamic_attribute": tftypes.DynamicPseudoType,
	// 	"dynamic_block":     tftypes.DynamicPseudoType,
	// 	"regular_att":       tftypes.String,
	// 	"regular_block": tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	// 		"bar": tftypes.String,
	// 		"foo": tftypes.Bool,
	// 	}},
	// }
	schema := tfprotov5.Schema{

		Version: 1,
		Block: &tfprotov5.SchemaBlock{
			Version: 1,
			Attributes: []*tfprotov5.SchemaAttribute{
				{
					Name:     "id",
					Type:     tftypes.String,
					Computed: true,
				},
				{
					Name:     "regular_attribute",
					Type:     tftypes.String,
					Optional: true,
				},
				{
					Name:     "dynamic_attribute",
					Type:     tftypes.DynamicPseudoType,
					Optional: true,
				},
			},
			BlockTypes: []*tfprotov5.SchemaNestedBlock{
				{
					TypeName: "dynamic_block",
					Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
					Block: &tfprotov5.SchemaBlock{
						Attributes: []*tfprotov5.SchemaAttribute{
							{
								Name:     "bar",
								Type:     tftypes.DynamicPseudoType,
								Optional: true,
							},
							{
								Name:     "foo",
								Type:     tftypes.Number,
								Optional: true,
							},
						},
					},
				},
				{
					TypeName: "regular_block",
					Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
					Block: &tfprotov5.SchemaBlock{
						Attributes: []*tfprotov5.SchemaAttribute{
							{
								Name:     "bar",
								Type:     tftypes.String,
								Optional: true,
							},
							{
								Name:     "foo",
								Type:     tftypes.Number,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
	return &schema

}

// dataSourceBlueprintConfig
type dataSourceDummy struct {
	// this struct can carry important elements for usage in the lifecycle, in this case the sdk is not needed
	// because we are not accessing anything from it
	// Leaving this as reference for the future
	//nolint:golint,unused
	autocloudClient interface{}
}

func NewDataSourceDummy() tfprotov5.DataSourceServer {
	return dataSourceDummy{}
}

func (d dataSourceDummy) ReadDataSource(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	objTypeDef, ok := GetDataDummyLowLevelSchema().ValueType().(tftypes.Object)
	if !ok {
		return nil, errors.New("Cant get lowlevel attributes")
	}

	/*
		// used as backup to read the dynamic value
		config := map[string]tftypes.Type{
			"id":                tftypes.String,
			"dynamic_attribute": tftypes.DynamicPseudoType,
			"dynamic_block": tftypes.List{
				ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
					"bar": tftypes.DynamicPseudoType,
					"foo": tftypes.Number,
				}},
			},
			"regular_attribute": tftypes.String,
			"regular_block": tftypes.List{
				ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
					"bar": tftypes.String,
					"foo": tftypes.Number,
				}},
			},
		}
		//objTypeDef.AttributeTypes["dynamic_block"] = objTypeDef.AttributeTypes["regular_block"]
	*/
	values, err := req.Config.Unmarshal(tftypes.Object{
		AttributeTypes: objTypeDef.AttributeTypes, // or config
	})
	if err != nil {
		return nil, fmt.Errorf("Cant unmarshall config input, %v", err.Error())
	}

	var input DataDummyState
	err = values.As(&input)
	if err != nil {
		return nil, fmt.Errorf("Cant convert config input, %v", err.Error())
	}

	outputObjectType := map[string]tftypes.Type{
		"bar": tftypes.String,
		"foo": tftypes.Number,
	}

	state, err := tfprotov5.NewDynamicValue(
		objTypeDef,
		tftypes.NewValue(tftypes.Object{
			AttributeTypes: objTypeDef.AttributeTypes, // or config
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
		}),
	)

	if err != nil {
		return nil, err
	}
	fmt.Println(state)

	return &tfprotov5.ReadDataSourceResponse{
		State: &state,
	}, nil
}

func (d dataSourceDummy) ValidateDataSourceConfig(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	return &tfprotov5.ValidateDataSourceConfigResponse{}, nil
}
