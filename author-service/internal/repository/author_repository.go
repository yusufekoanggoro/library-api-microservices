package repository

import (
	sharedDomain "author-service/pkg/shared/domain"
	"context"
	"log"

	"gorm.io/gorm"
)

type authorRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewAuthorRepository(readDB *gorm.DB, writeDB *gorm.DB) AuthorRepository {
	return &authorRepository{readDB: readDB, writeDB: writeDB}
}

func (r *authorRepository) GetAllAuthors(ctx context.Context, page, limit int) ([]*sharedDomain.Author, int64, error) {
	var authors []*sharedDomain.Author
	var totalRows int64

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	if err := r.readDB.Model(&sharedDomain.Author{}).Count(&totalRows).Error; err != nil {
		log.Println("GetAllAuthors count error:", err)
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err := r.readDB.Limit(limit).Offset(offset).Order("created_at DESC").Find(&authors).Error
	if err != nil {
		log.Println("GetAllAuthors query error:", err)
		return nil, 0, err
	}

	return authors, totalRows, nil
}

func (r *authorRepository) GetAuthorByID(ctx context.Context, id uint) (*sharedDomain.Author, error) {
	var author sharedDomain.Author

	err := r.readDB.First(&author, id).Error
	if err != nil {
		log.Println("GetAuthorByID error:", err)
		return nil, err
	}

	return &author, nil
}

func (r *authorRepository) SaveAuthor(ctx context.Context, author *sharedDomain.Author) error {
	return r.writeDB.Save(author).Error
}

func (r *authorRepository) DeleteAuthor(ctx context.Context, id uint) error {
	return r.writeDB.Delete(&sharedDomain.Author{}, id).Error
}

func (r *authorRepository) DeleteAuthorByID(ctx context.Context, id uint) error {
	return r.writeDB.Delete(&sharedDomain.Author{}, id).Error
}
