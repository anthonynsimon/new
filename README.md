# new

Create configurable project templates.

```
new -h

> render custom templates
> 
> Usage:
>   new [template path] [destination path] [flags]
> 
> Examples:
> new templates/team-project .
> 
> Flags:
>   -h, --help      help for new
>       --version   version for new
```

### Installation

Assuming you have [Go installed](https://golang.org/):

```bash
go get github.com/anthonynsimon/new
```


### Usage:

1. Create a file or folder with your template.
2. Add a `.new.yaml` file:

```yaml
# examples/example-project/.new.yaml

version: '1'

description: Example codebase project template

params:
    - name: name
      prompt: What's the name for the project?
      required: true
      kind: string
    
    - name: deployment
      prompt: What kind of deployment should be included? 
      kind: enum
      required: true
      enum:
        - kubernetes
        - docker
        - ec2
```

3. Template it at the destination directory.

```
new examples/example-project

> Example codebase project template

> What's the name for the project? my-new-project
> What kind of deployment should be included? ✔ kubernetes

> Rendering my-new-project
> Rendering my-new-project/Pipfile
> Rendering my-new-project/Pipfile.lock
> Rendering my-new-project/README.md
> Rendering my-new-project/bin
> Rendering my-new-project/bin/.gitkeep
> Rendering my-new-project/conf
> Rendering my-new-project/conf/local
> Rendering my-new-project/conf/local/app.yaml
> Rendering my-new-project/deploy
> Rendering my-new-project/deploy/.gitkeep
> Rendering my-new-project/deploy/deploy.yml
> Rendering my-new-project/docs
> Rendering my-new-project/docs/.gitkeep
> Rendering my-new-project/src
> Rendering my-new-project/src/.gitkeep
```

Done!



## Contributing
To get started take a look at our [contribution guidelines](CONTRIBUTING.md) for information about our process and where you can fit in!


## Changelog
You can check our [changelog](CHANGELOG.md) for information about releases.


## License
See [LICENSE](LICENSE).

