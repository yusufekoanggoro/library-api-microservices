package repository

import (
	sharedDomain "book-service/pkg/shared/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

type StockRepository interface {
	IncreaseStock(ctx context.Context, bookID uint, quantity int) error
	DecreaseStock(ctx context.Context, bookID uint, quantity int) error
}

type stockRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewStockRepository(readDB *gorm.DB, writeDB *gorm.DB) StockRepository {
	return &stockRepository{readDB: readDB, writeDB: writeDB}
}

func (r *stockRepository) IncreaseStock(ctx context.Context, bookID uint, quantity int) error {
	// var book *sharedDomain.Book

	// // FOR UPDATE (Row-level locking)
	// if err := r.writeDB.Model(&sharedDomain.Book{}).
	// 	Clauses(gorm.Expr("FOR UPDATE")).
	// 	Where("id = ?", bookID).
	// 	First(&book).Error; err != nil {
	// 	return err
	// }

	// // UPDATE
	// if err := r.writeDB.Model(&book).
	// 	Update("stock", gorm.Expr("stock + ?", quantity)).Error; err != nil {
	// 	return err
	// }

	// return nil
	return r.writeDB.Model(&sharedDomain.Book{}).
		Where("id = ?", bookID).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (r *stockRepository) DecreaseStock(ctx context.Context, bookID uint, quantity int) error {
	// var book *sharedDomain.Book

	// // FOR UPDATE (Row-level locking)
	// if err := r.writeDB.Model(&sharedDomain.Book{}).
	// 	Clauses(gorm.Expr("FOR UPDATE")).
	// 	Where("id = ? AND stock >= ?", bookID, quantity).
	// 	First(&book).Error; err != nil {
	// 	return err
	// }

	// // UPDATE
	// if err := r.writeDB.Model(&book).
	// 	Update("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
	// 	return err
	// }

	// return nil
	result := r.writeDB.Model(&sharedDomain.Book{}).
		Where("id = ? AND stock >= ?", bookID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("stok tidak mencukupi atau buku tidak ditemukan")
	}

	return nil
}
