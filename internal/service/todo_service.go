package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ujjvalkumar/todo-api/internal/domain"
	"github.com/ujjvalkumar/todo-api/internal/repository"
)

type TodoService struct {
	todoRepo *repository.TodoRepository
}

func NewTodoService(todoRepo *repository.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

func (s *TodoService) Create(ctx context.Context, userID uuid.UUID, req domain.CreateTodoRequest) (*domain.Todo, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	now := time.Now().UTC()
	todo := &domain.Todo{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.todoRepo.Create(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *TodoService) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Todo, error) {
	return s.todoRepo.GetByID(ctx, id, userID)
}

func (s *TodoService) List(ctx context.Context, userID uuid.UUID) ([]domain.Todo, error) {
	return s.todoRepo.ListByUser(ctx, userID)
}

func (s *TodoService) Update(ctx context.Context, id, userID uuid.UUID, req domain.UpdateTodoRequest) (*domain.Todo, error) {
	todo, err := s.todoRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Status != nil {
		todo.Status = *req.Status
	}
	todo.UpdatedAt = time.Now().UTC()

	if err := s.todoRepo.Update(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *TodoService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.todoRepo.Delete(ctx, id, userID)
}
