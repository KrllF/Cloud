package entity

// ListUserOptions список опций
type ListUserOptions struct {
	TokenSize int64
}

// ListUserOption функция для установки опций
type ListUserOption func(*ListUserOptions)

// WithTokenSize токен сайз
func WithTokenSize(size int64) ListUserOption {
	return func(opts *ListUserOptions) {
		opts.TokenSize = size
	}
}
