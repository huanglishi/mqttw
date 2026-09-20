// ========================
// 先将切片转为 map 集合，再反复查询，适合一个切片多次判断元素是否存在。
package gmap

// Set 基于map实现的集合，用于快速元素存在性判断
type Set[T comparable] map[T]struct{}

// NewSet 切片转为集合
func NewSet[T comparable](slice []T) Set[T] {
	s := make(Set[T], len(slice))
	for _, v := range slice {
		s[v] = struct{}{}
	}
	return s
}

// Exists 判断元素是否在集合中
func (s Set[T]) Exists(val T) bool {
	_, ok := s[val]
	return ok
}
