package client

import (
	"crypto/md5" // #nosec
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/internal/parallel"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
)

func (c *ascClient) UpdatePreviewsAndScreenshotsIfNeeded(ctx *context.Context, g parallel.Group, loc *asc.AppStoreVersionLocalization, config config.VersionLocalization) error {
	if loc.Relationships.AppPreviewSets != nil {
		var previewSets asc.AppPreviewSetsResponse

		_, err := c.client.FollowReference(ctx, loc.Relationships.AppPreviewSets.Links.Related, &previewSets)
		if err != nil {
			return err
		}

		if err := c.UpdatePreviewSets(ctx, g, previewSets.Data, loc.ID, config.PreviewSets); err != nil {
			return err
		}
	}

	if loc.Relationships.AppScreenshotSets != nil {
		var screenshotSets asc.AppScreenshotSetsResponse

		_, err := c.client.FollowReference(ctx, loc.Relationships.AppScreenshotSets.Links.Related, &screenshotSets)
		if err != nil {
			return err
		}

		if err := c.UpdateScreenshotSets(ctx, g, screenshotSets.Data, loc.ID, config.ScreenshotSets); err != nil {
			return err
		}
	}

	return nil
}

func (c *ascClient) UploadRoutingCoverage(ctx *context.Context, versionID string, config config.File) error {
	prepare := func(name string, checksum string) (shouldContinue bool, err error) {
		covResp, _, err := c.client.Apps.GetRoutingAppCoverageForAppStoreVersion(ctx, versionID, nil)
		if err != nil {
			log.Warn(err.Error())
		}

		if covResp == nil {
			return true, nil
		}

		if _, err := c.client.Apps.DeleteRoutingAppCoverage(ctx, covResp.Data.ID); err != nil {
			return false, err
		}

		return true, nil
	}

	create := func(name string, size int64) (id string, ops []asc.UploadOperation, err error) {
		resp, _, err := c.client.Apps.CreateRoutingAppCoverage(ctx, name, size, versionID)
		if err != nil {
			return "", nil, err
		}

		return resp.Data.ID, resp.Data.Attributes.UploadOperations, nil
	}

	commit := func(id string, checksum string) error {
		_, _, err := c.client.Apps.CommitRoutingAppCoverage(ctx, id, asc.Bool(true), &checksum)
		return err
	}

	return c.uploadFile(ctx, config.Path, prepare, create, commit)
}

//nolint:dupl // This is a false positive identified by dupl against UpdateScreenshotSets
func (c *ascClient) UpdatePreviewSets(ctx *context.Context, g parallel.Group, previewSets []asc.AppPreviewSet, appStoreVersionLocalizationID string, config config.PreviewSets) error {
	found := make(map[asc.PreviewType]bool)

	for i := range previewSets {
		previewSet := previewSets[i]
		previewType := *previewSet.Attributes.PreviewType
		found[previewType] = true
		previewsConfig := config.GetPreviews(previewType)

		if err := c.UploadPreviews(ctx, g, &previewSet, previewsConfig); err != nil {
			return err
		}
	}

	for previewType, previews := range config {
		t := *previewType.APIValue()
		if found[t] {
			continue
		}

		previewSetResp, _, err := c.client.Apps.CreateAppPreviewSet(ctx, t, appStoreVersionLocalizationID)
		if err != nil {
			return err
		}

		if err := c.UploadPreviews(ctx, g, &previewSetResp.Data, previews); err != nil {
			return err
		}
	}

	return nil
}

