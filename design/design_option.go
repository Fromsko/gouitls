package design

type User struct {
	Age  string
	Name string
}

func NewUser(opts ...UserOption) *User {
	user := new(User)
	for _, opt := range opts {
		opt.Apply(user)
	}
	return user
}

type UserOption interface {
	Apply(*User)
}

type UserName struct {
	Name string
}

func (u UserName) Apply(user *User) {
	user.Name = u.Name
}

func NewUserName(name string) *UserName {
	return &UserName{Name: name}
}

type UserAge func(*User)

func (op UserAge) Apply(u *User) {
	op(u)
}

// func main() {
// 	newAge := func(u *User) {
// 		u.Age = "18"
// 	}
// 	user := NewUser(
// 		NewUserName("小安"),
// 		UserAge(newAge),
// 	)
// 	fmt.Println(user)
// }
