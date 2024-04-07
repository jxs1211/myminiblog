// Copyright 2022 Innkeeper Jayflow <jxs121@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/marmotedu/miniblog/internal/pkg/model"
)

// UserStore defines the methods implemented by the user module in the store layer.
type UserStore interface {
	Create(ctx context.Context, user *model.UserM) error
	Get(ctx context.Context, username string) (*model.UserM, error)
	Update(ctx context.Context, user *model.UserM) error
	List(ctx context.Context, offset, limit int) (int64, []*model.UserM, error)
	Delete(ctx context.Context, username string) error
}

// users is the implementation of the UserStore interface.
type users struct {
	db *gorm.DB
}

// Ensure that users implements the UserStore interface.
var _ UserStore = (*users)(nil)

func newUsers(db *gorm.DB) *users {
	return &users{db}
}

// Create inserts a user record.
func (u *users) Create(ctx context.Context, user *model.UserM) error {
	return u.db.Create(&user).Error
}

// Get retrieves the specified user's database record by username.
func (u *users) Get(ctx context.Context, username string) (*model.UserM, error) {
	var user model.UserM
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates a user database record.
func (u *users) Update(ctx context.Context, user *model.UserM) error {
	return u.db.Save(user).Error
}

// List returns a list of users based on the offset and limit.
func (u *users) List(ctx context.Context, offset, limit int) (count int64, ret []*model.UserM, err error) {
	err = u.db.Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

// Delete deletes a database user record based on the username.
func (u *users) Delete(ctx context.Context, username string) error {
	err := u.db.Where("username = ?", username).Delete(&model.UserM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
