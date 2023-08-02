// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"
	"time"
)

type Timers struct {
	Start    time.Time     `json:"start"`
	Finish   time.Time     `json:"finish"`
	Duration time.Duration `json:"duration"`
}

type DelayInteractor struct {
	common
}

func NewDelay() *DelayInteractor {
	return &DelayInteractor{
		common: newCommon(),
	}
}

func (uc *DelayInteractor) Process(ctx context.Context, delay time.Duration) (*Timers, error) {
	start := uc.currentTime()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(delay):
	}

	finish := uc.currentTime()
	duration := finish.Sub(start)

	timers := Timers{
		Start:    start,
		Finish:   finish,
		Duration: duration,
	}

	return &timers, nil
}
