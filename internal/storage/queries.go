package storage

const (
	usersTableQuery = `CREATE TABLE IF NOT EXISTS users
	(
	 login    text NOT NULL,
	 password text NOT NULL,
	 CONSTRAINT PK_1_users PRIMARY KEY ( login )
	 )`

	ordersTableQuery = `
	CREATE TABLE IF NOT EXISTS orders
	(
	 order_num   text NOT NULL,
	 login       text NOT NULL,
	 change_time timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	 status      text NOT NULL DEFAULT 'NEW',
	 accrual     decimal DEFAULT 0,	 
	 CONSTRAINT PK_1_orders PRIMARY KEY ( order_num ),
	 CONSTRAINT REF_FK_1_orders FOREIGN KEY ( login ) REFERENCES users ( login )
	)`

	balanceTableQuery = `CREATE TABLE IF NOT EXISTS balance
	(
	 login           text NOT NULL UNIQUE,
	 current_balance decimal NOT NULL,
	 total_withdrawn decimal NOT NULL,
	 CONSTRAINT PK_1_balance PRIMARY KEY ( login ),
	 CONSTRAINT REF_FK_4_balance FOREIGN KEY ( login ) REFERENCES users ( login )
	)`

	withdrawalsTableQuery = `CREATE TABLE IF NOT EXISTS withdrawals
	(
	 new_order       text NOT NULL UNIQUE,
	 login           text NOT NULL,
	 "sum"             decimal NOT NULL,
	 withdrawal_time timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	 CONSTRAINT PK_1_withdrawals PRIMARY KEY ( new_order ),
	 CONSTRAINT REF_FK_3_withdrawals FOREIGN KEY ( login ) REFERENCES users ( login )
	)`

	createNewUserQuery     = `INSERT INTO users VALUES ($1, $2)`
	createUserBalanceQuery = `INSERT INTO balance VALUES ($1, 0, 0)`
	userPasswordQuery      = `SELECT password FROM users WHERE login = $1`

	createNewUserOrderQuery = `INSERT INTO orders (order_num, login) VALUES ($1, $2)`
	getUserOrdersListQuery  = `SELECT order_num, status, accrual, change_time FROM orders WHERE login = $1 ORDER BY change_time`
	getUserByOrderIDQuery   = `SELECT login FROM orders WHERE order_num = $1`
	updateUserOrderQuery    = `UPDATE orders SET status = $3, accrual = $4 WHERE login = $1 AND order_num = $2 AND status != $3`
	getUserBalanceQuery     = `SELECT current_balance FROM balance WHERE login = $1`
	updateUserBalanceQuery  = `UPDATE balance SET current_balance = $2 WHERE login = $1`

	getUserWithdrawnInfoQuery        = `SELECT current_balance, total_withdrawn FROM balance WHERE login = $1`
	getUserWithdrawalsQuery          = `SELECT new_order, "sum", withdrawal_time FROM withdrawals WHERE login = $1 ORDER BY withdrawal_time`
	updateUserWithdrawalBalanceQuery = `UPDATE balance SET current_balance = $2, total_withdrawn =$3 WHERE login = $1`
	createUserWithdrawalQuery        = `INSERT INTO withdrawals (new_order, login, "sum") VALUES ($1, $2, $3)`
)
