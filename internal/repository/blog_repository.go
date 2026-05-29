package repository

// BlogRepository handles all blog database queries.
type BlogRepository struct{}

func NewBlogRepository() *BlogRepository {
	return &BlogRepository{}
}

// TODO: FindAllPublished, FindBySlug, FindByID, Create, Update, Delete, CountAll, CountByStatus
