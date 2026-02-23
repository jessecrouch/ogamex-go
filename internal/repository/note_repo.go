package repository

import (
	"context"

	"gorm.io/gorm"

	"ogamex-go/internal/database"
	"ogamex-go/internal/schema"
)

type noteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) NoteRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &noteRepository{db: gdb}
}

func (r *noteRepository) Create(ctx context.Context, note *schema.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *noteRepository) GetByID(ctx context.Context, id uint) (*schema.Note, error) {
	var note schema.Note
	err := r.db.WithContext(ctx).First(&note, id).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *noteRepository) GetByUserID(ctx context.Context, userID uint) ([]*schema.Note, error) {
	var notes []*schema.Note
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notes).Error
	return notes, err
}

func (r *noteRepository) Update(ctx context.Context, note *schema.Note) error {
	return r.db.WithContext(ctx).Save(note).Error
}

func (r *noteRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.Note{}, id).Error
}
