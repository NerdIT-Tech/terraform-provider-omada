Feature: Manage Omada sites
  As a network engineer
  I want to create a site
  So that I can manage devices and clients for a given location

  Scenario: Create a site
    Given a configured Omada provider
    When I create a site named "tf-bdd-create"
    Then the site should exist
    And the site's name should be "tf-bdd-create"

  Scenario: Rename an existing site
    Given a configured Omada provider
    When I create a site named "tf-bdd-rename-before"
    And I rename the site to "tf-bdd-rename-after"
    Then the site should exist
    And the site's name should be "tf-bdd-rename-after"

  Scenario: Import an existing site by ID
    Given a configured Omada provider
    When I create a site named "tf-bdd-import"
    And I import the site by its ID
    Then the imported site's name should be "tf-bdd-import"
