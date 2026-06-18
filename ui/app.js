import { wasmReady } from "./loader.js";

const sampleProto = `syntax = "proto3";
package examples.trading.v1;

import "google/protobuf/timestamp.proto";

enum Side {
  SIDE_UNSPECIFIED = 0;
  SIDE_BUY = 1;
  SIDE_SELL = 2;
}

message Order {
  string id = 1;
  optional string venue = 2;
  Side side = 3;
  google.protobuf.Timestamp created_at = 4;
  map<string, string> tags = 5;
}
`;

const input = document.querySelector("#proto-input");
const output = document.querySelector("#xml-output");
const status = document.querySelector("#status");
const convertButton = document.querySelector("#convert-button");
const copyButton = document.querySelector("#copy-button");
const sampleButton = document.querySelector("#sample-button");
const themeToggle = document.querySelector("#theme-toggle");

let requestID = 0;
let convertTimer = 0;

const currentTheme = savedTheme();
const protoEditor = CodeMirror.fromTextArea(input, {
  mode: "protobuf",
  theme: editorTheme(currentTheme),
  lineNumbers: true,
  indentUnit: 2,
  tabSize: 2,
  lineWrapping: false,
  extraKeys: {
    Tab(editor) {
      editor.replaceSelection("\t");
    },
  },
});
const xmlEditor = CodeMirror.fromTextArea(output, {
  mode: "xml",
  theme: editorTheme(currentTheme),
  lineNumbers: true,
  indentUnit: 2,
  tabSize: 2,
  lineWrapping: false,
  readOnly: true,
});

applyTheme(currentTheme);
protoEditor.setValue(sampleProto);

convertButton.addEventListener("click", () => {
  convert();
});

copyButton.addEventListener("click", async () => {
  const xml = xmlEditor.getValue();
  if (!xml) {
    setStatus("Nothing to copy", "warn");
    return;
  }

  try {
    await navigator.clipboard.writeText(xml);
    setStatus("Copied XML", "ok");
  } catch {
    xmlEditor.focus();
    xmlEditor.execCommand("selectAll");
    document.execCommand("copy");
    setStatus("Copied XML", "ok");
  }
});

sampleButton.addEventListener("click", () => {
  protoEditor.setValue(sampleProto);
  convert();
});

themeToggle.addEventListener("change", () => {
  const theme = themeToggle.checked ? "dark" : "light";
  applyTheme(theme);
  try {
    localStorage.setItem("theme", theme);
  } catch {}
});

protoEditor.on("change", () => {
  clearTimeout(convertTimer);
  setStatus("Editing", "idle");
  convertTimer = window.setTimeout(convert, 450);
});

wasmReady
  .then(() => {
    setStatus("Ready", "ok");
    convert();
  })
  .catch((error) => {
    xmlEditor.setValue("");
    setStatus(error.message || String(error), "error");
  });

async function convert() {
  const currentRequestID = ++requestID;
  setBusy(true);
  setStatus("Generating", "idle");

  try {
    await wasmReady;
    const xml = await window.generateSBE(protoEditor.getValue());
    if (currentRequestID !== requestID) {
      return;
    }
    xmlEditor.setValue(xml);
    setStatus("Generated output.xml", "ok");
  } catch (error) {
    if (currentRequestID !== requestID) {
      return;
    }
    xmlEditor.setValue("");
    setStatus(error.message || String(error), "error");
  } finally {
    if (currentRequestID === requestID) {
      setBusy(false);
    }
  }
}

function setBusy(isBusy) {
  convertButton.disabled = isBusy;
  convertButton.textContent = isBusy ? "Generating" : "Generate";
}

function setStatus(message, tone) {
  status.textContent = message;
  status.dataset.tone = tone;
}

function savedTheme() {
  try {
    const theme = localStorage.getItem("theme");
    if (theme === "light" || theme === "dark") {
      return theme;
    }
  } catch {
    return "dark";
  }
  return "dark";
}

function applyTheme(theme) {
  document.documentElement.dataset.theme = theme;
  themeToggle.checked = theme === "dark";
  protoEditor.setOption("theme", editorTheme(theme));
  xmlEditor.setOption("theme", editorTheme(theme));
}

function editorTheme(theme) {
  return theme === "light" ? "sbexml-light" : "sbexml-dark";
}
