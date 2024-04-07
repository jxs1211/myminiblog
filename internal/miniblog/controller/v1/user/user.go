// Copyright 2022 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package user

import (
	"github.com/marmotedu/miniblog/internal/miniblog/biz"
	"github.com/marmotedu/miniblog/internal/miniblog/store"
	"github.com/marmotedu/miniblog/pkg/auth"
	pb "github.com/marmotedu/miniblog/pkg/proto/miniblog/v1"
)

// UserController is the implementation of the user module in the Controller layer, used to handle requests for the user module.
type UserController struct {
	a *auth.Authz
	b biz.IBiz
	pb.UnimplementedMiniBlogServer
}

// New creates a new user controller.
func New(ds store.IStore, a *auth.Authz) *UserController {
	return &UserController{a: a, b: biz.NewBiz(ds)}
}
