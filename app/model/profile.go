// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

import "megpoid.dev/go/go-skel/pkg/model"

type Profile struct {
	model.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func NewProfile(opts ...model.Option) *Profile {
	p := &Profile{
		Model: model.NewModel(opts...),
	}
	return p
}

type ProfileRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (p *ProfileRequest) Profile(opts ...model.Option) *Profile {
	profile := NewProfile(opts...)
	profile.FirstName = p.FirstName
	profile.LastName = p.LastName
	profile.Email = p.Email

	return profile
}
