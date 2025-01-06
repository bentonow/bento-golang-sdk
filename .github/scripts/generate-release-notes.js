name: Create Release

on:
  workflow_dispatch:
    inputs:
      version:
        description: 'Version number for the release (e.g., v1.0.0) must include v'
        required: true
        type: string
permissions:
  contents: write
  pull-requests: write

jobs:
  create-release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
      pull-requests: read

    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
          ref: main
          token: ${{ secrets.GITHUB_TOKEN }}

      - name: Validate version format
        run: |
          if [[ ! ${{ github.event.inputs.version }} =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            echo "Error: Version must be in format v1.0.0"
            exit 1
          fi
          echo "VERSION=${{ github.event.inputs.version }}" >> $GITHUB_ENV
          echo "PREVIOUS_TAG=$(git describe --tags --abbrev=0)" >> $GITHUB_ENV

      - name: Setup Bun
        uses: oven-sh/setup-bun@v1
        with:
          bun-version: latest

      - name: Generate Release Notes
        id: release_notes
        uses: actions/github-script@v6
        with:
          script: |
            const script = require('./.github/scripts/generate-release-notes.js')
            const notes = await script({github, context})
            core.setOutput('notes', notes)

      - name: Create Tag
        run: |
          git tag ${{ github.event.inputs.version }}
          git push origin ${{ github.event.inputs.version }}

      - name: Create Release
        uses: softprops/action-gh-release@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          tag_name: ${{ github.event.inputs.version }}
          name: Release ${{ github.event.inputs.version }}
          body: ${{ steps.release_notes.outputs.notes }}
          draft: false
          prerelease: false
