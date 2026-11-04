package service

import (
 "context"
 "net/http"
 "net/http/httptest"
 "testing"
)
func TestNewRejectsInsecure(t *testing.T) {
 for _,s:=range []string{"http://example.org","file:///tmp/x","https://user:pw@example.org","https://example.org/?token=abc","https://example.org/private","https://example.org/#fragment"}{
  if _,err:=New(s);err==nil{t.Errorf("accepted %q",s)}
 }
}
func TestMaintenance(t *testing.T) {
 server:=httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){_,_=w.Write([]byte(`{"status":"maintenance","message":"Actualizando mods"}`))}));defer server.Close()
 c,err:=New(server.URL);if err!=nil{t.Fatal(err)};c.http=server.Client()
 ok,msg,err:=c.CanPlay(context.Background(),"");if err!=nil||ok||msg!="Actualizando mods"{t.Fatalf("%v %q %v",ok,msg,err)}
}
func TestRedeem(t *testing.T) {
 server:=httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Header.Get("Authorization")!="Bearer verified-token"{t.Error("missing auth")}
  if r.URL.Path!="/api/v1/beta/redeem" || r.Method!="POST"{t.Error("unexpected request")}
  _,_=w.Write([]byte(`{"access":true}`))
 }));defer server.Close()
 c,_:=New(server.URL);c.http=server.Client()
 ok,err:=c.Redeem(context.Background(),"verified-token","MINERUX-12345678901234567890123456789012")
 if err!=nil||!ok{t.Fatalf("%v %v",ok,err)}
}
func TestFailClosed(t *testing.T) {
 server:=httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){http.Error(w,"offline",503)}));defer server.Close()
 c,_:=New(server.URL);c.http=server.Client()
 ok,_,err:=c.CanPlay(context.Background(),"token");if ok||err==nil{t.Fatalf("must fail closed: %v %v",ok,err)}
}

func TestBackendRedirectBlocked(t *testing.T) {
 redirected := false
 target := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
  redirected = true
  t.Error("redirect destination must never receive player bearer token")
 }))
 defer target.Close()
 origin := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
  http.Redirect(w,r,target.URL+"/steal",http.StatusTemporaryRedirect)
 }))
 defer origin.Close()
 c,err := New(origin.URL)
 if err!=nil { t.Fatal(err) }
 c.http.Transport = origin.Client().Transport
 _,err = c.CheckAccess(context.Background(),"verified-player-token")
 if err==nil || redirected { t.Fatalf("redirect must be rejected: err=%v redirected=%v",err,redirected) }
}
