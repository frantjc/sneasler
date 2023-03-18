// Handle signals e.g. from ctrl-c'ing `docker run`.
process.on("SIGINT", onSignal);
process.on("SIGTERM", onSignal);

const path = require("path");
const express = require("express");
const compression = require("compression");
const { createRequestHandler } = require("@remix-run/express");

const build = require("./build");
const mode = process.env.NODE_ENV;
const remixRequestHandler = createRequestHandler({
  build,
  mode,
});
const envRemixRequestHandler =
  mode === "development"
    ? (req, res, next) => {
        purgeRequireCache();
        remixRequestHandler(req, res, next);
      }
    : remixRequestHandler;

const app = express();

app.use(compression());

app.disable("x-powered-by");

// Remix fingerprints its assets so we can cache forever.
app.use(
  "/build",
  express.static("public/build", { immutable: true, maxAge: "1y" })
);

// Everything else (like favicon.ico) is cached for an hour. You may want to be
// more aggressive with this caching.
app.use(express.static("public", { maxAge: "1h" }));

app.get("/", envRemixRequestHandler);

app.listen(process.env.PORT || 3000);

const buildDir = path.join(process.cwd(), "./build");
function purgeRequireCache() {
  // Purge require cache on requests for "server side HMR" this won't let
  // you have in-memory objects between requests in development.
  // Alternatively you can set up nodemon/pm2-dev to restart the server on
  // file changes, but then you'll have to reconnect to databases/etc on each
  // change. We prefer the DX of this, so we've included it for you by default.
  for (const key in require.cache) {
    if (key.startsWith(buildDir)) {
      delete require.cache[key];
    }
  }
}

function onSignal() {
  process.exit(0);
}
