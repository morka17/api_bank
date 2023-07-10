package gapi

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/morka17/shiny_bank/v1/pb"
	mockdb "github.com/morka17/shiny_bank/v1/src/db/mock"
	db "github.com/morka17/shiny_bank/v1/src/db/sqlc"
	"github.com/morka17/shiny_bank/v1/src/utils"
	"github.com/morka17/shiny_bank/v1/src/worker"
	mockwk "github.com/morka17/shiny_bank/v1/src/worker/mock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
)

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (expected eqCreateUserTxParamsMatcher) Matches(X interface{}) bool {
	actualArg, ok := X.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(expected.password, actualArg.HashedPassword)
	if err != nil {
		return false
	}
	expected.arg.HashedPassword = actualArg.HashedPassword
	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	err = actualArg.AfterCreate(expected.user)
	if err != nil {
		return false
	}

	return true

}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	userTestCase := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "Created",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username:       user.Username,
						HashedPassword: user.HashedPassword,
						FullName:       user.FullName,
						Email:          user.Email,
					},
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{User: user}, nil)

				taskPayload := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)

			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res.GetUser())

				createdUser := res.GetUser()
				assert.Equal(t, user.Username, createdUser.Username)
				assert.Equal(t, user.FullName, createdUser.FullName)
				assert.Equal(t, user.Email, createdUser.Email)
			},
		},
		{
			name: "Internal",
			
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.CreateUserTxResult{}, sql.ErrConnDone)


				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0).
					Return(nil)

			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				assert.Error(t, err)
				_, ok := status.FromError(err)
				assert.True(t, ok)
				// assert.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range userTestCase {
		tc := userTestCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dCtrl := gomock.NewController(t)
			defer dCtrl.Finish()

			store := mockdb.NewMockStore(ctrl) // store 
			taskDistributor := mockwk.NewMockTaskDistributor(dCtrl) 
			tc.buildStubs(store, taskDistributor)

			// start test server and send request
			server := newTestServer(t, store, taskDistributor)
			res, err := server.CreateUser(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})

	}

}

func randomUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	user = db.User{
		Username:       utils.RandomOwner(),
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	assert.NoError(t, err)

	var getUser db.User
	err = json.Unmarshal(data, &getUser)

	assert.NoError(t, err)
	assert.Equal(t, user.Username, getUser.Username)
	assert.Equal(t, user.FullName, getUser.FullName)
	assert.Equal(t, user.Email, getUser.Email)
	assert.NotEmpty(t, user.HashedPassword)

}
