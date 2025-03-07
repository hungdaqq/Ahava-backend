package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	PlaceOrder(order models.PlaceOrder, final_price uint64) (models.Order, error)
	PlaceOrderItem(order_id uint, item models.CartItem) error
	GetOrderItems(order_id uint) ([]models.OrderItem, error)
	ListAllOrders(limit, offset int) (models.ListOrders, error)
	GetOrderDetails(user_id, order_id uint) (models.Order, error)
	GetOrderForWebhook(order_id uint) (models.Order, error)
	UpdateOrder(order_id uint, order models.Order) (models.Order, error)
}

type orderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (r *orderRepository) PlaceOrder(o models.PlaceOrder, final_price uint64) (models.Order, error) {
	// Define the order
	order := domain.Order{
		UserID:        o.UserID,
		Address:       o.Address,
		Name:          o.Name,
		Phone:         o.Phone,
		PaymentMethod: o.PaymentMethod,
		FinalPrice:    final_price,
		Coupon:        o.Coupon,
	}
	// Create the order
	err := r.DB.Create(&order).Error
	if err != nil {
		return models.Order{}, err
	}
	// Return the order
	return models.Order{
		ID:            order.ID,
		UserID:        order.UserID,
		Address:       order.Address,
		Name:          order.Name,
		Phone:         order.Phone,
		PaymentMethod: order.PaymentMethod,
		FinalPrice:    order.FinalPrice,
		Coupon:        order.Coupon,
		OrderStatus:   order.OrderStatus,
		PaymentStatus: order.PaymentStatus,
	}, nil
}

func (r *orderRepository) PlaceOrderItem(order_id uint, item models.CartItem) error {
	// Create the order item
	err := r.DB.Create(&domain.OrderItem{
		OrderID:           order_id,
		ProductID:         item.ProductID,
		Quantity:          item.Quantity,
		Size:              item.Size,
		OriginalPrice:     item.OriginalPrice,
		DiscountPrice:     item.DiscountPrice,
		ItemPrice:         item.ItemPrice,
		ItemDiscountPrice: item.ItemDiscountPrice},
	).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *orderRepository) GetOrderForWebhook(order_id uint) (models.Order, error) {
	// Define the order
	var order models.Order
	// Query to get the order details
	err := r.DB.Raw(`SELECT id,final_price FROM orders WHERE id=?`,
		order_id).Scan(&order).Error
	if err != nil {
		return models.Order{}, err
	}
	// Return the order details
	return order, nil
}

func (r *orderRepository) GetOrderDetails(user_id, order_id uint) (models.Order, error) {
	// Define the order
	var order models.Order
	// Query to get the order details
	err := r.DB.Where("order_id=? AND user_id=?").
		Find(&order).Error
	if err != nil {
		return models.Order{}, err
	}
	// Return the order details
	return order, nil
}

