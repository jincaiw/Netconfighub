package service

import (
	"context"
	"fmt"
	"log"

	hookExecutor "github.com/netconfighub/netconfighub/internal/hook"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
)

type HookService interface {
	Create(ctx context.Context, hook *model.Hook) (*model.Hook, error)
	GetByID(ctx context.Context, id uint) (*model.Hook, error)
	List(ctx context.Context, offset, limit int) ([]model.Hook, int64, error)
	Update(ctx context.Context, id uint, hook *model.Hook) (*model.Hook, error)
	Delete(ctx context.Context, id uint) error
	ExecuteHooks(ctx context.Context, event string, payload map[string]interface{})
}

type hookServiceImpl struct {
	hookRepo repository.HookRepository
	hookSem  chan struct{}
}

func NewHookService(hookRepo repository.HookRepository) HookService {
	return &hookServiceImpl{
		hookRepo: hookRepo,
		hookSem:  make(chan struct{}, 10),
	}
}

func (s *hookServiceImpl) Create(ctx context.Context, hook *model.Hook) (*model.Hook, error) {
	if hook.Type != "exec" && hook.Type != "git_push" && hook.Type != "webhook" {
		return nil, fmt.Errorf("无效的 Hook 类型，支持: exec, git_push, webhook")
	}
	if hook.Name == "" {
		return nil, fmt.Errorf("Hook 名称不能为空")
	}
	if err := s.hookRepo.Create(ctx, hook); err != nil {
		return nil, err
	}
	return s.hookRepo.FindByID(ctx, hook.ID)
}

func (s *hookServiceImpl) GetByID(ctx context.Context, id uint) (*model.Hook, error) {
	hook, err := s.hookRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Hook 不存在")
	}
	return hook, nil
}

func (s *hookServiceImpl) List(ctx context.Context, offset, limit int) ([]model.Hook, int64, error) {
	return s.hookRepo.List(ctx, offset, limit)
}

func (s *hookServiceImpl) Update(ctx context.Context, id uint, hook *model.Hook) (*model.Hook, error) {
	existing, err := s.hookRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Hook 不存在")
	}
	existing.Name = hook.Name
	existing.Type = hook.Type
	existing.Config = hook.Config
	existing.Enabled = hook.Enabled
	existing.Events = hook.Events
	if err := s.hookRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return s.hookRepo.FindByID(ctx, id)
}

func (s *hookServiceImpl) Delete(ctx context.Context, id uint) error {
	if _, err := s.hookRepo.FindByID(ctx, id); err != nil {
		return fmt.Errorf("Hook 不存在")
	}
	return s.hookRepo.Delete(ctx, id)
}

func (s *hookServiceImpl) ExecuteHooks(ctx context.Context, event string, payload map[string]interface{}) {
	hooks, err := s.hookRepo.FindEnabledByEvent(ctx, event)
	if err != nil {
		return
	}
	for _, hook := range hooks {
		h := hook
		switch h.Type {
		case "exec":
			go func() {
				s.hookSem <- struct{}{}
				defer func() { <-s.hookSem }()
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Exec Hook panic [hookID=%d]: %v", h.ID, r)
					}
				}()
				hookExecutor.ExecuteExecHook(h, payload)
			}()
		case "git_push":
			go func() {
				s.hookSem <- struct{}{}
				defer func() { <-s.hookSem }()
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Git Push Hook panic [hookID=%d]: %v", h.ID, r)
					}
				}()
				hookExecutor.ExecuteGitPushHook(h, payload)
			}()
		case "webhook":
			go func() {
				s.hookSem <- struct{}{}
				defer func() { <-s.hookSem }()
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Webhook Hook panic [hookID=%d]: %v", h.ID, r)
					}
				}()
				hookExecutor.ExecuteWebhookHook(h, payload)
			}()
		}
	}
}
