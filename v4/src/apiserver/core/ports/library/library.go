package library

import "context"

type Config struct {
	Host            string
	Port            string
	MaxErrorsTrying int64 `mapstructure:"max_errors_trying"`
}

type Client interface {
	GetLibraries(context.Context, string, uint64, uint64) (Infos, error)
	GetLibrariesByIDs(context.Context, []string) (Infos, error)
	GetBooks(context.Context, string, bool, uint64, uint64) (Books, error)
	GetBooksByIDs(context.Context, []string) (Books, error)
	TakeBook(context.Context, string, string) (ReservedBook, error)
	ReturnBook(context.Context, string, string) (Book, error)
}
