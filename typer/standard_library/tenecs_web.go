package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_web = packageWith(
	withStruct("WebApp", &tenecs_web_WebApp, tenecs_web_WebApp_Fields...),
	withStruct("HtmlElement", &tenecs_web_HtmlElement, tenecs_web_HtmlElement_Fields...),
)

var tenecs_web_WebApp = types.KnownType{
	Package: "tenecs.web",
	Name:    "WebApp",
	Generics: []types.VariableType{
		&types.TypeArgument{
			Name: "Model",
		},
		&types.TypeArgument{
			Name: "Event",
		},
	},
}

var tenecs_web_WebApp_Fields = []func(fields *StructWithFields){
	structField("init", &types.Function{
		Arguments: []types.FunctionArgument{},
		ReturnType: &types.TypeArgument{
			Name: "Model",
		},
	}),
	structField("update", &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name: "model",
				VariableType: &types.TypeArgument{
					Name: "Model",
				},
			},
			{
				Name: "event",
				VariableType: &types.TypeArgument{
					Name: "Event",
				},
			},
		},
		ReturnType: &types.TypeArgument{
			Name: "Model",
		},
	}),
	structField("view", &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name: "model",
				VariableType: &types.TypeArgument{
					Name: "Model",
				},
			},
		},
		ReturnType: &tenecs_web_HtmlElement,
	}),
}

var tenecs_web_HtmlElement = types.KnownType{
	Package: "tenecs.web",
	Name:    "HtmlElement",
}

var tenecs_web_HtmlElement_Fields = []func(fields *StructWithFields){
	structField("name", types.String()),
	structField("properties", types.List(types.Void())),
	structField("children", types.List(&tenecs_web_HtmlElement)),
}
