package updater

import (
 "context"
 "strings"
 "testing"
)

func TestFetchManifestRejectsNonHTTPS(t *testing.T) {
 for _,u:=range []string{
  "http://example.com/manifest.json",
  "file:///tmp/manifest.json",
  "https://user:password@example.com/manifest.json",
  "not a url",
 }{
  if _,err:=FetchManifest(context.Background(),u);err==nil{t.Errorf("accepted URL: %q",u)}
 }
}
func TestValidateManifestJSON(t *testing.T) {
 if _,err:=ValidateManifestJSON([]byte("not-json"));err==nil{t.Fatal("accepted invalid JSON")}
 m,err:=ValidateManifestJSON([]byte(`{"schema":1,"version":"dev","assets":[]}`))
 if err!=nil{t.Fatal(err)}
 if m.Version!="dev"{t.Fatalf("unexpected version %q",m.Version)}
 if _,err:=ValidateManifestJSON([]byte(strings.Repeat(" ",2<<20+1)));err==nil{t.Fatal("accepted oversized manifest")}
}
