package installer

import (
 "encoding/json"
 "encoding/hex"
 "errors"
 "fmt"
 "io"
 "os"
 "path/filepath"
 "regexp"
 "net/url"
 "strings"
)

const maxManifestBytes = 1 << 20
var safeFileName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)

func validWindowsName(name string) bool {
 if !safeFileName.MatchString(name) || strings.HasSuffix(name,".") {return false}
 stem:=strings.ToUpper(strings.SplitN(name,".",2)[0])
 switch stem {
 case "CON","PRN","AUX","NUL","COM1","COM2","COM3","COM4","COM5","COM6","COM7","COM8","COM9","LPT1","LPT2","LPT3","LPT4","LPT5","LPT6","LPT7","LPT8","LPT9":
  return false
 }
 return true
}

func validateManifest(m Manifest) error {
 if !validWindowsName(m.Version) || len(m.Files)==0 || len(m.Files)>256 {
  return errors.New("manifest version or file count is invalid")
 }
 names:=make(map[string]bool,len(m.Files))
 for _,file:=range m.Files {
  if !validWindowsName(file.Name) {return errors.New("unsafe artifact filename")}
  normalized:=strings.ToLower(file.Name)
  if names[normalized] {return errors.New("duplicate artifact filename")}
  names[normalized]=true
  hash,err:=hex.DecodeString(file.SHA256)
  if err!=nil || len(hash)!=32 || file.SHA256!=strings.ToLower(file.SHA256) {
   return errors.New("invalid artifact SHA-256")
  }
  parsed,err:=url.Parse(file.URL)
  if err!=nil || parsed.Scheme!="https" || parsed.Host=="" || parsed.User!=nil || parsed.Fragment!="" || parsed.Opaque!="" {
   return errors.New("artifact URL must be HTTPS")
  }
 }
 return nil
}

// Manifest is a locally supplied, trusted release manifest. Do not download
// and trust a manifest from the same untrusted origin as its artifacts.
type Manifest struct {
 Version string `json:"version"`
 Files []ManifestFile `json:"files"`
}
type ManifestFile struct {
 Name string `json:"name"`
 URL string `json:"url"`
 SHA256 string `json:"sha256"`
}

func ParseManifest(r io.Reader) (Manifest,error) {
 var m Manifest
 data,err:=io.ReadAll(io.LimitReader(r,maxManifestBytes+1))
 if err!=nil{return m,err}
 if len(data)>maxManifestBytes{return m,errors.New("manifest exceeds size limit")}
 decoder:=json.NewDecoder(strings.NewReader(string(data)))
 decoder.DisallowUnknownFields()
 if err:=decoder.Decode(&m);err!=nil{return m,fmt.Errorf("invalid manifest: %w",err)}
 var trailing any
 if err:=decoder.Decode(&trailing);err!=io.EOF{return m,errors.New("manifest contains trailing data")}
 if err:=validateManifest(m);err!=nil{return m,err}
 return m,nil
}

// InstallManifest downloads files into a dedicated, versioned directory.
// On failure, previously installed versions are left untouched. The new
// version is not marked active; the launcher must activate it separately.
func InstallManifest(m Manifest,root string,download func(Artifact,string)error) (string,error) {
 if err:=validateManifest(m);err!=nil{return "",err}
 versionDir:=filepath.Join(root,"versions",m.Version)
 if err:=os.MkdirAll(versionDir,0700);err!=nil{return "",err}
 for _,f:=range m.Files {
  if err:=download(Artifact{URL:f.URL,SHA256:f.SHA256},filepath.Join(versionDir,f.Name));err!=nil{
   return "",fmt.Errorf("install %s: %w",f.Name,err)
  }
 }
 return versionDir,nil
}
