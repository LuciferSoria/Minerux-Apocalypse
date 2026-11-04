package installer

import (
 "context"
 "crypto/sha256"
 "encoding/hex"
 "errors"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "os"
 "path/filepath"
 "strings"
 "time"
)

// Artifact is a downloadable, independently SHA-256-verified game dependency.
// Never trust a filename or a hash supplied by an untrusted remote manifest.
type Artifact struct {
 URL string
 SHA256 string
}

func Download(ctx context.Context, artifact Artifact, destination string) error {
 return downloadWithClient(ctx, artifact, destination, &http.Client{Timeout: 10*time.Minute})
}

func downloadWithClient(ctx context.Context, artifact Artifact, destination string, client *http.Client) error {
 u, err := url.Parse(artifact.URL)
 if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.Fragment != "" {
  return errors.New("artifact URL must be HTTPS without credentials or fragment")
 }
 expected, err := hex.DecodeString(artifact.SHA256)
 if err != nil || len(expected) != sha256.Size || strings.ToLower(artifact.SHA256) != artifact.SHA256 {
  return errors.New("invalid SHA-256 digest")
 }
 if destination == "" { return errors.New("destination is required") }
 // Re-running an installation must be safe: accept an existing regular file
 // only when its bytes already match the trusted digest.
 if info,statErr:=os.Lstat(destination);statErr==nil {
  if !info.Mode().IsRegular() {return errors.New("destination is not a regular file")}
  return verifyExisting(destination,artifact.SHA256)
 } else if !os.IsNotExist(statErr) {return statErr}
 parent := filepath.Dir(destination)
 if err := os.MkdirAll(parent,0700); err != nil { return err }
 tmp, err := os.CreateTemp(parent,".minerux-download-*")
 if err != nil { return err }
 defer os.Remove(tmp.Name())
 defer tmp.Close()
 // Use a copy so redirect policy cannot be weakened by a caller.
 localClient := *client
 localClient.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
  return errors.New("artifact redirects are not allowed")
 }
 req,err := http.NewRequestWithContext(ctx,http.MethodGet,u.String(),nil)
 if err != nil { return err }
 resp,err := localClient.Do(req)
 if err != nil { return err }
 defer resp.Body.Close()
 if resp.StatusCode != http.StatusOK { return fmt.Errorf("artifact download HTTP %d",resp.StatusCode) }
 h := sha256.New()
 // Hard limit: 2 GiB per artifact; prevents unbounded disk consumption.
 const maxBytes int64 = 2 << 30
 n,err := io.Copy(io.MultiWriter(tmp,h),io.LimitReader(resp.Body,maxBytes+1))
 if err != nil { return err }
 if n > maxBytes { return errors.New("artifact exceeds 2 GiB limit") }
 if !strings.EqualFold(hex.EncodeToString(h.Sum(nil)),artifact.SHA256) {
  return errors.New("artifact SHA-256 mismatch")
 }
 if err := tmp.Sync(); err != nil { return err }
 if err := tmp.Close(); err != nil { return err }
 // Refuse to replace an existing file: prevents silent replacement of
 // previously installed game files until a versioned update is implemented.
 if _,err := os.Lstat(destination); err == nil { return verifyExisting(destination,artifact.SHA256) } else if !os.IsNotExist(err) { return err }
 return os.Rename(tmp.Name(),destination)
}

// verifyExisting accepts an already-installed file only if its contents
// match the exact expected digest. It never follows a final-component symlink.
func verifyExisting(destination,expected string) error {
 info,err:=os.Lstat(destination)
 if err!=nil{return err}
 if !info.Mode().IsRegular(){return errors.New("destination is not a regular file")}
 file,err:=os.Open(destination)
 if err!=nil{return err}
 defer file.Close()
 hash:=sha256.New()
 if _,err:=io.Copy(hash,file);err!=nil{return err}
 if hex.EncodeToString(hash.Sum(nil))!=expected {return errors.New("existing artifact SHA-256 mismatch")}
 return nil
}
