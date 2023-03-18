const pkg = require("pkg");
const fs = require("fs");

function getPkgPlatform() {
  if (
    process.platform === "linux" &&
    fs.readFileSync("/etc/os-release").includes("alpine")
  ) {
    return "alpine";
  }

  return process.platform;
}

pkg.exec([
  "--target",
  `node18-${getPkgPlatform()}-${process.arch}`,
  "--compress",
  "GZip",
  "--output",
  "./pkg/sneasler",
  "./server.js",
]);
