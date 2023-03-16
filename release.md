## Release

Builds and releases are performed by [goreleaser](https://goreleaser.com/) using
a Github Action. The releases are not yet signed.

### How-to

The action is configured to run when a tag is created. To create a release:

1. Create a tag and push it:
    ```
    git tag -a v0.1.0 -m "v0.1.0"
    git push origin v0.1.0
    ```
1. Check [progress](https://github.com/jeremychase/wxgate/actions/workflows/release.yml).