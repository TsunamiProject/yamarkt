package storage

type Storage interface {
	RegisterUser()
	AuthUser()

	GetUserBalance()
	GetUserOrders()
	GetUserWithdrawals()

	CreateUserOrder()
	CreateUserWithdraw()
}
