package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/wenealves10/gobank/db/mocks"
	db "github.com/wenealves10/gobank/db/sqlc"
	"github.com/wenealves10/gobank/utils"
	"go.uber.org/mock/gomock"
)

type epCreateUserParamsMatcher struct {
    arg db.CreateUserParams
    password string
}

func (e epCreateUserParamsMatcher) Matches(x interface{}) bool {
    arg, ok := x.(db.CreateUserParams)
    if !ok {
        return false
    }

    err := utils.CheckPassword(e.password, arg.HashedPassword)
    if err != nil {
        return false
    }

    e.arg.HashedPassword = arg.HashedPassword
    return reflect.DeepEqual(e.arg, arg)
}

func (e epCreateUserParamsMatcher) String() string {
    return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EpCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
    return epCreateUserParamsMatcher{arg: arg, password: password}
}

func TestCreateUser(t *testing.T){
    user, password := randomUser(t)

    testCases := []struct {
        name string
        body gin.H
        buildStubs func(store *mocks.MockStore)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            body: gin.H{
                "username": user.Username,
                "full_name": user.FullName,
                "email": user.Email,
                "password": password,
            },
            buildStubs: func(store *mocks.MockStore) {

                arg := db.CreateUserParams{
                    Username: user.Username,
                    FullName: user.FullName,
                    Email: user.Email,
                }

                store.EXPECT().
                    CreateUser(gomock.Any(), EpCreateUserParams(arg, password)).
                    Times(1).
                    Return(user, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchUser(t, recorder.Body, user)
            },
        },
        {
            name: "InternalError",
            body: gin.H{
                "username": user.Username,
                "full_name": user.FullName,
                "email": user.Email,
                "password": password,
            },
            buildStubs: func(store *mocks.MockStore) {
                store.EXPECT().
                    CreateUser(gomock.Any(), gomock.Any()).
                    Times(1).
                    Return(db.User{}, sql.ErrConnDone)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name: "DuplicateUsername",
            body: gin.H{
                "username": user.Username,
                "full_name": user.FullName,
                "email": user.Email,
                "password": password,
            },
            buildStubs: func(store *mocks.MockStore) {
                store.EXPECT().
                    CreateUser(gomock.Any(), gomock.Any()).
                    Times(1).
                    Return(db.User{}, &pq.Error{Code: "23505"})
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusForbidden, recorder.Code)
            },
        },
        {
            name: "InvalidUsername",
            body: gin.H{
                "username": "invalid-user#",
                "full_name": user.FullName,
                "email": user.Email,
                "password": password,
            },
            buildStubs: func(store *mocks.MockStore) {
                store.EXPECT().
                    CreateUser(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "InvalidEmail",
            body: gin.H{
                "username": user.Username,
                "full_name": user.FullName,
                "email": "invalid-email",
                "password": password,
            },
            buildStubs: func(store *mocks.MockStore) {
                store.EXPECT().
                    CreateUser(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "TooShortPassword",
            body: gin.H{
                "username": user.Username,
                "full_name": user.FullName,
                "email": user.Email,
                "password": "123",
            },
            buildStubs: func(store *mocks.MockStore) {
                store.EXPECT().
                    CreateUser(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "TooLongPassword",
            body: gin.H{
                "username": user.Username,
                "full_name": user.FullName,
                "email": user.Email,
                "password": utils.RandomString(101),
            },
            buildStubs: func(store *mocks.MockStore) {
                store.EXPECT().
                    CreateUser(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
    }

    for i := range testCases {
        tc := testCases[i]
        t.Run(tc.name, func(t *testing.T){
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            store := mocks.NewMockStore(ctrl)
            tc.buildStubs(store)

            server := NewServer(store)
            recorder := httptest.NewRecorder()

            data, err := json.Marshal(tc.body)
            require.NoError(t, err)

            url := "/users"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
    
}

func randomUser(t *testing.T) (user db.User, password string) {
    password = utils.RandomString(6)
    hashedPassword, err := utils.HashPassword(password)
    require.NoError(t, err)

    user = db.User{
        Username: utils.RandomOwner(),
        FullName: utils.RandomOwner(),
        Email: utils.RandomEmail(),
        HashedPassword: hashedPassword,
    }
    return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotUser db.User
    err = json.Unmarshal(data, &gotUser)
    require.NoError(t, err)
    require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}