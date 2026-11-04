package installer

import (
 "archive/zip"
 "crypto/sha256"
 "encoding/hex"
 "errors"
 "fmt"
 "io"
 "os"
 "path/filepath"
 "strings"
)

// ImportModsZIP imports a local ZIP containing a mods/ folder. It never
// extracts arbitrary ZIP paths and refuses to overwrite installed mods.
// This stages mods only; it does not install NeoForge or Minecraft.
func ImportModsZIP(zipPath,root,version string) (string,int,error) {
 if !validWindowsName(version) {return "",0,errors.New("invalid version name")}
 archive,err:=zip.OpenReader(zipPath)
 if err!=nil{return "",0,err}
 defer archive.Close()
 if len(archive.File)>256 {return "",0,errors.New("too many ZIP entries")}
 destDir:=filepath.Join(root,"versions",version,"mods")
 type entry struct { file *zip.File; name string }
 var entries []entry
 seen:=map[string]bool{}
 var total uint64
 for _,f:=range archive.File {
  if f.FileInfo().IsDir() {if f.Name!="mods/" {return "",0,errors.New("unexpected ZIP directory")};continue}
  if f.Mode()&os.ModeSymlink!=0 {return "",0,errors.New("ZIP symlinks are forbidden")}
  if !strings.HasPrefix(f.Name,"mods/") {return "",0,errors.New("ZIP files must be under mods/")}
  name:=strings.TrimPrefix(f.Name,"mods/")
  if !validWindowsName(name) || !strings.EqualFold(filepath.Ext(name),".jar") {return "",0,fmt.Errorf("invalid mod filename: %q",name)}
  lower:=strings.ToLower(name)
  if seen[lower] {return "",0,errors.New("duplicate mod filename")}
  seen[lower]=true
  if f.UncompressedSize64>128<<20 {return "",0,errors.New("mod exceeds 128 MiB")}
  total+=f.UncompressedSize64
  if total>512<<20 {return "",0,errors.New("ZIP exceeds 512 MiB")}
  entries=append(entries,entry{f,name})
 }
 if len(entries)==0 {return "",0,errors.New("ZIP contains no mod JARs")}
 if err:=os.MkdirAll(destDir,0700);err!=nil{return "",0,err}
 for _,e:=range entries {
  destination:=filepath.Join(destDir,e.name)
  source,err:=e.file.Open()
  if err!=nil{return "",0,err}
  tmp,err:=os.CreateTemp(destDir,".minerux-mod-*")
  if err!=nil{source.Close();return "",0,err}
  hash:=sha256.New()
  n,copyErr:=io.Copy(io.MultiWriter(tmp,hash),io.LimitReader(source,128<<20+1))
  closeErr:=source.Close()
  syncErr:=tmp.Sync()
  tmpCloseErr:=tmp.Close()
  if copyErr!=nil || closeErr!=nil || syncErr!=nil || tmpCloseErr!=nil || n!=int64(e.file.UncompressedSize64) {
   os.Remove(tmp.Name())
   return "",0,fmt.Errorf("unable to verify ZIP entry %s",e.name)
  }
  expected:=hex.EncodeToString(hash.Sum(nil))
  if _,err:=os.Lstat(destination);err==nil {
   os.Remove(tmp.Name())
   if err:=verifyExisting(destination,expected);err!=nil{return "",0,fmt.Errorf("installed mod %s: %w",e.name,err)}
   continue
  } else if !os.IsNotExist(err) {os.Remove(tmp.Name());return "",0,err}
  if err:=os.Rename(tmp.Name(),destination);err!=nil{os.Remove(tmp.Name());return "",0,err}
 }
 return destDir,len(entries),nil
}
