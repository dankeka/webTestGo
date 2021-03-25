package types

type Users []User

type Sections []Section

type Products []Product

func (s Users) Append(o User) Users {
	return append(s, o)
}

func (s Sections) Append(o Section) Sections {
	return append(s, o)
}

func (s Products) Append(o Product) Products {
	return append(s, o)
}