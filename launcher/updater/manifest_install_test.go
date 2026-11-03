package updater

import (
 "context"
 "crypto/sha256"
 "encoding/hex"
 "errors"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestInstallManifestEmpty(t *testing.T) {
 var calls int
 m := Manifest{Schema:1, Version:"test"}
 err := InstallManifest(context.Background(),t.TempDir(),m,func(_, _ int, _ string){calls++})
 if err!=nil{t.Fatal(err)}
 if calls!=0{t.Fatalf("unexpected progress callbacks: %d",calls)}
}
func TestInstallManifestRejectsBadHeader(t *testing.T) {
 if err:=InstallManifest(context.Background(),t.TempDir(),Manifest{},nil);err==nil{t.Fatal("invalid header accepted")}
}
func TestInstallManifestCancelled(t *testing.T) {
 ctx,cancel:=context.WithCancel(context.Background());cancel()
 m:=Manifest{Schema:1,Version:"test",Assets:[]Asset{{Path:"mods/test.jar",URL:"https://example.org/test.jar",SHA256:strings.Repeat("a",64)}}}
 if err:=InstallManifest(ctx,t.TempDir(),m,nil);!errors.Is(err,context.Canceled){t.Fatalf("expected cancellation, got %v",err)}
}
func TestInstallManifestAlreadyVerified(t *testing.T) {
 dir:=t.TempDir()
 if err:=os.MkdirAll(filepath.Join(dir,"mods"),0755);err!=nil{t.Fatal(err)}
 content:=[]byte("licensed example asset")
 path:=filepath.Join(dir,"mods","example.jar")
 if err:=os.WriteFile(path,content,0644);err!=nil{t.Fatal(err)}
 hash:=sha256.Sum256(content)
 m:=Manifest{Schema:1,Version:"test",Assets:[]Asset{{Path:"mods/example.jar",URL:"https://example.org/never-download",SHA256:hex.EncodeToString(hash[:])}}}
 calls:=0
 if err:=InstallManifest(context.Background(),dir,m,func(done,total int,p string){calls++;if done!=1||total!=1||p!="mods/example.jar"{t.Fatalf("unexpected progress")}});err!=nil{t.Fatal(err)}
 if calls!=1{t.Fatalf("expected one callback, got %d",calls)}
}
