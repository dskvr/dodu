package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/tyutyutyu/dodu/pkg/scan"
)

// CSVKind selects which table to render.
type CSVKind string

// CSV table identifiers.
const (
	CSVImages     CSVKind = "images"
	CSVContainers CSVKind = "containers"
	CSVVolumes    CSVKind = "volumes"
	CSVBuildCache CSVKind = "build_cache"
)

// AllCSVKinds lists the supported tables.
var AllCSVKinds = []CSVKind{CSVImages, CSVContainers, CSVVolumes, CSVBuildCache}

// ToCSV writes a single CSV table to w.
func ToCSV(w io.Writer, snap *scan.Snapshot, kind CSVKind) error {
	if snap == nil {
		return fmt.Errorf("export: nil snapshot")
	}
	cw := csv.NewWriter(w)
	defer cw.Flush()

	switch kind {
	case CSVImages:
		if err := cw.Write([]string{"id", "repo_tags", "size_bytes", "shared_size_bytes", "containers", "created"}); err != nil {
			return err
		}
		for _, im := range snap.Images {
			if err := cw.Write([]string{
				im.ID,
				strings.Join(im.RepoTags, ";"),
				strconv.FormatInt(im.Size, 10),
				strconv.FormatInt(im.SharedSize, 10),
				strconv.FormatInt(im.Containers, 10),
				im.Created.UTC().Format(time.RFC3339),
			}); err != nil {
				return err
			}
		}
	case CSVContainers:
		if err := cw.Write([]string{"id", "names", "image", "image_id", "state", "size_rw_bytes", "size_root_fs_bytes", "log_size_bytes"}); err != nil {
			return err
		}
		for _, c := range snap.Containers {
			if err := cw.Write([]string{
				c.ID,
				strings.Join(c.Names, ";"),
				c.Image,
				c.ImageID,
				c.State,
				strconv.FormatInt(c.SizeRw, 10),
				strconv.FormatInt(c.SizeRootFs, 10),
				strconv.FormatInt(snap.LogSize(c.ID), 10),
			}); err != nil {
				return err
			}
		}
	case CSVVolumes:
		if err := cw.Write([]string{"name", "driver", "mountpoint", "usage_bytes", "ref_count"}); err != nil {
			return err
		}
		for _, v := range snap.Volumes {
			if err := cw.Write([]string{
				v.Name,
				v.Driver,
				v.Mountpoint,
				strconv.FormatInt(v.UsageBytes, 10),
				strconv.FormatInt(v.RefCount, 10),
			}); err != nil {
				return err
			}
		}
	case CSVBuildCache:
		if err := cw.Write([]string{"id", "type", "description", "size_bytes", "in_use", "shared"}); err != nil {
			return err
		}
		for _, b := range snap.BuildCache {
			if err := cw.Write([]string{
				b.ID,
				b.Type,
				b.Description,
				strconv.FormatInt(b.Size, 10),
				strconv.FormatBool(b.InUse),
				strconv.FormatBool(b.Shared),
			}); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("export: unknown csv kind %q", kind)
	}
	return cw.Error()
}
