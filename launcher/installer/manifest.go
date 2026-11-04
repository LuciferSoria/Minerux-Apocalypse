package installer

import (
 "encoding/json"
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
 if !safeFileName.MatchString(m.Version) || len(m.Files)==0 || len(m.Files)>256 {
  return m,errors.New("manifest version or file count is invalid")
 }
 names:=make(map[string]bool,len(m.Files))
 for _,file:=range m.Files {
  if !safeFileName.MatchString(file.Name) || file.Name=="." || file.Name==".." || strings.HasSuffix(file.Name,".") {
   return m,errors.New("unsafe artifact filename")
  }
  normalized:=strings.ToLower(file.Name)
  if names[normalized] {return m,errors.New("duplicate artifact filename")}
  names[normalized]=true
  if len(file.SHA256)!=64 {return m,errors.New("artifact SHA-256 is required")}
  parsed,err:=url.Parse(file.URL)
  if err!=nil || parsed.Scheme!="https" || parsed.Host=="" || parsed.User!=nil || parsed.Fragment!="" {return m,errors.New("artifact URL must be HTTPS")}
 }
 return m,nil
}

// InstallManifest downloads files into a dedicated, versioned directory.
// On failure, previously installed versions are left untouched. The new
// version is not marked active; the launcher must activate it separately.
func InstallManifest(m Manifest,root string,download func(Artifact,string)error) (string,error) {
 if !safeFileName.MatchString(m.Version) || len(m.Files)==0 || len(m.Files)>256 {return "",errors.New("invalid manifest")}
 seen:=map[string]bool{}
 for _,f:=range m.Files {
  if !safeFileName.MatchString(f.Name) || f.Name=="." || f.Name==".." || strings.HasSuffix(f.Name,".") {return "",errors.New("unsafe artifact filename")}
  name:=strings.ToLower(f.Name)
  if seen[name] {return "",errors.New("duplicate artifact filename")}
  seen[name]=true
 }
 versionDir:=filepath.Join(root,"versions",m.Version)
 if err:=os.MkdirAll(versionDir,0700);err!=nil{return "",err}
 for _,f:=range m.Files {
  if err:=download(Artifact{URL:f.URL,SHA256:f.SHA256},filepath.Join(versionDir,f.Name));err!=nil{
   return "",fmt.Errorf("install %s: %w",f.Name,err)
  }
 }
 return versionDir,nil
}
