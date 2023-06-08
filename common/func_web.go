package comm

import (
	"Lucky_Wheel/conf"
	"Lucky_Wheel/models"
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

func ClientIp(request *http.Request) string {
	host, _, _ := net.SplitHostPort(request.RemoteAddr)
	return host
}

func Redirect(writer http.ResponseWriter, url string) {
	writer.Header().Add("Location", url)
	writer.WriteHeader(http.StatusFound)
}

func GetLoginUser(request *http.Request) *models.ObjLoginUser {
	c, err := request.Cookie("lottery_loginuser")
	if err != nil {
		return nil
	}
	params, err := url.ParseQuery(c.Value)
	if err != nil {
		return nil
	}
	uid, err := strconv.Atoi(params.Get("uid"))
	if err != nil || uid < 1 {
		return nil
	}
	now, err := strconv.Atoi(params.Get("now"))
	if err != nil || NowUnix()-now > 86400*30 {
		return nil
	}
	loginuser := &models.ObjLoginUser{}
	loginuser.Uid = uid
	loginuser.Username = params.Get("username")
	loginuser.Now = now
	loginuser.Ip = ClientIp(request)
	loginuser.Sign = params.Get("sign")
	sign := createLoginUserSign(loginuser)
	if sign != loginuser.Sign {
		log.Println("func_web.GetLoginUser no right sign")
		return nil
	}
	return loginuser
}
func SetLoginUser(writer http.ResponseWriter, loginuser *models.ObjLoginUser) {
	if loginuser == nil || loginuser.Uid < 1 {
		c := &http.Cookie{
			Name:   "lottery_loginuser",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(writer, c)
		return
	}
	if loginuser.Sign == "" {
		loginuser.Sign = createLoginUserSign(loginuser)
	}
	params := url.Values{}
	params.Add("uid", strconv.Itoa(loginuser.Uid))
	params.Add("username", loginuser.Username)
	params.Add("now", strconv.Itoa(loginuser.Now))
	params.Add("ip", loginuser.Ip)
	params.Add("sign", loginuser.Sign)
	c := &http.Cookie{
		Name:  "lottery_loginuser",
		Value: params.Encode(),
		Path:  "/",
	}
	http.SetCookie(writer, c)
}

func createLoginUserSign(loginuser *models.ObjLoginUser) string {
	str := fmt.Sprintf("uid=%d&username=%s&secret=%s&now=%d",
		loginuser.Uid, loginuser.Username, conf.CookieSecret, loginuser.Now)
	sign := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return sign
}
