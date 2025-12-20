package core

import "context"

type requesterKeyType struct{}

var requesterKey = requesterKeyType{}

type Requester interface {
	GetUserID() int
	GetRole() string
}

type requesterData struct {
	UserID int
	Role   string
}

func NewRequester(userID int, role string) *requesterData {
	return &requesterData{
		UserID: userID,
		Role:   role,
	}
}

func (r *requesterData) GetUserID() int {
	return r.UserID
}

func (r *requesterData) GetRole() string {
	return r.Role
}

func GetRequester(ctx context.Context) Requester {
	if requester, ok := ctx.Value(requesterKey).(Requester); ok {
		return requester
	}
	return nil
}

func ContextWithRequester(ctx context.Context, requester Requester) context.Context {
	return context.WithValue(ctx, requesterKey, requester)
}
