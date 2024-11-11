package model

type Response struct {
	Status     int         `json:"-"`
	RC         string      `json:"rc,omitempty"`
	Message    string      `json:"message,omitempty"`
	Data       any         `json:"data,omitempty"`
	Errors     []ErrorData `json:"error,omitempty"`
	Extref     string      `json:"ext_ref,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type ErrorData struct {
	RC          string      `json:"-"`
	Status      int         `json:"-"`
	Description string      `json:"description,omitempty"`
	Title       string      `json:"title,omitempty"`
	Field       string      `json:"field,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	Tag         string      `json:"tag,omitempty"`
	TagValue    string      `json:"tag_value,omitempty"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Pagination struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Search string `json:"search,omitempty"`
	Sort   string `json:"sort,omitempty"`
	Total  int    `json:"total"`
}
