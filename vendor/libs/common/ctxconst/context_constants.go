package ctxconst

import (
	"context"
	"fmt"
)

type CtxKey string

const (
	CtxKeyRequestID         CtxKey = "ctx_request_id"
	CtxKeyUserID            CtxKey = "ctx_user_id"
	CtxKeyUserPhoneNumber   CtxKey = "ctx_user_phone_number"
	CtxKeyLanguage          CtxKey = "ctx_language"
	CtxKeyMethodName        CtxKey = "ctx_key_method_name"
	CtxKeyKeycloakRealmName CtxKey = "ctx_key_keycloak_realm_name"

	LanguageRU = "ru"
	LanguageEN = "en"
	LanguageKK = "kk"
)

func GetRequestID(ctx context.Context) interface{} {
	return ctx.Value(CtxKeyRequestID)
}

func SetRequestID(ctx context.Context, value interface{}) context.Context {
	return context.WithValue(ctx, CtxKeyRequestID, value)
}

func GetUserID(ctx context.Context) interface{} {
	return ctx.Value(CtxKeyUserID)
}

func SetUserID(ctx context.Context, value interface{}) context.Context {
	return context.WithValue(ctx, CtxKeyUserID, value)
}

func GetUserPhoneNumber(ctx context.Context) interface{} {
	return ctx.Value(CtxKeyUserPhoneNumber)
}

func SetUserPhoneNumber(ctx context.Context, value interface{}) context.Context {
	return context.WithValue(ctx, CtxKeyUserPhoneNumber, value)
}

func GetLanguage(ctx context.Context) string {
	if ctx.Value(CtxKeyLanguage) != nil {
		return fmt.Sprintf("%v", ctx.Value(CtxKeyLanguage))
	}

	return LanguageRU
}

func SetLanguage(ctx context.Context, value interface{}) context.Context {
	return context.WithValue(ctx, CtxKeyLanguage, value)
}

func GetMethodName(ctx context.Context) string {
	if ctx.Value(CtxKeyMethodName) != nil {
		return fmt.Sprintf("%v", ctx.Value(CtxKeyMethodName))
	}

	return ""
}

func SetMethodName(ctx context.Context, value interface{}) context.Context {
	return context.WithValue(ctx, CtxKeyMethodName, value)
}

func GetKeycloakRealmName(ctx context.Context) interface{} {
	return ctx.Value(CtxKeyKeycloakRealmName)
}

func SetKeycloakRealmName(ctx context.Context, value interface{}) context.Context {
	return context.WithValue(ctx, CtxKeyKeycloakRealmName, value)
}

func Copy(ctx context.Context) context.Context {
	newCtx := context.Background()
	newCtx = SetRequestID(newCtx, GetRequestID(ctx))
	newCtx = SetUserID(newCtx, GetUserID(ctx))
	newCtx = SetUserPhoneNumber(newCtx, GetUserPhoneNumber(ctx))
	newCtx = SetLanguage(newCtx, GetLanguage(ctx))

	return newCtx
}

func GetContextValues(ctx context.Context) map[string]interface{} {
	ctxValues := make(map[string]interface{})

	keys := []struct {
		key    CtxKey
		getter func(context.Context) interface{}
	}{
		{CtxKeyRequestID, GetRequestID},
		{CtxKeyUserID, GetUserID},
		{CtxKeyUserPhoneNumber, GetUserPhoneNumber},
		{CtxKeyMethodName, func(ctx context.Context) interface{} { return GetMethodName(ctx) }},
		{CtxKeyLanguage, func(ctx context.Context) interface{} { return GetLanguage(ctx) }},
	}

	for _, k := range keys {
		if value := k.getter(ctx); value != nil {
			ctxValues[string(k.key)] = value
		}
	}

	return ctxValues
}
