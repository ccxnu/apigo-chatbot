package domain

type ILogs interface {
	SaveRequest(request any, method string, className string)
	SaveException(request any, method string, className string, err error)
	SaveResponse(response any, method string, className string)
}
