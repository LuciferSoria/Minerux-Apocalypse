package updater

import (
 "context"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "net/http"
 "os"
 "path/filepath"
 "strings"
 "time"
)

type Asset struct {
 Path string `json:"path"`
 URL string `json:"url"`
 SHA256 string `json:"sha256"`
}
type Manifest struct {
 Schema int `json:"schema"`
 Version string `json:"version"`
 Assets []Asset `json:"assets"`
}
func Parse(data []byte) (Manifest, error) {
 var m Manifest
 if err:=json.Unmarshal(data,&m);err!=nil{return m,err}
 if m.Schema!=1{return m,errors.New("unsupported manifest schema")}
 if len(m.Version)==0{return m,errors.New("missing version")}
 for _,a:=range m.Assets {
  if err:=validate(a);err!=nil{return m,err}
 }
 return m,nil
}
func validate(a Asset) error {
 if a.Path=="" || filepath.IsAbs(a.Path) || strings.Contains(a.Path,"\\") || strings.Contains(a.Path,":") {return fmt.Errorf("unsafe asset path %q",a.Path)}
 for _,p:=range strings.Split(a.Path,"/"){if p==""||p=="."||p==".."{return fmt.Errorf("unsafe asset path %q",a.Path)}}
 if !strings.HasPrefix(a.URL,"https://"){return fmt.Errorf("asset %q requires HTTPS",a.Path)}
 if len(a.SHA256)!=64{return fmt.Errorf("invalid sha256 for %q",a.Path)}
 if _,err:=hex.DecodeString(a.SHA256);err!=nil{return err}
 return nil
}
func Verify(path,expected string) (bool,error) {
 f,err:=os.Open(path);if os.IsNotExist(err){return false,nil};if err!=nil{return false,err};defer f.Close()
 h:=sha256.New();if _,err=io.Copy(h,f);err!=nil{return false,err}
 return strings.EqualFold(hex.EncodeToString(h.Sum(nil)),expected),nil
}
// Install downloads into a temporary file, verifies SHA-256 and atomically replaces the target.
// It never executes downloaded files. The caller must use a trusted manifest.
func Install(ctx context.Context,root string,a Asset) error {
 if err:=validate(a);err!=nil{return err}
 dst:=filepath.Join(root,filepath.FromSlash(a.Path))
 // Reject symlink components to prevent escaping root.
 base:=filepath.Clean(root)
 cur:=base
 for _,part:=range strings.Split(a.Path,"/"){
  cur=filepath.Join(cur,part)
  info,err:=os.Lstat(cur)
  if err==nil && info.Mode()&os.ModeSymlink!=0{return fmt.Errorf("symlink in destination: %s",cur)}
  if err!=nil && !os.IsNotExist(err){return err}
 }
 if ok,err:=Verify(dst,a.SHA256);err!=nil{return err}else if ok{return nil}
 if err:=os.MkdirAll(filepath.Dir(dst),0755);err!=nil{return err}
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,a.URL,nil);if err!=nil{return err}
 client:=&http.Client{Timeout:10*time.Minute}
 resp,err:=client.Do(req);if err!=nil{return err};defer resp.Body.Close()
 if resp.StatusCode!=http.StatusOK{return fmt.Errorf("download %s: HTTP %d",a.Path,resp.StatusCode)}
 tmp,err:=os.CreateTemp(filepath.Dir(dst),".minerux-*.part");if err!=nil{return err}
 defer os.Remove(tmp.Name())
 h:=sha256.New()
 _,err=io.Copy(io.MultiWriter(tmp,h),io.LimitReader(resp.Body,1<<30+1))
 if err!=nil{tmp.Close();return err}
 if err=tmp.Close();err!=nil{return err}
 if !strings.EqualFold(hex.EncodeToString(h.Sum(nil)),a.SHA256){return fmt.Errorf("sha256 mismatch: %s",a.Path)}
 return os.Rename(tmp.Name(),dst)
}
