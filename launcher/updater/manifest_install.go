package updater

import (
 "context"
 "fmt"
 "os"
 "path/filepath"
 "strings"
)

// Progress reports the completed asset count after each verified installation.
type Progress func(completed, total int, path string)

// InstallManifest installs assets sequentially and stops on the first error.
// Only pass manifests from a trusted source. This function never launches the game.
func InstallManifest(ctx context.Context, root string, manifest Manifest, progress Progress) error {
 if manifest.Schema != 1 || manifest.Version == "" { return fmt.Errorf("invalid manifest header") }
 if root == "" { return fmt.Errorf("empty installation directory") }
 abs, err := filepath.Abs(root)
 if err != nil { return err }
 if err := os.MkdirAll(abs, 0755); err != nil { return err }
 seen:=make(map[string]bool,len(manifest.Assets))
 for i, asset := range manifest.Assets {
  if err := ctx.Err(); err != nil { return err }
  if err := validate(asset); err != nil { return fmt.Errorf("invalid asset %d: %w",i+1,err) }
  key:=strings.ToLower(asset.Path)
  if seen[key] { return fmt.Errorf("duplicate asset path %q",asset.Path) }
  seen[key]=true
  if err := Install(ctx, abs, asset); err != nil {
   return fmt.Errorf("asset %d/%d (%s): %w", i+1, len(manifest.Assets), asset.Path, err)
  }
  if progress != nil { progress(i+1, len(manifest.Assets), asset.Path) }
 }
 return nil
}
