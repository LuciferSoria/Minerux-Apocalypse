package installer

import (
 "archive/zip"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func makeModZIP(t *testing.T,files map[string]string) string {
 t.Helper()
 path:=filepath.Join(t.TempDir(),"mods.zip")
 out,err:=os.Create(path)
 if err!=nil{t.Fatal(err)}
 writer:=zip.NewWriter(out)
 for name,data:=range files {
  entry,err:=writer.Create(name)
  if err!=nil{t.Fatal(err)}
  if _,err:=entry.Write([]byte(data));err!=nil{t.Fatal(err)}
 }
 if err:=writer.Close();err!=nil{t.Fatal(err)}
 if err:=out.Close();err!=nil{t.Fatal(err)}
 return path
}
func TestImportModsZIP(t *testing.T) {
 archive:=makeModZIP(t,map[string]string{"mods/one.jar":"one","mods/two.jar":"two"})
 root:=t.TempDir()
 dir,count,err:=ImportModsZIP(archive,root,"26.1.2")
 if err!=nil || count!=2 {t.Fatalf("import: %v, %d",err,count)}
 if got,err:=os.ReadFile(filepath.Join(dir,"one.jar"));err!=nil || string(got)!="one"{t.Fatalf("file: %q %v",got,err)}
 if _,_,err:=ImportModsZIP(archive,root,"26.1.2");err!=nil{t.Fatalf("reimport: %v",err)}
}
func TestRejectUnsafeModZIP(t *testing.T) {
 for _,name:=range []string{"../evil.jar","mods/../evil.jar","mods/subdir/a.jar","mods/CON.jar","mods/a.txt","other/a.jar"} {
  archive:=makeModZIP(t,map[string]string{name:"bad"})
  if _,_,err:=ImportModsZIP(archive,t.TempDir(),"26.1.2");err==nil{t.Errorf("accepted %s",name)}
 }
}
func TestRejectChangedExistingMod(t *testing.T) {
 archive:=makeModZIP(t,map[string]string{"mods/a.jar":"good"})
 root:=t.TempDir()
 dir,_,err:=ImportModsZIP(archive,root,"26.1.2")
 if err!=nil{t.Fatal(err)}
 if err:=os.WriteFile(filepath.Join(dir,"a.jar"),[]byte("changed"),0600);err!=nil{t.Fatal(err)}
 if _,_,err:=ImportModsZIP(archive,root,"26.1.2");err==nil || !strings.Contains(err.Error(),"mismatch"){t.Fatalf("expected mismatch: %v",err)}
}
