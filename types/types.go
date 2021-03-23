package types

type IndexData struct {
	Sections []Section
	IsLogin  bool
}

type RegisterStruct struct {
	ErrRegister bool
	ErrText     string
	IsLogin 		bool
	Csrf        string
}

type LoginPostStruct struct {
	NameOrEmail string
	Password    string
}

type LoginGetStruct struct {
	ErrLogin bool
	ErrText  string
	IsLogin  bool
	Csrf     string
}

type MyUserProfilStruct struct {
	Access    bool
	IsLogin   bool
	User			User
}

type UpdateUserSettingsFormStruct struct {
	Age      string
	Site     string
	Email    string
	PubEmail string
	AboutMe  string
}

type AddProductPageStruct struct {
	Sections []Section
	IsLogin  bool
	UserId   int
	Csrf     string
}