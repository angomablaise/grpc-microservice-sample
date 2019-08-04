package server

import "context"

func ExportTokenAuthentication(ctx context.Context) (context.Context, error) {
	return tokenAuthentication(ctx)
}
