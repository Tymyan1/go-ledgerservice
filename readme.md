# Project in a nutshell
A simple ledger service which keeps track of transactions built in golang.
My first go project and a way to learn golang with practical application.

# Issues and pain points
## Lack of enums
The lack of native enum support is a bit surprising and I'm struggling to work with the custom string type and constant effectively, mostly due the the fact that this approach does not prevent assigning illegal values
```
type TransactionType string

const (
	CREDIT TransactionType = "CREDIT"
	DEBIT  TransactionType = "DEBIT"
)

func iOnlyWantTheLegalValues(txType TransactionType) {
	fmt.Println(txType)
}

func main() {
	iOnlyWantTheLegalValues("asdf") // compliles, runs and prints out "asdf"
}
```
I have also seen the usual definition of the enums using the iota and int values, however this approach requires functions translating the int values into the string form when unmarshalling a json (eg. in a request) or saving to database and still do not guard against illegal values.
```
const (
	CREDIT = iota
	DEBIT
)
...
iOnlyWantTheLegalValues(123) // compliles, runs and prints out 123
```
## Comparison of structs with nested pointers
The great thing about comparison in go is that `==` tends to work well for most cases. Unfortunatelly in the case of comparing two Transactions, which have nested decimal (using the [shopspring](https://pkg.go.dev/github.com/shopspring/decimal "shopspring")), failed miserably. After quite some time trying to understand why, it turns out the Decimal implementation is a struct with a pointer
```
type Decimal struct {
	value *big.Int
	exp int32
}
```
As a result, the `==` comparison of two decimals compare the pointers (ie. the memory addresses) rather than the values themselves and as a result `decimal.NewFromInt(10) != decimal.NewFromInt(10)`. As such transitively two transactions/structs with equal values (but different pointers) are not comparable using the `==`.
As a work-around, I decided to use the `reflect.DeepEqual(...)` which works, but seems like a seriously heavy-duty tool to use for something as simple as a struct with a few not very deeply nested fields.
Had a look into alternative implementations, but they just are not nearly as popular (shopspring has 5.8k stars on GitHub compared to eg. [ericlagergren](https://github.com/ericlagergren/decimal)'s 500 star), so I am mostly assuming I am misusing/misunderstanding how the comparisons should be done. An alternative here would be to use a different decimal implementation or write one myself, but that would be creep the scope of the project by a good amount.
## Error handling *everywhere*
Coming from the Java background, I am used to just throw an exception and not have to worry about it propagating up the stack myself, knowing that at some point it will get handled. Of course exceptions are not an answer to everything and sometimes a result object works much better (and with sealed interfaces and pattern matching in Java 21 this can be done very nicely), however from my experience more often than not when an error occurs, there is fairly little one can do to not just fail the whole operation and propagate the failure upwards. 
The way Go handles errors as a type of return of every function means that whenever a function is called, en evaluation must be done on whether or not it was successful, which leads to massive pollution of the code base and distracts from its intention (at least in my case). 
For example, in the following function I need a database transaction, lookup a transaction, do some validation, add some numbers, save the the transaction and commit. Without the error handling, the intention is pretty readable and obvious at first glance.
```
func (ts *TransactionService) ProcessTransaction(transaction Transaction, ctx context.Context) error {
	dbTx, _ := ts.txManager.BeginTx(ctx, nil)

	lastTxPt, _ := ts.transactionDb.QueryLatestForAccount(transaction.AccountKey)
	lastTx := getTransactionData(lastTxPt)

	validateTransaction(transaction, lastTx)
	
	newBalance, _ := calculateNewBalance(transaction, lastTx)
	transaction.Balance = newBalance
	
	ts.transactionDb.Save(transaction)

	return dbTx.Commit()
}
```
However after adding the error handling, it becomes very unobvious what is happening as the number of lines required doubles(!), the duplication of the code soars and if your goal is high-coverage, the amount of tests required to cover all the new branches, in my opinion, closes the point of where maintainability and readability of tests becomes close to impossible. And this is a simple function we're talking about, I don't want to think of what would happen if the function was more complicated than this!
```
unc (ts *TransactionService) ProcessTransaction(transaction Transaction, ctx context.Context) error {
	dbTx, err := ts.txManager.BeginTx(ctx, nil)
	if err != nil {
		return rollback(dbTx, err)
	}

	lastTxPt, err := ts.transactionDb.QueryLatestForAccount(transaction.AccountKey)
	if err != nil {
		return rollback(dbTx, err)
	}

	lastTx := getTransactionData(lastTxPt)

	if err := validateTransaction(transaction, lastTx); err != nil {
		return rollback(dbTx, err)
	}

	newBalance, err := calculateNewBalance(transaction, lastTx)
	if err != nil {
		return rollback(dbTx, err)
	}

	transaction.Balance = newBalance
	if err := ts.transactionDb.Save(transaction); err != nil {
		return rollback(dbTx, err)
	}

	return dbTx.Commit()
}
```
It is well possible I am missing something obvious to a seasoned Golang engineer and if that's the case, please let me know!

Meanwhile, I will dream of a way to at least aggregate and centralise the error handling within a function by using some kind of observer (where one would be able to declare an err variable and at the point it got populated, the execution would get paused and the error could be handled, either by propagation upwards, retrying or continuing the flow with some compensating action).
