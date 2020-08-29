package static

const ExampleConfig = `# This is an example applereleaser.yaml file with some sane defaults.
name: My Project
apps:
  ProjectApp:
    id: com.project.ProjectApp
	localizations: ~
    versions:
      platform: iOS 
      enablePhasedRelease: true
      localizations:
        en-US:
          whatsNew: ''
      idfaDeclaration: ~
`