func (c *ascClient) UploadPreviews(ctx *context.Context, g parallel.Group, previewSet *asc.AppPreviewSet, previewConfigs []config.Preview) error {
	previewsResp, _, err := c.client.Apps.ListAppPreviewsForSet(ctx, previewSet.ID, nil)
	if err != nil {
		return err
	}

	var previewsByName = make(map[string]*asc.AppPreview)

	for i := range previewsResp.Data {
		preview := previewsResp.Data[i]
		if preview.Attributes == nil || preview.Attributes.FileName == nil {
			continue
		}

		previewsByName[*preview.Attributes.FileName] = &preview
	}

	prepare := func(name string, checksum string) (shouldContinue bool, err error) {
		preview := previewsByName[name]
		if preview == nil {
			return true, nil
		}

		if preview.Attributes.SourceFileChecksum != nil &&
			*preview.Attributes.SourceFileChecksum == checksum {
			log.WithFields(log.Fields{
				"id":       preview.ID,
				"checksum": checksum,
			}).Debug("skip existing preview")

			return false, nil
		}

		log.WithFields(log.Fields{
			"name": name,
			"id":   preview.ID,
		}).Debug("delete preview")

		if _, err := c.client.Apps.DeleteAppPreview(ctx, preview.ID); err != nil {
			return false, err
		}

		return true, nil
	}

	create := func(name string, size int64) (id string, ops []asc.UploadOperation, err error) {
		log.WithFields(log.Fields{
			"name": name,
		}).Debug("create preview")

		resp, _, err := c.client.Apps.CreateAppPreview(ctx, name, size, previewSet.ID)
		if err != nil {
			return "", nil, err
		}

		return resp.Data.ID, resp.Data.Attributes.UploadOperations, nil
	}

	for i := range previewConfigs {
		previewConfig := previewConfigs[i]
		commit := func(id string, checksum string) error {
			log.WithFields(log.Fields{
				"id": id,
			}).Debug("commit preview")

			_, _, err := c.client.Apps.CommitAppPreview(ctx, id, asc.Bool(true), &checksum, &previewConfig.PreviewFrameTimeCode)

			return err
		}

		g.Go(func() error {
			return c.uploadFile(ctx, previewConfig.Path, prepare, create, commit)
		})
	}

	return nil
}

//nolint:dupl // This is a false positive identified by dupl against UpdatePreviewSets
func (c *ascClient) UpdateScreenshotSets(ctx *context.Context, g parallel.Group, screenshotSets []asc.AppScreenshotSet, appStoreVersionLocalizationID string, config config.ScreenshotSets) error {
	found := make(map[asc.ScreenshotDisplayType]bool)

	for i := range screenshotSets {
		screenshotSet := screenshotSets[i]
		screenshotType := *screenshotSet.Attributes.ScreenshotDisplayType
		found[screenshotType] = true
		screenshotConfig := config.GetScreenshots(screenshotType)

		if err := c.UploadScreenshots(ctx, g, &screenshotSet, screenshotConfig); err != nil {
			return err
		}
	}

	for screenshotType, screenshots := range config {
		t := *screenshotType.APIValue()
		if found[t] {
			continue
		}

		screenshotSetResp, _, err := c.client.Apps.CreateAppScreenshotSet(ctx, t, appStoreVersionLocalizationID)
		if err != nil {
			return err
		}

		if err := c.UploadScreenshots(ctx, g, &screenshotSetResp.Data, screenshots); err != nil {
			return err
		}
	}

	return nil
}

func (c *ascClient) UploadScreenshots(ctx *context.Context, g parallel.Group, screenshotSet *asc.AppScreenshotSet, config []config.File) error {
	shotsResp, _, err := c.client.Apps.ListAppScreenshotsForSet(ctx, screenshotSet.ID, nil)
	if err != nil {
		return err
	}

	var screenshotsByName = make(map[string]*asc.AppScreenshot)

	for i := range shotsResp.Data {
		shot := shotsResp.Data[i]
		if shot.Attributes == nil || shot.Attributes.FileName == nil {
			continue
		}

		screenshotsByName[*shot.Attributes.FileName] = &shot
	}

	prepare := func(name string, checksum string) (shouldContinue bool, err error) {
		shot := screenshotsByName[name]
		if shot == nil {
			return true, nil
		}

		if shot.Attributes.SourceFileChecksum != nil &&
			*shot.Attributes.SourceFileChecksum == checksum {
			log.WithFields(log.Fields{
				"id":       shot.ID,
				"checksum": checksum,
			}).Debug("skip existing screenshot")

			return false, nil
		}

		log.WithFields(log.Fields{
			"name": name,
			"id":   shot.ID,
		}).Debug("delete screenshot")

		if _, err := c.client.Apps.DeleteAppScreenshot(ctx, shot.ID); err != nil {
			return false, err
		}

		return true, nil
	}

	create := func(name string, size int64) (id string, ops []asc.UploadOperation, err error) {
		log.WithFields(log.Fields{
			"name": name,
		}).Debug("create screenshot")

		resp, _, err := c.client.Apps.CreateAppScreenshot(ctx, name, size, screenshotSet.ID)
		if err != nil {
			return "", nil, err
		}

		return resp.Data.ID, resp.Data.Attributes.UploadOperations, nil
	}

	commit := func(id string, checksum string) error {
		log.WithFields(log.Fields{
			"id": id,
		}).Debug("commit screenshot")

		_, _, err := c.client.Apps.CommitAppScreenshot(ctx, id, asc.Bool(true), &checksum)

		return err
	}

	for i := range config {
		screenshotConfig := config[i]

		g.Go(func() error {
			return c.uploadFile(ctx, screenshotConfig.Path, prepare, create, commit)
		})
	}

	return nil
}

