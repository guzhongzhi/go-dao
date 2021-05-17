package backend

import (
	"context"
	"fmt"
	"github.com/guzhongzhi/gmicro/client"
	"github.com/guzhongzhi/gmicro/logger"
	"github.com/guzhongzhi/gmicro/render"
	"github.com/guzhongzhi/gmicro/samples/bff/api"
	"google.golang.org/grpc"
	"net/http"
)

type ADD struct {
	Line1    string `json:"line1"`
	PostCode string `json:"post_code"`
}

type UserCreateMessage struct {
	ID      *string `json:"id"`
	Name    string  `json:"name"`
	Address struct {
		Line1    string `json:"line1"`
		PostCode string `json:"post_code"`
	} `json:"address"`
	Address2  *ADD `json:"address_2"`
	IsBlocked bool `json:"is_blocked"`
}

type UserCreateResponse struct {
	ID string `json:"id"`
}

type User struct {
}

func (s *User) Create(ctx context.Context, message UserCreateMessage) (render.Render, *UserCreateResponse) {
	c, err := client.NewGRPCClient("test", "127.0.0.1", 9000, nil)
	fmt.Println(err)
	err = c.Callback(func(conn *grpc.ClientConn, log logger.SuperLogger) error {
		c := api.NewSubEffectServiceClient(conn)
		in := &api.UpsertRequest{}
		_, err := c.Create(context.Background(), in)
		return err
	})
	fmt.Println(err)
	in := &api.UpsertRequest{}
	rsp := &api.UpsertResponse{}
	err = c.Call(context.Background(), "/api.SubEffectService/Create", in, rsp)
	return render.JSON{
		Status: http.StatusOK,
	}, &UserCreateResponse{}
}

func (s *User) Update(ctx context.Context, v UserCreateMessage) (render.Render, *UserCreateResponse) {
	return render.Text{}, &UserCreateResponse{}
}

func (s *User) Delete(ctx context.Context, v struct{ Id string `json:"id"` }) (render.Render, *UserCreateResponse) {
	return render.Text{}, &UserCreateResponse{}
}

func (s *User) Get(ctx context.Context, v struct{ Id string `json:"id"` }) (render.Render, string) {
	return render.Text{}, ""
}
