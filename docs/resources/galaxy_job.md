# galaxy_job Resource

Execute tools to load data. This is mainly intended for data managers or upload/download tools. Do not use this for data processing!

## Example Usage

```hcl
resource "galaxy_repository" "awkscript" {
  tool_shed = "toolshed.g2.bx.psu.edu"
  owner = "brinkmanlab"
  name = "awkscript"
  changeset_revision = "ceac6ffb3865"
  remove_from_disk = true
}

resource "galaxy_history" "test" {
  name = "test"
}

resource "galaxy_job" "example" {
  depends_on = [galaxy_repository.awkscript]
  tool_id = galaxy_repository.awkscript.tools[0].tool_guid
  history_id = galaxy_history.test.id
  params = {
    "code" = "BEGIN { print \"foo\" }"
  }
  wait_for_completion = true
}
```

## Argument Reference

* `hda` - &lt;List&gt; (Optional) Repeatable block of HDA inputs. Specify the same input id in multiple blocks to provide tool multiple HDAs per input.  
  Arguments:  
  * `id` - &lt;String&gt; (Required) HDA id  
  * `input` - &lt;String&gt; (Required) Input id as described in Galaxy tool wrapper XML  

* `hdca` - &lt;List&gt; (Optional) Repeatable block of HDCA inputs. Specify the same input id in multiple blocks to provide tool multiple HDCAs per input.  
  Arguments:  
  * `id` - &lt;String&gt; (Required) HDCA id  
  * `input` - &lt;String&gt; (Required) Input id as described in Galaxy tool wrapper XML  

* `history_id` - &lt;String&gt; (Required) Id of history where tool outputs are associated  
* `params` - &lt;Map&gt; (Optional) Map of parameter values keyed on input id  
  Element type: String
* `tool_guid` - &lt;String&gt; (Optional) UUID of tool as assigned by Galaxy instance  
  Exactly one of `tool_id` or `tool_guid`  
* `tool_id` - &lt;String&gt; (Optional) Id of the tool to execute in the form `toolshed hostname/repo owner/repo name/tool name/version`  
  Exactly one of `tool_id` or `tool_guid`  
* `wait_for_completion` - &lt;Bool&gt; (Optional) Wait for job to complete before creating dependant resources \[Default: true]  


## Attribute Reference

* `additional_jobs` - &lt;List&gt; If the input parameters spawn multiple jobs, the remaining jobs will be listed here  
  Attributes:  
  * `create_time` - &lt;String&gt; Job creation time  
  * `exit_code` - &lt;Int&gt; Exit code as returned by tool execution  
  * `history_id` - &lt;String&gt; Id of history where tool outputs are associated  
  * `state` - &lt;String&gt; Running state of job  
  * `tool_id` - &lt;String&gt; Id of the tool to execute in the form `toolshed hostname/repo owner/repo name/tool name/version`  
  * `update_time` - &lt;String&gt; Time job state lst updated  

* `create_time` - &lt;String&gt; Job creation time  
* `exit_code` - &lt;Int&gt; Exit code as returned by tool execution  
* `hda` - &lt;List&gt; Repeatable block of HDA inputs. Specify the same input id in multiple blocks to provide tool multiple HDAs per input.  
  Attributes:  
  * `id` - &lt;String&gt; HDA id  
  * `input` - &lt;String&gt; Input id as described in Galaxy tool wrapper XML  

* `hdca` - &lt;List&gt; Repeatable block of HDCA inputs. Specify the same input id in multiple blocks to provide tool multiple HDCAs per input.  
  Attributes:  
  * `id` - &lt;String&gt; HDCA id  
  * `input` - &lt;String&gt; Input id as described in Galaxy tool wrapper XML  

* `history_id` - &lt;String&gt; Id of history where tool outputs are associated  
* `params` - &lt;Map&gt; Map of parameter values keyed on input id  
  Element type: String
* `state` - &lt;String&gt; Running state of job  
* `tool_guid` - &lt;String&gt; UUID of tool as assigned by Galaxy instance  
* `tool_id` - &lt;String&gt; Id of the tool to execute in the form `toolshed hostname/repo owner/repo name/tool name/version`  
* `update_time` - &lt;String&gt; Time job state lst updated  
* `wait_for_completion` - &lt;Bool&gt; Wait for job to complete before creating dependant resources  

