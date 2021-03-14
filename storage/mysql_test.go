package storage

import (
	"github.com/jinzhu/gorm"
	"testing"
	"time"
)

type Order struct {
	OrderId     int32     `gorm:"primary_key;column:order_id"` // 订单号
	GoodsId     int32     `gorm:"column:goods_id"`             // 物品编号
	ProductId   int32     `gorm:"column:product_id"`           // 商品ID
	CategoryId  int32     `gorm:"column:category_id"`          // 商品分类ID
	UserId      int32     `gorm:"column:user_id"`              // 用户ID
	LogisticsId int32     `gorm:"column:logistics_id"`         // 物流ID
	Receiver    int32     `gorm:"column:receiver"`             // 收货人uid
	PayWays     int32     `gorm:"column:pay_ways"`             // 支付方式
	OrderStatus int32     `gorm:"column:order_status"`         // 订单状态
	IsValid     int32     `gorm:"column:is_valid"`             // 是否删除 1是，0否
	CreateTime  time.Time `gorm:"column:create_time"`          // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time"`          // 更新时间
}

type Warehouse struct {
	ProductId  int32     `gorm:"primary_key;column:product_id"` //商品ID
	CategoryId int32     `gorm:"column:category_id"`            //品类ID
	Stock      int32     `gorm:"column:stock"`                  //库存数量
	CreateTime time.Time `gorm:"column:create_time"`            //创建时间
	IsValid    int32     `gorm:"column:is_valid"`               //是否有效
	UpdateTime time.Time `gorm:"column:update_time"`            //记录更新时间
}

func TestMysqlTrans(t *testing.T) {
	serverAddr := "192.168.37.131:3306"
	userName := "devops"
	passWord := "devops@123"
	maxOpen := 10
	maxIdle := 20
	idleTimeout := 300
	dataBase := "devtest"
	err := InitDB(serverAddr, userName, passWord, dataBase, maxOpen, maxIdle, idleTimeout, true)
	if err != nil {
		t.Errorf("init database failed for:%+v", err)
		return
	}

	var productId int32 = 33
	tx := GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	rowsAffected := tx.Model(&Warehouse{}).Where("product_id=?", productId).Where("stock > 0").Update("stock", gorm.Expr("stock - 1")).RowsAffected
	if rowsAffected == 0 {
		tx.Rollback()
		return
	}

	var goodsId int32 = 1
	var categoryId int32 = 2
	var userId int32 = 123456
	var logisticsId int32 = 111
	var receiver int32 = 123456
	var payWays int32 = 1
	var orderStatus int32 = 1
	var isValid int32 = 1
	updateTime := time.Now()
	createTime := time.Now()
	order := &Order{
		GoodsId:     goodsId,
		ProductId:   productId,
		CategoryId:  categoryId,
		UserId:      userId,
		LogisticsId: logisticsId,
		Receiver:    receiver,
		PayWays:     payWays,
		OrderStatus: orderStatus,
		IsValid:     isValid,
		CreateTime:  createTime,
		UpdateTime:  updateTime,
	}
	err = tx.Create(&order).Error
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

}
