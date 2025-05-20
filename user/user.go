package user

import (
	"encoding/gob"
	"os"
)

var (
	filePath string
	users    []*User
)

type User struct {
	Name   string `yaml:"name"`
	Passwd string `yaml:"passwd"`
}

func Load(fp string) error {
	filePath = fp
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	dec.Decode(&users)
	return nil
}

func Write() error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	en := gob.NewEncoder(f)
	en.Encode(&users)
	return nil
}

func Get() []*User {
	return users
}

func CheckUser(name, passwd string) bool {
	for _, u := range users {
		if u.Name == name && u.Passwd == passwd {
			return true
		}
	}
	return false
}

func Remove(u string) {
	for i, user := range users {
		if user.Name == u {
			users = append(users[0:i], users[i+1:]...)
			break
		}
	}
	Write()
}

func Modify(o, u, p string) {
	re := false
	for _, user := range users {
		re = user.Name == o
		if re {
			user.Name = u
			user.Passwd = p
			break
		}
	}
	if !re {
		user := &User{};
		user.Name = u
		user.Passwd = p
		users = append(users, user)
	}
	Write()
}
