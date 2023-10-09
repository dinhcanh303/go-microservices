Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	Get(ctx context.Context, uuid string) (*domain.Comment, error)
	Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	Delete(ctx context.Context, uuid string) (bool, error)
	DeleteByPostID(ctx context.Context, postId string) (bool, error)
	DeleteByCommentID(ctx context.Context, commentId string) (bool, error)
	ListByPostID(ctx context.Context, postId string) ([]*domain.Comment, error)
	CountByPostID(ctx context.Context, postId string) (uint64, error)
	CountByCommentID(ctx context.Context, commentId string) (uint64, error)

-- name: Create :one
INSERT INTO 
    interaction.comments (
        id,
        
    )