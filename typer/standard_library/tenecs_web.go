package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_web = packageWith(
	withStruct(Tenecs_web_CssUrl),
	withStruct(Tenecs_web_WebApp),
	withStruct(Tenecs_web_HtmlElement),
	withStruct(Tenecs_web_HtmlElementProperty),
)

var Tenecs_web_CssUrl = structWithFields("CssUrl", tenecs_web_CssUrl, tenecs_web_CssUrl_Fields...)

var tenecs_web_CssUrl = types.Struct(
	"tenecs.web",
	"CssUrl",
	[]string{},
)

var tenecs_web_CssUrl_Fields = []func(fields *StructWithFields){
	structField("url", types.String()),
}

var Tenecs_web_WebApp = structWithFields("WebApp", tenecs_web_WebApp, tenecs_web_WebApp_Fields...)

var tenecs_web_WebApp = types.Struct(
	"tenecs.web",
	"WebApp",
	[]string{"Model", "Event"},
)

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
		ReturnType: tenecs_web_HtmlElement,
	}),
	structField("external", &types.OrVariableType{
		Elements: []types.VariableType{
			types.Void(),
			types.List(tenecs_web_CssUrl),
		},
	}),
}

var Tenecs_web_HtmlElement = structWithFields("HtmlElement", tenecs_web_HtmlElement, tenecs_web_HtmlElement_Fields...)

var tenecs_web_HtmlElement = types.Struct(
	"tenecs.web",
	"HtmlElement",
	[]string{"Event"},
)

var tenecs_web_HtmlElement_Fields = []func(fields *StructWithFields){
	structField("name", types.String()),
	structField("properties", types.List(tenecs_web_HtmlElementProperty)),
	structField("children", &types.OrVariableType{
		Elements: []types.VariableType{
			types.String(),
			types.List(tenecs_web_HtmlElement),
		},
	}),
}

var Tenecs_web_HtmlElementProperty = structWithFields("HtmlElementProperty", tenecs_web_HtmlElementProperty, tenecs_web_HtmlElementProperty_Fields...)

var tenecs_web_HtmlElementProperty = types.Struct(
	"tenecs.web",
	"HtmlElementProperty",
	[]string{"Event"},
)

var tenecs_web_HtmlElementProperty_Fields = []func(fields *StructWithFields){
	structField("name", types.String()),
	structField("value", &types.OrVariableType{
		Elements: []types.VariableType{
			types.String(),
			&types.Function{
				Arguments:  []types.FunctionArgument{},
				ReturnType: &types.TypeArgument{Name: "Event"},
			},
		},
	}),
}
