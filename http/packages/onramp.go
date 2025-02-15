package packages

import (
	"fmt"
	"sync"
)

type UserBalance struct {
	Balance int
	Locked  int
}

type BALANCE_INR struct {
	mutex sync.Mutex
	User  map[string]UserBalance
}

var (
	once     sync.Once
	instance *BALANCE_INR
)

func Init() *BALANCE_INR {
	once.Do(func() {
		instance = &BALANCE_INR{
			User: make(map[string]UserBalance),
		}
	})
	return instance
}

func (b *BALANCE_INR) Onramp(userId string) (*UserBalance, error) {
	fmt.Printf("📩 Received request for user: %s\n", userId)

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, exists := b.User[userId]; exists {
		fmt.Printf("🚨 User %s already exists!\n", userId)
		return nil, fmt.Errorf("user %s already exists", userId)
	}

	b.User[userId] = UserBalance{Balance: 1000, Locked: 0}
	fmt.Printf("✅ User %s added successfully!\n", userId)
	temp := b.User[userId]
	return &temp, nil
}

func (b *BALANCE_INR) GetAllUsers() (map[string]UserBalance ,int){
	b.mutex.Lock()
	defer b.mutex.Unlock()

	users := make(map[string]UserBalance)
	fmt.Println("🔍 Fetching all users...")

	for userId, userBalance := range b.User {
		users[userId] = userBalance
		fmt.Printf("👤 Found user: %s -> %+v\n", userId, userBalance)
	}

	fmt.Println("✅ Finished fetching all users")
	return users,len(users)
}