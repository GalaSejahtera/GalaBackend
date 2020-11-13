package constants

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidArgumentError    = status.Error(codes.InvalidArgument, "Invalid input.")
	InvalidPhoneNumberError = status.Error(codes.InvalidArgument, "Invalid phone number.")
	InvalidEmailError       = status.Error(codes.InvalidArgument, "Invalid email.")
	InvalidRoleError        = status.Error(codes.InvalidArgument, "Invalid role.")
	InvalidDateError        = status.Error(codes.InvalidArgument, "Invalid date.")
	ConsentNotSignedError   = status.Error(codes.InvalidArgument, "Please ensure that you have signed the consent form.")
	InvalidPasswordError    = status.Error(codes.InvalidArgument, "Invalid password, please ensure that password is more than 6 characters.")

	UserAlreadyExistError        = status.Error(codes.AlreadyExists, "User already exist!")
	ZoneAlreadyExistError        = status.Error(codes.AlreadyExists, "Zone already exist!")
	FaqAlreadyExistError         = status.Error(codes.AlreadyExists, "Faq already exist!")
	PhoneNumberAlreadyExistError = status.Error(codes.AlreadyExists, "Phone number already exist, please use another phone number.")
	EmailAlreadyExistError       = status.Error(codes.AlreadyExists, "Email already exist, please use another email.")

	UserNotFoundError     = status.Error(codes.NotFound, "User not found!")
	ZoneNotFoundError     = status.Error(codes.NotFound, "Zone not found!")
	FaqNotFoundError      = status.Error(codes.NotFound, "Faq not found!")
	ActivityNotFoundError = status.Error(codes.NotFound, "Activity not found!")
	MetadataNotFoundError = status.Error(codes.NotFound, "Metadata not found!")

	UserOperationError         = status.Error(codes.Internal, "Authentication Service failed. Might be due to invalid input.")
	UnauthorizedAccessError    = status.Error(codes.Unauthenticated, "User is not authorized to perform this action!")
	InvalidPasswordVerifyError = status.Error(codes.InvalidArgument, "Invalid password.")
	CreateTokenFailedError     = status.Error(codes.Internal, "Failed to create token.")
	VerifyTokenFailedError     = status.Error(codes.Internal, "Failed to verify token.")

	OperationUnsupportedError = status.Error(codes.Unimplemented, "Operation unsupported.")
	InternalError             = status.Error(codes.Internal, "Server unavailable, please try again.")
)
