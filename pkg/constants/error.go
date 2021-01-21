package constants

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidArgumentError = status.Error(codes.InvalidArgument, "Invalid input.")
	InvalidEmailError    = status.Error(codes.InvalidArgument, "Invalid email.")
	InvalidRoleError     = status.Error(codes.InvalidArgument, "Invalid role.")
	InvalidDateError     = status.Error(codes.InvalidArgument, "Invalid date.")
	InvalidPasswordError = status.Error(codes.InvalidArgument, "Invalid password, please ensure that password is more than 6 characters.")

	UserAlreadyExistError        = status.Error(codes.AlreadyExists, "User already exist!")
	ZoneAlreadyExistError        = status.Error(codes.AlreadyExists, "Zone already exist!")
	ReportAlreadyExistError      = status.Error(codes.AlreadyExists, "Report already exist!")
	PhoneNumberAlreadyExistError = status.Error(codes.AlreadyExists, "Phone number already exist, please use another phone number.")
	EmailAlreadyExistError       = status.Error(codes.AlreadyExists, "Email already exist, please use another email.")

	UserNotFoundError     = status.Error(codes.NotFound, "User not found!")
	DailyNotFoundError    = status.Error(codes.NotFound, "Covid-19 cases report not found!")
	ReportNotFoundError   = status.Error(codes.NotFound, "Report not found!")
	CovidNotFoundError    = status.Error(codes.NotFound, "Article not found!")
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
