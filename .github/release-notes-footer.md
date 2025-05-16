{{ release_notes }}

## Verify Release Artifacts

Releases include a build provenance attestation step that uses the [attest-build-provenance](https://github.com/actions/attest-build-provenance) GitHub action. Verification requires the [GitHub gh CLI](https://cli.github.com/) tool.

The attestation results for this release are at [{{ attestation_url }}]({{ attestation_url }}).  To verify this release, download a release artifact, install `gh` and run this command (change the artifact as needed for your use-case):

```shell
gh attestation verify {{ demo_artifact }} -R eliheady/wtutf
```

See the GitHub [documentation](https://docs.github.com/en/actions/security-for-github-actions/using-artifact-attestations/using-artifact-attestations-to-establish-provenance-for-builds#verifying-artifact-attestations-with-the-github-cli) for alternative verification steps.