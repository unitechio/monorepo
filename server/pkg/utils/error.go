package utils

import "fmt"

type ErrorApp struct {
	Code        int
	Description string
}

var (
	SUCCESS = ErrorApp{1, "SUCCESS"}

	INTERNAL_SERVER_ERROR = ErrorApp{500, "Internal Server Error"}
	BAD_REQUEST           = ErrorApp{400, "Bad Request"}
	UNAUTHORIZED          = ErrorApp{401, "Unauthorized"}
	ACCESS_DENIED         = ErrorApp{403, "Truy cập bị từ chối"}

	INVALID_CHECKSUM = ErrorApp{51, "Checksum error"}
	INVALID_MSISDN   = ErrorApp{51, "Unauthorized"}

	LINK_SAVE_ERROR = ErrorApp{50, "Lưu kết quả liên kết bại"}

	INFO_NOT_EXIST = ErrorApp{22, "Thông tin không tồn tại"}

	INVALID_PASSWORD = ErrorApp{-1, "Mật khẩu phải có độ dài 6-20 kí tự và không chứa khoảng trống"}

	XCG_NO_INFORMATION = ErrorApp{24, "Không lấy được thông tin đăng kiểm"}
	XCG_NOT_ALLOW      = ErrorApp{24, "Không được truy vấn thông tin đăng kiểm"}

	ERR_TYPE_DATA = ErrorApp{-1, "Dữ liệu không đúng định dạng"}

	FILE_UPLOAD_ERROR = ErrorApp{32, "MinIO upload error"}
	FILE_REMOVE_ERROR = ErrorApp{32, "MinIO remove error"}
	FILE_LOAD_ERROR   = ErrorApp{32, "MinIO load error"}
	FILE_BUCKET_ERROR = ErrorApp{32, "MinIO bucket error"}

	INVALID_PARTNER_INFO   = ErrorApp{401, "No valid partner_id+partner_secret pairs are provided"}
	PARTNER_NAME_NOT_MATCH = ErrorApp{401, "Partner name does not match"}

	PARTNER_CONFIG_NOT_EXIST = ErrorApp{63, "Cấu hình mã đối tác không hợp lệ"}

	AUTHORIZATION_CODE_EXPIRED = ErrorApp{64, "Authorization code expired"}

	CALL_REPORT_LOTS_FAILED = ErrorApp{67, "Has error when call api get lot map in report"}

	HAS_ERROR_WHEN_CALL_API_LINK = ErrorApp{70, "Has error when call api link"}

	DATA_LINK_IS_NULL_WHEN_CALL_API_LINK = ErrorApp{71, "Data link is null when call api link"}
)

func (e ErrorApp) WithArgs(args ...any) ErrorApp {
	return ErrorApp{
		Code:        e.Code,
		Description: fmt.Sprintf(e.Description, args...),
	}
}

func (e ErrorApp) Error() string {
	return e.Description
}

func (e ErrorApp) Is(target ErrorApp) bool {
	return e.Code == target.Code
}
