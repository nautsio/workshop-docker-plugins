# Workshop Build your own Docker plugin

## How to build/preview the slides
- Build your site using the `jekyll/jekyll:pages` Docker image

```
docker run --rm --name=jekyll --volume=$(pwd):/srv/jekyll -it -p 4000:4000 jekyll/jekyll:pages jekyll s
```

## Notes
- Reveal.js and the stying are pulled from the shared [cdn.nauts.io](https://github.com/nautsio/cdn) repo
- Markdown separators:
 - New slide: `^\n!SLIDE`
 - New vertical sub-slide: `^\n!SUB`
 - Presenter notes: `^\n!NOTE`
