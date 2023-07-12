
package gapi

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"bytes/google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/golang/mock/gomock"
	"github.com/morka17/shiny_bank/v1/pb"
	mockdb "github.com/morka17/shiny_bank/v1/src/db/mock"
	db "github.com/morka17/shiny_bank/v1/src/db/sqlc"
	"github.com/morka17/shiny_bank/v1/src/token"
	"github.com/morka17/shiny_bank/v1/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserAPI(t *testing.T) {
	user, _ := randomUser(t)

	newName := utils.RandomOwner()
	newEmail := utils.RandomEmail()
	invalidEmail := "invalid-email"

	userTestCase := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "Ok",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: sql.NullString{
						String: newName,
						Valid:  true,
					},
					Email: sql.NullString{
						String: newEmail,
						Valid:  true,
					},
				}
				updatedUser := db.User{
					Username:         user.Username,
					HashedPassword:   user.HashedPassword,
					FullName:         newName,
					PasswordChangeAt: user.PasswordChangeAt,
					CreatedAt:        user.CreatedAt,
					IsEmailVerified:  user.IsEmailVerified,
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedUser, nil)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {

				return newContextWithBearer(t, tokenMaker, user.Username, time.Minute)
			
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)

				createdUser := res.GetUser()
				assert.Equal(t, user.Username, createdUser.Username)
				assert.Equal(t, user.FullName, newEmail)
				assert.Equal(t, user.Email, newEmail)
			},
		},
		{
			name: "INVALIDEMAIL",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &invalidEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: sql.NullString{
						String: newName,
						Valid:  true,
					},
					Email: sql.NullString{
						String: newEmail,
						Valid:  true,
					},
				}


				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(0)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {

				return newContextWithBearer(t, tokenMaker, user.Username, time.Minute)
			
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				assert.NoError(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, st.Code())

			},
		},
		{
			name: "NOTFOUND",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, nil)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {

				return newContextWithBearer(t, tokenMaker, user.Username, time.Minute)
			
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				assert.Error(t, err)

				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.NotFound, st.Code)
			},
		},
		{
			name: "EXPIREDTOKEN",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.User{}, nil)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {

				return newContextWithBearer(t, tokenMaker, user.Username, -time.Minute)
			
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				assert.Error(t, err)

				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.Unauthenticated, st.Code)
			},
		},
		{
			name: "NOAUTHORIZATION",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, nil)

			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {

				return context.Background()
			
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				assert.Error(t, err)

				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.NotFound, st.Code)
			},
		},
		
	}

	for _, tc := range userTestCase {
		// tc := userTestCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl) // store
			tc.buildStubs(store)

			// start test server and send request
			server := newTestServer(t, store, nil)

			ctx := tc.buildContext(t, server.tokenMaker)
			res, err := server.UpdateUser(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})

	}

}
