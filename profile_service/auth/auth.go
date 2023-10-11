package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/Arovelti/identityhub/profile_service/models"
	"github.com/Arovelti/identityhub/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	protectedMethods = map[string]struct{}{
		"/pb.profile.Profiles/CreateProfile":     {},
		"/pb.profile.Profiles/UpdateProfile":     {},
		"/pb.profile.Profiles/DeleteProfileByID": {},

		// Non protected:
		// "/pb.profile.Profiles/ListProfiles": {},
	}

	ErrAuthFailed               = errors.New("authentication failed")
	ErrNoCredentials            = errors.New("no credentials provided")
	ErrInvalidCredentialsFormat = errors.New("invalid credentials format")
	ErrInvalidAuthHeader        = errors.New("invalid authorization header")
)

func BasicAuthInterceptor(repo repository.Repository) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		//	var username, password string

		username, password, err := extractCredentialsFromMetadata(ctx)
		if err != nil {
			return nil, fmt.Errorf("basic auth interceptor: %v", err)
		}

		// if md, ok := metadata.FromIncomingContext(ctx); ok {
		// 	authHeader := strings.Join(md.Get("authorization"), "")
		// 	if strings.HasPrefix(authHeader, "Basic ") {
		// 		creds, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		// 		if err == nil {
		// 			credsSlice := strings.SplitN(string(creds), ":", 2)
		// 			if len(credsSlice) == 2 {
		// 				username = credsSlice[0]
		// 				password = credsSlice[1]
		// 			}
		// 		}
		// 	}
		// }

		// Perform user authentication
		user, err := authenticate(repo, username, password)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Authentication failed")
		}

		// Check if the method is protected
		if _, ok := protectedMethods[info.FullMethod]; ok {
			if !user.Admin {
				return nil, status.Error(codes.PermissionDenied, "Permission denied")
			}
		}

		// Call the gRPC handler for the actual service method
		return handler(ctx, req)
	}
}

func authenticate(repo repository.Repository, username, password string) (*models.Profile, error) {
	p, err := repo.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %v", err)
	}

	if p.Name == username && p.Password == password {
		return p, nil
	}

	return nil, ErrAuthFailed
}

func extractCredentialsFromMetadata(ctx context.Context) (username, password string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", ErrNoCredentials
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return "", "", ErrNoCredentials
	}

	authHeader := authHeaders[0]

	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", "", ErrInvalidAuthHeader
	}

	creds, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
	if err != nil {
		return "", "", err
	}

	credsSlice := strings.SplitN(string(creds), ":", 2)
	if len(credsSlice) != 2 {
		return "", "", ErrInvalidCredentialsFormat
	}

	return credsSlice[0], credsSlice[1], nil
}
