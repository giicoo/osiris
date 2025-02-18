package repository

type Repo interface {
	Connection() error
	CloseConnection() error
}
