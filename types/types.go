package types

import "mime/multipart"

type IndexData struct {
	Sections []Section
	NewProducts []ProductIdAndTitleAndImg
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
	Csrf      string
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

type AddProductPost struct {
	Title       string
	SectionID   string
	Description string
	Files				[]*multipart.FileHeader
	Price 			string
}

type ProductAndImg struct {
	Product
	ImgUrl string
}

type OpenProductStruct struct {
	Product
	SellerInfo UserAvaAndName
	ImgUrls    []string
	IsLogin    bool
}

type MyProductsStruct struct {
	IsLogin  bool
	Products []ProductAndImg
}

type UserAvaAndName struct {
	Name   string
	AvaUrl string
}

type OpenUserAccStruct struct {
	IsLogin bool
	User    User
}