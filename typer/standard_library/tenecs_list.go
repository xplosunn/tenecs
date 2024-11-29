package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_list = packageWith(
	withFunction("append", tenecs_list_append),
	withStruct(Tenecs_list_Break),
	withFunction("filter", tenecs_list_filter),
	withFunction("flatMap", tenecs_list_flatMap),
	withFunction("fold", tenecs_list_fold),
	withFunction("forEach", tenecs_list_forEach),
	withFunction("length", tenecs_list_length),
	withFunction("map", tenecs_list_map),
	withFunction("mapUntil", tenecs_list_mapUntil),
	withFunction("mapNotNull", tenecs_list_mapNotNull),
	withFunction("repeat", tenecs_list_repeat),
)

var tenecs_list_append = functionFromSignature("<T>(list: List<T>, newElement: T): List<T>")

var tenecs_list_filter = functionFromSignature("<A>(list: List<A>, keep: (A) ~> Boolean): List<A>")

var tenecs_list_flatMap = functionFromSignature("<A, B>(list: List<A>, f: (A) ~> List<B>): List<B>")

var tenecs_list_fold = functionFromSignature("<A, Acc>(list: List<A>, zero: Acc, f: (Acc, A) ~> Acc): Acc")

var tenecs_list_forEach = functionFromSignature("<A>(list: List<A>, f: (A) ~> Void): Void")

var tenecs_list_length = functionFromSignature("<T>(list: List<T>): Int")

var tenecs_list_map = functionFromSignature("<A, B>(list: List<A>, f: (A) ~> B): List<B>")

var tenecs_list_mapUntil = functionFromSignature("<A, B, S>(list: List<A>, f: (A) ~> Break<S> | B): S | List<B>", Tenecs_list_Break)

var tenecs_list_mapNotNull = functionFromSignature("<A, B>(list: List<A>, f: (A) ~> B | Void): List<B>")

var tenecs_list_repeat = functionFromSignature("<A>(elem: A, times: Int): List<A>")

var Tenecs_list_Break = structWithFields("Break", tenecs_list_Break, tenecs_list_Break_Fields...)

var tenecs_list_Break = types.Struct("tenecs.list", "Break", []string{"S"})

var tenecs_list_Break_Fields = []func(fields *StructWithFields){
	structField("value", &types.TypeArgument{Name: "S"}),
}
