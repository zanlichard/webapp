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

func (Order) TableName() string {
	return "t_order"
}

type Warehouse struct {
	ProductId  int32     `gorm:"primary_key;product_id"` //商品ID
	CategoryId int32     `gorm:"category_id"`            //品类ID
	Stock      int32     `gorm:"stock"`                  //库存数量
	CreateTime time.Time `gorm:"create_time"`            //创建时间
	IsValid    int32     `gorm:"is_valid"`               //是否有效
	UpdateTime time.Time `gorm:"update_time"`            //记录更新时间
}

func (Warehouse) TableName() string {
	return "t_warehouse"
}

/*

CREATE TABLE `t_order` (
  `order_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `goods_id` int(11) NOT NULL COMMENT '品ID',
  `product_id` int(11) DEFAULT NULL COMMENT '品ID',
  `category_id` int(11) DEFAULT NULL COMMENT '品ID',
  `user_id` int(11) NOT NULL COMMENT '没id',
  `logistics_id` int(11) DEFAULT NULL COMMENT '鞯ズ',
  `receiver` int(11) DEFAULT NULL COMMENT '栈',
  `pay_ways` int(11) DEFAULT NULL COMMENT '支式',
  `order_status` int(11) DEFAULT NULL COMMENT '状态',
  `create_time` datetime DEFAULT NULL COMMENT '时',
  `update_time` datetime DEFAULT NULL COMMENT '时',
  `is_valid` int(11) DEFAULT '1' COMMENT '欠效',
  PRIMARY KEY (`order_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE `t_warehouse` (
  `product_id` int(11) NOT NULL COMMENT '品ID',
  `category_id` int(11) NOT NULL COMMENT '品ID',
  `stock` int(11) NOT NULL,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时',
  `is_valid` int(11) DEFAULT '1' COMMENT '欠效',
  PRIMARY KEY (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

  insert into t_warehouse('product_id','category_id','stock','is_valid') values('33','2','100','1');
  insert into t_warehouse('product_id','category_id','stock','is_valid') values('34','1','100','1');
*/

func TestMysqlTrans(t *testing.T) {
	serverAddr := "192.168.37.131:3306"
	userName := "devops"
	passWord := "devops@123"
	maxOpen := 10
	maxIdle := 20
	idleTimeout := 300
	dataBase := "devtest"
	err := InitDB(serverAddr, userName, passWord, dataBase, maxOpen, maxIdle, idleTimeout, false)
	if err != nil {
		t.Errorf("init database failed for:%+v", err)
		return
	}
	var productId int32 = 33
	tx := GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("error :%+v", r)
			tx.Rollback()
		}
	}()
	result := tx.Model(&Warehouse{}).Where("product_id=?", productId).Where("stock > 0").Update("stock", gorm.Expr("stock - 1"))
	rowsAffected := result.RowsAffected
	err = result.Error
	t.Logf("update warehouse")
	if err != nil || rowsAffected == 0 {
		t.Errorf("update t_warehouse failed for:%+v", err)
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
		t.Errorf("place order failed for:%+v", err)
		tx.Rollback()
	} else {
		tx.Commit()
	}

}
