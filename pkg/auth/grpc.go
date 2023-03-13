package auth

import (
	"context"
	"strings"

	grpc_permission "github.com/galexrt/arpanet/pkg/grpc/permission"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/proto/resources/users"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AuthAccIDCtxTag              = "auth.accid"
	AuthActiveCharIDCtxTag       = "auth.chrid"
	AuthActiveCharJobCtxTag      = "auth.chrjob"
	AuthActiveCharJobGradeCtxTag = "auth.chrjobg"
	AuthSubCtxTag                = "auth.sub"
)

func GRPCAuthFunc(ctx context.Context, fullMethod string) (context.Context, error) {
	token, err := GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}

	// Parse token only returns the token info when the token is still valid
	tokenInfo, err := Tokens.ParseWithClaims(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	grpc_ctxtags.Extract(ctx).Set(AuthAccIDCtxTag, tokenInfo.AccountID)
	grpc_ctxtags.Extract(ctx).Set(AuthActiveCharIDCtxTag, tokenInfo.ActiveCharID)
	grpc_ctxtags.Extract(ctx).Set(AuthActiveCharJobCtxTag, tokenInfo.ActiveCharJob)
	grpc_ctxtags.Extract(ctx).Set(AuthActiveCharJobGradeCtxTag, tokenInfo.ActiveCharJobGrade)
	grpc_ctxtags.Extract(ctx).Set(AuthSubCtxTag, tokenInfo.Subject)

	return ctx, nil
}

func GRPCPermissionUnaryFunc(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	// Check if the method is from a service otherwise the request must be invalid
	if strings.HasPrefix(info.FullMethod, "/services") {
		values := grpc_ctxtags.Extract(ctx).Values()
		if activeCharID, ok := values[AuthActiveCharIDCtxTag].(int32); ok {
			split := strings.Split(info.FullMethod[10:], ".")
			perm := strings.Join(split[1:], "-")

			if overrideSrv, ok := info.Server.(grpc_permission.GetPermsRemapFunc); ok {
				remap := overrideSrv.GetPermsRemap()
				if _, ok := remap[perm]; ok {
					perm = remap[perm]
				}
			}

			if perms.P.CanID(activeCharID, perm) {
				return ctx, nil
			}
		}
	}

	return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to do that! Permission: "+info.FullMethod)
}

func GRPCPermissionStreamFunc(ctx context.Context, srv interface{}, info *grpc.StreamServerInfo) (context.Context, error) {
	// Check if the method is from a service otherwise the request must be invalid
	if strings.HasPrefix(info.FullMethod, "/services") {
		values := grpc_ctxtags.Extract(ctx).Values()
		if activeCharID, ok := values[AuthActiveCharIDCtxTag].(int32); ok {
			split := strings.Split(info.FullMethod[10:], ".")
			perm := strings.Join(split[1:], "-")

			if overrideSrv, ok := srv.(grpc_permission.GetPermsRemapFunc); ok {
				remap := overrideSrv.GetPermsRemap()
				if _, ok := remap[perm]; ok {
					perm = remap[perm]
				}
			}

			if perms.P.CanID(activeCharID, perm) {
				return ctx, nil
			}
		}
	}

	return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to do that! Permission: "+info.FullMethod)
}

func GetTokenFromGRPCContext(ctx context.Context) (string, error) {
	return grpc_auth.AuthFromMD(ctx, "bearer")
}

func MustGetTokenFromGRPCContext(ctx context.Context) string {
	token, _ := GetTokenFromGRPCContext(ctx)
	return token
}

func GetUserIDFromContext(ctx context.Context) int32 {
	values := grpc_ctxtags.Extract(ctx).Values()

	return values[AuthActiveCharIDCtxTag].(int32)
}

func GetUserInfoFromContext(ctx context.Context) (int32, string, int32) {
	values := grpc_ctxtags.Extract(ctx).Values()

	return values[AuthActiveCharIDCtxTag].(int32), values[AuthActiveCharJobCtxTag].(string), values[AuthActiveCharJobGradeCtxTag].(int32)
}

func GetUserFromContext(ctx context.Context) (*users.UserShort, error) {
	return getUserByID(ctx, GetUserIDFromContext(ctx))
}

func getUserByID(ctx context.Context, userID int32) (*users.UserShort, error) {
	// Find user info for the new/old char index in the claim
	u := table.Users.AS("shortuser")
	stmt := u.SELECT(
		u.ID,
		u.Identifier,
		u.Job,
		u.JobGrade,
		u.Firstname,
		u.Lastname,
	).
		FROM(u).
		WHERE(u.ID.EQ(jet.Int32(userID))).
		LIMIT(1)

	var user users.UserShort
	if err := stmt.QueryContext(ctx, query.DB, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
