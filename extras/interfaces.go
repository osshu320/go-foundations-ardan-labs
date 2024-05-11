package main

import "log"

type User struct {
	Name  string
	Email string
}

func (u *User) Notify() error {
	log.Printf("User: Sending User Email To %s<%s>\n",
		u.Name,
		u.Email)
	return nil
}

type Notifier interface {
	Notify() error
}

func SendNotification(notify Notifier) error {
	return notify.Notify()
}

type Admin struct {
	User
	Level string
}

func (a *Admin) Notify() error {
	log.Printf("Admin: Sending Admin Email to %s<%s>\n",
		a.Name,
		a.Email)

	return nil
}

func interfaces_go_main() {
	// fmt.Println("Interfaces hello world")

	// bill := User{"Bill", "bill@email.com"}
	// _ = bill.Nofity()

	// jill := &User{"Jill", "jill@email.com"}
	// _ = jill.Nofity()

	// SendNotification(bill) // doesn't work
	// SendNotification(jill)

	admin := &Admin{
		User: User{
			Name:  "john smith",
			Email: "john@email.com",
		},
		Level: "super",
	}

	SendNotification(admin)
	admin.User.Notify()
	admin.Notify()
}
