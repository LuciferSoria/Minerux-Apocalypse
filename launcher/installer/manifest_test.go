package installer

import (
 "errors"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

const goodManifest = `{"version":"26.1.2-beta1","files":[{"name":"example.jar","url":"https://downloads.example.org/example.jar","sha256":"` + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + `"}]}`

func TestParseManifest(t *testing.T) {
 m,err:=ParseManifest(strings.NewReader(goodManifest))
 if err!=nil || m.Version!="26.1.2-beta1" || len(m.Files)!=1 {t.Fatalf("valid manifest: %+v %v",m,err)}
}
func TestRejectManifestTraversal(t *testing.T) {
 for _,name:=range []string{"../escape.jar","subdir/mod.jar","..","mod.jar."} {
  input:=strings.Replace(goodManifest,"example.jar",name,1)
  if _,err:=ParseManifest(strings.NewReader(input));err==nil{t.Errorf("accepted %q",name)}
 }
}
func TestRejectDuplicateWindowsName(t *testing.T) {
 m,err:=ParseManifest(strings.NewReader(goodManifest))
 if err!=nil{t.Fatal(err)}
 second:=m.Files[0]
 second.Name="EXAMPLE.JAR"
 m.Files=append(m.Files,second)
 called:=false
 _,err=InstallManifest(m,t.TempDir(),func(Artifact,string)error{called=true;return nil})
 if err==nil || called {t.Fatalf("duplicate accepted or download started: %v",err)}
}
func TestRejectOversizedManifest(t *testing.T) {
 if _,err:=ParseManifest(strings.NewReader(strings.Repeat(" ",maxManifestBytes+1)));err==nil{t.Fatal("oversized manifest accepted")}
}
func TestInstallStagesVersion(t *testing.T) {
 m,err:=ParseManifest(strings.NewReader(goodManifest))
 if err!=nil{t.Fatal(err)}
 root:=t.TempDir()
 folder,err:=InstallManifest(m,root,func(a Artifact,dest string)error{
  if a.URL!=m.Files[0].URL{return errors.New("wrong URL")}
  return os.WriteFile(dest,[]byte("staged"),0600)
 })
 if err!=nil{t.Fatal(err)}
 if folder!=filepath.Join(root,"versions","26.1.2-beta1"){t.Fatalf("wrong path %s",folder)}
 if _,err:=os.Stat(filepath.Join(folder,"example.jar"));err!=nil{t.Fatal(err)}
}

func TestRejectReservedWindowsNames(t *testing.T) {
 for _,name:=range []string{"CON","con.jar","NUL.txt","COM1.jar","lpt9.log","AUX","PRN"} {
  withName:=strings.Replace(goodManifest,"example.jar",name,1)
  if _,err:=ParseManifest(strings.NewReader(withName));err==nil{t.Errorf("accepted reserved filename %q",name)}
 }
}
func TestInstallRejectsInvalidManifestWithoutDownloading(t *testing.T) {
 m,err:=ParseManifest(strings.NewReader(goodManifest))
 if err!=nil{t.Fatal(err)}
 m.Files[0].SHA256="invalid"
 called:=false
 _,err=InstallManifest(m,t.TempDir(),func(Artifact,string)error{called=true;return nil})
 if err==nil || called {t.Fatalf("invalid manifest triggered download: %v",err)}
}
func TestRejectNonHexDigest(t *testing.T) {
 bad:=strings.Replace(goodManifest,strings.Repeat("a",64),strings.Repeat("z",64),1)
 if _,err:=ParseManifest(strings.NewReader(bad));err==nil{t.Fatal("accepted non-hex digest")}
}
