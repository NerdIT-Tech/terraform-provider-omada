package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// bddSiteResourceName is the Terraform resource address the site.feature
// scenarios exercise. Each scenario runs against its own isolated
// resource.Test call, so reusing one address across scenarios is safe.
const bddSiteResourceName = "omada_site.test"

// siteScenarioState accumulates the Terraform plan/apply steps a scenario's
// Given/When/Then steps describe. It is built up without touching the
// controller; the scenario's After hook runs the accumulated steps through
// resource.Test in a single acceptance test lifecycle, the same way
// TestAccSiteResource does.
type siteScenarioState struct {
	steps       []resource.TestStep
	currentName string
}

type siteScenarioStateKey struct{}

func siteScenarioStateFrom(ctx context.Context) *siteScenarioState {
	state, _ := ctx.Value(siteScenarioStateKey{}).(*siteScenarioState)
	return state
}

func (s *siteScenarioState) lastStep() (*resource.TestStep, error) {
	if len(s.steps) == 0 {
		return nil, fmt.Errorf("no Terraform step has been recorded yet; a When step must run first")
	}
	return &s.steps[len(s.steps)-1], nil
}

func appendCheck(existing resource.TestCheckFunc, add resource.TestCheckFunc) resource.TestCheckFunc {
	if existing == nil {
		return add
	}
	return resource.ComposeAggregateTestCheckFunc(existing, add)
}

func bddGivenConfiguredProvider(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, siteScenarioStateKey{}, &siteScenarioState{}), nil
}

func bddWhenICreateASiteNamed(ctx context.Context, name string) (context.Context, error) {
	state := siteScenarioStateFrom(ctx)
	if state == nil {
		return ctx, fmt.Errorf("scenario state missing; the Given step must run first")
	}
	state.currentName = name
	state.steps = append(state.steps, resource.TestStep{
		Config: testAccSiteResourceConfig(name),
	})
	return ctx, nil
}

func bddWhenIRenameTheSiteTo(ctx context.Context, name string) (context.Context, error) {
	state := siteScenarioStateFrom(ctx)
	if state == nil {
		return ctx, fmt.Errorf("scenario state missing; the Given step must run first")
	}
	state.currentName = name
	state.steps = append(state.steps, resource.TestStep{
		Config: testAccSiteResourceConfig(name),
	})
	return ctx, nil
}

func bddWhenIImportTheSiteByItsID(ctx context.Context) (context.Context, error) {
	state := siteScenarioStateFrom(ctx)
	if state == nil {
		return ctx, fmt.Errorf("scenario state missing; the Given step must run first")
	}
	expectedName := state.currentName
	state.steps = append(state.steps, resource.TestStep{
		ResourceName:      bddSiteResourceName,
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateCheck: func(states []*terraform.InstanceState) error {
			if len(states) != 1 {
				return fmt.Errorf("expected 1 imported instance state, got %d", len(states))
			}
			if got := states[0].Attributes["name"]; got != expectedName {
				return fmt.Errorf("imported site name = %q, want %q", got, expectedName)
			}
			return nil
		},
	})
	return ctx, nil
}

func bddThenTheSiteShouldExist(ctx context.Context) (context.Context, error) {
	state := siteScenarioStateFrom(ctx)
	if state == nil {
		return ctx, fmt.Errorf("scenario state missing; the Given step must run first")
	}
	step, err := state.lastStep()
	if err != nil {
		return ctx, err
	}
	step.Check = appendCheck(step.Check, resource.TestCheckResourceAttrSet(bddSiteResourceName, "id"))
	return ctx, nil
}

func bddThenTheSitesNameShouldBe(ctx context.Context, name string) (context.Context, error) {
	state := siteScenarioStateFrom(ctx)
	if state == nil {
		return ctx, fmt.Errorf("scenario state missing; the Given step must run first")
	}
	step, err := state.lastStep()
	if err != nil {
		return ctx, err
	}
	step.Check = appendCheck(step.Check, resource.TestCheckResourceAttr(bddSiteResourceName, "name", name))
	return ctx, nil
}

func bddThenTheImportedSitesNameShouldBe(ctx context.Context, name string) (context.Context, error) {
	state := siteScenarioStateFrom(ctx)
	if state == nil {
		return ctx, fmt.Errorf("scenario state missing; the Given step must run first")
	}
	if state.currentName != name {
		return ctx, fmt.Errorf("scenario inconsistency: expected imported site name %q, but the site created/renamed earlier in the scenario is named %q", name, state.currentName)
	}
	return ctx, nil
}

// TestFeaturesSite runs the Gherkin scenarios in features/site.feature
// against a real Omada Controller, gated by the same OMADA_* environment
// variables and TF_ACC as TestAccSiteResource.
func TestFeaturesSite(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			s.Given(`^a configured Omada provider$`, bddGivenConfiguredProvider)
			s.When(`^I create a site named "([^"]*)"$`, bddWhenICreateASiteNamed)
			s.When(`^I rename the site to "([^"]*)"$`, bddWhenIRenameTheSiteTo)
			s.When(`^I import the site by its ID$`, bddWhenIImportTheSiteByItsID)
			s.Then(`^the site should exist$`, bddThenTheSiteShouldExist)
			s.Then(`^the site's name should be "([^"]*)"$`, bddThenTheSitesNameShouldBe)
			s.Then(`^the imported site's name should be "([^"]*)"$`, bddThenTheImportedSitesNameShouldBe)

			s.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				if err != nil {
					return ctx, err
				}
				state := siteScenarioStateFrom(ctx)
				if state == nil || len(state.steps) == 0 {
					return ctx, fmt.Errorf("scenario %q recorded no Terraform steps to run", sc.Name)
				}

				resource.Test(t, resource.TestCase{
					PreCheck:                 func() { testAccPreCheck(t) },
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					Steps:                    state.steps,
				})
				return ctx, nil
			})
		},
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features/site.feature"},
			// TestingT is deliberately not set: godog would run each
			// scenario as a t.Run subtest, but the After hook below drives
			// resource.Test using the outer *testing.T, and a subtest
			// goroutine calling Skip/FailNow on its parent T panics ("may
			// have called FailNow on a parent test"). Running scenarios in
			// TestFeaturesSite's own goroutine keeps that call valid; godog's
			// pretty formatter still reports pass/fail per scenario.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run site.feature scenarios")
	}
}
