package model

import (
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	mongoModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

func UserRepositoryToUsersRepositoryMapper(usersRepository []UserRepository) UsersRepository {
	return NewUsersRepository(
		append([]UserRepository{}, usersRepository...),
	)
}

func UsersRepositoryToUsersMapper(usersRepository UsersRepository) userModel.Users {
	users := make([]userModel.User, len(usersRepository.Users))
	for index, userRepository := range usersRepository.Users {
		users[index] = UserRepositoryToUserMapper(userRepository)
	}

	return userModel.NewUsers(
		users,
		usersRepository.PaginationResponse,
	)
}

func UserRepositoryToUserMapper(userRepository UserRepository) userModel.User {
	return userModel.NewUser(
		userRepository.ID.Hex(),
		userRepository.CreatedAt,
		userRepository.UpdatedAt,
		userRepository.Name,
		userRepository.Email,
		userRepository.Password,
		userRepository.Role,
		userRepository.Verified,
	)
}

func UserResetExpiryRepositoryToUserResetExpiryMapper(userResetExpiryRepository UserResetExpiryRepository) userModel.UserResetExpiry {
	return userModel.NewUserResetExpiry(
		userResetExpiryRepository.ResetExpiry,
	)
}

func UserCreateToUserCreateRepositoryMapper(userCreate userModel.UserCreate) UserCreateRepository {
	return NewUserCreateRepository(
		userCreate.Name,
		userCreate.Email,
		userCreate.Password,
		userCreate.Role,
		userCreate.Verified,
		userCreate.VerificationCode,
		userCreate.CreatedAt,
		userCreate.UpdatedAt,
	)
}

func UserUpdateToUserUpdateRepositoryMapper(logger model.Logger, location string, userUpdate userModel.UserUpdate) common.Result[UserUpdateRepository] {
	userObjectID := mongoModel.HexToObjectIDMapper(logger, location+".UserUpdateToUserUpdateRepositoryMapper", userUpdate.ID)
	if validator.IsError(userObjectID.Error) {
		return common.NewResultOnFailure[UserUpdateRepository](userObjectID.Error)
	}

	return common.NewResultOnSuccess(NewUserUpdateRepository(
		userObjectID.Data,
		userUpdate.Name,
		userUpdate.UpdatedAt,
	))
}

func UserForgottenPasswordToUserForgottenPasswordRepositoryMapper(userForgottenPassword userModel.UserForgottenPassword) UserForgottenPasswordRepository {
	return NewUserForgottenPasswordRepository(
		userForgottenPassword.ResetToken,
		userForgottenPassword.ResetExpiry,
	)
}

func UserResetPasswordToUserResetPasswordRepositoryMapper(userResetPassword userModel.UserResetPassword) UserResetPasswordRepository {
	return NewUserResetPasswordRepository(
		userResetPassword.Password,
	)
}
