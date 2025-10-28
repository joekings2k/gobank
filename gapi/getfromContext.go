package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
)



func getFromContext(ctx context.Context) (userAgent string, clientIP string){
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if ua := md.Get("user-agent"); len(ua) > 0 {
			userAgent = ua[0]
		}
		if ci := md.Get("x-client-ip"); len(ci) > 0 {
			clientIP = ci[0]
		}
	}
	return
}