package grpcws

import "net/http"

type Trailer struct {
	trailer
}

func HTTPTrailerToGrpcWebTrailer(httpTrailer http.Header) Trailer {
	return Trailer{trailer{httpTrailer}}
}
