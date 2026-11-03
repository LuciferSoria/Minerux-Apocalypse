package updater

import (
 "context"
 "fmt"
 "os"
 "path/filepath"
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
 for i, asset := range manifest.Assets {
  if err := ctx.Err(); err != nil { return err }
  if err := Install(ctx, abs, asset); err != nil {
   return fmt.Errorf("asset %d/%d (%s): %w", i+1, len(manifest.Assets), asset.Path, err)
  }
  if progress != nil { progress(i+1, len(manifest.Assets), asset.Path) }
 }
 return nil
}
