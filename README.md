## Installation

If you have a Go development environment, installing this command is as simple as typing the following line in your command line:

```
go install github.com/parhamdoustdar/godium
```

## Usage

### Obtaining an Integration Token

Before you can do anything with this application, you need to generate an integration token on the [settings page in Medium](https://medium.com/me/settings). Once that is done, use the `set-token` as follows to set that token to be used by godium:

```
godium set-token <token>
```

Alternatively, if you try running the `info` or `publish` commands, godium will try to automatically open the settings page for you. If that is successful, you can do the previous steps; namely generate an integration token and add it with the command above to godium.

### Getting Your Information

To double-check that you have the right integration token set, you can use the `info` command to get the information for the current user:

```
godium info
```

If everything goes fine, you should see the current user's information. If not, you will probably receive an error like this:

```
medium: Token was invalid. (6003)
```

If you get this error, this means something went wrong with the copying process. To fix this, open the file godium uses to store your integration token (you get the file path when you run `godium set-token <token>`), and compare this with what you see on your settings page.

### Publishing to Medium

Publishing to Medium is a simple, three-step process:

- Write your article in [markdown syntax][].
- Provide some information about the article using the frontmatter (explained below).
- Use `godium publish <filename>` to publish this to Medium. This command will publish the article with the `draft` status and open your browser window to make any final changes and confirm by clicking the `Publish` button to finalize.

As mentioned before, godium requires some information about your article to be present as a [frontmatter][]. Here is a list of what you'll need:

- title: self-explanatory -- the title of your article
- tags: a list of tags in yaml format (refer to the example below)

#### Example Article

To better illustrate the publishing process, here is a well-formed article.

```
---
title: Example Post
tags: [tech, development, markdown]
---

Here is some post, with a link to [Google][].

[Google]: https://www.google.com
```

Much simpler than dealing with Medium's editor, don't you think?

[markdown syntax]: https://daringfireball.net/projects/markdown/syntax
[frontmatter]: http://assemble.io/docs/YAML-front-matter.html