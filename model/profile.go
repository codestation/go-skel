// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

type Profile struct {
	Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserToken string `json:"user_token"`
}

func NewProfile(opts ...Option) *Profile {
	p := &Profile{
		Model: NewModel(opts...),
	}
	return p
}

type ProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (p *ProfileRequest) Profile(opts ...Option) *Profile {
	profile := NewProfile(opts...)
	profile.FirstName = p.FirstName
	profile.LastName = p.LastName

	return profile
}
