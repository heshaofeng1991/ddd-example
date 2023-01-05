/*
	@Author  johnny
	@Author  heshaofeng1991@gmail.com
	@Version v1.0.0
	@File    cqrs
	@Date    2022/4/14 15:38
	@Desc
*/

package inbound

import (
	"context"

	applicationInbound "github.com/heshaofeng1991/ddd-johnny/application/order/inbound"
)

func NewApplication(ctx context.Context) applicationInbound.Application {
	return applicationInbound.Application{
		Commands: applicationInbound.Commands{},
		Queries:  applicationInbound.Queries{},
	}
}
