package module

type Module[T IWithID[TID], TID IUintID, TDTO any, TQ any] struct {
	Handler IHandler[T, TID, TDTO, TQ]
}

func NewModule[T IWithID[TID], TID IUintID, TDTO any, TQ any](
	namespace string,
	hYields HandlerYields[T, TDTO],
	sYields ServiceYields[T, TID],
) *Module[T, TID, TDTO, TQ] {
	repo := NewRepository[T](namespace)
	svc := NewService(repo, sYields)
	handler := NewHandler[T, TID, TDTO, TQ](
		NewServiceAdapter(svc),
		hYields,
	)

	return &Module[T, TID, TDTO, TQ]{
		Handler: handler,
	}
}