func (r *orderRepository) UpdateOrder(order_id uint, o models.Order) (models.Order, error) {
	// Define the order
	var order models.Order
	// Update the order
	result := r.DB.Model(&domain.Order{}).
		Where("id = ?", order_id).
		Updates(domain.Order{
			PaymentMethod: o.PaymentMethod,
			OrderStatus:   o.OrderStatus,
			PaymentStatus: o.PaymentStatus,
		}).
		Scan(&order)
	if result.Error != nil {
		return models.Order{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.Order{}, models.ErrEntityNotFound
	}
	// Return the updated order
	return order, nil
}

func (r *orderRepository) ListAllOrders(limit, offset int) (models.ListOrders, error) {
	// Define the list of orders
	var orders []models.Order
	var total int64
	// Define the query
	query := r.DB.Model(&domain.Order{})
	if err := query.Count(&total).Error; err != nil {
		return models.ListOrders{}, err
	}
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return models.ListOrders{}, err
	}
	var orderDetails []models.OrderDetails
	for _, order := range orders {
		orderDetails = append(orderDetails, models.OrderDetails{Order: order})
	}
	// Return the list of orders
	return models.ListOrders{
		Orders: orderDetails,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (r *orderRepository) GetOrderItems(order_id uint) ([]models.OrderItem, error) {
	// Define the order items
	var orderItems []models.OrderItem
	// Query to get the order items
	err := r.DB.Where("order_id=?", order_id).
		Find(&orderItems).Error
	if err != nil {
		return nil, err
	}
	// Return the order items
	return orderItems, nil
}

// func (r *orderRepository) GetOrders(order models.) ([]domain.Order, error) {

// 	var orders []domain.Order

// 	if err := r.DB.Raw("select * from orders where user_id=?", id).Scan(&orders).Error; err != nil {
// 		return []domain.Order{}, err
// 	}

// 	return orders, nil

// }

// func (ad *orderRepository) GetCart(id uint) ([]models.GetCart, error) {

// 	var cart []models.GetCart

// 	if err := ad.DB.Raw("SELECT products.name,cart_products.quantity,cart_products.total_price AS Total FROM cart_products JOIN products ON cart_products.product_id=products.id WHERE user_id=$1", id).Scan(&cart).Error; err != nil {
// 		return []models.GetCart{}, err
// 	}
// 	return cart, nil

// }

// func (i *orderRepository) OrderItems(userID, addressid, paymentid uint, total float64, coupon string) (int, error) {

// 	var id uint
// 	query := `
//     INSERT INTO orders (user_id,address_id, payment_method_id, final_price,coupon_used)
//     VALUES (?, ?, ?, ?, ?)
//     RETURNING id
//     `
// 	i.DB.Raw(query, userID, addressid, paymentid, total, coupon).Scan(&id)

// 	return id, nil

// }

// func (i *orderRepository) AddOrderProducts(order_id uint, cart []models.GetCart) error {

// 	query := `
//     INSERT INTO order_items (order_id,product_id,quantity,total_price)
//     VALUES (?, ?, ?, ?)
//     `

// 	for _, v := range cart {
// 		var inv int
// 		if err := i.DB.Raw("select id from products where name=$1", v.ProductName).Scan(&inv).Error; err != nil {
// 			return err
// 		}

// 		if err := i.DB.Exec(query, order_id, inv, v.Quantity, v.Total).Error; err != nil {
// 			return err
// 		}
// 	}

// 	return nil

// }

// func (i *orderRepository) CancelOrder(id uint) error {

// 	if err := i.DB.Exec("update orders set order_status='CANCELED' where id=$1", id).Error; err != nil {
// 		return err
// 	}

// 	return nil

// }

// func (i *orderRepository) EditOrderStatus(status string, id uint) error {

// 	if err := i.DB.Exec("update orders set order_status=$1 where id=$2", status, id).Error; err != nil {
// 		return err
// 	}

// 	return nil

// }

// func (r *orderRepository) AdminOrders(status string) ([]domain.OrderDetails, error) {

// 	var orders []domain.OrderDetails
// 	if err := r.DB.Raw("SELECT orders.id AS id, users.name AS username, CONCAT('House Name:',addresses.house_name, ',', 'Street:', addresses.street, ',', 'Province:', addresses.province, ',', 'State', addresses.state, ',', 'Phone:', addresses.phone) AS address, payment_methods.payment_name AS payment_method, orders.final_price As total FROM orders JOIN users ON users.id = orders.user_id JOIN payment_methods ON payment_methods.id = orders.payment_method_id JOIN addresses ON orders.address_id = addresses.id WHERE order_status = $1", status).Scan(&orders).Error; err != nil {
// 		return []domain.OrderDetails{}, err
// 	}

// 	return orders, nil

// }

// func (o *orderRepository) CheckOrder(orderID string, userID uint) error {

// 	var count int
// 	err := o.DB.Raw("select count(*) from orders where order_id = ?", orderID).Scan(&count).Error
// 	if err != nil {
// 		return err
// 	}
// 	if count < 0 {
// 		return errors.New("no such order exist")
// 	}
// 	var checkUser int
// 	err = o.DB.Raw("select user_id from orders where order_id = ?", orderID).Scan(&checkUser).Error
// 	if err != nil {
// 		return err
// 	}

// 	if userID != checkUser {
// 		return errors.New("the order is not did by this user")
// 	}

// 	return nil
// }

// func (o *orderRepository) GetOrderDetail(orderID string) (domain.Order, error) {

// 	var orderDetails domain.Order
// 	err := o.DB.Raw("select * from orders where order_id = ?", orderID).Scan(&orderDetails).Error
// 	if err != nil {
// 		return domain.Order{}, err
// 	}

// 	return orderDetails, nil

// }

// func (i *orderRepository) ReturnOrder(id uint) error {

// 	if err := i.DB.Exec("update orders set order_status='RETURNED' where id=$1", id).Error; err != nil {
// 		return err
// 	}

// 	return nil

// }

// func (o *orderRepository) CheckOrderStatusByID(id uint) (string, error) {

// 	var status string
// 	err := o.DB.Raw("select order_status from orders where id = ?", id).Scan(&status).Error
// 	if err != nil {
// 		return "", err
// 	}

// 	return status, nil
// }

// func (o *orderRepository) FindAmountFromOrderID(id uint) (float64, error) {

// 	var amount float64
// 	err := o.DB.Raw("select final_price from orders where id = ?", id).Scan(&amount).Error
// 	if err != nil {
// 		return 0, err
// 	}

// 	return amount, nil
// }

// func (i *orderRepository) CreditToUserWallet(amount float64, walletid uint) error {

// 	if err := i.DB.Exec("update wallets set amount=$1 where id=$2", amount, walletId).Error; err != nil {
// 		return err
// 	}

// 	return nil

// }

// func (o *orderRepository) FindUserIdFromOrderID(id uint) (int, error) {

// 	var userID uint
// 	err := o.DB.Raw("select user_id from orders where id = ?", id).Scan(&userID).Error
// 	if err != nil {
// 		return 0, err
// 	}

// 	return userID, nil
// }

// func (o *orderRepository) FindWalletIdFromUserID(userID uint) (int, error) {

// 	var count int
// 	err := o.DB.Raw("select count(*) from wallets where user_id = ?", userId).Scan(&count).Error
// 	if err != nil {
// 		return 0, err
// 	}

// 	var walletid uint
// 	if count > 0 {
// 		err := o.DB.Raw("select id from wallets where user_id = ?", userId).Scan(&walletID).Error
// 		if err != nil {
// 			return 0, err
// 		}
// 	}

// 	return walletID, nil

// }

// func (o *orderRepository) CreateNewWallet(userID uint) (int, error) {

// 	var walletid uint
// 	err := o.DB.Exec("Insert into wallets(user_id,amount) values($1,$2)", userID, 0).Error
// 	if err != nil {
// 		return 0, err
// 	}

// 	if err := o.DB.Raw("select id from wallets where user_id=$1", userID).Scan(&walletID).Error; err != nil {
// 		return 0, err
// 	}

// 	return walletID, nil
// }

// func (o *orderRepository) MakePaymentStatusAsPaid(id uint) error {

// 	err := o.DB.Exec("UPDATE orders SET payment_status = 'PAID' WHERE id = $1", id).Error
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (o *orderRepository) GetProductImagesInAOrder(id uint) ([]string, error) {

// 	var images []string
// 	err := o.DB.Raw(`SELECT products.image
// 	FROM order_items
// 	JOIN products ON products.id = order_items.product_id
// 	JOIN orders ON orders.id = order_items.order_id
// 	WHERE orders.id = $1`, id).Scan(&images).Error
// 	if err != nil {
// 		return []string{}, err
// 	}

// 	return images, nil
// }

// func (o *orderRepository) GetIndividualOrderDetails(id uint) (models.IndividualOrderDetails, error) {

// 	var details models.IndividualOrderDetails
// 	err := o.DB.Raw(`SELECT orders.id AS order_id,
// 	CONCAT('House Name:',addresses.house_name, ' ', 'Street:', addresses.street, ' ', 'Province:', addresses.province, ' ', 'State', addresses.state) AS address,
// 	addresses.phone AS phone,
// 	orders.coupon_used,
// 	payment_methods.payment_name AS payment_method,
// 	orders.final_price As total_amount ,
// 	orders.order_status,
// 	orders.payment_status
// 	FROM orders
// 	 JOIN payment_methods ON payment_methods.id = orders.payment_method_id
// 	JOIN addresses ON orders.address_id = addresses.id
// 	WHERE orders.id = $1`, id).Scan(&details).Error
// 	if err != nil {
// 		return models.IndividualOrderDetails{}, err
// 	}

// 	return details, nil
// }

// func (o *orderRepository) GetProductDetailsInOrder(id uint) ([]models.ProductDetails, error) {

// 	var products []models.ProductDetails
// 	err := o.DB.Raw(`SELECT  products.name,
// 	products.image,
// 	order_items.quantity,
// 	order_items.total_price AS amount
// 	FROM order_items
// 	JOIN products ON products.id = order_items.product_id
// 	JOIN orders ON order_items.order_id = orders.id
// 	WHERE orders.id = $1`, id).Scan(&products).Error
// 	if err != nil {
// 		return []models.ProductDetails{}, err
// 	}

// 	return products, nil
// }

// func (o *orderRepository) FindPaymentMethodOfOrder(id uint) (string, error) {

// 	var payment string

// 	if err := o.DB.Raw(`select payment_methods.payment_name
// 	 from payment_methods
// 	  join orders on orders.payment_method_id = payment_methods.id
// 	   where orders.id = $1`, id).Scan(&payment).Error; err != nil {
// 		return "", err
// 	}
// 	return payment, nil
// }
