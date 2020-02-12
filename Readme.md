#goProjectValidator

GPV is a CLI tool that takes three primary arguments:

1. Validation JSON : The validation JSON that defines the entire specification
2. Directory in which test JSON is placed
3. Directory into which you would like the three markdown documents to be written 

## Validation Json
```json
{
  "project" : "Babylon",
  "scope" : "Provide a tool for programmatically reading in validation files as JSON to establish project requirements, parse JSON Results from tests, and build validation documents",
  "markdown" : [
    {
      "source" : "https://raw.githubusercontent.com/metrumresearchgroup/babylontest/master/markdown/expectations.md"
    }
  ],
  "stories" : [
    {
      "name" : "Babylon should allow users to run NonMem locally",
      "tags" : [
        "TestBabylonCompletesLocalExecution",
        "TestBabylonParallelExecution"
      ],
      "risk" : "high",
      "markdown" : [
        {
          "source" : "https://raw.githubusercontent.com/metrumresearchgroup/babylontest/master/markdown/stories/run_jobs_locally.md"
        }
      ]
    }
  ]
}
```

The structure is pretty simple. The Top level object defines the specification. The markdown key contains an array of 
"source" objects. These are http / https accessible MD files we want to render into the output files. It was easier
and more sane for an editing approach to do it this way than to have people editing overly large and complex JSON files.

Underneath the Specification is the `stories` key, which contains user issues expressed as stories. Each "tag" listed here
should correlate to a test name in the go test output. GPV will stitch tests matching those names into the story such that
it's self-contained. Each story also has a markdown section indicating any remote markdowns that should be rendered in