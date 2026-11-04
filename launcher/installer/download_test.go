package installer

import (
 "context"
 "crypto/sha256"
 "encoding/hex"
 "net/http"
 "net/http/httptest"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestRejectInsecureURL(t *testing.T) {
 for _,url := range []string{"http://example.org/a","file:///tmp/a","https://user:pw@example.org/a","https://example.org/a#fragment"} {
  err := Download(context.Background(),Artifact{URL:url,SHA256:strings.Repeat("0",64)},filepath.Join(t.TempDir(),"mod.jar"))
  if err==nil { t.Errorf("accepted unsafe URL: %s",url) }
 }
}
func TestRejectBadDigest(t *testing.T) {
 err:=Download(context.Background(),Artifact{URL:"https://example.org/mod.jar",SHA256:"bad"},filepath.Join(t.TempDir(),"mod.jar"))
 if err==nil {t.Fatal("accepted invalid digest")}
}
func TestRejectRedirectBeforeDownload(t *testing.T) {
 targetHit:=false
 target:=httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){targetHit=true}))
 defer target.Close()
 // Redirects are refused even for a trusted HTTPS origin. This test only
 // asserts the transport policy; the TLS test server is not publicly trusted.
 origin:=httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  http.Redirect(w,r,target.URL,http.StatusTemporaryRedirect)
 }))
 defer origin.Close()
 // A normal TLS client rejects the self-signed test cert first; either way,
 // the destination must never be contacted.
 err:=Download(context.Background(),Artifact{URL:origin.URL,SHA256:strings.Repeat("0",64)},filepath.Join(t.TempDir(),"mod.jar"))
 if err==nil || targetHit {t.Fatalf("unsafe redirect followed: %v",err)}
}
func TestExistingFileIsPreserved(t *testing.T) {
 path:=filepath.Join(t.TempDir(),"existing.jar")
 if err:=os.WriteFile(path,[]byte("original"),0600);err!=nil{t.Fatal(err)}
 // Invalid digest causes an early failure, and must never touch the file.
 err:=Download(context.Background(),Artifact{URL:"https://example.org/a",SHA256:"invalid"},path)
 if err==nil{t.Fatal("expected failure")}
 got,err:=os.ReadFile(path)
 if err!=nil||string(got)!="original"{t.Fatalf("existing file changed: %q %v",got,err)}
}
func TestDigestFormat(t *testing.T) {
 sum:=sha256.Sum256([]byte("minerux"))
 if len(hex.EncodeToString(sum[:]))!=64{t.Fatal("bad SHA-256 format")}
}

func TestVerifiedDownload(t *testing.T) {
 payload:=[]byte("verified minerux mod bytes")
 sum:=sha256.Sum256(payload)
 server:=httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
  _,_=w.Write(payload)
 }))
 defer server.Close()
 destination:=filepath.Join(t.TempDir(),"mods","verified.jar")
 err:=downloadWithClient(context.Background(),Artifact{URL:server.URL+"/verified.jar",SHA256:hex.EncodeToString(sum[:])},destination,server.Client())
 if err!=nil {t.Fatal(err)}
 got,err:=os.ReadFile(destination)
 if err!=nil || string(got)!=string(payload) {t.Fatalf("wrong installed bytes: %v",err)}
 // A second install must not replace the verified file.
 if err:=downloadWithClient(context.Background(),Artifact{URL:server.URL, SHA256:hex.EncodeToString(sum[:])},destination,server.Client());err==nil {
  t.Fatal("existing destination was overwritten")
 }
}

func TestHashMismatchLeavesNoInstalledFile(t *testing.T) {
 server:=httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
  _,_=w.Write([]byte("unexpected payload"))
 }))
 defer server.Close()
 destination:=filepath.Join(t.TempDir(),"mods","bad.jar")
 err:=downloadWithClient(context.Background(),Artifact{URL:server.URL,SHA256:strings.Repeat("0",64)},destination,server.Client())
 if err==nil || !strings.Contains(err.Error(),"mismatch") {t.Fatalf("expected hash mismatch: %v",err)}
 if _,err:=os.Stat(destination);!os.IsNotExist(err){t.Fatalf("invalid artifact was installed: %v",err)}
}

func TestRedirectNeverReachesTarget(t *testing.T) {
 targetHit:=false
 target:=httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){targetHit=true}))
 defer target.Close()
 origin:=httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
  http.Redirect(w,r,target.URL,http.StatusTemporaryRedirect)
 }))
 defer origin.Close()
 err:=downloadWithClient(context.Background(),Artifact{URL:origin.URL,SHA256:strings.Repeat("0",64)},filepath.Join(t.TempDir(),"mod.jar"),origin.Client())
 if err==nil || targetHit {t.Fatalf("redirect was followed: err=%v targetHit=%v",err,targetHit)}
}
