package model

import (
	"context"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/utility"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DisableInactiveUsers disable inactive users
func (m *Model) DisableInactiveUsers(ctx context.Context) error {
	// get all users
	_, users, err := m.userDAO.Query(ctx, nil, nil, &dto.FilterData{
		Item:  constants.Role,
		Value: constants.User,
	})
	if err != nil {
		return err
	}

	// check user last updated time, set users whose updated time > 20 minutes to inactive
	for _, user := range users {
		if utility.TimeToMilli(utility.MalaysiaTime(time.Now()))-user.LastUpdated > 10*60*1000 {
			user.IsActive = false
			_, err = m.userDAO.Update(ctx, user)
			if err != nil {
				logger.Log.Error(err.Error())
			}
		}
	}
	return nil
}

// CreateUser creates new user
func (m *Model) CreateUser(ctx context.Context, user *dto.User) (*dto.User, error) {

	// check if User exist
	_, err := m.userDAO.Get(ctx, user.ID)

	// only can create user if not found
	if err != nil && status.Code(err) == codes.Unknown {
		return m.userDAO.Create(ctx, user)
	}

	if err != nil {
		return nil, err
	}

	return nil, status.Error(codes.AlreadyExists, "User already exist!")
}

// ClientUpdateUser updates user by client
func (m *Model) ClientUpdateUser(ctx context.Context, user *dto.User) (*dto.User, error) {

	// check if user exists
	u, err := m.userDAO.Get(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// patch user
	u.IsActive = user.IsActive
	u.Lat = user.Lat
	u.Long = user.Long
	u.LastUpdated = user.LastUpdated
	u.Name = user.Name

	_, err = m.userDAO.Update(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// UpdateUser updates user
func (m *Model) UpdateUser(ctx context.Context, user *dto.User) (*dto.User, error) {

	// check if user exists
	u, err := m.userDAO.Get(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// patch user
	oldEmail := u.Email
	u.Role = user.Role
	u.Email = user.Email
	u.IsActive = user.IsActive
	u.Lat = user.Lat
	u.Long = user.Long
	u.Name = user.Name

	_, err = m.userDAO.Update(ctx, u)
	if err != nil {
		return nil, err
	}

	// revoke all tokens by user id if email is changed
	if oldEmail != u.Email {
		err = m.authDAO.DeleteByID(ctx, u.ID)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

// UpdateUserPassword updates user password only
func (m *Model) UpdateUserPassword(ctx context.Context, user *dto.User) (*dto.User, error) {

	// check if user exists
	u, err := m.userDAO.Get(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// patch user
	u.Password = user.Password

	// hash password
	u.Password, err = utility.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}

	_, err = m.userDAO.Update(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// UpdateUsers update users
func (m *Model) UpdateUsers(ctx context.Context, user *dto.User, ids []string) ([]string, error) {
	if len(ids) > 1 {
		return nil, constants.OperationUnsupportedError
	}
	user.ID = ids[0]
	u, err := m.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return []string{u.ID}, err
}

// GetUser gets user by ID
func (m *Model) GetUser(ctx context.Context, id string) (*dto.User, error) {
	return m.userDAO.Get(ctx, id)
}

// CreateToken creates token with custom ttl
func (m *Model) CreateToken(ctx context.Context, auth *dto.AuthObject) (*dto.AuthObject, error) {

	// Create Token
	var err error
	atClaims := jwt.MapClaims{}
	atClaims[constants.Authorized] = true
	atClaims[constants.AccessUuid] = uuid.NewV4().String()
	atClaims[constants.UserId] = auth.UserId
	atClaims[constants.DisplayName] = auth.DisplayName
	atClaims[constants.Exp] = auth.TTL
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	auth.Token = token

	_, err = m.authDAO.Create(ctx, auth)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

// GetUserIDByToken gets userID by token
func (m *Model) GetUserIDByToken(ctx context.Context, token string) (string, error) {
	auth, err := m.authDAO.Get(ctx, token)
	if err != nil {
		return "", err
	}
	return auth.UserId, nil
}

// RevokeTokensByUserID revokes all tokens by UserID
func (m *Model) RevokeTokensByUserID(ctx context.Context, id string) error {
	err := m.authDAO.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// BatchGetUsers get users by slice of IDs
func (m *Model) BatchGetUsers(ctx context.Context, ids []string) ([]*dto.User, error) {
	return m.userDAO.BatchGet(ctx, ids)
}

// QueryUsers queries users by sort, range, filter
func (m *Model) QueryUsers(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.User, error) {
	return m.userDAO.Query(ctx, sort, itemsRange, filter)
}

// DeleteUser deletes user by ID
func (m *Model) DeleteUser(ctx context.Context, id string) (*dto.User, error) {
	// check if user exist
	u, err := m.userDAO.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = m.userDAO.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	// revoke all tokens
	err = m.authDAO.DeleteByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// RevokeUserTokens revoke all user tokens
func (m *Model) RevokeUserTokens(ctx context.Context) error {
	// force revoke all user tokens
	_, users, err := m.QueryUsers(ctx, nil, nil, nil)
	for _, user := range users {
		err = m.authDAO.DeleteByID(ctx, user.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteUsers delete users by IDs
func (m *Model) DeleteUsers(ctx context.Context, ids []string) ([]string, error) {
	var deletedIDs []string
	for _, id := range ids {
		u, err := m.DeleteUser(ctx, id)
		if err != nil {
			return nil, err
		}
		deletedIDs = append(deletedIDs, u.ID)
	}

	return deletedIDs, nil
}

// Login verifies user by email and password and return tokens
func (m *Model) Login(ctx context.Context, email string, password string) (*dto.User, error) {
	// get users by email
	_, users, err := m.userDAO.Query(ctx, nil, nil, &dto.FilterData{
		Item:  constants.Email,
		Value: email,
	})
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, constants.UserNotFoundError
	}
	user := users[0]

	// verify password
	if !m.verifyPassword(user.Password, password) {
		return nil, constants.InvalidPasswordVerifyError
	}

	// create token
	ts, err := m.createToken(user)
	if err != nil {
		return nil, constants.CreateTokenFailedError
	}

	// create auth in db
	err = m.createAuth(user.ID, ts)
	if err != nil {
		return nil, err
	}
	user.AccessToken = ts.AccessToken
	user.RefreshToken = ts.RefreshToken

	return user, nil
}

// VerifyUser verifies user by header
func (m *Model) VerifyUser(ctx context.Context, header string) (*dto.User, error) {
	tokenAuth, err := m.extractTokenMetadata(header)
	if err != nil {
		return nil, err
	}
	id, err := m.fetchAuth(tokenAuth)
	if err != nil {
		return nil, err
	}

	// get user by id
	user, err := m.userDAO.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Logout logs user out from the system by header
func (m *Model) Logout(ctx context.Context, header string) error {
	au, err := m.extractTokenMetadata(header)
	if err != nil {
		return err
	}
	err = m.authDAO.DeleteByID(ctx, au.ID)
	if err != nil {
		return err
	}
	return nil
}

// Refresh returns new token to authorized user by header
func (m *Model) Refresh(ctx context.Context, header string) (*dto.User, error) {
	refreshToken := m.extractToken(header)
	// verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constants.VerifyTokenFailedError
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, constants.VerifyTokenFailedError
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims[constants.RefreshUuid].(string)
		if !ok {
			return nil, constants.VerifyTokenFailedError
		}
		userId, ok := claims[constants.UserId].(string)
		if !ok {
			return nil, constants.VerifyTokenFailedError
		}
		// Delete the previous Refresh Token
		err := m.deleteAuth(refreshUuid)
		if err != nil {
			return nil, err
		}

		// get user
		user, err := m.GetUser(ctx, userId)
		if err != nil {
			return nil, err
		}

		// Create new pairs of refresh and access tokens
		ts, err := m.createToken(user)
		if err != nil {
			return nil, err
		}
		// save the tokens metadata to mongodb
		err = m.createAuth(userId, ts)
		if err != nil {
			return nil, err
		}

		user.AccessToken = ts.AccessToken
		user.RefreshToken = ts.RefreshToken

		return user, nil
	} else {
		return nil, constants.VerifyTokenFailedError
	}
}

// verifyPassword is used to verify plain password to database hash password
func (m *Model) verifyPassword(hashPassword string, plainPassword string) bool {
	byteHashPassword := []byte(hashPassword)
	bytePlainPassword := []byte(plainPassword)
	err := bcrypt.CompareHashAndPassword(byteHashPassword, bytePlainPassword)
	if err != nil {
		return false
	}
	return true
}

func (m *Model) createToken(user *dto.User) (*dto.User, error) {
	td := &dto.User{}
	td.AtExpires = time.Now().Add(time.Minute*constants.AccessTokenTTLMinutes).Unix()*1000 - 1000
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour*24*constants.RefreshTokenTTLDays).Unix()*1000 - 1000
	td.RefreshUuid = uuid.NewV4().String()

	// Create Access Token
	var err error
	atClaims := jwt.MapClaims{}
	atClaims[constants.Authorized] = true
	atClaims[constants.AccessUuid] = td.AccessUuid
	atClaims[constants.UserId] = user.ID
	atClaims[constants.Email] = user.Email
	atClaims[constants.Exp] = td.AtExpires
	atClaims[constants.Name] = user.Name
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// Create Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims[constants.RefreshUuid] = td.RefreshUuid
	rtClaims[constants.UserId] = user.ID
	rtClaims[constants.Email] = user.Email
	rtClaims[constants.Exp] = td.RtExpires
	rtClaims[constants.Name] = user.Name
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (m *Model) createAuth(id string, td *dto.User) error {
	accessToken := &dto.AuthObject{
		Token:  td.AccessUuid,
		UserId: id,
		TTL:    utility.MilliToTime(td.AtExpires),
		Type:   constants.Access,
	}
	refreshToken := &dto.AuthObject{
		Token:  td.RefreshUuid,
		UserId: id,
		TTL:    utility.MilliToTime(td.RtExpires),
		Type:   constants.Refresh,
	}
	ctx := context.TODO()
	_, err := m.authDAO.Create(ctx, accessToken)
	if err != nil {
		return err
	}
	_, err = m.authDAO.Create(ctx, refreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) verifyToken(header string) (*jwt.Token, error) {
	tokenString := m.extractToken(header)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constants.VerifyTokenFailedError
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (m *Model) extractTokenMetadata(header string) (*dto.User, error) {
	token, err := m.verifyToken(header)
	if err != nil {
		return nil, err
	}

	// Validate token
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, constants.VerifyTokenFailedError
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims[constants.AccessUuid].(string)
		if !ok {
			return nil, constants.VerifyTokenFailedError
		}
		userId, ok := claims[constants.UserId].(string)
		if !ok {
			return nil, constants.VerifyTokenFailedError
		}
		return &dto.User{
			AccessUuid: accessUuid,
			ID:         userId,
		}, nil
	}
	return nil, constants.VerifyTokenFailedError
}

func (m *Model) fetchAuth(authD *dto.User) (string, error) {
	authObject, err := m.authDAO.Get(context.TODO(), authD.AccessUuid)
	if err != nil {
		return "", err
	}
	return authObject.UserId, nil
}

func (m *Model) extractToken(header string) string {
	strArr := strings.Split(header, " ")
	if len(strArr) == 1 {
		return strArr[0]
	}
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (m *Model) deleteAuth(givenUuid string) error {
	err := m.authDAO.Delete(context.TODO(), givenUuid)
	if err != nil {
		logger.Log.Warn("DeleteAuth", zap.String("error", err.Error()), zap.String("token", givenUuid))
	}
	return nil
}

// GetNearbyUsers get nearby users count given user
func (m *Model) GetNearbyUsers(ctx context.Context, user *dto.User) (int64, []*dto.User, error) {
	total, users, err := m.userDAO.GetNearbyUsers(ctx, user)
	if err != nil {
		return 0, nil, err
	}
	return total, users, nil
}
