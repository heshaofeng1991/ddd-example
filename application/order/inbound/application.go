/*
	@Author  johnny
	@Author  heshaofeng1991@gmail.com
	@Version v1.0.0
	@File    cqrs
	@Date    2022/4/14 15:36
	@Desc
*/

package inbound

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct{}

type Queries struct{}
