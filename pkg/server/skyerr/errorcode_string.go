// Code generated by "stringer -type=ErrorCode"; DO NOT EDIT

package skyerr

import "fmt"

const (
	_ErrorCode_name_0 = "NotAuthenticatedPermissionDeniedAccessKeyNotAcceptedAccessTokenNotAcceptedInvalidCredentialsInvalidSignatureBadRequestInvalidArgumentDuplicatedResourceNotFoundNotSupportedNotImplementedConstraintViolatedIncompatibleSchemaAtomicOperationFailurePartialOperationFailureUndefinedOperationPluginUnavailablePluginTimeoutRecordQueryInvalidPluginInitializing"
	_ErrorCode_name_1 = "UnexpectedErrorUnexpectedUserInfoNotFoundUnexpectedUnableToOpenDatabaseUnexpectedPushNotificationNotConfiguredInternalQueryInvalid"
)

var (
	_ErrorCode_index_0 = [...]uint16{0, 16, 32, 52, 74, 92, 108, 118, 133, 143, 159, 171, 185, 203, 221, 243, 266, 284, 301, 314, 332, 350}
	_ErrorCode_index_1 = [...]uint8{0, 15, 41, 71, 110, 130}
)

func (i ErrorCode) String() string {
	switch {
	case 101 <= i && i <= 121:
		i -= 101
		return _ErrorCode_name_0[_ErrorCode_index_0[i]:_ErrorCode_index_0[i+1]]
	case 10000 <= i && i <= 10004:
		i -= 10000
		return _ErrorCode_name_1[_ErrorCode_index_1[i]:_ErrorCode_index_1[i+1]]
	default:
		return fmt.Sprintf("ErrorCode(%d)", i)
	}
}
