package selectors

type indexPage struct {
	BaseHTML        string
	AuthComponent   string
	LoginComponent  string
	SignUpComponent string
}

var IndexPage = indexPage{
	BaseHTML:        "base-html",
	AuthComponent:   "auth-component",
	LoginComponent:  "login-component",
	SignUpComponent: "signup-component",
}

type usersPage struct {
	UsersView         string
	AddUserButton     string
	GetUserHTMLButton string
	GetUserJSONButton string
	DeleteUsersButton string
	UsersList         string
}

var UsersPage = usersPage{
	UsersView:         "users-view",
	AddUserButton:     "index-add-user-button",
	GetUserHTMLButton: "index-get-user-html-button",
	GetUserJSONButton: "index-get-user-json-button",
	DeleteUsersButton: "index-delete-users-button",
	UsersList:         "index-user-list",
}
