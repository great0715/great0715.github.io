# Jorge Lopez Calero · Senior Software Engineer

**Frontend · Backend · DevOps**

[Live portfolio](https://great0715.github.io/) · [GitHub profile](https://github.com/great0715)

A responsive, no-build portfolio with interactive demonstrations of frontend interfaces, operational dashboards, web editing, and delivery workflows. The **demos are original examples**, not past client work.

## Featured demos

| Demo | What's interactive | Technology |
| --- | --- | --- |
| [Atelier Storefront](https://great0715.github.io/demos/atelier.html) | Search, filters, sorting, cart | Vanilla JavaScript, CSS, local storage |
| [Signal Dashboard](https://great0715.github.io/demos/signal.html) | Sample metrics, date ranges, CSV export | JavaScript, SVG |
| [Canvas Studio](https://great0715.github.io/demos/canvas.html) | Content, layout, and appearance controls | JavaScript, HTML, CSS |
| [Release Lab](https://great0715.github.io/demos/release-lab.html) | Simulated pipeline execution, failures, rollback | JavaScript, static state machine |

## Backend reference

[`examples/go-api/`](examples/go-api/) is a **real, locally runnable Go API** with unit tests and a Dockerfile. The Release Lab browser demo is **a simulation and is not connected to this API**. GitHub Pages hosts static files; it does not execute the Go server.

```bash
cd examples/go-api
go test ./...
go run .
# Test in another terminal: curl http://localhost:8080/healthz
```

## Run locally

```bash
python -m http.server 8000
```

Open `http://localhost:8000/`. No JavaScript package installation or build step is required.

## Quality checks

GitHub Actions checks local HTML links and runs the Go example's tests and `go vet`. To run the link checker locally:

```bash
python scripts/check_site.py
```

## Hosting

This repository is designed for **GitHub Pages** at `https://great0715.github.io/`:

1. Make the repository public and ensure `index.html` is at the repository root.
2. Open **Settings → Pages**.
3. Under **Build and deployment**, set **Source** to **Deploy from a branch**.
4. Select **main** and **/(root)**, then save.

The `.nojekyll` file tells Pages to serve the existing static site without a Jekyll build.

## Notes

- Review biographical details, dates, technologies, and the public contact email before publishing.
- All storefront, dashboard, and pipeline data is demonstration data.
- Private client code is not included.
