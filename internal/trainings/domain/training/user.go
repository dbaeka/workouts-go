package training

import (
	"fmt"

	commonerrors "github.com/dbaeka/workouts-go/internal/common/errors"
	"github.com/pkg/errors"
)

// UserType is enum-like type
// By using struct with hidden field, we ensure immutability instead of declaring string as a UserType and using constants
type UserType struct {
	s string
}

func (u UserType) IsZero() bool {
	return u == UserType{}
}

func (u UserType) String() string {
	return u.s
}

var (
	Trainer  = UserType{"trainer"}
	Attendee = UserType{"attendee"}
)

func NewUserTypeFromString(userType string) (UserType, error) {
	switch userType {
	case "trainer":
		return Trainer, nil
	case "attendee":
		return Attendee, nil
	}

	return UserType{}, commonerrors.NewSlugError(
		fmt.Sprintf("invalid '%s' role", userType),
		"invalid-role",
	)
}

type User struct {
	userUUID string
	userType UserType
}

func (u User) UUID() string {
	return u.userUUID
}

func (u User) Type() UserType {
	return u.userType
}

func (u User) IsEmpty() bool {
	return u == User{}
}

func NewUser(userUUID string, userType UserType) (User, error) {
	if userUUID == "" {
		return User{}, errors.New("missing user UUID")
	}
	if userType.IsZero() {
		return User{}, errors.New("missing user type")
	}
	return User{userUUID: userUUID, userType: userType}, nil
}

func MustNewUser(userUUID string, userType UserType) User {
	u, err := NewUser(userUUID, userType)
	if err != nil {
		panic(err)
	}

	return u
}

type ForbiddenToSeeTrainingError struct {
	RequestingUserUUID string
	TrainingOwnerUUID  string
}

func (f ForbiddenToSeeTrainingError) Error() string {
	return fmt.Sprintf(
		"user %q can't see user %q training",
		f.RequestingUserUUID, f.TrainingOwnerUUID,
	)
}

func CanUserSeeTraining(user User, training Training) error {
	if user.Type() == Trainer {
		return nil
	}
	if user.UUID() == training.UserUUID() {
		return nil
	}

	return ForbiddenToSeeTrainingError{user.UUID(), training.UserUUID()}
}
