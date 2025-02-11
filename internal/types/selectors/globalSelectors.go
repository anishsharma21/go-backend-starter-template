package selectors

type indexPage struct {
	AddUserButton     string
	GetUserHTMLButton string
	GetUserJSONButton string
	DeleteUsersButton string
	UsersList         string
}

var IndexPage = indexPage{
	AddUserButton:     "index-add-user-button",
	GetUserHTMLButton: "index-get-user-html-button",
	GetUserJSONButton: "index-get-user-json-button",
	DeleteUsersButton: "index-delete-users-button",
	UsersList:         "index-user-list",
}
