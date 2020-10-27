// Package template is a pipe that processes a configuration's template fields and stores it in the context
package template

import (
	"github.com/cidertool/cider/internal/template"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
	"github.com/hashicorp/go-multierror"
)

// Pipe is a global hook pipe.
type Pipe struct{}

// String is the name of this pipe.
func (Pipe) String() string {
	return "applying template values"
}

// Run executes the hooks.
func (p Pipe) Run(ctx *context.Context) error {
	var template = template.New(ctx)

	project, err := ctx.RawConfig.Copy()

	if err != nil {
		return err
	}

	var errors *multierror.Error

	for appName := range project {
		app := project[appName]
		if err := updateApp(&app, template); err != nil {
			errors = multierror.Append(errors, err)
		}

		project[appName] = app
	}

	ctx.Config = project

	return errors.ErrorOrNil()
}

func updateApp(app *config.App, tmpl *template.Template) error {
	var errors error

	for locName := range app.Localizations {
		loc := app.Localizations[locName]
		if err := updateAppLocalization(&loc, tmpl); err != nil {
			errors = multierror.Append(errors, err)
		}

		app.Localizations[locName] = loc
	}

	if err := updateAppTestflight(&app.Testflight, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := updateAppVersions(&app.Versions, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	return errors
}

func updateAppLocalization(loc *config.AppLocalization, tmpl *template.Template) error {
	var errors error
	if err := applyTemplateVar(&loc.Name, loc.Name, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := applyTemplateVar(&loc.PrivacyPolicyText, loc.PrivacyPolicyText, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := applyTemplateVar(&loc.PrivacyPolicyURL, loc.PrivacyPolicyURL, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := applyTemplateVar(&loc.Subtitle, loc.Subtitle, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	return errors
}

func updateAppTestflight(tf *config.Testflight, tmpl *template.Template) error {
	var errors error
	if err := applyTemplateVar(&tf.LicenseAgreement, tf.LicenseAgreement, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	for locName := range tf.Localizations {
		loc := tf.Localizations[locName]
		if err := updateTestflightLocalization(&loc, tmpl); err != nil {
			errors = multierror.Append(errors, err)
		}

		tf.Localizations[locName] = loc
	}

	if err := updateReviewDetails(tf.ReviewDetails, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	return errors
}

func updateAppVersions(version *config.Version, tmpl *template.Template) error {
	var errors error
	if err := applyTemplateVar(&version.Copyright, version.Copyright, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	for locName := range version.Localizations {
		loc := version.Localizations[locName]
		if err := updateVersionLocalization(&loc, tmpl); err != nil {
			errors = multierror.Append(errors, err)
		}

		version.Localizations[locName] = loc
	}

	if version.RoutingCoverage != nil {
		if err := applyTemplateVar(&version.RoutingCoverage.Path, version.RoutingCoverage.Path, tmpl); err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	if err := updateReviewDetails(version.ReviewDetails, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	return errors
}

func updateTestflightLocalization(loc *config.TestflightLocalization, tmpl *template.Template) error {
	var errors error
	if err := applyTemplateVar(&loc.Description, loc.Description, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := applyTemplateVar(&loc.FeedbackEmail, loc.FeedbackEmail, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := applyTemplateVar(&loc.MarketingURL, loc.MarketingURL, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := applyTemplateVar(&loc.PrivacyPolicyURL, loc.PrivacyPolicyURL, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := applyTemplateVar(&loc.TVOSPrivacyPolicy, loc.TVOSPrivacyPolicy, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := applyTemplateVar(&loc.WhatsNew, loc.WhatsNew, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	return errors
}

func updateVersionLocalization(loc *config.VersionLocalization, tmpl *template.Template) error {
	var errors error
	if err := applyTemplateVar(&loc.Description, loc.Description, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := applyTemplateVar(&loc.Keywords, loc.Keywords, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := applyTemplateVar(&loc.MarketingURL, loc.MarketingURL, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := applyTemplateVar(&loc.PromotionalText, loc.PromotionalText, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := applyTemplateVar(&loc.SupportURL, loc.SupportURL, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	if err := applyTemplateVar(&loc.WhatsNewText, loc.WhatsNewText, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	for previewType, set := range loc.PreviewSets {
		var previews = make([]config.Preview, len(set))

		for i, preview := range set {
			if err := applyTemplateVar(&preview.Path, preview.Path, tmpl); err != nil {
				errors = multierror.Append(errors, err)
			}

			previews[i] = preview
		}

		loc.PreviewSets[previewType] = previews
	}

	for screenshotType, set := range loc.ScreenshotSets {
		var screenshots = make([]config.File, len(set))

		for i, screenshot := range set {
			if err := applyTemplateVar(&screenshot.Path, screenshot.Path, tmpl); err != nil {
				errors = multierror.Append(errors, err)
			}

			screenshots[i] = screenshot
		}

		loc.ScreenshotSets[screenshotType] = screenshots
	}

	return errors
}

func updateReviewDetails(details *config.ReviewDetails, tmpl *template.Template) error {
	var errors error

	if details == nil {
		return errors
	}

	if details.Contact != nil {
		if err := applyTemplateVar(&details.Contact.Email, details.Contact.Email, tmpl); err != nil {
			errors = multierror.Append(errors, err)
		}

		if err := applyTemplateVar(&details.Contact.FirstName, details.Contact.FirstName, tmpl); err != nil {
			errors = multierror.Append(errors, err)
		}

		if err := applyTemplateVar(&details.Contact.LastName, details.Contact.LastName, tmpl); err != nil {
			errors = multierror.Append(errors, err)
		}

		if err := applyTemplateVar(&details.Contact.Phone, details.Contact.Phone, tmpl); err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	if details.DemoAccount != nil {
		if err := applyTemplateVar(&details.DemoAccount.Name, details.DemoAccount.Name, tmpl); err != nil {
			errors = multierror.Append(errors, err)
		}

		if err := applyTemplateVar(&details.DemoAccount.Password, details.DemoAccount.Password, tmpl); err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	if err := applyTemplateVar(&details.Notes, details.Notes, tmpl); err != nil {
		errors = multierror.Append(errors, err)
	}

	var attachments = make([]config.File, len(details.Attachments))

	for i, attachment := range details.Attachments {
		if err := applyTemplateVar(&attachment.Path, attachment.Path, tmpl); err != nil {
			errors = multierror.Append(errors, err)
		}

		attachments[i] = attachment
	}

	details.Attachments = attachments

	return errors
}

func applyTemplateVar(v *string, s string, tmpl *template.Template) error {
	new, err := tmpl.Apply(s)
	if err != nil {
		return err
	}

	*v = new

	return nil
}
