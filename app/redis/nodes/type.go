package nodes

type IServer interface {
	ListenAndServe() error
}
