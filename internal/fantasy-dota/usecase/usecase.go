package usecase

type Usecase struct {
	repo repository
}

func NewUsecase(repo repository) *Usecase {
	return &Usecase{repo: repo}
}
