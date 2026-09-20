package gf

const (
	StatusOK  int   = 0 //状态正常码
	StatusNo  int   = 1 //状态禁用码
	ZeroInt64 int64 = 0 //比较int64零值
	OneInt64  int64 = 1 //比较int64的1值
	ZeroInt8  int8  = 0 //比较int8零值
	OneInt8   int8  = 1 //比较int8的1值
)

type (
	Map        = map[string]interface{}      // Map is alias of frequently-used map type map[string]interface{}.
	MapAnyAny  = map[interface{}]interface{} // MapAnyAny is alias of frequently-used map type map[interface{}]interface{}.
	MapAnyStr  = map[interface{}]string      // MapAnyStr is alias of frequently-used map type map[interface{}]string.
	MapAnyInt  = map[interface{}]int         // MapAnyInt is alias of frequently-used map type map[interface{}]int.
	MapStrStr  = map[string]string           // MapStrStr is alias of frequently-used map type map[string]string.
	MapStrInt  = map[string]int              // MapStrInt is alias of frequently-used map type map[string]int.
	MapIntAny  = map[int]interface{}         // MapIntAny is alias of frequently-used map type map[int]interface{}.
	MapIntStr  = map[int]string              // MapIntStr is alias of frequently-used map type map[int]string.
	MapIntInt  = map[int]int                 // MapIntInt is alias of frequently-used map type map[int]int.
	MapAnyBool = map[interface{}]bool        // MapAnyBool is alias of frequently-used map type map[interface{}]bool.
	MapStrBool = map[string]bool             // MapStrBool is alias of frequently-used map type map[string]bool.
	MapIntBool = map[int]bool                // MapIntBool is alias of frequently-used map type map[int]bool.
)

type (
	List        = []Map        // List is alias of frequently-used slice type []Map.
	ListAnyAny  = []MapAnyAny  // ListAnyAny is alias of frequently-used slice type []MapAnyAny.
	ListAnyStr  = []MapAnyStr  // ListAnyStr is alias of frequently-used slice type []MapAnyStr.
	ListAnyInt  = []MapAnyInt  // ListAnyInt is alias of frequently-used slice type []MapAnyInt.
	ListStrStr  = []MapStrStr  // ListStrStr is alias of frequently-used slice type []MapStrStr.
	ListStrInt  = []MapStrInt  // ListStrInt is alias of frequently-used slice type []MapStrInt.
	ListIntAny  = []MapIntAny  // ListIntAny is alias of frequently-used slice type []MapIntAny.
	ListIntStr  = []MapIntStr  // ListIntStr is alias of frequently-used slice type []MapIntStr.
	ListIntInt  = []MapIntInt  // ListIntInt is alias of frequently-used slice type []MapIntInt.
	ListAnyBool = []MapAnyBool // ListAnyBool is alias of frequently-used slice type []MapAnyBool.
	ListStrBool = []MapStrBool // ListStrBool is alias of frequently-used slice type []MapStrBool.
	ListIntBool = []MapIntBool // ListIntBool is alias of frequently-used slice type []MapIntBool.
)

type (
	Slice    = []interface{} // Slice is alias of frequently-used slice type []interface{}.
	SliceStr = []string      // SliceStr is alias of frequently-used slice type []string.
	SliceInt = []int         // SliceInt is alias of frequently-used slice type []int.
)
