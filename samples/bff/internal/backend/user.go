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
	"time"
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
	Address2 *ADD `json:"address_2"`
}

type User struct {
}

func (s *User) Create(ctx context.Context, message UserCreateMessage) render.Render {
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
		Data:   rsp,
		Status: http.StatusOK,
	}
}

func (s *User) Update(ctx context.Context, v UserCreateMessage) render.Render {
	return render.Text{
		Content: "fdsafsad",
	}
}

func (s *User) Delete(ctx context.Context, v struct{ Id string `json:"id"` }) render.Render {
	return render.Text{
		Content: fmt.Sprintf("%v", time.Now().UnixNano()),
	}
}
