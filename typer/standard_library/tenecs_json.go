package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_json = packageWith(
	withFunction("toJson", tenecs_json_toJson),
)

var tenecs_json_toJson = &types.Function{
	Generics: []string{
		"T",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "input",
			VariableType: &types.TypeArgument{
				Name: "T",
			},
		},
	},
	ReturnType: types.String(),
}