func (c *ascClient) UploadReviewAttachments(ctx *context.Context, reviewDetailID string, config []config.File) error {
	if len(config) == 0 {
		return nil
	}

	var g = parallel.New(ctx.MaxProcesses)

	attachmentsResp, _, err := c.client.Submission.ListAttachmentsForReviewDetail(ctx, reviewDetailID, nil)
	if err != nil {
		return err
	}

	var attachmentsByName = make(map[string]*asc.AppStoreReviewAttachment)

	for i := range attachmentsResp.Data {
		attachment := attachmentsResp.Data[i]
		if attachment.Attributes == nil || attachment.Attributes.FileName == nil {
			continue
		}

		attachmentsByName[*attachment.Attributes.FileName] = &attachment
	}

	prepare := func(name string, checksum string) (shouldContinue bool, err error) {
		attachment := attachmentsByName[name]
		if attachment == nil {
			return true, nil
		}

		if attachment.Attributes.SourceFileChecksum != nil &&
			*attachment.Attributes.SourceFileChecksum == checksum {
			log.WithFields(log.Fields{
				"id":       attachment.ID,
				"checksum": checksum,
			}).Debug("skip existing attachment")

			return false, nil
		}

		log.WithFields(log.Fields{
			"name": name,
			"id":   attachment.ID,
		}).Debug("delete attachment")

		if _, err := c.client.Submission.DeleteAttachment(ctx, attachment.ID); err != nil {
			return false, err
		}

		return true, nil
	}

	create := func(name string, size int64) (id string, ops []asc.UploadOperation, err error) {
		log.WithFields(log.Fields{
			"name": name,
		}).Debug("create attachment")

		resp, _, err := c.client.Submission.CreateAttachment(ctx, name, size, reviewDetailID)
		if err != nil {
			return "", nil, err
		}

		return resp.Data.ID, resp.Data.Attributes.UploadOperations, nil
	}

	commit := func(id string, checksum string) error {
		log.WithFields(log.Fields{
			"id": id,
		}).Debug("commit attachment")

		_, _, err := c.client.Submission.CommitAttachment(ctx, id, asc.Bool(true), &checksum)

		return err
	}

	for i := range config {
		attachmentConfig := config[i]

		g.Go(func() error {
			return c.uploadFile(ctx, attachmentConfig.Path, prepare, create, commit)
		})
	}

	return g.Wait()
}

type prepareFunc func(name string, checksum string) (shouldContinue bool, err error)
type createFunc func(name string, size int64) (id string, ops []asc.UploadOperation, err error)
type commitFunc func(id string, checksum string) error

func (c *ascClient) uploadFile(ctx *context.Context, path string, prepare prepareFunc, create createFunc, commit commitFunc) (err error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}

	defer func() {
		closeErr := f.Close()
		if closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				log.Fatal(closeErr.Error())
			}
		}
	}()

	fstat, err := os.Stat(path)
	if err != nil {
		return err
	}

	checksum, err := md5Checksum(f)
	if err != nil {
		return err
	}

	shouldContinue, err := prepare(fstat.Name(), checksum)
	if err != nil {
		return err
	} else if !shouldContinue {
		return nil
	}

	id, ops, err := create(fstat.Name(), fstat.Size())
	if err != nil {
		return err
	}

	if err = c.client.Upload(ctx, ops, f); err != nil {
		return err
	}

	return commit(id, checksum)
}

func md5Checksum(f io.Reader) (string, error) {
	/* #nosec */
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
