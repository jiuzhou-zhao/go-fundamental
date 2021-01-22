package meta

import (
	"context"
	"fmt"
	"math/rand"

	"google.golang.org/grpc/metadata"
)

const (
	// RequestIdOnMetaData unique request id
	RequestIdOnMetaData = "ymi_micro_srv_req_id"
)

func randomNumber() uint64 {
	return uint64(rand.Int63())
}

func getRandomID() string {
	return fmt.Sprintf("%x", randomNumber())
}

func GetRequestIDFromMD(md metadata.MD) string {
	ids := md.Get(RequestIdOnMetaData)
	if len(ids) > 0 {
		return ids[0]
	}
	return ""
}

func IdFromIncomingContext(ctx context.Context) string {
	id := ""

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		id = GetRequestIDFromMD(md)
	}

	if id == "" {
		id = getRandomID()
	}
	return id
}

func IdToOutgoingContext(ctx context.Context, id string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, RequestIdOnMetaData, id)
}

func IdFromOutgoingContext(ctx context.Context) string {
	id := ""

	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		id = GetRequestIDFromMD(md)
	}

	if id == "" {
		id = getRandomID()
	}
	return id
}
