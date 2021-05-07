package data

import (
	"database/sql"
	"fmt"
	"github.com/guzhongzhi/gmicro/dao"
)

type TOrderFindOptions struct {
	OID *string
	ID  *int
}

func (s *TOrderFindOptions) Filter() (interface{}, error) {
	return "", nil
}

func (s *TOrderFindOptions) Options() (interface{}, error) {
	panic("implement me")
}

func (s *TOrderFindOptions) Pagination() *dao.Pagination {
	return dao.NewPagination()
}

type TOrder struct {
	Id        int    `sql:"id"`
	OID       string `sql:"oid"`
	TradeNo   string `sql:"trade_no"`
	Sandbox   string `sql:"sandbox"`
	BuyerInfo string `sql:"buyer_info:"`
}

func (s *TOrder) IsNew() bool {
	return s.Id == 0
}

func (s *TOrder) ID() interface{} {
	return s.Id
}

func (s *TOrder) SetID(v interface{}) {
	s.Id = v.(int)
}

func (s *TOrder) String() string {
	return fmt.Sprintf("%v", s.Id)
}

type TOrderDAO interface {
	dao.MysqlDAO
}

func NewTOrderDAO(db *sql.DB) TOrderDAO {
	return &tOrderDAO{
		dao.NewMysqlDAO(db, "t_order", "id", dao.MysqlDAOOptions{
			FindSQL: "SELECT o.id AS id,o.oid,o.trade_no FROM t_order AS  o " +
				"LEFT JOIN t_iap_order AS t ON t.t_order_id = o.id " +
				"WHERE 1=1 LIMIT 1 %s",
		}),
	}
}

type tOrderDAO struct {
	dao.MysqlDAO
}
