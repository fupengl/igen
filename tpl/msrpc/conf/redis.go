package conf

import "fmt"

func RNamespace(s string) string {
	return fmt.Sprintf("igen.msdemo.%s", s)
}

func RIdentifyCode(app string, typ int, phone string) string {
	return RNamespace(fmt.Sprintf("sms.%s.%d.%s", app, typ, phone))
}

func RToken(id string, platform string) string {
	return RNamespace(fmt.Sprintf("token.%s.%s", id, platform))
}
