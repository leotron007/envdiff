// Package ci provides CI/CD platform detection and annotation support
// for envdiff.
//
// It detects well-known CI environments (GitHub Actions, GitLab CI,
// CircleCI) by inspecting environment variables, and emits
// platform-native annotations (e.g. GitHub workflow commands) so that
// diff results surface directly in pull-request checks and pipeline
// logs.
//
// Basic usage:
//
//	detector := ci.NewDetector()
//	if detector.IsCI() {
//		annotator := ci.NewAnnotator(detector.Detect(), os.Stdout)
//		annotator.Annotate(diffEntries)
//	}
package ci
