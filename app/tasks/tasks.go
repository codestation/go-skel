// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"megpoid.dev/go/go-skel/app/repository/uow"
	"megpoid.dev/go/go-skel/app/usecase"
	"megpoid.dev/go/go-skel/pkg/sql"
)

const (
	TypeSayHello     = "say:hello"
	TypeProfileCheck = "profile:check"
)

type HelloPayload struct {
	Message string
}

type ProfilePayload struct {
	ID int64
}

func NewSayHelloTask(message string) (*asynq.Task, error) {
	payload, err := json.Marshal(HelloPayload{Message: message})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSayHello, payload), nil
}

func NewProfileCheckTask(id int64) (*asynq.Task, error) {
	payload, err := json.Marshal(ProfilePayload{ID: id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeProfileCheck, payload), nil
}

func HandleSayHelloTask(ctx context.Context, t *asynq.Task) error {
	var p HelloPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sending Message: %s", p.Message)
	fmt.Printf("Hello %s", p.Message)
	return nil
}

type ProfileChecker struct {
	pool sql.Connector
}

func (process *ProfileChecker) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ProfilePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	unitOfWork := uow.New(process.pool)
	uc := usecase.NewProfile(unitOfWork)

	profile, err := uc.GetProfile(ctx, p.ID)
	if err != nil {
		return fmt.Errorf("failed to find profile: %w", err)
	}

	log.Printf("Profile found: email=%s", profile.Email)

	return nil
}

func NewProfileProcessor(pool sql.Connector) *ProfileChecker {
	return &ProfileChecker{pool: pool}
}
