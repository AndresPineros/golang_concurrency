package main

import (
	"context"
	"fmt"
)

/*
To write to a context, each package needs to define their own types.
This way there can't be collisions with other packages.

Imagine a package writes to the "user" key and then another package also writes to "user",
but they meant different users. It would break the program.
*/

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func UserID(c context.Context) string {
	return c.Value(ctxUserID).(string)
}

func AuthToken(c context.Context) string {
	return c.Value(ctxAuthToken).(string)
}

func ProcessRequest(userID, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"handling response for %v (auth: %v)",
		UserID(ctx),
		AuthToken(ctx),
	)
}

func main() {
	ProcessRequest("jane", "abc123")
}
