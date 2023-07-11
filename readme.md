# Trello Hello

## Overview

This is a small command line interface client to show me my trello cards in a compact format. Since trello is generic, it does not understand the semantics of particular boards - for example the fact that the "Done" board should be skipped when showing my to do list.

## Features

* Show cards from multiple boards and lists in a compact format
* Only show cards from boards with names meaning they are pending to dos

## Purpose

This is mostly an exercise to learn some Go stuff. But it is useful to me, too. I'm using the "github.com/go-playground/validator/v10" package because we use it at work and I wanted to understand its basic functionality in a simpler context.

I wanted to understand the pattern of a library exporting a struct and the application code that uses the library defines an interface for the subset of the external API actually used in this application, and that interface supports mocking for unit tests.

## Is This Really for General Use?

No, not really. It's too tied to my personal workflow and naming conventions. But you might be able to use it as an example when making something tailored to your workflow.

## Install and Run

- clone this repo
- Follow [the trello docs](https://developer.atlassian.com/cloud/trello/guides/rest-api/authorization/) to get an API Key and a token
- set your API key as environment variable `KEY`
- set your API token as environment variable `TOKEN`
- Run with `go run .` or build a binary with `go build .` then `./trello-hello`
- Copy the binary into your PATH if you like `sudo cp trello-hello /usr/local/bin`

## Sample Output

```
📋 Trip Planning
  📃Doing
    🪧Shop for travel umbrella
  📃To Do Soon
    🪧Passport Renewal
    🪧International Driver's Permit
  📃To Do
    🪧Get a phrasebook
📋Personal
  📃Doing
    🪧Clean out garage
  📃To Do
  📃To Do: Low Priority
    🪧Fix fence door hinge
    🪧Send thank you card to Walter
```
