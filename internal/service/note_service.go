package service

import (
	"context"
	"time"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type NoteService struct {
	noteRepo    repository.NoteRepository
	planetRepo  repository.PlanetRepository
}

func NewNoteService(noteRepo repository.NoteRepository, planetRepo repository.PlanetRepository) *NoteService {
	return &NoteService{
		noteRepo:   noteRepo,
		planetRepo: planetRepo,
	}
}

type CreateNoteInput struct {
	Galaxy   int
	System   int
	Position int
	Type     int
	Subject  string
	Text     string
}

func (s *NoteService) CreateNote(ctx context.Context, userID uint, input CreateNoteInput) (*schema.Note, error) {
	note := &schema.Note{
		UserID:   userID,
		Galaxy:   input.Galaxy,
		System:   input.System,
		Position: input.Position,
		Type:     input.Type,
		Subject:  input.Subject,
		Text:     input.Text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.noteRepo.Create(ctx, note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *NoteService) GetNotes(ctx context.Context, userID uint) ([]*schema.Note, error) {
	return s.noteRepo.GetByUserID(ctx, userID)
}

func (s *NoteService) GetNote(ctx context.Context, noteID uint) (*schema.Note, error) {
	return s.noteRepo.GetByID(ctx, noteID)
}

func (s *NoteService) UpdateNote(ctx context.Context, noteID uint, subject, text string) (*schema.Note, error) {
	note, err := s.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return nil, err
	}

	note.Subject = subject
	note.Text = text
	note.UpdatedAt = time.Now()

	err = s.noteRepo.Update(ctx, note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (s *NoteService) DeleteNote(ctx context.Context, noteID uint) error {
	return s.noteRepo.Delete(ctx, noteID)
}
