package isaerrors

import "errors"

var ErrDifferentFundNotAllowed = errors.New("customers can only invest in one fund at this time")
