---
page_title: Install a repository and run its data manager
---
# Install a repository and run its data manager

```hcl
# Install the tool repository, including its data manager
resource "galaxy_repository" "rgi" {
  tool_shed = "toolshed.g2.bx.psu.edu"
  owner = "card"
  name = "rgi"
}

# Create a history to store the data managers output
resource "galaxy_history" "terraform" {
  name = "terraform"
}

# Run the data manager to load the tools data
resource "galaxy_job" "load_data" {
  depends_on = [galaxy_repository.rgi]  # An issue with the Galaxy API requires this. https://github.com/galaxyproject/galaxy/issues/10378
  tool_id = "toolshed.g2.bx.psu.edu/repos/card/rgi/rgi_database_builder/1.0.0"
  history_id = galaxy_history.terraform.id
  params = {
    "name" = "5.1.1"
    "url" = "https://card.mcmaster.ca/download/1/software-v5.1.1.tar.bz2"
  }
  wait_for_completion = true
}
```