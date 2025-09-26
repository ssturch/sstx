## **Описание:**
Пакет для управления транзакциями полученными от pgx.Pool.
## **Пример:**

```

func main() {
ctx := context.Background()

	pool, _ := pgx.Connect(ctx, "postgres://user:pass@localhost:5432/db?sslmode=disable")

	txManager := New(pool)
	tx, _ := txManager.GetTx(ctx)
	//tx, _ := txManager.GetTxWithOpt(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})

	ur := new(testdata.UserRepo)
	sr := new(testdata.ScopeRepo)

	userRepo := NewUnit[testdata.UserRepo](tx, ur)
	scopeRepo := NewUnit[testdata.ScopeRepo](tx, sr)

	var err error
	defer tx.StopTxByErr(ctx, err)

	_, err = userRepo.Exec().CreateUser(ctx, "user_data")
	if err != nil {
		return
	}
	_, err = scopeRepo.Exec().CreateScope(ctx, "scope_data")

	// Вместо tx.StopTxByErr можно использовать tx.Rollback() или tx.Commit()
}
```



