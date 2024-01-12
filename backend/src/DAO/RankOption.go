package DAO

type RankOption struct {
	Mode       int64
	R18        bool
	Likelimit  int64
	ShowSingle bool
}

//
//var (
//	ByPid    = int64(0)
//	ByAuthor = int64(1)
//	ByRank   = int64(2)
//)
//
//type rankoption func(*Option)
//
//func NewOption(op ...option) *Option {
//	Op := &Option{
//		Mode:       ByPid,
//		R18:        false,
//		Likelimit:  0,
//		ShowSingle: false,
//	}
//	for _, O := range op {
//		O(Op)
//	}
//	return Op
//}
//func WithR18(r18 bool) option {
//	return func(o *Option) {
//		o.R18 = r18
//	}
//}
//func WithLikeLimit(num int64) option {
//	return func(o *Option) {
//		o.Likelimit = num
//	}
//}
//func WithShowSingle(show bool) option {
//	return func(o *Option) {
//		o.ShowSingle = show
//	}
//}
//func WithMode(mode int64) option {
//	return func(o *Option) {
//		o.Mode = int64(1) << mode
//	}
//}
