package dotest

import "fmt"

/**
**
**
** 一种多态
**
**
 */
type Duck interface {
	LayEgg(amount int)
	BirthEgg(amount int)
}

type MaleDuck struct {
}

func NewMaleDuck() Duck {
	return &MaleDuck{}
}

func (d *MaleDuck) LayEgg(amount int) {
	fmt.Println("maleDuck LayEgg", amount)
}

func (d *MaleDuck) BirthEgg(amount int) {
	fmt.Println("maleDuck BirthEgg", amount)
}

type FemaleDuck struct {
}

func NewFemaleDuck() Duck {
	return &FemaleDuck{}
}

func (d *FemaleDuck) LayEgg(amount int) {
	fmt.Println("femaleDuck LayEgg", amount)
}

func (d *FemaleDuck) BirthEgg(amount int) {
	fmt.Println("femaleDuck BirthEgg", amount)
}

/**
**
**
** 一种多态
**
**
 */
type NewDuck interface {
	LayEgg(int)
	BirthEgg(int)
}

func LayEgg(a NewDuck, amount int) {
	a.LayEgg(amount)
}

func BirthEgg(a NewDuck, amount int) {
	a.BirthEgg(amount)
}

