package usecases

import (
	"Task_7-Clean_Architecture_Task_Management_REST_API/domain"
	"context"
	"time"
)

type taskUsecase struct {
	taskRepository domain.TaskRepository
	contextTimeout time.Duration
}

func NewTaskUsecase(taskRepository domain.TaskRepository, timeout time.Duration) domain.TaskUsecase {
	return &taskUsecase{
		taskRepository: taskRepository,
		contextTimeout: timeout,
	}
}

func (taskUC *taskUsecase) Create(c context.Context, task *domain.Task) error {
	ctx, cancel := context.WithTimeout(c, taskUC.contextTimeout)
	defer cancel()
	return taskUC.taskRepository.Create(ctx, task)
}

func (taskUC *taskUsecase) GetTasks(c context.Context) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, taskUC.contextTimeout)
	defer cancel()
	return taskUC.taskRepository.GetTasks(ctx)
}

func (taskUC *taskUsecase) GetTaskByID(c context.Context, taskID string) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, taskUC.contextTimeout)
	defer cancel()
	return taskUC.taskRepository.GetTaskByID(ctx, taskID)
}

func (taskUC *taskUsecase) UpdateTask(c context.Context, taskID string, updated_task *domain.Task) error {
	ctx, cancel := context.WithTimeout(c, taskUC.contextTimeout)
	defer cancel()
	return taskUC.taskRepository.UpdateTask(ctx, taskID, updated_task)
}

func (taskUC *taskUsecase) DeleteTask(c context.Context, taskID string) error {
	ctx, cancel := context.WithTimeout(c, taskUC.contextTimeout)
	defer cancel()
	return taskUC.taskRepository.DeleteTask(ctx, taskID)
}
