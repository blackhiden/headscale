# Topology

The Topology page (`/topology` in Headplane) renders your tailnet as an
interactive graph.

## What is shown

- **Nodes** — each machine is a circle, colored by its owner (user or
  "Tag-owned").
- **Subnet routers** — advertised subnet routes appear as small nodes below
  their router, connected by a dashed line.
- **Exit nodes** — highlighted with an amber ring.
- **Status** — a green dot marks online nodes; expired nodes are dimmed.

## Interacting with the graph

- **Pan** — click and drag the canvas.
- **Zoom** — use the mouse wheel or the `+` / `−` buttons in the corner.
- **Reset** — the `⟲` button restores the default view.
- **Inspect** — hover a node for a tooltip with its owner, online state, and
  role; click a machine node to open its details page.

## Filters

Use the filter tabs to show only **Online**, **Offline**, or **Expired**
machines, and the owner dropdown to focus on a single user or tag group.