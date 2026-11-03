package updater

import (
 "crypto/sha256"
 "encoding/hex"
 "os"
 "path/filepath"
 "strings"
 "testing"
)
func TestParseValid(t *testing.T) {
 b:=[]byte(`{"schema":1,"version":"dev","assets":[{"path":"mods/example.jar","url":"https://example.org/example.jar","sha256":"`+strings.Repeat("a",64)+`"}]}`)
 m,err:=Parse(b);if err!=nil{t.Fatal(err)};if len(m.Assets)!=1{t.Fatal("expected one asset")}
}
func TestRejectUnsafe(t *testing.T) {
 for _,path:=range []string{"../escape","/absolute","mods/../escape","mods\\bad.jar","C:/bad","mods//bad"}{
  a:=Asset{Path:path,URL:"https://example.org/file",SHA256:strings.Repeat("a",64)}
  if err:=validate(a);err==nil{t.Errorf("accepted unsafe path %q",path)}
 }
}
func TestRejectHTTP(t *testing.T) {
 a:=Asset{Path:"mods/a.jar",URL:"http://example.org/a.jar",SHA256:strings.Repeat("a",64)}
 if validate(a)==nil{t.Fatal("accepted HTTP")}
}
func TestVerify(t *testing.T) {
 path:=filepath.Join(t.TempDir(),"test.jar")
 data:=[]byte("hello")
 if err:=os.WriteFile(path,data,0600);err!=nil{t.Fatal(err)}
 h:=sha256.Sum256(data)
 ok,err:=Verify(path,hex.EncodeToString(h[:]));if err!=nil||!ok{t.Fatalf("verify: %v %v",ok,err)}
 ok,err=Verify(path,strings.Repeat("0",64));if err!=nil||ok{t.Fatal("wrong digest accepted")}
}
