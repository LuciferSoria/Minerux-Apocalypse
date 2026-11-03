package updater

import (
 "context"
 "encoding/json"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "time"
)

const maxManifestBytes = 2 << 20

// FetchManifest retrieves a manifest from an HTTPS URL, with a strict size limit.
// IMPORTANT: HTTPS is transport security, not proof that the manifest is authorized.
// The caller must configure a trusted URL and ideally pin a signature or digest.
func FetchManifest(ctx context.Context, rawURL string) (Manifest, error) {
 var empty Manifest
 u,err:=url.Parse(rawURL)
 if err!=nil || u.Scheme!="https" || u.Host=="" || u.User!=nil {
  return empty,fmt.Errorf("manifest URL must be HTTPS with no credentials")
 }
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,rawURL,nil)
 if err!=nil{return empty,err}
 client:=&http.Client{
  Timeout:30*time.Second,
  CheckRedirect:func(req *http.Request,via []*http.Request) error {
   if req.URL.Scheme!="https" {return fmt.Errorf("insecure redirect rejected")}
   if len(via)>=5{return fmt.Errorf("too many redirects")}
   return nil
  },
 }
 resp,err:=client.Do(req)
 if err!=nil{return empty,err}
 defer resp.Body.Close()
 if resp.StatusCode!=http.StatusOK{return empty,fmt.Errorf("manifest HTTP %d",resp.StatusCode)}
 data,err:=io.ReadAll(io.LimitReader(resp.Body,maxManifestBytes+1))
 if err!=nil{return empty,err}
 if len(data)>maxManifestBytes{return empty,fmt.Errorf("manifest too large")}
 // Parse also validates asset paths and checksums.
 m,err:=Parse(data)
 if err!=nil{return empty,fmt.Errorf("invalid manifest: %w",err)}
 return m,nil
}

// ValidateManifestJSON checks a local manifest before any network downloads.
func ValidateManifestJSON(data []byte) (Manifest,error) {
 var m Manifest
 if !json.Valid(data){return m,fmt.Errorf("invalid JSON")}
 if len(data)>maxManifestBytes{return m,fmt.Errorf("manifest too large")}
 return Parse(data)
}
