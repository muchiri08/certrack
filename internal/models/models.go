package models

import "database/sql"

type Models struct {
	Users UserModel
	Certs CertModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Users: UserModel{db},
		Certs: CertModel{db},
	}
}
